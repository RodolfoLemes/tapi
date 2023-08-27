package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"tapi/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/xid"
)

func dial() (*amqp.Connection, *amqp.Channel, error) {
	amqpURL := config.Env.Broker.RabbitMQ.URL
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, err
}

func createQueue(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

type rabbitMQProducer[T any] struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	queueName *string
}

func newRabbitMQProducer[T any]() (*rabbitMQProducer[T], error) {
	conn, ch, err := dial()
	if err != nil {
		return nil, err
	}

	return &rabbitMQProducer[T]{conn, ch, nil}, nil
}

func (r *rabbitMQProducer[T]) AttachQueue(queueName string) error {
	err := createQueue(r.ch, queueName)
	if err != nil {
		return err
	}

	r.queueName = &queueName
	return nil
}

func (r *rabbitMQProducer[T]) Publish(ctx context.Context, msg T) error {
	if r.queueName == nil {
		return errors.New("no queue attached")
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.ch.PublishWithContext(
		ctx,
		"",           // since we are working with working queues, the exchange is not need
		*r.queueName, // queueName  is the routing key for working queues
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

type rabbitMQConsumer[T any] struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	id        string
	queueName *string
}

func newRabbitMQConsumer[T any]() (*rabbitMQConsumer[T], error) {
	conn, ch, err := dial()
	if err != nil {
		return nil, err
	}

	id := xid.New().String()

	return &rabbitMQConsumer[T]{conn, ch, id, nil}, nil
}

func (r *rabbitMQConsumer[T]) AttachQueue(queueName string) error {
	err := createQueue(r.ch, queueName)
	if err != nil {
		return err
	}

	r.queueName = &queueName
	return nil
}

func (r *rabbitMQConsumer[T]) Consume(ctx context.Context) (<-chan Message[T], error) {
	err := r.ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}

	msgs, err := r.ch.Consume(
		*r.queueName,
		r.id,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	channel := make(chan Message[T], 10)

	go func() {
		for d := range msgs {
			var data T
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				// deal later
				log.Println("failed to unmarshal rabbitmq message", d)
			}

			channel <- Message[T]{
				ID:        d.MessageId,
				Timestamp: d.Timestamp,
				Data:      data,
				ackTag:    d.DeliveryTag,
			}
		}
	}()

	return channel, nil
}

func (r *rabbitMQConsumer[T]) Ack(msg Message[T]) error {
	return r.ch.Ack(msg.ackTag, false)
}
