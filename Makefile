.PHONY: build run clean docker-build docker-run test

# Vars
APP_NAME=telegram-bot
DOCKER_IMAGE=telegram-torrent-bot

# TODO: Change if needed for your system
GOOS?=linux
GOARCH?=amd64


build:
	@echo "Building application..."
	@go build -o $(APP_NAME) .

run: build
	@echo "Running application..."
	@./$(APP_NAME)

clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)

# Docker 
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker-compose up -d

docker-stop:
	@echo "Stopping Docker container..."
	@docker-compose down

# Logs
docker-logs:
	@echo "Showing container logs..."
	@docker-compose logs -f

# Testing
test:
	@echo "Running tests..."
	@go test -v ./...

# Mkdir
setup:
	@echo "Setting up project..."
	@mkdir -p downloads logs
	@cp -n .env.example .env || true
	@echo "Setup complete. Remember to edit your .env file with your Telegram bot token."
