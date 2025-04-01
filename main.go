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
	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init Logger
	logger, err := server.NewLogger(cfg.LogPath, true)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.LogInfo("Starting Telegram Torrent Bot")
	logger.LogInfo("Download path: %s", cfg.DownloadPath)
	logger.LogInfo("Log path: %s", cfg.LogPath)

	// Init Bot
	botCfg, err := bot.NewBotConfig(cfg, logger)
	if err != nil {
		logger.LogError("Failed to initialize bot: %v", err)
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Start Bot
	telegramBot := bot.NewBot(botCfg, logger)
	logger.LogInfo("Bot initialized. Starting...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := telegramBot.Start(); err != nil {
			logger.LogError("Bot error: %v", err)
			log.Fatalf("Bot error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	logger.LogInfo("Shutdown signal received, closing bot...")
	# TODO: MVP-2 Add cleanup 
	logger.LogInfo("Bot shutdown complete")
}
