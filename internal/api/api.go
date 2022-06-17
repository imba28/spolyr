package api

import (
	"github.com/gorilla/mux"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"net/http"
)

func NewOAPI(db *db.Repositories, geniusAPIToken string) http.Handler {
	fetcher := lyrics.New(geniusAPIToken, 3)
	syncer := lyrics.NewSyncer(fetcher, db.Tracks)

	AuthApiController := openapi.NewAuthApiController(NewAuthApiService())
	ImportController := openapi.NewImportApiController(NewImportApiService(db.Tracks, syncer, fetcher))
	TracksApiController := openapi.NewTracksApiController(NewTracksApiService(db.Tracks))

	return openapi.NewRouter(AuthApiController, TracksApiController, ImportController)
}

func New(db *db.Repositories, geniusAPIToken string) http.Handler {
	r := mux.NewRouter()

	// fetcher := lyrics.New(geniusAPIToken, 3)
	// syncer := lyrics.NewSyncer(fetcher, db.Tracks)

	r.PathPrefix("/api").Handler(NewOAPI(db, geniusAPIToken))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	return r
}
