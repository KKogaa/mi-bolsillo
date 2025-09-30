package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl        string
	Port               string
	EmailProviderUrl   string
	EmailProviderToken string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		Port:               os.Getenv("PORT"),
		EmailProviderUrl:   os.Getenv("EMAIL_PROVIDER_URL"),
		EmailProviderToken: os.Getenv("EMAIL_PROVIDER_TOKEN"),
	}
}
