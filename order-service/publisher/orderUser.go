package publisher

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type OrderCreatedMessage struct {
	Email string `json:"email"`
}

type OrderCreatedPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewOrderCreatedPublisher(amqpURL string) (*OrderCreatedPublisher, error) {
	const op = "NewPublisher"

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Printf("[%s] error dialing amqp: %s", op, err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[%s] error opening channel: %s", op, err)
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"order_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("[%s] Failed to declare exchange: %s", op, err)
		return nil, err
	}

	return &OrderCreatedPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *OrderCreatedPublisher) OrderCreatedPublish(email string) error {
	const op = "OrderCreatedPublish"

	body, err := json.Marshal(OrderCreatedMessage{
		Email: email,
	})
	if err != nil {
		log.Printf("[%s] error marshaling message: %s", op, err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		"order_exchange",
		"order.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		log.Printf("[%s] Failed to publish message: %s", op, err)
		return err
	}

	log.Printf("[%s] Published order.created message for email: %s", op, email)
	return nil
}

func (p *OrderCreatedPublisher) Close() {
	if err := p.channel.Close(); err != nil {
		log.Printf("[Publisher] Failed to close channel: %s", err)
	}
	if err := p.conn.Close(); err != nil {
		log.Printf("[Publisher] Failed to close connection: %s", err)
	}
}
