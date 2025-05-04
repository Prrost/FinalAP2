package main

import (
	"context"
	"user-service/Storage/Sqlite"
	"user-service/config"
	"user-service/internal/RabbitMQ/Consumers"
	"user-service/internal/server"
	"user-service/useCase"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	Storage := Sqlite.NewSqliteStorage(cfg)

	UseCase := useCase.NewUseCase(Storage, cfg)

	OrderUserConsumerRMQ, err := Consumers.NewConsumer(cfg.RMQ, UseCase)
	if err == nil {
		_ = OrderUserConsumerRMQ.StartConsuming(ctx)
		defer OrderUserConsumerRMQ.Close()
	}

	server.RunGRPCServer(cfg, UseCase)
}
