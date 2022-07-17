package api

import (
	"context"
	"github.com/imba28/spolyr/pkg/openapi"
	spotify2 "github.com/zmb3/spotify/v2"
	"net/http"
)

type playlistApiService struct {
}

func (p playlistApiService) PlaylistsGet(ctx context.Context, page int32, limit int32) (openapi.ImplResponse, error) {
	if !isAuthenticated(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	c := oauthClientFromContext(ctx)
	pp, err := c.CurrentUsersPlaylists(ctx, spotify2.Limit(int(limit)), spotify2.Offset(int((page-1)*limit)))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	playlists := make([]openapi.PlaylistInfo, len(pp.Playlists))
	for i, playlist := range pp.Playlists {
		playlists[i].Name = playlist.Name
		playlists[i].Owner = playlist.Owner.DisplayName
		if len(playlist.Images) > 0 {
			playlists[i].CoverImage = playlist.Images[0].URL
		}
		playlists[i].SpotifyId = playlist.ID.String()
		playlists[i].TrackCount = int32(playlist.Tracks.Total)
		playlists[i].IsCollaborative = playlist.Collaborative
		playlists[i].IsPublic = playlist.IsPublic
	}

	res := openapi.PlaylistsGet200Response{
		Meta: openapi.PaginationMetadata{
			Limit: limit,
			Page:  page,
			Total: int32(pp.Total),
		},
		Data: playlists,
	}
	return openapi.Response(http.StatusOK, res), nil
}

var _ openapi.PlaylistsApiServicer = &playlistApiService{}

func newPlaylistApiService() playlistApiService {
	return playlistApiService{}
}
