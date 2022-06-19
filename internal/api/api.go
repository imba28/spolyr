package api

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const authKey = "auth"

type customClaims struct {
	Token oauth2.Token `json:"oauth_token"`
	jwt.RegisteredClaims
}

func tokenFromContext(ctx context.Context) *oauth2.Token {
	t, _ := ctx.Value(authKey).(oauth2.Token)
	return &t
}

func AuthenticationMiddleware(secret []byte) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if r.Method != http.MethodOptions {
				if c, err := r.Cookie("jwt"); err == nil {
					token, _ := jwt.ParseWithClaims(c.Value, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
						return secret, nil
					})
					if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
						log.Println(claims.Token)
						ctx = context.WithValue(ctx, authKey, claims.Token)
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewOAPI(db *db.Repositories, oauthClientId, geniusAPIToken string, secret []byte) http.Handler {
	fetcher := lyrics.New(geniusAPIToken, 3)
	syncer := lyrics.NewSyncer(fetcher, db.Tracks)

	AuthApiController := openapi.NewAuthApiController(NewAuthApiService(oauthClientId, secret))
	ImportController := openapi.NewImportApiController(NewImportApiService(db.Tracks, syncer, fetcher))
	TracksApiController := openapi.NewTracksApiController(NewTracksApiService(db.Tracks))
	r := openapi.NewRouter(AuthApiController, TracksApiController, ImportController)

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"https://localhost:8081", "http://localhost:8081", "http://127.0.0.1:8081"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type"},
		Debug:            true,
	})
	return AuthenticationMiddleware(secret)(c.Handler(r))
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

func New(db *db.Repositories, oauthClientId, geniusAPIToken string, secret []byte) http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/api").Handler(NewOAPI(db, geniusAPIToken, oauthClientId, secret))
	r.PathPrefix("/").Handler(spaFileHandler("public"))
	return r
}
