package api

import (
	"github.com/gorilla/mux"
	"github.com/imba28/spolyr/pkg/db"
	jwt2 "github.com/imba28/spolyr/pkg/jwt"
	"github.com/imba28/spolyr/pkg/language"
	"github.com/imba28/spolyr/pkg/lyrics"
	"github.com/imba28/spolyr/pkg/openapi/openapi"
	"github.com/rs/cors"
	"net/http"
	"os"
	"path/filepath"
)

type languageDetector interface {
	Detect(string) (string, error)
}

func NewOAPI(db *db.Repositories, oauthClientId, geniusAPIToken string, secret []byte, d languageDetector) http.Handler {
	fetcher := lyrics.New(geniusAPIToken, 3, d)
	syncer := lyrics.NewSyncer(fetcher, db.Tracks)

	authApiController := openapi.NewAuthApiController(newAuthApiService(oauthClientId, secret))
	importController := openapi.NewImportApiController(newImportApiService(db.Tracks, syncer, fetcher, d))
	tracksApiController := openapi.NewTracksApiController(newTracksApiService(db.Tracks, d))
	playlistController := openapi.NewPlaylistsApiController(newPlaylistApiService())

	r := openapi.NewRouter(authApiController, tracksApiController, importController, playlistController)

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"https://localhost:8081", "https://127.0.0.1:8081"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
		MaxAge:           3600,
		Debug:            true,
	})
	return AuthenticationMiddleware(jwt2.New(secret))(c.Handler(r))
}

func spaFileHandler(publicFolder string) http.HandlerFunc {
	const indexPath = "index.html"
	staticHandler := http.FileServer(http.Dir(publicFolder))

	return func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		path = filepath.Join(publicFolder, path)

		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(publicFolder, indexPath))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		staticHandler.ServeHTTP(w, r)
	}
}

func New(db *db.Repositories, oauthClientId, geniusAPIToken string, secret []byte, d language.Detector) http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/api").Handler(NewOAPI(db, geniusAPIToken, oauthClientId, secret, d))
	r.PathPrefix("/").Handler(spaFileHandler("public"))
	return r
}