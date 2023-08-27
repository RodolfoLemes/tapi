package services

import (
	"context"

	"tapi/internal/entities"
	"tapi/internal/repositories"
)

type ScheduleService interface {
	Schedule(ctx context.Context, dto ScheduleDTO) error
	ListSchedules(ctx context.Context, dto ListSchedulesDTO) ([]*entities.Schedule, error)
}

func NewScheduleService(
	scheduleRepository repositories.ScheduleRepository,
	userService UserService,
	commandService CommandService,
) ScheduleService {
	return &applicationScheduleService{scheduleRepository, userService, commandService}
}

type applicationScheduleService struct {
	scheduleRepository repositories.ScheduleRepository
	userService        UserService
	commandService     CommandService
}

func (s *applicationScheduleService) Schedule(ctx context.Context, dto ScheduleDTO) error {
	user, err := s.userService.GetOrCreateUserByTelegramID(ctx, dto.User)
	if err != nil {
		return err
	}

	dto.Command.UserID = user.ID
	command, err := s.commandService.CreateCommand(ctx, dto.Command)
	if err != nil {
		return err
	}

	schedule := entities.NewSchedule(
		dto.Name,
		dto.ScheduledAt,
		command.ID,
		user.ID,
		dto.Command.TelegramChatID,
	)
	err = s.scheduleRepository.Create(ctx, schedule)
	if err != nil {
		return err
	}

	// Schedule
	// maybe a rabbitmq for another service to run a cron job inside it

	return nil
}

func (s *applicationScheduleService) ListSchedules(
	ctx context.Context,
	dto ListSchedulesDTO,
) ([]*entities.Schedule, error) {
	return nil, nil
}
