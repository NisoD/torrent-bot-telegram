.PHONY: build run clean docker-build docker-run test

# Variables
APP_NAME=telegram-bot
DOCKER_IMAGE=telegram-torrent-bot

# Go build flags
GOOS?=linux
GOARCH?=amd64

# Build the application
build:
	@echo "Building application..."
	@go build -o $(APP_NAME) .

# Run the application
run: build
	@echo "Running application..."
	@./$(APP_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker-compose up -d

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	@docker-compose down

# Show container logs
docker-logs:
	@echo "Showing container logs..."
	@docker-compose logs -f

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Create required directories
setup:
	@echo "Setting up project..."
	@mkdir -p downloads logs
	@cp -n .env.example .env || true
	@echo "Setup complete. Remember to edit your .env file with your Telegram bot token."
