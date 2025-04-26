package server

import (
	userpb "github.com/Prrost/protoFinalAP2/user"
	"google.golang.org/grpc"
	"log"
	"net"
	"user-service/config"
	"user-service/internal/handlers"
	"user-service/internal/middleware"
	"user-service/useCase"
)

func RunGRPCServer(cfg *config.Config, uc *useCase.UseCase) {
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.LoggingInterceptor))

	userServer := handlers.NewUserServer(cfg, uc)

	userpb.RegisterUserServiceServer(grpcServer, userServer)

	log.Printf("gRPC server listening on %s\n", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error on server start: %v", err)
	}
}
