package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized as %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! I'm your bot.")
			bot.Send(msg)
		}
	}
}
