local M = {}

---@class SeismicConfig
---@field api_key string
---@field api_url string
---@field enabled boolean

local defaults = {
  api_key = '',
  api_url = 'https://correct-wolverine-majoramari-6049fd71.koyeb.app',
  enabled = true,
  -- status_bar_enabled = true, -- Not implemented yet
}

---@type SeismicConfig
local config = {}

--- Merges the user-provided configuration with the defaults.
---@param user_config SeismicConfig
function M.setup(user_config)
  config = vim.tbl_deep_extend('force', defaults, user_config or {})
end

--- Returns the current configuration.
---@return SeismicConfig
function M.get()
  return config
end

--- Sets and persists the API key.
---@param key string
function M.set_api_key(key)
  config.api_key = key
  -- Persist the key to a file so the user doesn't have to
  -- enter it every time.
  local data_dir = vim.fn.stdpath('data')
  local key_file = data_dir .. '/seismic_api_key'
  local file = io.open(key_file, 'w')
  if file then
    file:write(key)
    file:close()
  end
end

-- On startup, try to load the API key from the persisted file
local function load_api_key()
  local data_dir = vim.fn.stdpath('data')
  local key_file = data_dir .. '/seismic_api_key'
  local file = io.open(key_file, 'r')
  if file then
    local key = file:read('*a')
    file:close()
    if key and key ~= '' then
      config.api_key = key:gsub('^%s*(.-)%s*$', '%1') -- trim whitespace
    end
  end
end

-- Initialize the config with defaults
M.setup({})
load_api_key()

return M
