local M = {}

local config = require('seismic.config')
local heartbeat = require('seismic.heartbeat')
local commands = require('seismic.commands')

--- The main setup function for the plugin.
---@param user_config SeismicConfig The user's configuration, to be merged with the defaults.
function M.setup(user_config)
  config.setup(user_config)
  commands.setup()

  -- Create an augroup to ensure the autocmds are cleared on reload
  local group = vim.api.nvim_create_augroup('Seismic', { clear = true })

  -- Event: text changed in insert mode
  vim.api.nvim_create_autocmd('TextChangedI', {
    group = group,
    pattern = '*',
    callback = function()
      heartbeat.handle_activity(false)
    end,
  })

  -- Event: text changed in normal mode (e.g. paste)
  vim.api.nvim_create_autocmd('TextChanged', {
    group = group,
    pattern = '*',
    callback = function()
      heartbeat.handle_activity(false)
    end,
  })

  -- Event: switched to a new buffer
  vim.api.nvim_create_autocmd('BufEnter', {
    group = group,
    pattern = '*',
    callback = function()
      heartbeat.handle_activity(true)
    end,
  })

  -- Event: saved a buffer
  vim.api.nvim_create_autocmd('BufWritePost', {
    group = group,
    pattern = '*',
    callback = function()
      heartbeat.handle_activity(true)
    end,
  })

  -- Event: gained focus
  vim.api.nvim_create_autocmd('FocusGained', {
    group = group,
    pattern = '*',
    callback = function()
      heartbeat.handle_activity(false)
    end,
  })

  -- Set up a timer to flush the queue periodically
  vim.fn.timer_start(5 * 60 * 1000, function()
    heartbeat.flush_queue()
  end, { ['repeat'] = -1 })

  if config.get().enabled then
    vim.notify('Seismic is tracking your coding activity.', vim.log.levels.INFO)
  end
end

M.config = M.setup

return M
