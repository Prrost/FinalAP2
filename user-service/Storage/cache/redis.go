package cache

import (
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedis() *redis.Client {
	log.Println("[NewRedis] Creating new redis client")

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
