package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Command struct {
	ID             string    `json:"id" bson:"_id"`
	Name           string    `json:"name" bson:"name"`
	Text           string    `json:"text" bson:"text"`
	SentAt         time.Time `json:"sent_at" bson:"sent_at"`
	UserID         string    `json:"user_id" bson:"user_id"`
	TelegramChatID uint64    `json:"telegram_chat_id" bson:"telegram_chat_id"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

func NewCommand(
	name,
	text string,
	sendAt time.Time,
	userID string,
	telegramChatID uint64,
) *Command {
	return &Command{
		ID:             primitive.NewObjectID().Hex(),
		Name:           name,
		Text:           text,
		SentAt:         sendAt,
		UserID:         userID,
		TelegramChatID: telegramChatID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
