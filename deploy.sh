#!/bin/bash

# Telegram Torrent Bot deployment script
set -e

# Print colored messages
print_message() {
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    NC='\033[0m' # No Color
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    YELLOW='\033[1;33m'
    NC='\033[0m' # No Color
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker compose &> /dev/null; then
    echo "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Setup directories
print_message "Setting up directories..."
mkdir -p downloads logs

# Check if .env file exists, create if it doesn't
if [ ! -f .env ]; then
    print_message "Creating .env file..."
    if [ -f .env.example ]; then
        cp .env.example .env
        print_warning "Please edit the .env file with your Telegram bot token before continuing."
        exit 0
    else
        echo "TELEGRAM_BOT_TOKEN=your_bot_token_here" > .env
        echo "DOWNLOAD_PATH=/app/downloads" >> .env
        echo "LOG_PATH=/app/logs" >> .env
        print_warning "Please edit the .env file with your Telegram bot token before continuing."
        exit 0
    fi
fi

# Build and start the container
print_message "Building and starting the Docker container..."
docker compose up --build -d

# Show logs
print_message "Container started. Showing logs (press Ctrl+C to exit logs):"
docker compose logs -f

# Exit message
print_message "The bot is now running in the background."
print_message "Use 'docker compose logs -f' to view logs again."
print_message "Use 'docker compose down' to stop the bot."
