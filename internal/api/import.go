package api

import (
	"context"
	"encoding/json"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"github.com/imba28/spolyr/internal/spotify"
	"golang.org/x/oauth2"
	"net/http"
)

func NewImportApiService(repo db.TrackRepository, syncer *lyrics.Syncer, fetcher lyrics.AsyncFetcher) ImportApiServicer {
	return ImportApiServicer{
		repo:    repo,
		syncer:  syncer,
		fetcher: fetcher,
	}
}

type ImportApiServicer struct {
	repo    db.TrackRepository
	syncer  *lyrics.Syncer
	fetcher lyrics.Fetcher
}

func (i ImportApiServicer) ImportLibraryPost(ctx context.Context) (openapi.ImplResponse, error) {
	token := ctx.Value(spotifyTokenKey).(string)
	var tok oauth2.Token
	err := json.Unmarshal([]byte(token), &tok)
	if err != nil {
		return openapi.Response(http.StatusUnauthorized, nil), nil
	}

	err = spotify.SyncTracks(spotify.NewSpotifyTrackProvider(auth.NewClient(&tok)), i.repo)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	return openapi.Response(http.StatusOK, nil), nil
}

func (i ImportApiServicer) ImportLyricsPost(ctx context.Context) (openapi.ImplResponse, error) {
	panic("implement me")
}

func (i ImportApiServicer) ImportPlaylistIdPost(ctx context.Context, s string) (openapi.ImplResponse, error) {
	panic("implement me")
}

var _ openapi.ImportApiServicer = &ImportApiServicer{}
