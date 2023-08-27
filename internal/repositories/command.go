package repositories

import (
	"context"

	"tapi/internal/entities"
	"tapi/internal/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommandRepository interface {
	Create(ctx context.Context, command *entities.Command) error
}

func NewCommandRepository(db *mongodb.MongoDatabase) CommandRepository {
	return &mongoCommandRepository{
		coll: db.GetCollection("commands"),
	}
}

type mongoCommandRepository struct {
	coll *mongo.Collection
}

func (m *mongoCommandRepository) Create(ctx context.Context, command *entities.Command) error {
	_, err := m.coll.InsertOne(ctx, command)

	return err
}
