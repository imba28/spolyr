package spotify

import (
	"context"
	"errors"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zmb3/spotify/v2"
	"io"
	"testing"
)

type userProviderMock struct {
	mock.Mock
}

func (c userProviderMock) Tracks(ctx context.Context) ([]*db.Track, error) {
	args := c.Called(ctx)
	return args.Get(0).([]*db.Track), args.Error(1)
}
func (c userProviderMock) Next(ctx context.Context) error {
	return c.Called(ctx).Error(0)
}

type trackSaverMock struct {
	mock.Mock
}

func (t *trackSaverMock) Save(track *db.Track) error {
	args := t.Called(track)
	return args.Error(0)
}

var _ userTrackProvider = &userProviderMock{}
var _ trackSaver = &trackSaverMock{}

func TestSyncTracks__saves_tracks_to_store(t *testing.T) {
	result := []*db.Track{
		{Name: "track 1", Artist: "an artist, another artist"},
		{Name: "track 2", Artist: "an artist, another artist"},
	}
	ctx := context.Background()

	client := new(userProviderMock)
	client.On("Tracks", ctx).Return(result, nil)
	client.On("Next", ctx).Return(spotify.ErrNoMorePages)

	store := new(trackSaverMock)
	store.
		On("Save", mock.AnythingOfType("*db.Track")).
		Times(len(result)).
		Return(nil)

	_ = SyncTracks(ctx, client, store)

	store.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestSyncTracks__returns_error_if_fetching_tracks_results_in_error(t *testing.T) {
	expectedError := errors.New("unexpected error")
	ctx := context.Background()
	client := new(userProviderMock)
	client.On("Tracks", ctx).Times(1).Return([]*db.Track{}, expectedError)

	store := new(trackSaverMock)

	err := SyncTracks(ctx, client, store)

	assert.EqualError(t, err, expectedError.Error())
	store.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestSyncTracks__returns_error_if_fetching_next_page_results_in_error(t *testing.T) {
	ctx := context.Background()
	client := new(userProviderMock)
	client.On("Tracks", ctx).Times(1).Return([]*db.Track{}, nil)
	client.On("Next", ctx).Times(1).Return(io.ErrUnexpectedEOF)

	store := new(trackSaverMock)

	err := SyncTracks(ctx, client, store)

	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	store.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestSyncTracks__returns_error_if_database_returns_error(t *testing.T) {
	expectedError := errors.New("could not write to database")
	ctx := context.Background()

	client := new(userProviderMock)
	client.On("Tracks", ctx).Times(1).Return([]*db.Track{
		{Name: "track 1", Artist: "an artist, another artist"},
		{Name: "track 2", Artist: "an artist, another artist"},
	}, nil)

	store := new(trackSaverMock)
	store.On("Save", mock.Anything).Times(1).Return(expectedError)

	err := SyncTracks(ctx, client, store)

	assert.EqualError(t, err, expectedError.Error())
	store.AssertExpectations(t)
	client.AssertExpectations(t)
}
