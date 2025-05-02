package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
}

func LoadConfig() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		telegramToken = "8075731455:AAG_rbshHzEbIoFu-qjKWbWy_VVbH6P710c"
	}

	return &Config{
		TelegramBotToken: telegramToken,
	}, nil
}
