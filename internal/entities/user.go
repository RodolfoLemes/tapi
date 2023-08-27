package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         string    `json:"id" bson:"_id"`
	TelegramID uint64    `json:"telegram_id" bson:"telegram_id"`
	Username   string    `json:"username" bson:"username"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

func NewUser(telegramID uint64, username string) *User {
	return &User{
		ID:         primitive.NewObjectID().Hex(),
		TelegramID: telegramID,
		Username:   username,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
