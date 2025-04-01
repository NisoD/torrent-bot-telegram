FROM golang:1.22-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/telegram-bot

# Create a smaller final image
FROM alpine:3.18

# Install CA certificates and other runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/telegram-bot /app/telegram-bot

# Create directories for downloads and logs
RUN mkdir -p /app/downloads /app/logs

# Set environment variables
ENV DOWNLOAD_PATH=/app/downloads
ENV LOG_PATH=/app/logs

# Set the command to run the application
ENTRYPOINT ["/app/telegram-bot"]
