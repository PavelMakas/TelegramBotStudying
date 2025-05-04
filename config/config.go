package config

import (
	"log"
)

type Config struct {
	TelegramBotToken string
	OpenAIAPIKey     string
}

func LoadConfig() (*Config, error) {
	log.Printf("Config loaded successfully")
	return &Config{
		TelegramBotToken: APIKeys.TelegramBotToken,
		OpenAIAPIKey:     APIKeys.OpenAIAPIKey,
	}, nil
}
