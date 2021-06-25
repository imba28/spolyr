package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	Tracks TrackRepository
	client *mongo.Client
}

func New(username, password, databaseName, host string, maxLyricsImportErrorCount int) (*Repositories, error) {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", host)).SetAuth(options.Credential{
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

	db := client.Database(databaseName)
	trackStore, err := newMongoTrackStore(db)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		NewMongoTrackRepository(trackStore, maxLyricsImportErrorCount),
		client,
	}, nil
}
