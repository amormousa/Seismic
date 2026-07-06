local M = {}

--- Checks if a buffer should be tracked.
---@param bufnr number The buffer number.
---@return boolean
function M.should_track(bufnr)
  -- For now, track everything. In the future, we can add checks
  -- for things like filetype, buftype, etc.
  local buftype = vim.api.nvim_get_buf_option(bufnr, 'buftype')
  if buftype ~= '' and buftype ~= 'acwrite' then
    return false
  end
  return true
end

--- Detects the project name from the current buffer.
--- Traverses up from the file's directory to find a .git folder.
---@param bufnr number The buffer number.
---@return string The project name, or 'unknown'.
function M.detect_project(bufnr)
  local file_path = vim.api.nvim_buf_get_name(bufnr)
  local dir = vim.fn.fnamemodify(file_path, ':h')
  local project_root = vim.fn.finddir('.git', dir .. ';')
  if project_root ~= '' and project_root ~= nil then
    return vim.fn.fnamemodify(project_root, ':h:t')
  end
  return 'unknown'
end

--- Detects the current git branch.
---@return string|nil The branch name, or nil if not in a git repo.
function M.detect_branch()
  if vim.fn.executable('git') == 1 then
    local branch = vim.fn.system('git rev-parse --abbrev-ref HEAD')
    if vim.v.shell_error == 0 then
      return vim.fn.trim(branch)
    end
  end
  return nil
end

--- Detects the operating system.
---@return string
function M.detect_os()
  return vim.loop.os_uname().sysname
end

--- Detects the machine hostname.
---@return string
function M.detect_machine()
  return vim.loop.os_uname().nodename
end

--- Detects the user's timezone.
---@return nil
function M.detect_timezone()
  -- This is tricky in a cross-platform way without external dependencies.
  -- For now, we'll leave it out. The backend can infer it from the IP.
  return nil
end

return M
