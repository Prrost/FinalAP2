package main

import (
	"user-service/Storage/Sqlite"
	"user-service/config"
	"user-service/internal/server"
	"user-service/useCase"
)

func main() {
	cfg := config.LoadConfig()

	Storage := Sqlite.NewSqliteStorage(cfg)

	UseCase := useCase.NewUseCase(Storage, cfg)

	server.RunGRPCServer(cfg, UseCase)
}
