package services

import "time"

type CreateCommandDTO struct {
	Name           string
	Text           string
	SentAt         time.Time
	UserID         string
	TelegramChatID uint64
}
