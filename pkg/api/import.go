package api

import (
	"context"
	"errors"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/lyrics"
	"github.com/imba28/spolyr/pkg/openapi"
	"github.com/imba28/spolyr/pkg/spotify"
	"log"
	"net/http"
)

var (
	errLyricsNotFound = errors.New("no lyrics found")
)

func newImportApiService(repo db.TrackRepository, syncer *lyrics.Syncer, fetcher lyrics.AsyncFetcher, d languageDetector) ImportApiServicer {
	return ImportApiServicer{
		repo:             repo,
		syncer:           syncer,
		fetcher:          fetcher,
		languageDetector: d,
	}
}

type ImportApiServicer struct {
	repo             db.TrackRepository
	syncer           *lyrics.Syncer
	fetcher          lyrics.Fetcher
	languageDetector languageDetector
}

func (i ImportApiServicer) ImportLyricsTrackIdPost(ctx context.Context, id string) (openapi.ImplResponse, error) {
	if !isAuthenticated(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	t, err := i.repo.FindTrack(id)
	if err != nil {
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	// do not import lyrics more than once
	if t.Loaded {
		return openapi.Response(http.StatusOK, toTrackDetail(*t)), nil
	}

	err = i.fetcher.Fetch(t)
	if err != nil {
		return openapi.Response(http.StatusNotFound, nil), errLyricsNotFound
	}

	languageOfLyrics, err := i.languageDetector.Detect(t.Lyrics)
	if err != nil {
		log.Printf("Could not detect language of track %d. Falling back to english.", t.ID)
		t.Language = "english"
	} else {
		t.Language = languageOfLyrics
	}

	err = i.repo.Save(t)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, toTrackDetail(*t)), nil
}

func (i ImportApiServicer) ImportLyricsGet(ctx context.Context) (openapi.ImplResponse, error) {
	if !isAuthenticated(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	if !i.syncer.Syncing() {
		return openapi.Response(http.StatusOK, openapi.LyricsImportStatus{
			Running: false,
		}), nil
	}

	return openapi.Response(http.StatusOK, openapi.LyricsImportStatus{
		Running:          true,
		TracksCompleted:  int32(i.syncer.SyncedTracks()),
		TracksTotal:      int32(i.syncer.TotalTracks()),
		TracksError:      int32(i.syncer.TracksFailed()),
		TracksSuccessful: int32(i.syncer.TracksSuccess()),
		Log:              i.syncer.Logs(),
	}), nil
}

func (i ImportApiServicer) ImportLibraryPost(ctx context.Context) (openapi.ImplResponse, error) {
	if !isAuthenticated(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	err := spotify.SyncTracks(ctx, spotify.NewSpotifyTrackProvider(oauthClientFromContext(ctx)), i.repo)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, nil), nil
}

func (i ImportApiServicer) ImportLyricsPost(ctx context.Context) (openapi.ImplResponse, error) {
	if !isAuthenticated(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	_, err := i.syncer.Sync()
	if err == lyrics.ErrBusy {
		return openapi.Response(http.StatusTooManyRequests, nil), nil
	}

	return openapi.Response(http.StatusOK, nil), nil
}

func (i ImportApiServicer) ImportPlaylistIdPost(ctx context.Context, playlistId string) (openapi.ImplResponse, error) {
	if !isAuthenticated(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	c := oauthClientFromContext(ctx)
	err := spotify.NewPlaylistProvider(c, i.repo).Download(ctx, playlistId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	return openapi.Response(http.StatusOK, nil), nil

}

var _ openapi.ImportApiServicer = &ImportApiServicer{}
