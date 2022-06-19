package api

import (
	"context"
	"errors"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

var _ openapi.TracksApiServicer = &TracksApiService{}

type TracksApiService struct {
	repo db.TrackRepository
}

func (s *TracksApiService) TracksIdPatch(ctx context.Context, s2 string, lyrics openapi.Lyrics) (openapi.ImplResponse, error) {
	// TODO - update TracksIdPost with the required logic for this service method.
	// Add api_tracks_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, TrackDetail{}) or use other options such as http.Ok ...
	//return openapi.Response(200, TrackDetail{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return openapi.Response(404, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return openapi.Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(500, {}) or use other options such as http.Ok ...
	//return openapi.Response(500, nil),nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("TracksIdPost method not implemented")
}

// NewTracksApiService creates a default api service
func NewTracksApiService(repo db.TrackRepository) *TracksApiService {
	return &TracksApiService{
		repo: repo,
	}
}

func (s *TracksApiService) TracksGet(ctx context.Context, page int32, limit int32, query string) (openapi.ImplResponse, error) {
	var tracks []*db.Track
	var err error

	if query != "" {
		if strings.Index(query, " ") > -1 && strings.Index(query, "\"") == -1 && strings.Index(query, "-") == -1 {
			qs := strings.Split(query, " ")
			for i := range qs {
				qs[i] = "\"" + qs[i] + "\""
			}
			query = strings.Join(qs, " ")
		}

		tracks, err = s.repo.Search(query)
	} else {
		tracks, err = s.repo.LatestTracks(int64(limit))
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return openapi.Response(http.StatusNotFound, openapi.TracksGet200Response{}), nil
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
		}
	}

	res := openapi.TracksGet200Response{
		Data: data,
	}

	return openapi.Response(http.StatusOK, res), nil
}

// TracksIdGet - Returns a track
func (s *TracksApiService) TracksIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	// TODO - update TracksIdGet with the required logic for this service method.
	// Add api_tracks_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, TrackDetail{}) or use other options such as http.Ok ...
	//return openapi.Response(200, TrackDetail{}), nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return openapi.Response(404, nil),nil

	//TODO: Uncomment the next line to return response Response(500, {}) or use other options such as http.Ok ...
	//return openapi.Response(500, nil),nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("TracksIdGet method not implemented")
}
