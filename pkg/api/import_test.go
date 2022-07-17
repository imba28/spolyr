package api

import (
	"context"
	"errors"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/lyrics"
	"github.com/imba28/spolyr/pkg/openapi"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zmb3/spotify/v2"
	"net/http"
	"testing"
)

type fetcherMock struct {
	mock.Mock
}

func (f *fetcherMock) Fetch(track *db.Track) error {
	return f.Called(track).Error(0)
}
func (f *fetcherMock) FetchAll(tracks []*db.Track) (<-chan lyrics.Result, error) {
	args := f.Called(tracks)
	return args.Get(0).(<-chan lyrics.Result), args.Error(1)
}

var _ lyrics.Fetcher = &fetcherMock{}

func TestImportApiServicer_ImportLibraryPost(t *testing.T) {
	t.Run("denies unauthenticated access", func(t *testing.T) {
		service := ImportApiServicer{}
		res, err := service.ImportLibraryPost(context.Background())

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Error(t, err)
	})

	t.Run("imports tracks returned from spotify library API", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "=~/me/tracks",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
					"limit":  10,
					"offset": 0,
					"total":  2,
					"items": []map[string]interface{}{
						{
							"track": map[string]interface{}{
								"id":   "1",
								"name": "Track A",
							},
						},
						{
							"track": map[string]interface{}{
								"id":   "2",
								"name": "Track B",
							},
						},
					},
				})
			})

		repoMock := new(trackRepoMock)
		repoMock.On("Save", mock.AnythingOfType("*db.Track")).
			Times(2).
			Return(nil)
		service := ImportApiServicer{repo: repoMock}

		c := spotify.New(http.DefaultClient)
		ctx := context.WithValue(context.Background(), spotifyOauthClientKey, c)
		ctx = context.WithValue(ctx, jwtAccessKey, "a-valid-token")

		res, err := service.ImportLibraryPost(ctx)

		repoMock.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Nil(t, err)
	})

	t.Run("imports tracks returned from spotify library API", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "=~/me/tracks",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
					"limit":  10,
					"offset": 0,
					"total":  1,
					"items": []map[string]interface{}{
						{
							"track": map[string]interface{}{
								"id":   "1",
								"name": "Track A",
							},
						},
					},
				})
			})

		repoMock := new(trackRepoMock)
		repoMock.On("Save", mock.AnythingOfType("*db.Track")).
			Return(errors.New("database error"))
		service := ImportApiServicer{repo: repoMock}

		c := spotify.New(http.DefaultClient)
		ctx := context.WithValue(context.Background(), spotifyOauthClientKey, c)
		ctx = context.WithValue(ctx, jwtAccessKey, "a-valid-token")

		res, err := service.ImportLibraryPost(ctx)

		repoMock.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Error(t, err)
	})
}

func TestImportApiServicer_ImportLyricsTrackIdPost(t *testing.T) {
	t.Run("denies unauthenticated access", func(t *testing.T) {
		t.Run("denies unauthenticated access", func(t *testing.T) {
			service := ImportApiServicer{}
			res, err := service.ImportLyricsTrackIdPost(context.Background(), "1234")

			assert.Equal(t, http.StatusUnauthorized, res.Code)
			assert.Error(t, err)
		})
	})

	t.Run("imports and saves lyrics of a track", func(t *testing.T) {
		requestedId := "1234"
		track := &db.Track{
			SpotifyID: requestedId,
			Artist:    "Eminem",
			Name:      "Lose Yourself",
			Lyrics:    "There's vomit on his sweater already, mom's spaghetti",
		}
		repoMock := new(trackRepoMock)
		repoMock.
			On("FindTrack", requestedId).
			Return(track, nil).
			On("Save", track).
			Return(nil)
		lm := new(languageDetectorMock)
		lm.On("Detect", track.Lyrics).Return("english", nil)
		lyricsFetcherMock := new(fetcherMock)
		lyricsFetcherMock.On("Fetch", track).Return(nil)

		service := ImportApiServicer{repo: repoMock, languageDetector: lm, fetcher: lyricsFetcherMock}
		ctx := context.WithValue(context.Background(), jwtAccessKey, "a-valid-token")
		res, err := service.ImportLyricsTrackIdPost(ctx, requestedId)

		assert.IsType(t, openapi.TrackDetail{}, res.Body)
		assert.Equal(t, requestedId, res.Body.(openapi.TrackDetail).SpotifyId)
		assert.Nil(t, err)
		repoMock.AssertExpectations(t)
		lm.AssertExpectations(t)
		lyricsFetcherMock.AssertExpectations(t)
	})

	t.Run("does not import lyrics if track contains lyrics already", func(t *testing.T) {
		requestedId := "1234"
		track := &db.Track{
			SpotifyID: requestedId,
			Artist:    "Eminem",
			Name:      "Lose Yourself",
			Lyrics:    "There's vomit on his sweater already, mom's spaghetti",
			Loaded:    true,
		}
		repoMock := new(trackRepoMock)
		repoMock.
			On("FindTrack", requestedId).
			Return(track, nil).
			On("Save", track).
			Return(nil)
		lm := new(languageDetectorMock)
		lm.On("Detect", track.Lyrics).Return("english", nil)
		lyricsFetcherMock := new(fetcherMock)
		lyricsFetcherMock.On("Fetch", track).Return(nil)

		service := ImportApiServicer{repo: repoMock, languageDetector: lm, fetcher: lyricsFetcherMock}
		ctx := context.WithValue(context.Background(), jwtAccessKey, "a-valid-token")
		res, err := service.ImportLyricsTrackIdPost(ctx, requestedId)

		assert.IsType(t, openapi.TrackDetail{}, res.Body)
		assert.Nil(t, err)

		repoMock.AssertNotCalled(t, "Save")
		lm.AssertNotCalled(t, "Detect")
	})
}
