package repositories

import (
	"context"
	"errors"

	"tapi/internal/entities"
	"tapi/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, User *entities.User) error
	FindByTelegramID(ctx context.Context, telegramID uint64) (*entities.User, error)
}

func NewUserRepository(db *mongodb.MongoDatabase) UserRepository {
	return &mongoUserRepository{
		coll: db.GetCollection("users"),
	}
}

type mongoUserRepository struct {
	coll *mongo.Collection
}

func (m *mongoUserRepository) Create(
	ctx context.Context,
	User *entities.User,
) error {
	_, err := m.coll.InsertOne(ctx, User)

	return err
}

func (m *mongoUserRepository) FindByTelegramID(
	ctx context.Context,
	telegramID uint64,
) (*entities.User, error) {
	result := m.coll.FindOne(ctx, bson.D{
		{Key: "telegram_id", Value: telegramID},
	})

	var User *entities.User

	if err := result.Decode(&User); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound{m.coll.Name(), telegramID}
		}
		return nil, err
	}

	return User, nil
}
