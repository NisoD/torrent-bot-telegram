#!/bin/bash

# Telegram Torrent Bot health check script
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is running
if ! docker info &>/dev/null; then
    echo -e "${RED}[ERROR]${NC} Docker is not running!"
    exit 1
fi

# Check if the container is running
if [ "$(docker ps -q -f name=telegram-torrent-bot)" ]; then
    echo -e "${GREEN}[OK]${NC} Container is running"
    
    # Get container stats in one call
    STATS=$(docker stats --no-stream --format "{{.CPUPerc}} {{.MemUsage}}" telegram-torrent-bot)
    CPU=$(echo "$STATS" | awk '{print $1}')
    MEM=$(echo "$STATS" | awk '{print $2, $3}')
    
    echo -e "${GREEN}[INFO]${NC} CPU Usage: $CPU"
    echo -e "${GREEN}[INFO]${NC} Memory Usage: $MEM"
    
    # Check logs for errors in the last hour
    ERROR_COUNT=$(docker logs --since 1h telegram-torrent-bot 2>&1 | grep -c "ERROR" || true)
    if [ "$ERROR_COUNT" -gt 0 ]; then
        echo -e "${YELLOW}[WARNING]${NC} Found $ERROR_COUNT errors in the logs in the last hour"
    else
        echo -e "${GREEN}[OK]${NC} No errors found in recent logs"
    fi
    
    # Check disk space (handle missing directories)
    DOWNLOADS_SIZE=$(du -sh downloads 2>/dev/null | cut -f1 || echo "0B")
    LOGS_SIZE=$(du -sh logs 2>/dev/null | cut -f1 || echo "0B")
    
    echo -e "${GREEN}[INFO]${NC} Downloads size: $DOWNLOADS_SIZE"
    echo -e "${GREEN}[INFO]${NC} Logs size: $LOGS_SIZE"
    
    exit 0
else
    echo -e "${RED}[ERROR]${NC} Container is not running!"
    echo "Use 'docker compose up -d' to start the container."
    exit 1
fi
