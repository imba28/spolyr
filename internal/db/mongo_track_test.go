package db

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type MongoTrackStoreMock struct {
	mock.Mock
}

func (t *MongoTrackStoreMock) Save(filter interface{}, data interface{}) error {
	args := t.Called(filter, data)
	return args.Error(0)
}

func (t *MongoTrackStoreMock) FindOne(filter interface{}) (*Track, error) {
	args := t.Called(filter)
	return args.Get(0).(*Track), args.Error(1)
}

func (t *MongoTrackStoreMock) Find(filter interface{}, opts ...*options.FindOptions) ([]*Track, error) {
	args := t.Called(filter, opts)
	return args.Get(0).([]*Track), args.Error(1)
}

func (t *MongoTrackStoreMock) Count(filter interface{}) (int64, error) {
	args := t.Called(filter)
	return args.Get(0).(int64), args.Error(1)
}

var _ mongoTrackStore = &MongoTrackStoreMock{}

func TestSomething_Count__successful(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(42), nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	res, err := r.Count()

	mocksStore.AssertExpectations(t)
	assert.Equal(t, res, int64(42))
	assert.Nil(t, err)
}

func TestSomething_Count__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(-1), errors.New("some error"))

	r := NewMongoTrackRepository(mocksStore, 3)
	res, err := r.Count()

	mocksStore.AssertExpectations(t)
	assert.Equal(t, res, int64(-1))
	assert.Error(t, err)
}

func TestSomething_FindOne__successful(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("FindOne", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return(&Track{Name: "a track", SpotifyID: "id-123"}, nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	track, err := r.FindTrack("id-123")

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, track.SpotifyID, "id-123")
}

func TestSomething_FindOne__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("FindOne", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return(&Track{}, errors.New("db error"))

	r := NewMongoTrackRepository(mocksStore, 3)
	track, err := r.FindTrack("id-123")

	mocksStore.AssertExpectations(t)
	assert.Error(t, err)
	assert.Empty(t, track.Name)
}

func TestSomething_TracksWithoutLyrics(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.M"), mock.Anything).
		Return([]*Track{{Name: "track 1"}, {Name: "track 2"}}, nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	tracks, err := r.TracksWithoutLyrics()

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, tracks[0].Name, "track 1")
}

func TestSomething_TracksWithoutLyrics__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.M"), mock.Anything).
		Return([]*Track{}, errors.New("db error"))

	r := NewMongoTrackRepository(mocksStore, 3)
	tracks, err := r.TracksWithoutLyrics()

	mocksStore.AssertExpectations(t)
	assert.Len(t, tracks, 0)
	assert.Error(t, err)
}

func TestSomething_CountWithoutLyrics__successful(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(10), nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	res, err := r.CountWithoutLyrics()

	mocksStore.AssertExpectations(t)
	assert.Equal(t, res, int64(10))
	assert.Nil(t, err)
}

func TestSomething_CountWithoutLyrics__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(-1), errors.New("database error"))

	r := NewMongoTrackRepository(mocksStore, 3)
	_, err := r.CountWithoutLyrics()

	mocksStore.AssertExpectations(t)
	assert.Error(t, err)
}

func TestSomething_CountWithLyrics__successful(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(10), nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	res, err := r.CountWithLyrics()

	mocksStore.AssertExpectations(t)
	assert.Equal(t, res, int64(10))
	assert.Nil(t, err)
}

func TestSomething_CountWithLyrics__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(-1), errors.New("database error"))

	r := NewMongoTrackRepository(mocksStore, 3)
	_, err := r.CountWithLyrics()

	mocksStore.AssertExpectations(t)
	assert.Error(t, err)
}

func TestSomething_LatestTracks__successful__no_results(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return([]*Track{}, nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	res, err := r.LatestTracks(10)

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, res, 0)
}

func TestSomething_LatestTracks__successful__two_results(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return([]*Track{{}, {}}, nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	res, err := r.LatestTracks(10)

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, res, 2)
}

func TestSomething_LatestTracks__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return([]*Track{}, errors.New("database error"))

	r := NewMongoTrackRepository(mocksStore, 3)
	_, err := r.LatestTracks(10)

	mocksStore.AssertExpectations(t)
	assert.Error(t, err)
}

func TestSomething_Search__successful(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	res := []*Track{{}, {}}
	mocksStore.On("Find", mock.AnythingOfType("primitive.M"), mock.Anything).
		Return(res, nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	resActual, err := r.Search("foobar")

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, res, resActual)
}

func TestSomething_Search__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.M"), mock.Anything).
		Return([]*Track{}, errors.New("database error"))

	r := NewMongoTrackRepository(mocksStore, 3)
	_, err := r.Search("foobar")

	mocksStore.AssertExpectations(t)
	assert.Error(t, err)
}

func TestSomething_Save__calls_store_save(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Save", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D")).Return(nil)

	track := Track{}
	r := NewMongoTrackRepository(mocksStore, 3)
	err := r.Save(&track)

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestSomething_Save__does_not_override_with_empty_lyrics(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Save", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D")).Return(nil)

	track := Track{
		Loaded: false,
		Lyrics: "",
	}
	r := NewMongoTrackRepository(mocksStore, 3)
	_ = r.Save(&track)

	mocksStore.AssertExpectations(t)
	updateFields := mocksStore.Calls[0].Arguments.Get(1).(primitive.D).Map()["$set"].(primitive.D)
	if updateFields.Map()["lyrics"] != nil {
		t.Error("should not override lyrics if loaded property is set")
	}
}

func TestSomething_Save__saves_lyrics_if_loaded_flag_is_set(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Save", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D")).Return(nil)

	track := Track{
		Loaded: true,
		Lyrics: "new lyrics",
	}
	r := NewMongoTrackRepository(mocksStore, 3)
	_ = r.Save(&track)

	mocksStore.AssertExpectations(t)
	updateFields := mocksStore.Calls[0].Arguments.Get(1).(primitive.D).Map()["$set"].(primitive.D)
	if updateFields.Map()["lyrics"] != track.Lyrics {
		t.Error("should set lyrics if loaded flag is explicitly set")
	}
}

func TestSomething_Save__limits_the_update_to_one_specific_tracks(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Save", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D")).Return(nil)

	track := Track{
		SpotifyID: "spotify-123-id",
	}
	r := NewMongoTrackRepository(mocksStore, 3)
	_ = r.Save(&track)

	mocksStore.AssertExpectations(t)
	filterFields := mocksStore.Calls[0].Arguments.Get(0).(primitive.D)
	if filterFields.Map()["spotify_id"] != track.SpotifyID {
		t.Error("track should apply the update query only to his own track document")
	}
}

func TestSomething_Save__error(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Save", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D")).Return(errors.New("database error"))

	track := Track{}
	r := NewMongoTrackRepository(mocksStore, 3)
	err := r.Save(&track)

	mocksStore.AssertExpectations(t)
	assert.Error(t, err)
}

func TestSomething_TracksWithoutLyricsError(t *testing.T) {
	mocksStore := new(MongoTrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.M"), mock.Anything).
		Return([]*Track{{Name: "track 1", LyricsImportErrorCount: 1}, {Name: "track 2", LyricsImportErrorCount: 0}, {Name: "track 3", LyricsImportErrorCount: 2}}, nil)

	r := NewMongoTrackRepository(mocksStore, 3)
	tracks, err := r.TracksWithoutLyricsError()

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, tracks, 3)
	assert.Equal(t, tracks[0].Name, "track 1")

}
