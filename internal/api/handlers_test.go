package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"github.com/imba28/spolyr/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type lyricsFetcherMock struct {
	mock.Mock
}

func (l lyricsFetcherMock) Fetch(t *model.Track) error {
	return l.Called(t).Error(0)
}
func (l lyricsFetcherMock) FetchAll(ts []*model.Track) (<-chan lyrics.Result, error) {
	args := l.Called(ts)
	if args.Get(1) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(chan lyrics.Result), args.Error(1)
}

type trackServiceMock struct {
	mock.Mock
}

func (t trackServiceMock) TracksWithLyricsError() ([]*model.Track, error) {
	panic("implement me")
}

func (t trackServiceMock) FindTrack(id string) (*model.Track, error) {
	args := t.Called(id)

	var r0 *model.Track
	var r1 error
	if args.Get(1) == nil {
		r0 = args.Get(0).(*model.Track)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (t trackServiceMock) TracksWithoutLyrics() ([]*model.Track, error) {
	return testReturnValues(t.Called())
}

func (t trackServiceMock) TracksWithoutLyricsError() ([]*model.Track, error) {
	panic("not implemented")
}

func (t trackServiceMock) CountWithoutLyrics() (int64, error) {
	args := t.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (t trackServiceMock) CountWithLyrics() (int64, error) {
	args := t.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (t trackServiceMock) Count() (int64, error) {
	args := t.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (t trackServiceMock) LatestTracks(limit int64) ([]*model.Track, error) {
	return testReturnValues(t.Called(limit))
}

func (t trackServiceMock) Search(query string) ([]*model.Track, error) {
	return testReturnValues(t.Called(query))
}

func (t trackServiceMock) Save(track *model.Track) error {
	return t.Called(track).Error(0)
}

var _ db.TrackRepository = trackServiceMock{}

func testReturnValues(args mock.Arguments) ([]*model.Track, error) {
	var r0 []*model.Track
	var r1 error
	if args.Get(1) == nil {
		r0 = args.Get(0).([]*model.Track)
	} else {
		r1 = args.Error(1)
	}
	return r0, r1
}

func setUp() *gin.Engine {
	gin.SetMode(gin.TestMode)

	store := memstore.NewStore([]byte("secret"))
	router := gin.Default()
	router.Use(sessions.Sessions("session", store))

	return router
}

func TestHomePageHandler(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("Count").Return(int64(10), nil)
	mockTrackService.On("CountWithLyrics").Return(int64(10), nil)
	mockTrackService.On("LatestTracks", mock.AnythingOfType("int64")).Return([]*model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)

	router := setUp()
	router.Use(HomePageHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, 200, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackDetailHandler(t *testing.T) {
	trackId := "foobar"

	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(&model.Track{
		SpotifyID: trackId,
		Artist:    "artist",
		Name:      "test track",
	}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	router := setUp()
	router.Use(TrackDetailHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, 200, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackDetailHandler_returns_404_of_track_does_not_exist(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(nil, mongo.ErrNoDocuments)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	router := setUp()
	router.Use(TrackDetailHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackDetailHandler_returns_500_if_something_goes_wrong(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(nil, errors.New("something went wrong"))

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	router := setUp()

	router.Use(TrackDetailHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackUpdateHandler(t *testing.T) {
	trackId := "foobar"
	track := model.Track{
		SpotifyID: trackId,
		Artist:    "artist",
		Name:      "test track",
	}
	newLyrics := "It's my life..."

	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", trackId).Return(&track, nil)
	mockTrackService.On("Save", mock.AnythingOfType("*model.Track")).Return(nil)

	rr := httptest.NewRecorder()
	form := url.Values{}
	form.Set("lyrics", newLyrics)
	request, _ := http.NewRequest(http.MethodPost, "/test/"+trackId, strings.NewReader(form.Encode()))
	request.PostForm = form

	router := setUp()
	router.POST("/test/:spotifyID", TrackUpdateHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, rr.Header().Get("Location"), "/tracks/id/"+trackId)
	assert.Equal(t, track.Lyrics, newLyrics)
	assert.True(t, track.Loaded)
	mockTrackService.AssertExpectations(t)
}

func TestTrackUpdateHandler_returns_401_if_lyrics_are_missing(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(&model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)

	router := setUp()
	router.POST("/", TrackUpdateHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackUpdateHandler_returns_404_if_tracks_does_not_exist(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(nil, mongo.ErrNoDocuments)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)

	router := setUp()
	router.POST("/", TrackUpdateHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackUpdateHandler_returns_500_if_track_cannot_be_saved(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(&model.Track{}, nil)
	mockTrackService.On("Save", mock.AnythingOfType("*model.Track")).Return(errors.New("something went wrong"))

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	form := url.Values{}
	form.Set("lyrics", "It's my life...")
	request.PostForm = form

	router := setUp()
	router.POST("/", TrackUpdateHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackUpdateHandler_sets_lyrics_error_counter_to_0(t *testing.T) {
	track := model.Track{LyricsImportErrorCount: 5}

	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(&track, nil)
	mockTrackService.On("Save", mock.AnythingOfType("*model.Track")).Return(nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	form := url.Values{}
	form.Set("lyrics", "It's my life...")
	request.PostForm = form

	router := setUp()
	router.POST("/", TrackUpdateHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, track.LyricsImportErrorCount, 0)
	mockTrackService.AssertExpectations(t)
}

func TestTrackEditFormHandler(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(&model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	router := setUp()
	router.GET("/", TrackEditFormHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackEditFormHandler_returns_404_if_track_does_not_exist(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("FindTrack", mock.AnythingOfType("string")).Return(nil, mongo.ErrNoDocuments)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	router := setUp()
	router.GET("/", TrackEditFormHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackMissingLyricsHandler(t *testing.T) {
	mockTrackService := trackServiceMock{}
	mockTrackService.On("TracksWithoutLyrics").Return([]*model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	router := setUp()
	router.GET("/", TrackMissingLyricsHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackSearchHandler(t *testing.T) {
	query := "test"
	mockTrackService := trackServiceMock{}
	mockTrackService.On("Search", query).Return([]*model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/?q="+query, nil)

	router := setUp()
	router.GET("/", TrackSearchHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockTrackService.AssertExpectations(t)
}

func TestTrackSearchHandler__transforms_multiple_keywords_to_AND_query_string(t *testing.T) {
	query := "it's my life"
	resultingQuery := "\"it's\" \"my\" \"life\""

	mockTrackService := trackServiceMock{}
	mockTrackService.On("Search", resultingQuery).Return([]*model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/?q="+query, nil)

	router := setUp()
	router.GET("/", TrackSearchHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	mockTrackService.AssertExpectations(t)
}

func TestTrackSearchHandler__does_not_modify_query_if_at_least_one_contains_an_exclude_instruction(t *testing.T) {
	query := "house snow -summer"

	mockTrackService := trackServiceMock{}
	mockTrackService.On("Search", query).Return([]*model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/?q="+query, nil)

	router := setUp()
	router.GET("/", TrackSearchHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	mockTrackService.AssertExpectations(t)
}

func TestTrackSearchHandler__does_not_modify_query_if_at_least_one_contains_an_explicit_match_instruction(t *testing.T) {
	query := "it's my \"life\""

	mockTrackService := trackServiceMock{}
	mockTrackService.On("Search", query).Return([]*model.Track{}, nil)

	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/?q="+query, nil)

	router := setUp()
	router.GET("/", TrackSearchHandler(mockTrackService))
	router.ServeHTTP(rr, request)

	mockTrackService.AssertExpectations(t)
}

func TestLyricsTrackSyncHandler(t *testing.T) {
	t.Run("calls lyrics fetcher service", func(t *testing.T) {
		track := model.Track{
			SpotifyID: "foobar",
		}
		mockTrackService := trackServiceMock{}
		mockTrackService.On("FindTrack", track.SpotifyID).Return(&track, nil)
		mockTrackService.On("Save", &track).Return(nil)

		mockFetcher := lyricsFetcherMock{}
		mockFetcher.On("Fetch", &track).Return(nil)

		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/tracks/id/"+track.SpotifyID, nil)

		router := setUp()
		router.GET("/tracks/id/:spotifyID", LyricsTrackSyncHandler(mockTrackService, mockFetcher))
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, "/tracks/id/foobar", rr.Header().Get("Location"))
		mockTrackService.AssertExpectations(t)
	})
}

func TestNoRouteHandle(t *testing.T) {
	t.Run("always returns 404 status", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		router := setUp()
		router.GET("/", NoRouteHandle)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Greater(t, rr.Body.Len(), 25, "should render some kind of html interface")
	})
}
