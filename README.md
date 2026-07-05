# Seismic

Seismic is a developer time tracking platform. It records editor activity from a VS Code extension, sends heartbeat events to a Go API, and displays coding activity, stats, and leaderboard views in an Angular web app.

## Project Structure

```text
apps/
  api/       Go/Fiber API, PostgreSQL storage, auth, stats, leaderboard, Swagger docs
  web/       Angular dashboard, login, settings, leaderboard, and stats UI
  vscode/    VS Code extension that tracks coding activity and sends heartbeats
bruno/       Bruno API request collections for local API testing
```

## Stack

- API: Go, Fiber, PostgreSQL, pgx, JWT, Swagger
- Web: Angular, Bun, RxJS, lucide-angular
- VS Code extension: TypeScript, Bun, VS Code Extension API
- API testing: Bruno collections

## Local URLs

| Service | URL |
| --- | --- |
| API | `http://localhost:5024` |
| Web app | `http://localhost:4200` |
| API health | `http://localhost:5024/health` |
| Swagger docs | `http://localhost:5024/api/docs/` |

The published API URL currently used by the web app and extension is:

```text
https://correct-wolverine-majoramari-6049fd71.koyeb.app
```

## Requirements

- Go 1.25+
- Bun 1.3+
- PostgreSQL
- VS Code, for extension development
- Bruno, optional for API collection testing

## API Setup

```bash
cd apps/api
cp .env.example .env
```

Update `.env` with your local values:

```env
DATABASE_URL=postgresql://user:password@localhost:5432/seismic
JWT_SECRET=replace-this-with-a-long-random-secret
PORT=5024
APP_URL=http://localhost:4200
ALLOWED_ORIGINS=http://localhost:4200
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=user@example.com
SMTP_PASS=password
```

Run the API:

```bash
go run .
```

The API runs migrations on startup.

## Web Setup

```bash
cd apps/web
bun install
bun run start
```

Development configuration points the web app at:

```text
http://localhost:5024
```

Production configuration points the web app at:

```text
https://correct-wolverine-majoramari-6049fd71.koyeb.app
```

Build the web app:

```bash
bun run build
```

## VS Code Extension Setup

```bash
cd apps/vscode
bun install
bun run build
```

Open `apps/vscode` in VS Code and launch the extension host from the debugger.

The extension sends heartbeats to:

- `http://localhost:5024` while running in VS Code extension development mode
- `https://correct-wolverine-majoramari-6049fd71.koyeb.app` when installed/published

Users can override the API URL with the `seismic.apiUrl` VS Code setting.

### Extension Commands

- `Seismic: Set API Key`
- `Seismic: Open Dashboard`
- `Seismic: Enable Tracking`
- `Seismic: Disable Tracking`

### Extension Settings

- `seismic.apiKey`: API key used to authenticate heartbeat requests
- `seismic.apiUrl`: API server URL
- `seismic.enabled`: enables or disables tracking
- `seismic.statusBarEnabled`: shows or hides today's tracked time in the status bar

The extension also stores shared editor configuration in `~/.seismic.cfg` so other editor integrations can reuse the same API key and API URL.

## Main API Routes

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/api/auth/magic-link` | Request a login magic link |
| `GET` | `/api/auth/verify` | Verify magic link |
| `POST` | `/api/auth/complete-signup` | Complete signup |
| `POST` | `/api/auth/refresh` | Refresh access token |
| `GET` | `/api/auth/apikey` | Get the user's editor API key |
| `POST` | `/api/auth/apikey/regenerate` | Regenerate editor API key |
| `POST` | `/api/heartbeat` | Receive editor heartbeat events |
| `GET` | `/api/stats/summary` | Summary stats |
| `GET` | `/api/stats/languages` | Language stats |
| `GET` | `/api/stats/heatmap` | Activity heatmap |
| `GET` | `/api/leaderboard` | Leaderboard data |
| `POST` | `/api/admin/process-sessions` | Manually process sessions |

## VS Code Extension Publishing

Install the VS Code extension publishing tool:

```bash
npm install -g @vscode/vsce
```

Build and package the extension:

```bash
cd apps/vscode
bun run build
vsce package
```

Test the generated `.vsix` locally:

```bash
code --install-extension seismic-0.1.0.vsix
```

Publish after logging in with the Marketplace publisher ID:

```bash
vsce login <publisher-id>
vsce publish
```

The Marketplace publisher display name can be `Muhannad H.`, but the `publisher` field in `apps/vscode/package.json` must be the publisher ID from Visual Studio Marketplace. Publisher IDs are stable identifiers and usually do not contain spaces or punctuation.

## Testing

API:

```bash
cd apps/api
go test ./...
```

Web:

```bash
cd apps/web
bun run test
```

VS Code extension:

```bash
cd apps/vscode
bun test
```

## Notes

- Keep secrets out of Git. Use `.env` locally and platform environment variables in production.
- The extension only tracks real file documents and ignores untitled or non-file VS Code documents.
- Failed heartbeat requests are queued locally by the extension and retried later.
- The web and extension API URLs should stay aligned when the production API host changes.
