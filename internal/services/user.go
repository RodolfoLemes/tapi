package services

import (
	"context"
	"errors"

	"tapi/internal/entities"
	"tapi/internal/repositories"
)

type UserService interface {
	GetOrCreateUserByTelegramID(
		ctx context.Context,
		dto GetOrCreateUserByTelegramIDDTO,
	) (*entities.User, error)
}

func NewUserService(
	userRepository repositories.UserRepository,
) UserService {
	return &applicationUserService{userRepository}
}

type applicationUserService struct {
	userRepository repositories.UserRepository
}

func (s applicationUserService) GetOrCreateUserByTelegramID(
	ctx context.Context,
	dto GetOrCreateUserByTelegramIDDTO,
) (*entities.User, error) {
	user, err := s.userRepository.FindByTelegramID(ctx, dto.TelegramID)
	if err != nil {
		if !errors.As(err, &repositories.ErrNotFound{}) {
			return nil, err
		}
		user := entities.NewUser(dto.TelegramID, dto.Username)

		err = s.userRepository.Create(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
