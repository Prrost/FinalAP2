package clients

import (
	"api-gateway/config"
	"api-gateway/gateway/grpc/clients/subClients"
	"context"
	userpb "github.com/Prrost/protoFinalAP2/user"
)

type Client struct {
	UserClient userpb.UserServiceClient
}

func NewMainClient(ctx context.Context, cfg *config.Config) (*Client, []error) {
	errArr := make([]error, 0)

	userClient, err := subClients.InitUserClient(ctx, *cfg)
	if err != nil {
		errArr = append(errArr, err)
	}

	return &Client{
		UserClient: userClient,
	}, errArr
}
