package api

import (
	"context"
	"github.com/imba28/spolyr/pkg/openapi/openapi"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
	"net/http"
	"testing"
)

func TestPlaylistApiService_PlaylistsGet(t *testing.T) {
	t.Run("deny unauthenticated access", func(t *testing.T) {
		service := playlistApiService{}
		res, err := service.PlaylistsGet(context.Background(), 1, 5)

		assert.Equal(t, res.Code, http.StatusUnauthorized)
		assert.Error(t, err)
	})

	t.Run("load playlists from Spotify api", func(t *testing.T) {
		requestedPage, requestedLimit := int32(1), int32(5)

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "=~/me/playlists",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"limit":  requestedLimit,
					"offset": (requestedPage - 1) * requestedLimit,
					"total":  100,
					"items": []map[string]interface{}{
						{
							"id":   "1",
							"name": "Playlist A",
							"tracks": map[string]interface{}{
								"total": 10,
							},
						},
						{
							"id":   "2",
							"name": "Playlist B",
							"tracks": map[string]interface{}{
								"total": 1,
							},
						},
					},
				})
			},
		)

		service := playlistApiService{}

		// inject spotify oauth client using the default transport so http responses can be mocked
		c := spotify.New(http.DefaultClient)
		ctx := context.WithValue(context.Background(), spotifyOauthClientKey, c)
		ctx = context.WithValue(ctx, jwtAccessKey, "foo")
		res, err := service.PlaylistsGet(ctx, requestedPage, requestedLimit)

		assert.Nil(t, err)
		b, ok := res.Body.(openapi.PlaylistsGet200Response)
		if !ok {
			t.Errorf("body should be openapi.PlaylistsGet200Response")
		} else {
			assert.Equal(t, b.Meta.Limit, requestedLimit)
			assert.Equal(t, b.Meta.Page, requestedPage)
			assert.Equal(t, b.Meta.Total, int32(100))
		}
	})

	t.Run("returns API errors", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "=~/me/playlists",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusBadRequest, nil)
			})

		service := playlistApiService{}

		c := spotify.New(http.DefaultClient)
		ctx := context.WithValue(context.Background(), spotifyOauthClientKey, c)
		ctx = context.WithValue(ctx, jwtAccessKey, "a-valid-token")
		res, err := service.PlaylistsGet(ctx, 1, 1)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Error(t, err)
	})
}
