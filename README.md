# TG-Notes Bot

TG-Notes Bot is a Telegram bot that allows you to quickly save notes to a Git repository. It's a simple and efficient way to keep track of your ideas, thoughts, and to-do lists.

## Features

- **Simple and Fast:** Just send a message to the bot, and it will be saved as a new note in your Git repository.
- **Secure:** The bot uses SSH to authenticate with your Git repository, ensuring that your notes are stored securely.
- **Configurable:** You can easily configure the bot to use your own Git repository and customize its behavior.
- **Dockerized:** The bot can be easily deployed using Docker.

## How it Works

The bot uses the Telegram Bot API to receive messages from users. When a new message is received, the bot creates a new file in a local Git repository, commits the file, and then pushes the changes to a remote repository.

The bot's workflow is as follows:

1. A user sends a message to the Telegram bot.
2. The bot receives the message via a webhook.
3. The message content is used to create a new file in a local Git repository. The filename is generated based on the current date and time.
4. The new file is committed to the local repository with a commit message "note from tg-notes".
5. The changes are pushed to the remote Git repository.

## Getting Started

To get started with the TG-Notes Bot, you'll need to have the following prerequisites:

- Go (version 1.22 or later)
- A Telegram bot token
- A Git repository

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/protomorphine/tg-notes.git
   ```
2. Create a new Telegram bot using the [BotFather](https://core.telegram.org/bots#6-botfather) and get your bot token.
3. Create a new Git repository (e.g., on GitHub) to store your notes.
4. Generate an SSH key pair to allow the bot to authenticate with your Git repository.

### Configuration

The bot is configured using a YAML file. You can create a `config/local.yaml` file based on the `config/local.yaml` example:

```yaml
environment: "local"
logger:
  minLevel: "INFO"
bot:
  initTimeout: 5s
  allowedUserID: YOUR_TELEGRAM_USER_ID
httpServer:
  addr: ":2000"
gitRepository:
  url: "YOUR_GIT_REPOSITORY_URL"
  path: "/path/to/your/local/repository"
  key: "YOUR_SSH_PRIVATE_KEY"
  keyPassword: "YOUR_SSH_KEY_PASSWORD"
  saveTo: "notes"
  branch: "main"
  committer:
    name: "tg-notes-bot"
```

Replace the placeholder values with your own settings:

- `YOUR_TELEGRAM_USER_ID`: Your Telegram user ID. You can get this by sending a message to the `@userinfobot` bot on Telegram.
- `YOUR_GIT_REPOSITORY_URL`: The SSH URL of your Git repository.
- `/path/to/your/local/repository`: The local path where the Git repository will be cloned.
- `YOUR_SSH_PRIVATE_KEY`: Your SSH private key.
- `YOUR_SSH_KEY_PASSWORD`: The password for your SSH key (if any).

### Running the Bot

Once you've configured the bot, you can run it using the following command:

```bash
go run main.go --config config/local.yaml
```

## Docker

You can also run the bot using Docker. First, build the Docker image:

```bash
docker build -t tg-notes-bot .
```

Then, run the Docker container with your configuration file:

```bash
docker run -v $(pwd)/config/local.yaml:/app/config/local.yaml tg-notes-bot
```

## Dependencies

The TG-Notes Bot uses the following main dependencies:

- [go-telegram/bot](https://github.com/go-telegram/bot): A Telegram bot library for Go.
- [go-git/go-git](https://github.com/go-git/go-git): A pure Go implementation of Git.
- [ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv): A library for reading configuration from files, environment variables, and command-line arguments.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.