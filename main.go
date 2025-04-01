package main

import (
	"BotTelegram/bot"
	"BotTelegram/config"
	"BotTelegram/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := server.NewLogger(cfg.LogPath, true)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.LogInfo("Starting Telegram Torrent Bot")
	logger.LogInfo("Download path: %s", cfg.DownloadPath)
	logger.LogInfo("Log path: %s", cfg.LogPath)

	// Initialize bot
	botCfg, err := bot.NewBotConfig(cfg, logger)
	if err != nil {
		logger.LogError("Failed to initialize bot: %v", err)
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Create and start bot
	telegramBot := bot.NewBot(botCfg, logger)
	logger.LogInfo("Bot initialized. Starting...")

	// Setup signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start bot in a goroutine
	go func() {
		if err := telegramBot.Start(); err != nil {
			logger.LogError("Bot error: %v", err)
			log.Fatalf("Bot error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	logger.LogInfo("Shutdown signal received, closing bot...")

	// Perform any cleanup here
	logger.LogInfo("Bot shutdown complete")
}
