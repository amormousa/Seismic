# Seismic for Neovim

Your code speaks. We listen. Seismic is a code-time tracking plugin that sends
heartbeats to the Seismic API.

## Installation

Install with your favorite plugin manager.

### [lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
{
  'Majoramari/Seismic',
  dir = 'apps/nvim',
  config = function()
    require('seismic').setup({
      -- your config here
    })
  end,
}
```

### [packer.nvim](https://github.com/wbthomason/packer.nvim)

```lua
use {
  'Majoramari/Seismic',
  run = 'cd apps/nvim && make', -- Or manually build
  config = function()
    require('seismic').setup({
      -- your config here
    })
  end,
}
```

## Setup

The only required configuration is your `api_key`. You can get one from
[seismic.icu/settings](https://seismic.icu/settings).

You can either set it in the `setup` function:

```lua
require('seismic').setup({
  api_key = 'your-api-key-goes-here',
})
```

Or you can use the `:SeismicSetApiKey` command after installation. The key
will be persisted locally so you don't have to set it every time.

## Configuration

| Option    | Type    | Default                                                    | Description                                 |
| --------- | ------- | ---------------------------------------------------------- | ------------------------------------------- |
| `api_key` | string  | `''`                                                       | Your Seismic API key.                       |
| `api_url` | string  | `https://correct-wolverine-majoramari-6049fd71.koyeb.app` | The URL of the Seismic API.                 |
| `enabled` | boolean | `true`                                                     | Whether to enable or disable time tracking. |

## Commands

-   `:SeismicSetApiKey`: Prompts you to enter and save your API key.
-   `:SeismicOpenDashboard`: Opens your Seismic dashboard in your browser.
-   `:SeismicEnable`: Enables time tracking.
-   `:SeismicDisable`: Disables time tracking.
