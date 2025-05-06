package subClients

import (
	"context"
	"log"
	"time"

	"api-gateway/config"
	bookpb "github.com/Prrost/protoFinalAP2/books"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitBookClient(ctx context.Context, cfg config.Config) (bookpb.BookServiceClient, error) {
	const op = "InitBookClient"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	clientConn, err := grpc.DialContext(
		ctx,
		cfg.BookService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("[%s] Cannot init client: %v", op, err)
		return nil, err
	}

	return bookpb.NewBookServiceClient(clientConn), nil
}
