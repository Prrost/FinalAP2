package main

import (
	"log"
	"net"

	bookpb "github.com/Prrost/protoFinalAP2/books"
	"google.golang.org/grpc"

	"github.com/Prrost/FinalAP2/book-service/Storage"
	"github.com/Prrost/FinalAP2/book-service/config"
	"github.com/Prrost/FinalAP2/book-service/internal/handlers"
	"github.com/Prrost/FinalAP2/book-service/useCase"
)

func main() {
	cfg := config.LoadConfig()
	storage := Storage.NewSQLiteStorage(cfg)
	uc := useCase.NewUseCase(storage)
	srv := handlers.NewServer(uc)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	bookpb.RegisterBookServiceServer(grpcServer, srv)
	log.Printf("BookService listening on %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
