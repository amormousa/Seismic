local M = {}

local config = require('seismic.config')
local queue_file = vim.fn.stdpath('data') .. '/seismic_queue.json'

--- Reads the queue from the filesystem.
---@return HeartbeatPayload[]
local function read_queue()
  local file = io.open(queue_file, 'r')
  if not file then
    return {}
  end
  local content = file:read('*a')
  file:close()
  local ok, decoded = pcall(vim.fn.json_decode, content)
  if not ok or type(decoded) ~= 'table' then
    return {}
  end
  return decoded
end

--- Writes the queue to the filesystem.
---@param queue_data HeartbeatPayload[]
local function write_queue(queue_data)
  local encoded = vim.fn.json_encode(queue_data)
  local file = io.open(queue_file, 'w')
  if file then
    file:write(encoded)
    file:close()
  end
end

--- Adds a payload to the offline queue.
---@param payload HeartbeatPayload
function M.enqueue(payload)
  local queue_data = read_queue()
  table.insert(queue_data, payload)
  write_queue(queue_data)
end

--- Attempts to send all heartbeats in the offline queue.
function M.flush()
  local conf = config.get()
  local api_key = conf.api_key
  local api_url = conf.api_url

  if not api_key or api_key == '' then
    return
  end

  local queue_data = read_queue()
  if #queue_data == 0 then
    return
  end

  local remaining_payloads = {}

  for _, payload in ipairs(queue_data) do
    local body = vim.fn.json_encode(payload)
    local cmd = {
      'curl',
      '-s',
      '-w',
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

    local result = vim.fn.system(cmd)
    local status_code = tonumber(result)

    if not status_code or status_code < 200 or status_code >= 300 then
      -- If it's not a success, put it back in the queue
      table.insert(remaining_payloads, payload)
    end
  end

  write_queue(remaining_payloads)
end

return M
