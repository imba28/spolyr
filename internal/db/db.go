package db

import (
	"context"
	"fmt"
	"github.com/imba28/spolyr/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Access struct {
	db  *mongo.Database
	ctx context.Context
}

func (a Access) FindTrack(spotifyID string) (*model.Track, error) {
	filter := bson.D{primitive.E{Key: "spotify_id", Value: spotifyID}}
	var t model.Track
	return &t, a.db.Collection(model.TrackCollection).FindOne(a.ctx, filter).Decode(&t)
}

func (a Access) TracksWithoutLyrics() ([]*model.Track, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}}
	res, err := a.db.Collection(model.TrackCollection).Find(a.ctx, filter)
	if err != nil {
		return nil, err
	}
	return a.decodeTracks(res)
}

func (a Access) TracksWithLyricsCount() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$eq": true}}
	return a.db.Collection(model.TrackCollection).CountDocuments(a.ctx, filter)
}

func (a Access) TrackCount() (int64, error) {
	return a.db.Collection(model.TrackCollection).EstimatedDocumentCount(a.ctx)
}

func (a Access) LatestTracks(limit int64) ([]*model.Track, error) {
	opts := options.Find().SetLimit(limit)
	cur, err := a.db.Collection(model.TrackCollection).Find(a.ctx, bson.D{{}}, opts)
	if err != nil {
		return nil, err
	}

	return a.decodeTracks(cur)
}

func (a Access) FindTracks(filter interface{}) ([]*model.Track, error) {
	if filter == nil {
		filter = bson.D{{}}
	}

	cur, err := a.db.Collection(model.TrackCollection).Find(a.ctx, filter)
	if err != nil {
		return nil, err
	}

	return a.decodeTracks(cur)
}

func (a Access) decodeTracks(cur *mongo.Cursor) ([]*model.Track, error) {
	var tracks []*model.Track

	for cur.Next(a.ctx) {
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

	cur.Close(a.ctx)
	return tracks, nil
}

func (a Access) SaveTrack(t *model.Track) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"spotify_id", t.SpotifyID}}
	fieldsToUpdate := bson.D{
		{"spotify_id", t.SpotifyID},
		{"name", t.Name},
		{"artist", t.Artist},
		{"album_name", t.AlbumName},
		{"preview_url", t.PreviewURL},
		{"image_url", t.ImageURL},
	}
	if t.Loaded {
		fieldsToUpdate = append(fieldsToUpdate, bson.E{"lyrics", t.Lyrics}, bson.E{"loaded", t.Loaded})
	}

	update := bson.D{
		{"$set", fieldsToUpdate},
	}

	_, err := a.db.Collection(model.TrackCollection).UpdateOne(a.ctx, filter, update, opts)
	return err
}

func New(username, password string) (*Access, error) {
	ctx := context.TODO()
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
	idxOptions := options.Index()
	idxOptions.SetName("fulltext_index").SetWeights(map[string]int{
		"name":       9,
		"artist":     5,
		"album_name": 4,
		"lyrics":     2,
	})
	idx := mongo.IndexModel{
		Keys:    bson.M{"name": "text", "artist": "text", "album_name": "text"},
		Options: idxOptions,
	}
	_, err = db.Collection(model.TrackCollection).Indexes().CreateOne(ctx, idx)
	if err != nil {
		fmt.Sprintln(err)
		return nil, err
	}
	return &Access{
		db,
		context.TODO(),
	}, nil
}
