# Telegram Notes Bot

Single Go service: saves Telegram messages as Markdown files, commits/pushes to Git asynchronously.

## Commands

```bash
go test ./...           # Run all tests
go test ./... -race    # Run with race detector
go mod tidy            # Update dependencies
go fmt ./...           # Format code
mockery               # Generate mocks (via .mockery.yml)
```

## Architecture

- **Async flow is critical**: user message → immediate write → confirmation. Git commit/push happens in background `Processor` goroutine (see `internal/storage/git/`).
- Config: YAML file (`--config` flag) + env vars override (`TG_API_KEY`, `WEBHOOK_URL`, `KEY`, `KEY_PASSWD`).
- Entry points: `main.go` (HTTP server), `bot.go` (Telegram webhook setup).

## Key files

- `internal/storage/git/storage.go` – `NoteAdder` interface + background processor
- `internal/bot/handlers/notesaving/` – message handler → usecase → storage
- `internal/config/config.go` – config struct + env overrides
- `internal/log/` – custom logger with request ID context

## Testing

- Mocks generated via `mockery` into `mocks/` subdirectories
- Tests alongside code (`*_test.go`)
- Config env vars: create `.env` file for local testing