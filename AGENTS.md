# AGENTS.md

This document provides context and guidelines for AI agents (such as Gemini CLI) working on the Telegram Notes Bot project.

## Project Overview

Telegram Notes Bot is a Go application that allows users to save notes via Telegram. Notes are stored as Markdown files in a local Git repository and periodically pushed to a remote repository. The bot is designed with an asynchronous workflow: user messages are saved immediately, while a background processor handles Git operations (commit/push) to ensure responsiveness.

## Key Components

- **Bot (`internal/bot`)**: Contains Telegram bot handlers, middleware, and the main bot setup.
  - `handlers/help`: Handles `/help` command, renders a template.
  - `handlers/notesaving`: Handles text messages, saves notes via `NoteAdder` interface.
  - `middleware`: Provides authentication, logging, panic recovery, and request ID.
- **Storage (`internal/storage/git`)**: Implements `NoteAdder` and background processor for Git operations using `go-git`.
- **Config (`internal/config`)**: Loads configuration from YAML and environment variables via `cleanenv`.
- **Logger (`internal/log`)**: Custom logger with context attributes, used throughout the app.
- **Main Entry Points**:
  - `main.go`: Parses flags, initializes config, logger, storage, bot, and starts HTTP server.
  - `bot.go`: Sets up the Telegram bot, registers handlers, and configures webhook.

## Technologies

- **Go 1.25.2** – The project uses modern Go features and follows idiomatic Go style.
- **Libraries**:
  - `github.com/go-telegram/bot` – Telegram Bot API bindings.
  - `github.com/go-git/go-git/v6` – Git repository manipulation.
  - `github.com/ilyakaznacheev/cleanenv` – Configuration management.
  - `github.com/stretchr/testify` – Testing utilities.
  - `github.com/stretchr/testify/mock` (indirectly via `mockery`) – Mock generation.

## Project Structure

```
.
├── cmd/                      (not used; main.go at root)
├── config/                   – YAML config files (local.yaml, prod.yaml)
├── internal/
│   ├── app/                   – Usecase layer 
│   ├── bot/                   – Telegram bot handlers & middleware
│   ├── config/                – Config loader
│   ├── log/                   – Logger package
│   └── storage/               – Git storage implementation
├── .mockery.yml               – Mockery configuration
├── bot.go                     – Bot setup
├── logger.go                  – Global logger initialisation
├── main.go                    – Entry point
└── README.md                  – Project documentation
```

## Coding Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go) and general Go best practices.
- Use `gofmt` for formatting.
- Write clear comments for exported functions and types.
- Handle errors explicitly; avoid panics (except in middleware recovery).
- Use interfaces to decouple components (e.g., `NoteAdder`, `MessageSender`).
- Write tests using `testify`; generate mocks with `mockery` for interfaces.

## Development Workflow

1. **Always create a plan** before modifying any files. Describe the changes, affected components, and testing strategy. The agent must propose a plan and **explicitly ask for user confirmation** before proceeding with implementation. Do not start writing code until approval is received.
2. Use Go modules: `go mod tidy` after dependency changes.
3. Run tests: `go test ./...` (and ensure mocks are up-to-date with `mockery`).
4. Follow the existing code structure – place new handlers in `internal/bot/handlers/`, storage logic in `internal/storage/`, etc.
5. Update configuration schemas in `internal/config/config.go` and examples in `README.md` if you change config fields.
6. Keep the asynchronous note-saving flow in mind: `notesaving` handler should only call `Add()` on the storage, which writes locally and returns quickly. The background processor handles Git sync.

## Testing

- Unit tests reside alongside the code they test (e.g., `notesaving_test.go`).
- Use `mockery` to generate mocks for interfaces. The configuration is in `.mockery.yml`.
- Aim for good test coverage; test error paths and edge cases.
- Run tests with `-race` to detect data races.

## Configuration

- The app reads a YAML file specified via `--config` flag.
- Environment variables override YAML values (see `README.md` for the list).
- Sensitive data (API keys, SSH keys) must be passed via environment variables, never committed.

## Important Notes

- The Git storage uses a background goroutine (`Processor`) that periodically commits and pushes changes. This is crucial for performance.
- Authentication middleware restricts access to a single user ID (`allowedUserID`). Unauthorized users receive a template error.
- Templates are embedded in the binary using `//go:embed` (see `resources` folders).
- All components should receive a logger with context (request ID) to enable tracing.

## Agent Instructions

- When asked to implement a new feature or fix a bug, first provide a **detailed plan** outlining the changes. Include which files will be modified, any new interfaces or types, and how you will test it.
- After presenting the plan, **explicitly ask the user for confirmation** to proceed. Use a clear prompt such as: "Please confirm that you approve this plan so I can start implementing."
- Do not begin writing any code or making changes until you receive an explicit approval message from the user.
- Adhere to the coding guidelines above.
- Write tests for new functionality and update existing tests if needed.
- Keep documentation (especially `README.md`) in sync with code changes.
- If you need to add a dependency, explain why and run `go mod tidy`.
- Use the provided mockery configuration to generate mocks for any new interfaces.
- Ensure that concurrency patterns (goroutines, channels) are correctly handled and race-free.

By following these guidelines, you will help maintain code quality and consistency across the project.
