package api

import (
	"context"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"github.com/imba28/spolyr/internal/spotify"
	spotify2 "github.com/zmb3/spotify/v2"
	"log"
	"net/http"
)

func newImportApiService(repo db.TrackRepository, syncer *lyrics.Syncer, fetcher lyrics.AsyncFetcher) ImportApiServicer {
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
	c := spotifyClientFromContext(ctx)
	if c == nil {
		return openapi.Response(http.StatusForbidden, nil), nil
	}

	err := spotify.SyncTracks(ctx, spotify.NewSpotifyTrackProvider(c), i.repo)
	if err != nil {
		log.Println(err)
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	return openapi.Response(http.StatusOK, nil), nil
}

func (i ImportApiServicer) ImportLyricsPost(ctx context.Context) (openapi.ImplResponse, error) {
	panic("implement me")
}

func (i ImportApiServicer) ImportPlaylistIdPost(ctx context.Context, playlistId string) (openapi.ImplResponse, error) {
	token := tokenFromContext(ctx)
	if token == nil {
		return openapi.Response(http.StatusUnauthorized, nil), nil
	}

	c := spotify2.New(auth.Client(ctx, token))
	err := spotify.NewPlaylistProvider(c, i.repo).Download(ctx, playlistId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	return openapi.Response(http.StatusOK, nil), nil

}

var _ openapi.ImportApiServicer = &ImportApiServicer{}
