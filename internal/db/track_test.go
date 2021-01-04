package db

import (
	"errors"
	"github.com/imba28/spolyr/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type TrackStoreMock struct {
	mock.Mock
}

func (t *TrackStoreMock) Save(filter interface{}, data interface{}) error {
	args := t.Called(filter, data)
	return args.Error(0)
}

func (t *TrackStoreMock) FindOne(filter interface{}) (*model.Track, error) {
	args := t.Called(filter)
	return args.Get(0).(*model.Track), args.Error(1)
}

func (t *TrackStoreMock) Find(filter interface{}, opts ...*options.FindOptions) ([]*model.Track, error) {
	args := t.Called(filter, opts)
	return args.Get(0).([]*model.Track), args.Error(1)
}

func (t *TrackStoreMock) Count(filter interface{}) (int64, error) {
	args := t.Called(filter)
	return args.Get(0).(int64), args.Error(1)
}

var _ TrackStore = &TrackStoreMock{}

func TestSomething_Count__successful(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(42), nil)

	r := NewTrackRepository(mocksStore)
	res, err := r.Count()

	mocksStore.AssertExpectations(t)
	assert.Equal(t, res, int64(42))
	assert.Nil(t, err)
}

func TestSomething_Count__error(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("Count", mock.AnythingOfType("primitive.M")).Return(int64(-1), errors.New("some error"))

	r := NewTrackRepository(mocksStore)
	res, err := r.Count()

	mocksStore.AssertExpectations(t)
	assert.Equal(t, res, int64(-1))
	assert.Error(t, err)
}

func TestSomething_FindOne__successful(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("FindOne", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return(&model.Track{Name: "a track", SpotifyID: "id-123"}, nil)

	r := NewTrackRepository(mocksStore)
	track, err := r.FindTrack("id-123")

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, track.SpotifyID, "id-123")
}

func TestSomething_FindOne__error(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("FindOne", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return(&model.Track{}, errors.New("db error"))

	r := NewTrackRepository(mocksStore)
	track, err := r.FindTrack("id-123")

	mocksStore.AssertExpectations(t)
	assert.Error(t, err)
	assert.Empty(t, track.Name)
}

func TestSomething_TracksWithoutLyrics(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.M"), mock.Anything).
		Return([]*model.Track{{Name: "track 1"}, {Name: "track 2"}}, nil)

	r := NewTrackRepository(mocksStore)
	tracks, err := r.TracksWithoutLyrics()

	mocksStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, tracks[0].Name, "track 1")
}

func TestSomething_TracksWithoutLyrics__error(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("Find", mock.AnythingOfType("primitive.M"), mock.Anything).
		Return([]*model.Track{}, errors.New("db error"))

	r := NewTrackRepository(mocksStore)
	tracks, err := r.TracksWithoutLyrics()

	mocksStore.AssertExpectations(t)
	assert.Len(t, tracks, 0)
	assert.Error(t, err)
}
