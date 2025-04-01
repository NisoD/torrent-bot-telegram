package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
)


type Config struct {
	TelegramToken string
	DownloadPath  string
	LogPath       string
	MaxFileSize   int64 
}

func LoadConfig() (*Config, error) {
	
	godotenv.Load()
	
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}

	downloadPath := os.Getenv("DOWNLOAD_PATH")
	if downloadPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		downloadPath = filepath.Join(homeDir, "Downloads", "BotTelegram")
	}

	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = filepath.Join(downloadPath, "logs")
	}

	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, err
	}

	// Telegram's max file size is 50MB by default, but can be increased to 2GB for bots in channels/groups
	// TODO: Check Group guidlines for larger files
	maxFileSizeStr := os.Getenv("MAX_FILE_SIZE")
	var maxFileSize int64 = 50 * 1024 * 1024 // 50 MB default
	if maxFileSizeStr != "" {
		var err error
		maxFileSize, err = strconv.ParseInt(maxFileSizeStr, 10, 64)
		if err != nil {
			maxFileSize = 50 * 1024 * 1024
		}
	}
	return &Config{
		TelegramToken: telegramToken,
		DownloadPath:  downloadPath,
		LogPath:       logPath,
		MaxFileSize:   maxFileSize,
	}, nil
}
