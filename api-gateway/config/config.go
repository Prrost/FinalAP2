package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port             string
	InventoryService string
	OrderService     string
	UserService      string
	BookService      string // ← добавлено
	JWTSecret        string
	DBPath           string
	Front            string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	if os.Getenv("ENV") == "doc" {
		log.Println("Running in docker")
		return &Config{
			Port:        os.Getenv("PORT"),
			JWTSecret:   os.Getenv("JWT_SECRET"),
			Front:       os.Getenv("FRONT"),
			BookService: os.Getenv("BOOK_SERVICE"), // ← добавлено
		}
	} else {
		return &Config{
			Port:             os.Getenv("PORT"),
			JWTSecret:        os.Getenv("JWT_SECRET"),
			UserService:      os.Getenv("USER_SERVICE"),
			InventoryService: os.Getenv("INVENTORY_SERVICE"), // оставлено без изменений
			OrderService:     os.Getenv("ORDER_SERVICE"),     // оставлено без изменений
			BookService:      os.Getenv("BOOK_SERVICE"),      // ← добавлено
			Front:            os.Getenv("FRONT"),
		}
	}
}
