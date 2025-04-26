package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s" + err.Error())
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		DBPath:    os.Getenv("DB_PATH"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
