package template

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomePage__sets_status_code(t *testing.T) {
	var tests = []int{http.StatusOK, http.StatusBadRequest, http.StatusUnauthorized}

	for _, statusCode := range tests {
		t.Run(fmt.Sprintf("testing response code %d", statusCode), func(t *testing.T) {
			r := httptest.NewRecorder()

			_ = HomePage(r, gin.H{}, statusCode)

			assert.Equal(t, statusCode, r.Code)
		})
	}
}

func TestHomePage(t *testing.T) {
	r := httptest.NewRecorder()

	err := HomePage(r, gin.H{
		"TrackCount":            0,
		"TracksWithLyricsCount": 5,
		"TracksLatest":          []model.Track{},
	}, http.StatusOK)

	assert.Nil(t, err)
}

func TestTrackPage(t *testing.T) {
	r := httptest.NewRecorder()

	err := TrackPage(r, gin.H{"Track": model.Track{}, "MaxLyricsImportErrorCount": 3}, http.StatusOK)

	assert.Nil(t, err)
}

func TestTrackEditPage(t *testing.T) {
	r := httptest.NewRecorder()

	err := TrackEditPage(r, gin.H{"Track": model.Track{}}, http.StatusOK)

	assert.Nil(t, err)
}

func TestTrackLyricsSyncPage(t *testing.T) {
	r := httptest.NewRecorder()

	err := LyricsSyncLogPage(r, gin.H{"Syncing": true, "SyncLog": "foobar"}, http.StatusOK)

	assert.Nil(t, err)
}

func TestTrackLyricsSyncPage__not_syncing(t *testing.T) {
	r := httptest.NewRecorder()

	err := LyricsSyncLogPage(r, gin.H{"Syncing": false}, http.StatusOK)

	assert.Nil(t, err)
}

func TestTracksPage(t *testing.T) {
	r := httptest.NewRecorder()

	err := TracksPage(r, gin.H{"Tracks": []model.Track{}}, http.StatusOK)

	assert.Nil(t, err)
}

func TestErrorPage(t *testing.T) {
	r := httptest.NewRecorder()

	err := ErrorPage(r, gin.H{"Message": "error!", "Status": http.StatusNotFound}, http.StatusNotFound)

	assert.Nil(t, err)
}

func TestImportPage(t *testing.T) {
	r := httptest.NewRecorder()

	err := ImportPage(r, gin.H{}, http.StatusOK)

	assert.Nil(t, err)
}
