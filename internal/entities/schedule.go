package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	ID             string    `json:"id" bson:"_id"`
	Name           string    `json:"name" bson:"name"`
	ScheduledAt    time.Time `json:"scheduled_at" bson:"scheduled_at"`
	CommandID      string    `json:"command_id" bson:"command_id"`
	UserID         string    `json:"user_id" bson:"user_id"`
	TelegramChatID uint64    `json:"telegram_chat_id" bson:"telegram_chat_id"`
	IsDone         bool      `json:"is_done" bson:"is_done"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

func NewSchedule(
	name string,
	scheduledAt time.Time,
	commandID string,
	userID string,
	telegramChatID uint64,
) *Schedule {
	return &Schedule{
		ID:             primitive.NewObjectID().Hex(),
		Name:           name,
		ScheduledAt:    scheduledAt,
		CommandID:      commandID,
		UserID:         userID,
		TelegramChatID: telegramChatID,
		IsDone:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
