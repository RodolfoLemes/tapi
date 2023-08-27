package services

import "time"

type ScheduleDTO struct {
	Command     CreateCommandDTO
	User        GetOrCreateUserByTelegramIDDTO
	Name        string
	ScheduledAt time.Time
}

type ListSchedulesDTO struct{}
