package nats

import (
	"encoding/json"

	"github.com/Prrost/FinalAP2/order-service/infra/logger"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc *nats.Conn
}

func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Publisher{nc: nc}, nil
}

func (p *Publisher) Publish(subject string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	logger.Log.Infof("→ %s: %s", subject, data)
	return p.nc.Publish(subject, data)
}
