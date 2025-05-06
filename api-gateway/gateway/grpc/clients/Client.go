package clients

import (
	"context"

	"api-gateway/config"
	"api-gateway/gateway/grpc/clients/subClients"

	bookpb "github.com/Prrost/protoFinalAP2/books"
	userpb "github.com/Prrost/protoFinalAP2/user"
)

type Client struct {
	UserClient userpb.UserServiceClient
	BookClient bookpb.BookServiceClient
}

func NewMainClient(ctx context.Context, cfg *config.Config) (*Client, []error) {
	errs := []error{}

	uc, err := subClients.InitUserClient(ctx, *cfg)
	if err != nil {
		errs = append(errs, err)
	}
	bc, err := subClients.InitBookClient(ctx, *cfg)
	if err != nil {
		errs = append(errs, err)
	}

	return &Client{
		UserClient: uc,
		BookClient: bc,
	}, errs
}
