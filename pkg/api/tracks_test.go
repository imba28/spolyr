package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"testing"
)

type trackRepoMock struct {
	mock.Mock
}

func (t *trackRepoMock) Count() (int64, error) {
	args := t.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (t *trackRepoMock) CountWithLyrics() (int64, error) {
	args := t.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (t *trackRepoMock) AllTracks(page, limit int) ([]*db.Track, int, error) {
	args := t.Called(page, limit)
	return args.Get(0).([]*db.Track), args.Int(1), args.Error(2)
}
func (t *trackRepoMock) FindTrack(s string) (*db.Track, error) {
	args := t.Called(s)
	return args.Get(0).(*db.Track), args.Error(1)
}
func (t *trackRepoMock) LatestTracks(limit int64) ([]*db.Track, error) {
	args := t.Called(limit)
	return args.Get(0).([]*db.Track), args.Error(1)
}
func (t *trackRepoMock) TracksWithoutLyricsError() ([]*db.Track, error) {
	args := t.Called()
	return args.Get(0).([]*db.Track), args.Error(1)
}
func (t *trackRepoMock) Search(query string, page, limit int, language string) ([]*db.Track, int, error) {
	args := t.Called(query, page, limit, language)
	return args.Get(0).([]*db.Track), args.Int(1), args.Error(2)
}
func (t *trackRepoMock) Save(track *db.Track) error {
	return t.Called(track).Error(0)
}

var _ db.TrackRepository = &trackRepoMock{}

type languageDetectorMock struct {
	mock.Mock
}

func (l *languageDetectorMock) Detect(s string) (string, error) {
	args := l.Called(s)
	return args.String(0), args.Error(1)
}

var _ languageDetector = &languageDetectorMock{}

func TestTracksApiService_TracksIdGet(t *testing.T) {
	spotifyId := "1234"
	track := db.Track{
		SpotifyID: spotifyId,
	}
	m := new(trackRepoMock)
	m.On("FindTrack", spotifyId).Return(&track, nil)
	trackApi := TracksApiService{repo: m}

	res, err := trackApi.TracksIdGet(context.Background(), spotifyId)

	assert.Nil(t, err)
	assert.IsType(t, openapi.TrackDetail{}, res.Body)

	td, _ := res.Body.(openapi.TrackDetail)
	assert.Equal(t, track.SpotifyID, td.SpotifyId)
}

func TestTracksApiService_TracksGet(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		tracks := []*db.Track{{SpotifyID: "1"}, {SpotifyID: "2"}}
		query, limit, page := "foo", int32(5), int32(2)
		totalResults := 10

		m := new(trackRepoMock)
		m.On("Search", query, int(page), int(limit), mock.Anything).Return(tracks, totalResults, nil)
		lm := new(languageDetectorMock)
		lm.On("Detect", query).Return("english", nil)
		trackApi := TracksApiService{repo: m, languageDetector: lm}

		res, err := trackApi.TracksGet(context.Background(), page, limit, query)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, openapi.TracksGet200Response{}, res.Body)

		tr, _ := res.Body.(openapi.TracksGet200Response)
		assert.Equal(t, int32(totalResults), tr.Meta.Total)
		assert.Equal(t, page, tr.Meta.Page)
		assert.Equal(t, limit, tr.Meta.Limit)
		assert.Len(t, tr.Data, len(tracks))

		assert.Equal(t, tr.Data[0].SpotifyId, tracks[0].SpotifyID)
		assert.Equal(t, tr.Data[1].SpotifyId, tracks[1].SpotifyID)
	})

	t.Run("transforms query keywords", func(t *testing.T) {
		tests := []struct {
			query         string
			expectedQuery string
		}{
			{
				query:         "hello world",
				expectedQuery: "\"hello\" \"world\"",
			},
			{
				query:         "\"hello\" world",
				expectedQuery: "\"hello\" world",
			},
			{
				query:         "world",
				expectedQuery: "world",
			},
			{
				query:         "united states -america",
				expectedQuery: "united states -america",
			},
			{
				query:         "united \"states\" -america",
				expectedQuery: "united \"states\" -america",
			},
		}

		for _, testCase := range tests {
			t.Run(fmt.Sprintf("query %s", testCase.query), func(t *testing.T) {
				var tracks []*db.Track

				m := new(trackRepoMock)
				m.On("Search", testCase.expectedQuery, 1, 10, mock.Anything).Return(tracks, 1, nil)
				lm := new(languageDetectorMock)
				lm.On("Detect", testCase.query).Return("english", nil)
				trackApi := TracksApiService{repo: m, languageDetector: lm}

				_, _ = trackApi.TracksGet(context.Background(), 1, 10, testCase.query)

				m.AssertExpectations(t)
			})

		}

	})

	t.Run("no results found", func(t *testing.T) {
		var tracks []*db.Track

		m := new(trackRepoMock)
		m.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tracks, 0, mongo.ErrNoDocuments)
		lm := new(languageDetectorMock)
		lm.On("Detect", mock.Anything).Return("english", nil)
		trackApi := TracksApiService{repo: m, languageDetector: lm}

		res, err := trackApi.TracksGet(context.Background(), 1, 10, "query")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, openapi.TracksGet200Response{}, res.Body)

		tr, _ := res.Body.(openapi.TracksGet200Response)
		assert.Equal(t, int32(0), tr.Meta.Total)
		assert.Len(t, tr.Data, 0)
	})

	t.Run("database error", func(t *testing.T) {
		var tracks []*db.Track
		databaseErr := errors.New("database error")

		m := new(trackRepoMock)
		m.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tracks, 0, databaseErr)
		lm := new(languageDetectorMock)
		lm.On("Detect", mock.Anything).Return("english", nil)
		trackApi := TracksApiService{repo: m, languageDetector: lm}

		res, err := trackApi.TracksGet(context.Background(), 1, 10, "query")

		assert.Equal(t, databaseErr, err)
		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})
}

func TestTracksApiService_TracksIdPatch(t *testing.T) {
	t.Run("unauthenticated access", func(t *testing.T) {
		m := new(trackRepoMock)
		trackApi := TracksApiService{repo: m}

		res, err := trackApi.TracksIdPatch(context.Background(), "id", openapi.Lyrics{Lyrics: "new lyrics"})

		m.AssertNotCalled(t, "FindTrack")
		m.AssertNotCalled(t, "Save")

		assert.Equal(t, ErrNotAuthenticated, err)
		assert.Equal(t, res.Code, http.StatusUnauthorized)
	})

	t.Run("authenticated access", func(t *testing.T) {
		newLyrics := "neue lyrics"
		track := db.Track{SpotifyID: "id"}
		m := new(trackRepoMock)
		lm := new(languageDetectorMock)
		trackApi := TracksApiService{repo: m, languageDetector: lm}
		m.On("FindTrack", "id").Return(&track, nil)
		m.On("Save", &track).Return(nil)
		lm.On("Detect", newLyrics).Return("german", nil)

		ctx := context.WithValue(context.Background(), jwtAccessKey, "valid-token")
		res, err := trackApi.TracksIdPatch(ctx, "id", openapi.Lyrics{Lyrics: newLyrics})

		m.AssertExpectations(t)
		lm.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, res.Code, http.StatusOK)

		td, _ := res.Body.(openapi.TrackDetail)
		assert.Equal(t, td.Lyrics, newLyrics)
		assert.True(t, td.HasLyrics)
		assert.Equal(t, td.Language, "german")
	})

	t.Run("database error", func(t *testing.T) {
		track := db.Track{SpotifyID: "id"}
		databaseErr := errors.New("database error")
		m := new(trackRepoMock)
		lm := new(languageDetectorMock)
		trackApi := TracksApiService{repo: m, languageDetector: lm}
		m.On("FindTrack", "id").Return(&track, nil)
		m.On("Save", &track).Return(databaseErr)
		lm.On("Detect", mock.Anything).Return("", nil)

		ctx := context.WithValue(context.Background(), jwtAccessKey, "valid-token")
		res, err := trackApi.TracksIdPatch(ctx, "id", openapi.Lyrics{Lyrics: "new lyrics"})

		m.AssertExpectations(t)
		lm.AssertExpectations(t)

		assert.Equal(t, databaseErr, err)
		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})

	t.Run("track not found", func(t *testing.T) {
		var track db.Track
		m := new(trackRepoMock)
		lm := new(languageDetectorMock)
		trackApi := TracksApiService{repo: m}
		m.On("FindTrack", "id").Return(&track, errors.New("not found"))
		lm.On("Detect", mock.Anything).Return("english", nil)

		ctx := context.WithValue(context.Background(), jwtAccessKey, "valid-token")
		res, err := trackApi.TracksIdPatch(ctx, "id", openapi.Lyrics{Lyrics: "new lyrics"})

		m.AssertExpectations(t)
		m.AssertNotCalled(t, "Save")

		assert.Nil(t, err)
		assert.Equal(t, res.Code, http.StatusNotFound)
	})
}
