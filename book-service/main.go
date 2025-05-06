package main

import (
	"log"
	"net"

	sqlite "github.com/Prrost/FinalAP2/book-service/Storage"
	"github.com/Prrost/FinalAP2/book-service/config"
	"github.com/Prrost/FinalAP2/book-service/internal/handlers"
	"github.com/Prrost/FinalAP2/book-service/useCase"

	bookpb "github.com/Prrost/protoFinalAP2/books"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Загружаем конфиг (.env)
	cfg := config.LoadConfig()

	// Инициализируем хранилище SQLite
	storage := sqlite.NewSQLiteStorage(cfg)

	// Инициализируем бизнес-логику
	uc := useCase.NewUseCase(storage)

	// Создаём gRPC-сервер и регистрируем наш сервис
	grpcServer := grpc.NewServer()
	handlersrv := handlers.NewServer(uc)
	bookpb.RegisterBookServiceServer(grpcServer, handlersrv)

	// Включаем поддержку reflection для grpcurl и отладки
	reflection.Register(grpcServer)

	// Начинаем слушать
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.Port, err)
	}
	log.Printf("BookService listening on %s", cfg.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC Serve error: %v", err)
	}
}
