package main

import (
	"api-gateway/config"
	"api-gateway/gateway/grpc/clients"
	"api-gateway/gateway/http/handlers"
	"context"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	grpcClient, errors := clients.NewMainClient(ctx, cfg)
	if len(errors) > 0 {
		fmt.Println("Errors in grpc clients creating:")
		for i, err := range errors {
			fmt.Printf("Error #%d: %s\n", i+1, err)
		}
	}

	r := handlers.SetupRouter(cfg, grpcClient)

	err := r.Run(cfg.Port)
	if err != nil {
		panic(err)
	}
}
