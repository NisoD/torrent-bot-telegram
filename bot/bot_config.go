package bot

import (
	"BotTelegram/config"
	"BotTelegram/server"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: Extend bot config for further usage
type BotConfig struct {
	Token     string
	API       *tgbotapi.BotAPI
	AppConfig *config.Config
}

func NewBotConfig(cfg *config.Config, logger *server.Logger) (*BotConfig, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	// optional debugging
	bot.Debug = true
	logger.LogInfo("Authorized on account %s", bot.Self.UserName)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &BotConfig{
		Token:     cfg.TelegramToken,
		API:       bot,
		AppConfig: cfg,
	}, nil
}
