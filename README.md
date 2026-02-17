# Telegram Notes Bot

This is a simple Telegram bot that allows you to save your notes to a Git repository.

## Features

- Saves notes as Markdown files in a Git repository.
- Periodically pushes changes to a remote repository.
- Authentication middleware to restrict access to the bot.
- Supports `/help` command to display a help message.
- Configurable via a YAML file and environment variables.
- Dockerized for easy deployment.

## How It Works

The bot provides a fast, asynchronous workflow for saving notes. When a user sends a message, it's immediately saved to a local file on the server, and a "Saved!" confirmation is sent back.

A background process then handles the synchronization with the remote Git repository. This ensures that the user experience is fast and responsive, as it doesn't have to wait for slow network operations.

The following diagram illustrates the complete workflow, including the names of the internal components and functions involved.

```
+------+   +----------+   +--------------------+   +--------------+           +-------------------+
| User |   | Telegram |   | NoteSaving Handler |   |  GitStorage  |           |  Remote Git Repo  |
+------+   +----------+   +--------------------+   +--------------+           +-------------------+
   |           |                  |                      |                             |
   |---Note--->|                  |                      |                             |
   |           |---Webhook------->|                      |                             |
   |           |                  |                      |                             |
   |           |                  |----[ Add() ]-------->|                             |
   |           |                  |                      |--[ createFile() ]           |
   |           |                  |                      |                             |
   |           |                  |<---[ Success ]-------|                             |
   |           |                  |                      |                             |
   |           |<--Send "Saved!"--|                      |                             |
   |           |                  |                      |                             |
   |<--Confirm-|                  |                      |                             |
   |           |                  |                      |                             |
...|...........|..................|....................................................|...
   .           .                  .                      .                             .
   .           .                  .  BACKGROUND `Processor` JOB                        .
   .           .                  .                      .                             .
   |           |                  |                      |   [ 1. Wait for Trigger ]   |
   |           |                  |                      |    - Timer (Ticker) fires   |
   |           |                  |                      |    - Note Buffer is full    |
   |           |                  |                      |                             |
   |           |                  |                      |---[ 2. prepareStorage() ]-->|
   |           |                  |                      |      (git pull)             |
   |           |                  |                      |<----------------------------|
   |           |                  |                      |                             |
   |           |                  |                      |-------[ 3. save() ]-------->|
   |           |                  |                      |        (git commit)         |
   |           |                  |                      |         (git push)          |
   |           |                  |                      |---------------------------->|
   |           |                  |                      |                             |
```

## Configuration

The application is configured via a YAML file and environment variables. The configuration file path must be provided as a command-line argument.

```yaml
# config/local.yaml
environment: "local"

logger:
  minLevel: "DEBUG"

bot:
  initTimeout: "1m"
  webHookURL: "https://example.com" # Should be redefined via environment variable
  allowedUserID: 123456789

httpServer:
  addr: ":8080"

gitRepository:
  url: "git@github.com:user/repo.git" # Should be redefined
  path: "/app/notes"
  auth:
    key: "" # Should be redefined via environment variable
    keyPassword: "" # Should be redefined via environment variable
  saveTo: "notes"
  branch: "main"
  remoteName: "origin"
  committer:
    name: "tg-notes bot"
  bufSize: 10
  updateDuration: "5m"
```

The following environment variables can be used to override the configuration:

- `TG_API_KEY`: The Telegram bot API key.
- `WEBHOOK_URL`: The URL where the bot will receive updates.
- `KEY`: The SSH private key to access the Git repository.
- `KEY_PASSWD`: The password for the SSH key.

## Installation and Usage

The application can be built and run using Docker.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/protomorphine/tg-notes.git
    cd tg-notes
    ```

2.  **Create a `config/local.yaml` file** with your configuration.

3.  **Create a `.env` file** with your environment variables:

    ```bash
    # .env
    TG_API_KEY=your_telegram_api_key
    WEBHOOK_URL=your_webhook_url
    KEY=your_ssh_private_key
    KEY_PASSWD=your_ssh_key_password
    ```

4.  **Build and run the Docker container:**

    ```bash
    docker build -t tg-notes .
    docker run --rm -it --env-file .env tg-notes --config ./config/local.yaml
    ```

## Commands

- `/help`: Shows a help message.
- Any other text message will be saved as a new note.
