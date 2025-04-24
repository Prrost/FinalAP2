package clients

import (
	"api-gateway/config"
	"context"
)

type Client struct {
}

func NewMainClient(ctx context.Context, cfg *config.Config) (*Client, []error) {
	errArr := make([]error, 0)

	return &Client{}, errArr
}
