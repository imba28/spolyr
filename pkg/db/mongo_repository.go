package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrTrackNotFound  = errors.New("track not found")
	ErrTracksNotFound = errors.New("tracks not found")
)

type TrackRepository interface {
	FindTrack(string) (*Track, error)
	LatestTracks(limit int64) ([]*Track, error)
	TracksWithoutLyricsError() ([]*Track, error)
	AllTracks(page, limit int) ([]*Track, int, error)
	Search(query string, page, limit int, language string) ([]*Track, int, error)
	Save(track *Track) error

	Count() (int64, error)
	CountWithLyrics() (int64, error)
}

type MongoTrackRepository struct {
	maxLyricsImportError int
	db                   *mongo.Database
}

func decodeTracks(cur *mongo.Cursor) ([]*Track, error) {
	var tracks []*Track
	ctx := context.Background()
	for cur.Next(ctx) {
		var t Track
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

func (r MongoTrackRepository) findOneByQuery(filter interface{}, o ...*options.FindOneOptions) (*Track, error) {
	var t Track
	err := r.db.Collection(TrackCollection).FindOne(context.Background(), filter, o...).Decode(&t)
	if err != nil {
		return nil, ErrTrackNotFound
	}
	return &t, nil
}
func (r MongoTrackRepository) findByQuery(filter interface{}, o ...*options.FindOptions) ([]*Track, error) {
	c, err := r.db.Collection(TrackCollection).Find(context.Background(), filter, o...)
	if err != nil {
		return nil, ErrTracksNotFound
	}
	return decodeTracks(c)
}

func (r MongoTrackRepository) count(filter interface{}) (int64, error) {
	return r.db.Collection(TrackCollection).CountDocuments(context.Background(), filter)
}

func (t MongoTrackRepository) FindTrack(spotifyID string) (*Track, error) {
	filter := bson.D{primitive.E{Key: "spotify_id", Value: spotifyID}}

	return t.findOneByQuery(filter)
}

func (t MongoTrackRepository) TracksWithoutLyrics() ([]*Track, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}}
	return t.findByQuery(filter)
}

func (t MongoTrackRepository) TracksWithoutLyricsError() ([]*Track, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}, "lyrics_import_error_count": bson.M{"$lt": t.maxLyricsImportError}}
	return t.findByQuery(filter)
}

func (t MongoTrackRepository) TracksWithLyricsError() ([]*Track, error) {
	filter := bson.M{"lyrics_import_error_count": bson.M{"$gte": t.maxLyricsImportError}}
	return t.findByQuery(filter)
}

func (t MongoTrackRepository) CountWithoutLyrics() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}}
	return t.count(filter)
}

func (t MongoTrackRepository) CountWithLyrics() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$eq": true}}
	return t.count(filter)
}

func (t MongoTrackRepository) Count() (int64, error) {
	return t.count(bson.M{})
}

func (t MongoTrackRepository) LatestTracks(limit int64) ([]*Track, error) {
	opts := options.Find().SetLimit(limit).
		SetSort(bson.D{{"_id", -1}})
	return t.findByQuery(bson.D{{}}, opts)
}

func (t MongoTrackRepository) AllTracks(page, limit int) ([]*Track, int, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64((page - 1) * limit))

	filter := bson.D{}
	total, err := t.count(filter)
	if err != nil {
		return nil, 0, err
	}

	tracks, err := t.findByQuery(filter, opts)

	return tracks, int(total), err
}

func (t MongoTrackRepository) Search(query string, page, limit int, language string) ([]*Track, int, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64((page - 1) * limit))
	filter := bson.M{
		"$text": bson.M{
			"$search":   query,
			"$language": language,
		},
	}

	total, err := t.count(filter)
	if err != nil {
		return nil, 0, err
	}

	tracks, err := t.findByQuery(filter, opts)

	return tracks, int(total), err
}

func (t MongoTrackRepository) Save(track *Track) error {
	filter := bson.D{{"spotify_id", track.SpotifyID}}
	fieldsToUpdate := bson.D{
		{"spotify_id", track.SpotifyID},
		{"name", track.Name},
		{"artist", track.Artist},
		{"album_name", track.AlbumName},
		{"preview_url", track.PreviewURL},
		{"image_url", track.ImageURL},
		{"lyrics_import_error_count", track.LyricsImportErrorCount},
	}

	if track.Loaded {
		fieldsToUpdate = append(fieldsToUpdate, bson.E{"lyrics", track.Lyrics}, bson.E{"loaded", track.Loaded})
	}

	if track.Language != "" {
		fieldsToUpdate = append(fieldsToUpdate, bson.E{"language", track.Language})
	}

	return t.save(filter, bson.D{
		{"$set", fieldsToUpdate},
	})
}

func (r MongoTrackRepository) save(filter, update interface{}) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.db.Collection(TrackCollection).UpdateOne(context.Background(), filter, update, opts)
	return err
}

func NewMongoTrackRepository(db *mongo.Database, maxLyricsImportError int) (MongoTrackRepository, error) {
	err := migrateDatabase(db)
	return MongoTrackRepository{
		db:                   db,
		maxLyricsImportError: maxLyricsImportError,
	}, err
}
