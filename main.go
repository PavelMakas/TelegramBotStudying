package main

import (
	"log"
	"strings"

	"telegram-bot/config"
	"telegram-bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal("Error creating bot:", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			// Handle regular messages
			reversedText := handlers.HandleMessage(update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reversedText)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		// Handle commands
		switch update.Message.Command() {
		case "story":
			args := update.Message.CommandArguments()
			style := "sci-fi"
			if args != "" {
				style = strings.ToLower(args)
			}

			story, err := handlers.GenerateStory(style, "")
			if err != nil {
				errorMsg := "Sorry, I couldn't generate a story at this time. "
				if strings.Contains(err.Error(), "insufficient_quota") {
					errorMsg += "The API quota has been exceeded. Please try again later."
				} else {
					errorMsg += "Please try again later."
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, story)
			bot.Send(msg)
		}
	}
}
