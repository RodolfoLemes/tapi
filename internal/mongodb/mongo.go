package mongodb

import (
	"context"
	"time"

	"tapi/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	client   *mongo.Client
	database *mongo.Database
}

func New() (*MongoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	bsonOpts := &options.BSONOptions{
		// UseJSONStructTags: true,
	}

	mongoURI := config.Env.MongoDB.URI

	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetBSONOptions(bsonOpts)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	db := client.Database("db")

	return &MongoDatabase{client, db}, nil
}

func (db *MongoDatabase) GetCollection(name string) *mongo.Collection {
	return db.database.Collection(name)
}
