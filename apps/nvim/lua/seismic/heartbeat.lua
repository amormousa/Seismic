local M = {}

local config = require('seismic.config')
local detector = require('seismic.detector')
local queue = require('seismic.queue')

local last_heartbeat_time = 0
local last_file = ''
local has_shown_invalid_key_warning = false

local TWO_MINUTES = 2 * 60 * 1000 -- in milliseconds

---@class HeartbeatPayload
---@field file string
---@field project string
---@field language string
---@field editor string
---@field branch string|nil
---@field os string
---@field machine string
---@field lines number
---@field cursorLine number
---@field timezone nil
---@field time number

--- Builds the heartbeat payload.
---@return HeartbeatPayload
local function build_payload()
  local bufnr = vim.api.nvim_get_current_buf()
  local file = vim.api.nvim_buf_get_name(bufnr)
  local filetype = vim.api.nvim_get_buf_option(bufnr, 'filetype')
  local line_count = vim.api.nvim_buf_line_count(bufnr)
  local cursor_pos = vim.api.nvim_win_get_cursor(0)

  return {
    file = file,
    project = detector.detect_project(bufnr),
    language = filetype,
    editor = 'neovim',
    branch = detector.detect_branch(),
    os = detector.detect_os(),
    machine = detector.detect_machine(),
    lines = line_count,
    cursorLine = cursor_pos[1],
    timezone = detector.detect_timezone(),
    time = vim.loop.now(),
  }
end

--- Sends the heartbeat to the API.
---@param payload HeartbeatPayload
local function send(payload)
  local conf = config.get()
  local api_key = conf.api_key
  local api_url = conf.api_url

  local body = vim.fn.json_encode(payload)

  local cmd = {
    'curl',
    '-s', -- silent
    '-w', -- write-out
    '%{http_code}',
    '-X',
    'POST',
    '-H',
    'Content-Type: application/json',
    '-H',
    'Authorization: Bearer ' .. api_key,
    '--data',
    body,
    api_url .. '/api/heartbeat',
  }

  vim.fn.jobstart(cmd, {
    on_exit = function(_, data, _)
      local status_code = tonumber(data[1])
      if status_code == 401 then
        if not has_shown_invalid_key_warning then
          vim.schedule(function()
            vim.notify('Seismic: Invalid API key. Run ":SeismicSetApiKey" to update it.', vim.log.levels.WARN)
          end)
          has_shown_invalid_key_warning = true
        end
      elseif status_code and status_code >= 200 and status_code < 300 then
        -- We're online, so flush the queue
        queue.flush()
      else
        -- Network error, queue the heartbeat
        queue.enqueue(payload)
      end
    end,
  })
end

--- Handles an editor activity event.
--- Decides whether to send a heartbeat based on throttling rules.
---@param forced boolean If true, sends the heartbeat regardless of the 2-minute rule.
function M.handle_activity(forced)
  local conf = config.get()
  if not conf.enabled then
    return
  end
  if not conf.api_key or conf.api_key == '' then
    if not has_shown_invalid_key_warning then
      vim.schedule(function()
        vim.notify('Seismic: API key not set. Run ":SeismicSetApiKey".', vim.log.levels.WARN)
      end)
      has_shown_invalid_key_warning = true
    end
    return
  end

  local bufnr = vim.api.nvim_get_current_buf()
  if not detector.should_track(bufnr) then
    return
  end

  local now = vim.loop.now()
  local current_file = vim.api.nvim_buf_get_name(bufnr)
  local file_changed = current_file ~= last_file

  if not forced and not file_changed and (now - last_heartbeat_time) < TWO_MINUTES then
    return -- too soon
  end

  last_heartbeat_time = now
  last_file = current_file

  local payload = build_payload()
  send(payload)
end

--- Flushes the offline heartbeat queue.
function M.flush_queue()
  queue.flush()
end

return M
