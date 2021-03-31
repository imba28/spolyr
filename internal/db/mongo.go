package db

import (
	"context"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/imba28/spolyr/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

const TrackCollection = "tracks"

//go:embed migrations
var migrationFiles embed.FS

type MongoTrackStore struct {
	conn *mongo.Database
}

func (m MongoTrackStore) Save(filter interface{}, update interface{}) error {
	opts := options.Update().SetUpsert(true)
	_, err := m.conn.Collection(TrackCollection).UpdateOne(context.Background(), filter, update, opts)
	return err
}

func (m MongoTrackStore) FindOne(filter interface{}) (*model.Track, error) {
	var t model.Track
	err := m.conn.Collection(TrackCollection).FindOne(context.Background(), filter).Decode(&t)
	return &t, err
}

func (m MongoTrackStore) Find(filter interface{}, opts ...*options.FindOptions) ([]*model.Track, error) {
	cur, err := m.conn.Collection(TrackCollection).Find(context.Background(), filter, opts...)
	if err != nil {
		return nil, err
	}
	return m.decodeTracks(cur)
}

func (m MongoTrackStore) Count(filter interface{}) (int64, error) {
	return m.conn.Collection(TrackCollection).CountDocuments(context.Background(), filter)
}

func (m MongoTrackStore) decodeTracks(cur *mongo.Cursor) ([]*model.Track, error) {
	var tracks []*model.Track
	ctx := context.Background()
	for cur.Next(ctx) {
		var t model.Track
		err := cur.Decode(&t)
		if err != nil {
			return tracks, err
		}

		tracks = append(tracks, &t)
	}

	if err := cur.Err(); err != nil {
		return tracks, err
	}

	err := cur.Close(context.Background())
	return tracks, err
}

func createIndices(db *mongo.Database) error {
	driver, err := mongodb.WithInstance(db.Client(), &mongodb.Config{
		DatabaseName: db.Name(),
	})
	if err != nil {
		return err
	}

	source, err := httpfs.New(http.FS(migrationFiles), "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("httpfs", source, db.Name(), driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func newMongoTrackStore(db *mongo.Database) (TrackStore, error) {
	err := createIndices(db)
	return MongoTrackStore{
		conn: db,
	}, err
}

var _ TrackStore = MongoTrackStore{}
