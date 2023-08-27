package repositories

import (
	"context"

	"tapi/internal/entities"
	"tapi/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entities.Schedule) error
	List(ctx context.Context, userID string) ([]*entities.Schedule, error)
}

func NewScheduleRepository(db *mongodb.MongoDatabase) ScheduleRepository {
	return &mongoScheduleRepository{
		coll: db.GetCollection("schedules"),
	}
}

type mongoScheduleRepository struct {
	coll *mongo.Collection
}

func (m *mongoScheduleRepository) Create(ctx context.Context, schedule *entities.Schedule) error {
	_, err := m.coll.InsertOne(ctx, schedule)

	return err
}

func (m *mongoScheduleRepository) List(ctx context.Context, userID string) ([]*entities.Schedule, error) {
	cursor, err := m.coll.Find(ctx, bson.D{
		{Key: "telegram_user_id", Value: userID},
	})
	if err != nil {
		return nil, err
	}

	var schedules []*entities.Schedule

	for cursor.Next(ctx) {
		var schedule *entities.Schedule

		err = cursor.Decode(&schedule)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
