# Seismic

Track your coding time from VS Code and view your activity in Seismic.

## Features

- Tracks coding activity while you work in real files.
- Sends heartbeat events to Seismic with project, language, branch, cursor, line count, OS, machine, and timezone metadata.
- Queues failed heartbeat requests and retries them later.
- Shows today's coding time in the VS Code status bar.
- Lets you enable or disable tracking from the Command Palette.

## Setup

1. Install the extension.
2. Open the Command Palette.
3. Run `Seismic: Set API Key`.
4. Paste your API key from `https://seismic.icu/settings`.

## Commands

- `Seismic: Set API Key`
- `Seismic: Open Dashboard`
- `Seismic: Enable Tracking`
- `Seismic: Disable Tracking`

## Settings

- `seismic.apiKey`: API key used to authenticate heartbeat requests.
- `seismic.apiUrl`: Seismic API server URL.
- `seismic.enabled`: Enables or disables tracking.
- `seismic.statusBarEnabled`: Shows or hides today's coding time in the status bar.