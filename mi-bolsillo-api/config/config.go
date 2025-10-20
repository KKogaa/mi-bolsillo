package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl        string
	DatabaseToken      string
	Port               string
	EmailProviderUrl   string
	EmailProviderToken string
	ClerkJWKSUrl       string
	GrokAPIKey         string
	TelegramBotToken   string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		DatabaseToken:      os.Getenv("DATABASE_TOKEN"),
		Port:               os.Getenv("PORT"),
		EmailProviderUrl:   os.Getenv("EMAIL_PROVIDER_URL"),
		EmailProviderToken: os.Getenv("EMAIL_PROVIDER_TOKEN"),
		ClerkJWKSUrl:       os.Getenv("CLERK_JWKS_URL"),
		GrokAPIKey:         os.Getenv("GROK_API_KEY"),
		TelegramBotToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}
