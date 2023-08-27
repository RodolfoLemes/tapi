package server

import (
	"context"
	"time"

	"tapi/internal/bot"
	"tapi/internal/broker"
	"tapi/internal/services"
)

type Server struct {
	consumer broker.Consumer[bot.Command]
	producer broker.Producer[bot.Command]

	scheduleService services.ScheduleService
}

func New(
	consumer broker.Consumer[bot.Command],
	producer broker.Producer[bot.Command],
	scheduleService services.ScheduleService,
) *Server {
	return &Server{consumer, producer, scheduleService}
}

func (s *Server) Init(ctx context.Context) error {
	msgs, err := s.consumer.Consume(ctx)
	if err != nil {
		return err
	}

	go s.listenerUpdates(msgs)

	return nil
}

func (s *Server) listenerUpdates(msgs <-chan broker.Message[bot.Command]) {
	for msg := range msgs {
		s.handleBrokerMessage(msg)
	}
}

func (s *Server) handleBrokerMessage(msg broker.Message[bot.Command]) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.scheduleService.Schedule(ctx, services.ScheduleDTO{
		Command: services.CreateCommandDTO{
			Name:           string(msg.Data.Name),
			Text:           msg.Data.Text,
			SentAt:         msg.Timestamp,
			TelegramChatID: uint64(msg.Data.TelegramChatID),
		},
		User: services.GetOrCreateUserByTelegramIDDTO{
			TelegramID: uint64(msg.Data.TelegramUserID),
			Username:   msg.Data.TelegramUsername,
		},
		Name:        msg.Data.Text,
		ScheduledAt: msg.Timestamp.Add(5 * time.Minute), // TODO
	})
	if err != nil {
		return err
	}

	err = s.consumer.Ack(msg)
	if err != nil {
		return err
	}

	err = s.producer.Publish(ctx, bot.Command{
		Name:             msg.Data.Name,
		Text:             "Awesome. Command successfully processed, have a nice day!",
		TelegramUserID:   msg.Data.TelegramUserID,
		TelegramChatID:   msg.Data.TelegramChatID,
		TelegramUsername: msg.Data.TelegramUsername,
	})
	if err != nil {
		return err
	}

	return nil
}
