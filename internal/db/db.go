package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	Tracks TrackRepository
}

func New(username, password string) (*Repositories, error) {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/").SetAuth(options.Credential{
		Username: username,
		Password: password,
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database("spolyr")
	trackStore, err := newMongoTrackStore(db)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		NewTrackRepository(trackStore),
	}, nil
}
