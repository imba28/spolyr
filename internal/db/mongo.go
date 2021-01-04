package db

import (
	"context"
	"github.com/imba28/spolyr/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TrackCollection = "tracks"

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
	_, err := db.Collection(TrackCollection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "spotify_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("spotify_id_index"),
		},
	)
	if err != nil {
		return err
	}

	fullTextIndex := mongo.IndexModel{
		Keys: bson.M{"name": "text", "artist": "text", "album_name": "text"},
		Options: options.Index().SetName("fulltext_index").SetWeights(map[string]int{
			"name":       9,
			"artist":     5,
			"album_name": 4,
			"lyrics":     2,
		}),
	}
	_, err = db.Collection(TrackCollection).Indexes().CreateOne(context.Background(), fullTextIndex)
	return err
}

func newMongoTrackStore(db *mongo.Database) (TrackStore, error) {
	return MongoTrackStore{
		conn: db,
	}, nil
}

var _ TrackStore = MongoTrackStore{}
