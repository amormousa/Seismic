local M = {}

local config = require('seismic.config')

--- Sets up the user-facing commands.
function M.setup()
  --- Command to set the Seismic API key.
  vim.api.nvim_create_user_command('SeismicSetApiKey', function()
    vim.ui.input({ prompt = 'Enter your Seismic API key:' }, function(key)
      if key and key ~= '' then
        config.set_api_key(key)
        vim.notify('Seismic: API key saved!', vim.log.levels.INFO)
      end
    end)
  end, {})

  --- Command to open the Seismic dashboard.
  vim.api.nvim_create_user_command('SeismicOpenDashboard', function()
    local url = 'https://seismic.icu/dashboard'
    local os = vim.loop.os_uname().sysname
    local cmd
    if os == 'Darwin' then
      cmd = 'open'
    elseif os == 'Linux' then
      cmd = 'xdg-open'
    elseif os == 'Windows_NT' then
      cmd = 'start'
    else
      vim.notify('Seismic: Unsupported OS for opening dashboard.', vim.log.levels.WARN)
      return
    end
    vim.fn.system(cmd .. ' ' .. url)
  end, {})

  --- Command to enable Seismic tracking.
  vim.api.nvim_create_user_command('SeismicEnable', function()
    config.get().enabled = true
    vim.notify('Seismic: Tracking enabled', vim.log.levels.INFO)
  end, {})

  --- Command to disable Seismic tracking.
  vim.api.nvim_create_user_command('SeismicDisable', function()
    config.get().enabled = false
    vim.notify('Seismic: Tracking disabled', vim.log.levels.INFO)
  end, {})
end

return M
