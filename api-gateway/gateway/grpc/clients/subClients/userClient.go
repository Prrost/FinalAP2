package subClients

import (
	"api-gateway/config"
	"context"
	userpb "github.com/Prrost/protoFinalAP2/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"time"
)

func InitUserClient(ctx context.Context, cfg config.Config) (userpb.UserServiceClient, error) {
	const op = "InitUserClient"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	clientConn, err := grpc.NewClient(
		cfg.UserService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("[%s] Cannot init client: %v", op, err)
		return nil, err
	}

	return userpb.NewUserServiceClient(clientConn), nil
}
