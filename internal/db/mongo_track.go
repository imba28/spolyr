package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTrackStore interface {
	Save(filter interface{}, data interface{}) error
	FindOne(filter interface{}) (*Track, error)
	Find(filter interface{}, opts ...*options.FindOptions) ([]*Track, error)
	Count(filter interface{}) (int64, error)
}

type TrackRepository interface {
	FindTrack(string) (*Track, error)
	TracksWithoutLyrics() ([]*Track, error)
	TracksWithoutLyricsError() ([]*Track, error)
	TracksWithLyricsError() ([]*Track, error)
	CountWithoutLyrics() (int64, error)
	CountWithLyrics() (int64, error)
	Count() (int64, error)
	LatestTracks(limit int64) ([]*Track, error)
	Search(query string) ([]*Track, error)
	Save(track *Track) error
}

type MongoTrackRepository struct {
	store                mongoTrackStore
	maxLyricsImportError int
}

func (t MongoTrackRepository) FindTrack(spotifyID string) (*Track, error) {
	filter := bson.D{primitive.E{Key: "spotify_id", Value: spotifyID}}
	return t.store.FindOne(filter)
}

func (t MongoTrackRepository) TracksWithoutLyrics() ([]*Track, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}}
	return t.store.Find(filter)
}

func (t MongoTrackRepository) TracksWithoutLyricsError() ([]*Track, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}, "lyrics_import_error_count": bson.M{"$lt": t.maxLyricsImportError}}
	return t.store.Find(filter)
}

func (t MongoTrackRepository) TracksWithLyricsError() ([]*Track, error) {
	filter := bson.M{"lyrics_import_error_count": bson.M{"$gte": t.maxLyricsImportError}}
	return t.store.Find(filter)
}

func (t MongoTrackRepository) CountWithoutLyrics() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}}
	return t.store.Count(filter)
}

func (t MongoTrackRepository) CountWithLyrics() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$eq": true}}
	return t.store.Count(filter)
}

func (t MongoTrackRepository) Count() (int64, error) {
	return t.store.Count(bson.M{})
}

func (t MongoTrackRepository) LatestTracks(limit int64) ([]*Track, error) {
	opts := options.Find().SetLimit(limit).
		SetSort(bson.D{{"_id", -1}})
	return t.store.Find(bson.D{{}}, opts)
}

func (t MongoTrackRepository) Search(query string) ([]*Track, error) {
	return t.store.Find(bson.M{
		"$text": bson.M{
			"$search": query,
		},
	})
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

	return t.store.Save(filter, bson.D{
		{"$set", fieldsToUpdate},
	})
}

func NewMongoTrackRepository(s mongoTrackStore, maxLyricsImportError int) MongoTrackRepository {
	return MongoTrackRepository{
		store:                s,
		maxLyricsImportError: maxLyricsImportError,
	}
}
