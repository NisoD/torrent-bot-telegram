package bot

import (
	"BotTelegram/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the Telegram bot instance
type Bot struct {
	Config *BotConfig
	Logger *server.Logger
}

// NewBot creates a new Bot instance
func NewBot(cfg *BotConfig, logger *server.Logger) *Bot {
	return &Bot{
		Config: cfg,
		Logger: logger,
	}
}

// Start starts the bot and listens for incoming messages
func (b *Bot) Start() error {
	// Set up updates configuration
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Get updates channel
	updates := b.Config.API.GetUpdatesChan(updateConfig)

	// Log bot start
	b.Logger.LogInfo("Bot started successfully. Waiting for messages...")

	// Handle updates
	for update := range updates {
		// Skip any non-message updates
		if update.Message == nil {
			continue
		}

		// Log received message
		b.Logger.LogInfo("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Handle the message
		b.handleMessage(update.Message)
	}

	return nil
}
