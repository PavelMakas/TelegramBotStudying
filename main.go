package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"telegram-bot/config"
	"telegram-bot/handlers"
	"telegram-bot/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	bot := utils.NewBot(cfg.TelegramBotToken)
	updates := bot.GetUpdatesChan()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					bot.SendMessage(update.Message.Chat.ID, "Welcome! I can help you generate stories. Use /story <style> to get started.")
				case "story":
					args := update.Message.CommandArguments()
					if args == "" {
						bot.SendMessage(update.Message.Chat.ID, "Please specify a story style. For example: /story fantasy")
						continue
					}

					story, err := handlers.GenerateStory(args, cfg.OpenAIAPIKey)
					if err != nil {
						log.Printf("Error generating story: %v", err)
						bot.SendMessage(update.Message.Chat.ID, "Sorry, I couldn't generate a story at this time. Please try again later.")
						continue
					}

					bot.SendMessage(update.Message.Chat.ID, story)
				default:
					bot.SendMessage(update.Message.Chat.ID, "Unknown command. Use /start to see available commands.")
				}
			} else {
				response := handlers.HandleMessage(update.Message.Text)
				bot.SendMessage(update.Message.Chat.ID, response)
			}

		case <-sigChan:
			log.Println("Shutting down...")
			return
		}
	}
}
