# Telegram Torrent Bot

A professional Telegram bot that downloads torrents from magnet links using the [Rain](https://github.com/cenkalti/rain) library.

## Features

- Download torrents directly from Telegram using magnet links
- Select specific files from torrents to download
- Real-time download progress updates
- File upload back to Telegram when complete
- Comprehensive logging system
- Containerized with Docker for easy deployment

## Project Structure

```
.
├── bot/                 # Bot logic
├── config/              # Configuration handling
├── server/              # Core functionality
├── downloads/           # Downloaded files storage
├── logs/                # Log files
├── .env                 # Environment variables
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Docker build instructions
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Application entry point
└── Makefile             # Build and run automation
```

## Requirements

- Go 1.16+ (for local development)
- Docker and Docker Compose (for containerized deployment)
- Telegram Bot Token (from [@BotFather](https://t.me/BotFather))

## Quick Start

### Local Development

1. Clone the repository
2. Set up environment:
   ```bash
   make setup
   ```
3. Edit the `.env` file with your Telegram bot token:
   ```
   TELEGRAM_BOT_TOKEN=your_bot_token_here
   ```
4. Build and run:
   ```bash
   make run
   ```

### Docker Deployment

1. Set up environment:
   ```bash
   make setup
   ```
2. Edit the `.env` file with your Telegram bot token
3. Build and run with Docker:
   ```bash
   make docker-run
   ```
4. View logs:
   ```bash
   make docker-logs
   ```
5. Stop the container:
   ```bash
   make docker-stop
   ```

## Environment Variables

| Variable             | Description                          | Default          |
| -------------------- | ------------------------------------ | ---------------- |
| `TELEGRAM_BOT_TOKEN` | Your Telegram bot token              | (Required)       |
| `DOWNLOAD_PATH`      | Path to store downloads              | `/app/downloads` |
| `LOG_PATH`           | Path to store logs                   | `/app/logs`      |
| `MAX_FILE_SIZE`      | Maximum file size for upload (bytes) | 50MB (52428800)  |

## Using the Bot

1. Start a conversation with your bot on Telegram
2. Send `/start` to begin
3. Send a magnet link to fetch torrent information
4. Select which files to download:
   - Send specific numbers separated by commas (e.g., "1,3,5")
   - Or send "all" to download everything
5. The bot will download your files and upload them back to you

## Advanced Configuration

For advanced configuration options, you can modify the Docker and application settings:

- Change volume mount points in `docker-compose.yml`
- Adjust application parameters in the `.env` file
- Modify bot behavior in the source code

## Troubleshooting

Common issues:

- **Bot not responding**: Check your `TELEGRAM_BOT_TOKEN` and ensure the bot is running
- **Download failures**: Some torrents may have few or no seeds
- **Upload failures**: Telegram has a 50MB file size limit for bots

Check the logs folder for detailed error information.

## License

MIT

## Disclaimer

This tool is intended for downloading and sharing legal content only. Users are responsible for complying with all applicable copyright laws.
