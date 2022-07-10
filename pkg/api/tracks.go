package api

import (
	"context"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/openapi/openapi"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
)

var _ openapi.TracksApiServicer = &TracksApiService{}

type TracksApiService struct {
	repo             db.TrackRepository
	languageDetector languageDetector
}

func (s *TracksApiService) TracksIdPatch(ctx context.Context, id string, lyrics openapi.Lyrics) (openapi.ImplResponse, error) {
	if !isAuthenticated(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	t, err := s.repo.FindTrack(id)
	if err != nil {
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	t.Loaded = true
	t.Lyrics = lyrics.Lyrics

	ll, err := s.languageDetector.Detect(t.Lyrics)
	if err != nil {
		log.Printf("Could not detect language of track(%d). Using english as default", t.ID)
		t.Language = "english"
	} else {
		t.Language = ll
	}

	err = s.repo.Save(t)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, toTrackDetail(*t)), nil
}

func toTrackDetail(t db.Track) openapi.TrackDetail {
	return openapi.TrackDetail{
		SpotifyId:              t.SpotifyID,
		Title:                  t.Name,
		Album:                  t.AlbumName,
		CoverImage:             t.ImageURL,
		PreviewURL:             t.PreviewURL,
		Artists:                strings.Split(t.Artist, ","),
		HasLyrics:              t.Loaded,
		Lyrics:                 t.Lyrics,
		LyricsImportErrorCount: int32(t.LyricsImportErrorCount),
		Language:               t.Language,
	}
}

// newTracksApiService creates a default api service
func newTracksApiService(repo db.TrackRepository, languageDetector languageDetector) *TracksApiService {
	return &TracksApiService{
		repo:             repo,
		languageDetector: languageDetector,
	}
}

func (s *TracksApiService) TracksGet(ctx context.Context, page int32, limit int32, query string) (openapi.ImplResponse, error) {
	var tracks []*db.Track
	var err error
	var total int

	if query != "" {
		var queryLanguage string
		queryLanguage, err = s.languageDetector.Detect(query)
		if err != nil {
			queryLanguage = "english"
		}

		if strings.Index(query, " ") > -1 && strings.Index(query, "\"") == -1 && strings.Index(query, "-") == -1 {
			qs := strings.Split(query, " ")
			for i := range qs {
				qs[i] = "\"" + qs[i] + "\""
			}
			query = strings.Join(qs, " ")
		}

		tracks, total, err = s.repo.Search(query, int(page), int(limit), queryLanguage)
	} else {
		total = 10
		tracks, err = s.repo.LatestTracks(int64(limit))
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	data := make([]openapi.TrackInfo, len(tracks))
	for i, track := range tracks {
		data[i] = openapi.TrackInfo{
			SpotifyId:  track.SpotifyID,
			Title:      track.Name,
			Album:      track.AlbumName,
			CoverImage: track.ImageURL,
			PreviewURL: track.PreviewURL,
			Artists:    strings.Split(track.Artist, ", "),
			HasLyrics:  track.Loaded,
			Language:   track.Language,
		}
	}

	res := openapi.TracksGet200Response{
		Data: data,
		Meta: openapi.PaginationMetadata{
			Page:  page,
			Limit: limit,
			Total: int32(total),
		},
	}

	return openapi.Response(http.StatusOK, res), nil
}

// TracksIdGet - Returns a track
func (s *TracksApiService) TracksIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	t, err := s.repo.FindTrack(id)
	if err != nil {
		return openapi.Response(404, nil), nil
	}

	return openapi.Response(200, openapi.TrackDetail{
		SpotifyId:              t.SpotifyID,
		Title:                  t.Name,
		Album:                  t.AlbumName,
		CoverImage:             t.ImageURL,
		PreviewURL:             t.PreviewURL,
		Artists:                strings.Split(t.Artist, ","),
		HasLyrics:              t.Loaded,
		Lyrics:                 t.Lyrics,
		LyricsImportErrorCount: int32(t.LyricsImportErrorCount),
		Language:               t.Language,
	}), nil
}
