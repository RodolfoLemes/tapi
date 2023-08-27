package container

import (
	"tapi/internal/mongodb"
	"tapi/internal/repositories"
	"tapi/internal/services"
)

type Container struct {
	UserService     services.UserService
	CommandService  services.CommandService
	ScheduleService services.ScheduleService
}

func New() (*Container, error) {
	db, err := mongodb.New()
	if err != nil {
		return nil, err
	}

	userRepository := repositories.NewUserRepository(db)
	commandRepository := repositories.NewCommandRepository(db)
	scheduleRepository := repositories.NewScheduleRepository(db)

	userService := services.NewUserService(userRepository)
	commandService := services.NewCommandService(commandRepository)
	scheduleService := services.NewScheduleService(scheduleRepository, userService, commandService)

	return &Container{userService, commandService, scheduleService}, nil
}
