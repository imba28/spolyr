package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/db"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

const (
	envUsername = "DATABASE_USER"
	envPassword = "DATABASE_PASSWORD"
	envHost     = "DATABASE_HOST"

	testDatabaseName = "tests_functional"
)

var testDb *db.Repositories

func init() {
	dbConn, err := db.New(os.Getenv(envUsername), os.Getenv(envPassword), testDatabaseName, os.Getenv(envHost), 3)
	if err != nil {
		panic(err)
	}

	tracks := []*db.Track{
		{Name: "test1", SpotifyID: "1", Lyrics: "Lorem ipsum dolor sit amet", Loaded: true},
		{Name: "test2", SpotifyID: "2", Lyrics: "You know the rules and so do I.", Loaded: true},
		{Name: "test3", SpotifyID: "3", Lyrics: "His palms are sweaty, knees weak, arms are heavy", Loaded: true},
		{Name: "test4", SpotifyID: "4", Lyrics: "Fed to the rules and I hit the ground running", Loaded: true},
		{Name: "test5", SpotifyID: "5", Lyrics: "Annie are you ok?", Loaded: false, LyricsImportErrorCount: 3},
	}

	for i := range tracks {
		if err := dbConn.Tracks.Save(tracks[i]); err != nil {
			panic(err)
		}
	}

	testDb = dbConn
}

func newTestApp() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	return New(testDb, "")
}

func TestApp(t *testing.T) {
	t.Run("test if homepage returns 200", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		app := newTestApp()

		app.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("test if track page returns 200", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/tracks/id/1", nil)
		app := newTestApp()

		app.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("test if track page of not existing track returns 404", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/tracks/id/42", nil)
		app := newTestApp()

		app.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("test if search function finds tracks", func(t *testing.T) {
		tests := []struct {
			query string
			wants []string
		}{
			{"palms", []string{"test3"}},
			{"rule", []string{"test2", "test4"}},
			{"amet", []string{"test1"}},
			{"sweaty palms", []string{"test3"}},
		}

		for _, test := range tests {
			rr := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, "/search?q="+test.query, nil)
			app := newTestApp()

			app.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusOK, rr.Code)
			for _, title := range test.wants {
				assert.Contains(t, rr.Body.String(), title, "should find "+title+" because its lyrics contain the word '"+test.query+"'")
			}
		}
	})

	t.Run("Prevents an unauthenticated user from editing a track", func(t *testing.T) {
		form := url.Values{}
		form.Add("lyrics", "la la la")
		request, _ := http.NewRequest(http.MethodPost, "/tracks/id/1/edit", strings.NewReader(form.Encode()))

		app := newTestApp()
		rr := httptest.NewRecorder()
		app.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Prevents an unauthenticated user from view the edit form", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/tracks/id/1/edit", nil)
		app := newTestApp()

		app.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Lists all tracks where the import error has surpassed the threshold", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/tracks/no-lyrics-found", nil)
		app := newTestApp()

		app.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "test5", "should find 'test5' because the lyrics import error count is greater than the threshold")
	})
}
