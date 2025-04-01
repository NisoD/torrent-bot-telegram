package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds application configuration
type Config struct {
	TelegramToken string
	DownloadPath  string
	LogPath       string
	MaxFileSize   int64 // Maximum file size in bytes that can be uploaded to Telegram
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	// Get required environment variables
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}

	// Set download path
	downloadPath := os.Getenv("DOWNLOAD_PATH")
	if downloadPath == "" {
		// Use a default path if not specified
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		downloadPath = filepath.Join(homeDir, "Downloads", "BotTelegram")
	}

	// Set log path
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = filepath.Join(downloadPath, "logs")
	}

	// Ensure directories exist
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, err
	}

	// Telegram's max file size is 50MB by default, but can be increased to 2GB for bots in channels/groups
	// We'll set a default of 50MB
	maxFileSizeStr := os.Getenv("MAX_FILE_SIZE")
	var maxFileSize int64 = 50 * 1024 * 1024 // 50 MB default
	if maxFileSizeStr != "" {
		var err error
		maxFileSize, err = strconv.ParseInt(maxFileSizeStr, 10, 64)
		if err != nil {
			// If there's an error parsing, stick with the default
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
