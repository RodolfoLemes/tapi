package services

import (
	"context"

	"tapi/internal/entities"
	"tapi/internal/repositories"
)

type CommandService interface {
	CreateCommand(ctx context.Context, dto CreateCommandDTO) (*entities.Command, error)
}

func NewCommandService(
	commandRepository repositories.CommandRepository,
) CommandService {
	return &applicationCommandService{commandRepository}
}

type applicationCommandService struct {
	commandRepository repositories.CommandRepository
}

func (s *applicationCommandService) CreateCommand(
	ctx context.Context,
	dto CreateCommandDTO,
) (*entities.Command, error) {
	command := entities.NewCommand(
		dto.Name,
		dto.Text,
		dto.SentAt,
		dto.UserID,
		dto.TelegramChatID,
	)
	return command, s.commandRepository.Create(ctx, command)
}
