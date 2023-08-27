package broker

import (
	"context"
	"time"

	"tapi/config"
)

var (
	CommandRequestQueue  string = "command.request"
	CommandResponseQueue string = "command.response"
)

type Message[T any] struct {
	ID        string
	Data      T
	Timestamp time.Time

	ackTag uint64
}

type Producer[T any] interface {
	AttachQueue(queueName string) error
	Publish(ctx context.Context, msg T) error
}

type Consumer[T any] interface {
	AttachQueue(queueName string) error
	Consume(ctx context.Context) (<-chan Message[T], error)
	Ack(msg Message[T]) error
}

func NewProducer[T any]() (Producer[T], error) {
	if config.Env.Broker.Provider == "rabbitmq" {
		return newRabbitMQProducer[T]()
	}

	return nil, nil
}

func NewConsumer[T any]() (Consumer[T], error) {
	if config.Env.Broker.Provider == "rabbitmq" {
		return newRabbitMQConsumer[T]()
	}

	return nil, nil
}
