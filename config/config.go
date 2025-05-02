package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	OpenAIAPIKey     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		OpenAIAPIKey:     os.Getenv("OPENAI_API_KEY"),
	}, nil
}
