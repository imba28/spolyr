package db

import (
	"github.com/imba28/spolyr/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrackStore interface {
	Save(filter interface{}, data interface{}) error
	FindOne(filter interface{}) (*model.Track, error)
	Find(filter interface{}, opts ...*options.FindOptions) ([]*model.Track, error)
	Count(filter interface{}) (int64, error)
}

type TrackRepository struct {
	store TrackStore
}

func (t TrackRepository) FindTrack(spotifyID string) (*model.Track, error) {
	filter := bson.D{primitive.E{Key: "spotify_id", Value: spotifyID}}
	return t.store.FindOne(filter)
}

func (t TrackRepository) TracksWithoutLyrics() ([]*model.Track, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}}
	return t.store.Find(filter)
}

func (t TrackRepository) CountWithoutLyrics() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$ne": true}}
	return t.store.Count(filter)
}

func (t TrackRepository) CountWithLyrics() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$eq": true}}
	return t.store.Count(filter)
}

func (t TrackRepository) Count() (int64, error) {
	filter := bson.M{"loaded": bson.M{"$eq": true}}
	return t.store.Count(filter)
}

func (t TrackRepository) LatestTracks(limit int64) ([]*model.Track, error) {
	opts := options.Find().SetLimit(limit)
	return t.store.Find(bson.D{{}}, opts)
}

func (t TrackRepository) Search(query string) ([]*model.Track, error) {
	return t.store.Find(bson.M{
		"$text": bson.M{
			"$search": query,
		},
	})
}

func (t TrackRepository) Save(track *model.Track) error {
	filter := bson.D{{"spotify_id", track.SpotifyID}}
	fieldsToUpdate := bson.D{
		{"spotify_id", track.SpotifyID},
		{"name", track.Name},
		{"artist", track.Artist},
		{"album_name", track.AlbumName},
		{"preview_url", track.PreviewURL},
		{"image_url", track.ImageURL},
	}

	if track.Loaded {
		fieldsToUpdate = append(fieldsToUpdate, bson.E{"lyrics", track.Lyrics}, bson.E{"loaded", track.Loaded})
	}

	return t.store.Save(filter, bson.D{
		{"$set", fieldsToUpdate},
	})
}

func NewTrackRepository(s TrackStore) TrackRepository {
	return TrackRepository{
		store: s,
	}
}
