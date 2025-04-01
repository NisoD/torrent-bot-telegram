FROM golang:1.22-alpine AS builder

#dependencies 
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/telegram-bot
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/telegram-bot /app/telegram-bot

RUN mkdir -p /app/downloads /app/logs

ENV DOWNLOAD_PATH=/app/downloads
ENV LOG_PATH=/app/logs

ENTRYPOINT ["/app/telegram-bot"]
