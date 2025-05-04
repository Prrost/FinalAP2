package Consumers

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"

	"user-service/domain"
	"user-service/useCase"
)

type OrderCreatedMessage struct {
	Email string `json:"email"`
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	useCase *useCase.UseCase
}

func NewConsumer(amqpURL string, useCase *useCase.UseCase) (*Consumer, error) {
	const op = "NewConsumer"
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Printf("[%s]Failed to connect to RabbitMQ: %s", op, err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[%s]Failed to open a channel RabbitMQ: %s", op, err)
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
		log.Printf("[%s]Failed to declare an exchange: %s", op, err)
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"user_service_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("[%s]Failed to declare a queue: %s", op, err)
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"order.created",
		"order_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Printf("[%s]Failed to bind a queue: %s", op, err)
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		useCase: useCase,
	}, nil
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	const op = "StartConsuming"
	msgs, err := c.channel.Consume(
		"user_service_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	log.Printf("[%s]Consumer started", op)

	go func() {
		for {
			select {
			case msg := <-msgs:
				var event OrderCreatedMessage
				err = json.Unmarshal(msg.Body, &event)
				if err != nil {
					log.Printf("[%s]Failed to unmarshal a message: %s", op, err)
					_ = msg.Nack(false, false)
					continue
				}
				log.Printf("[%s]Resived message: %s", op, event)

				_, err = c.useCase.CreateUser(domain.User{
					Email:   event.Email,
					IsAdmin: false,
				})
				if err != nil {
					st, ok := status.FromError(err)
					if ok {
						switch st.Code() {
						case codes.InvalidArgument:
							log.Printf("[%s]Invalid credentials: %s", op, st.Message())
							_ = msg.Nack(false, false)
						case codes.Internal:
							log.Printf("[%s]Internal error: %s", op, st.Message())
							_ = msg.Nack(false, false)
						case codes.AlreadyExists:
							err = msg.Ack(false)
							if err != nil {
								log.Printf("[%s]Failed to ack a message: %s", op, err)
							}
							log.Printf("[%s]User already exists, not creating: %s", op, st.Message())
						}
					} else {
						log.Printf("[%s]Unknown error: %s", op, err.Error())
						_ = msg.Nack(false, false)
					}
				}
				err = msg.Ack(false)
				if err != nil {
					log.Printf("[%s]Failed to ack a message: %s", op, err)
				}
			case <-ctx.Done():
				log.Printf("[%s]Consumer stopped", op)
				return
			}
		}
	}()
	return nil
}

func (c *Consumer) Close() {
	_ = c.channel.Close()
	_ = c.conn.Close()
}
