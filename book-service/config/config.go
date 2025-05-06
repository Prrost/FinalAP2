package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port   string
	DBPath string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env: %v", err)
	}
	return &Config{
		Port:   os.Getenv("PORT"),
		DBPath: os.Getenv("DB_PATH"),
	}
}
