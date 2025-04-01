package bot

import (
	"BotTelegram/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Config *BotConfig
	Logger *server.Logger
}

func NewBot(cfg *BotConfig, logger *server.Logger) *Bot {
	return &Bot{
		Config: cfg,
		Logger: logger,
	}
}

// Start and listen
func (b *Bot) Start() error {

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.Config.API.GetUpdatesChan(updateConfig)

	b.Logger.LogInfo("Bot started successfully. Waiting for messages...")

	for update := range updates {
		
		if update.Message == nil {
			continue
		}

		b.Logger.LogInfo("[%s] %s", update.Message.From.UserName, update.Message.Text)
		b.handleMessage(update.Message)
	}

	return nil
}
