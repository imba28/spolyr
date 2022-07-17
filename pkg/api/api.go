package api

import (
	"github.com/gorilla/mux"
	"github.com/imba28/spolyr/pkg/db"
	jwt2 "github.com/imba28/spolyr/pkg/jwt"
	"github.com/imba28/spolyr/pkg/lyrics"
	"github.com/imba28/spolyr/pkg/openapi"
	"github.com/rs/cors"
	"net/http"
	"sync"
)

type languageDetector interface {
	Detect(string) (string, error)
}

func (s *Server) apiHandler() http.Handler {
	fetcher := lyrics.New(s.geniusAPIToken, 3, s.languageDetector)
	syncer := lyrics.NewSyncer(fetcher, s.db.Tracks)

	authApiController := openapi.NewAuthApiController(newAuthApiService(s.oauthClientID, s.oauthClientSecret, s.secret, s.publicProtocol, s.publicDomain, s.publicHttpPort))
	importController := openapi.NewImportApiController(newImportApiService(s.db.Tracks, syncer, fetcher, s.languageDetector))
	tracksApiController := openapi.NewTracksApiController(newTracksApiService(s.db.Tracks, s.languageDetector))
	playlistController := openapi.NewPlaylistsApiController(newPlaylistApiService())

	r := openapi.NewRouter(authApiController, tracksApiController, importController, playlistController)

	var handler http.Handler = r

	if s.env == Dev {
		c := cors.New(cors.Options{
			AllowCredentials: true,
			AllowedOrigins:   []string{"https://localhost:8081", "https://127.0.0.1:8081"},
			AllowedHeaders:   []string{"User-Agent", "Content-Type"},
			AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
			MaxAge:           3600,
			Debug:            true,
		})
		handler = c.Handler(r)
	}

	return AuthenticationMiddleware(jwt2.New(s.secret))(handler)
}

type Server struct {
	db                *db.Repositories
	oauthClientID     string
	oauthClientSecret string
	geniusAPIToken    string
	secret            []byte
	languageDetector  languageDetector

	env    Env
	router *mux.Router

	publicHttpPort int
	publicDomain   string
	publicProtocol string

	sync.Once
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Do(func() {
		s.init()
	})
	s.router.ServeHTTP(w, r)
}

func (s *Server) init() {
	s.router.PathPrefix("/api").Handler(s.apiHandler())
	s.router.PathPrefix("/").Handler(spaFileHandler("public"))
}

func NewServer(options ...ServerOptions) *Server {
	s := Server{
		secret: []byte("not so secret. change me"),
		router: mux.NewRouter(),
		env:    Prod,

		publicDomain:   "localhost",
		publicProtocol: "http",
		publicHttpPort: 8080,
	}

	for i := range options {
		options[i](&s)
	}

	return &s
}

type ServerOptions func(s *Server)

func WithOAuth(clientId string, clientSecret string) ServerOptions {
	return func(s *Server) {
		s.oauthClientID = clientId
		s.oauthClientSecret = clientSecret
	}
}
func WithGeniusAPI(token string) ServerOptions {
	return func(s *Server) {
		s.geniusAPIToken = token
	}
}
func WithSecret(secret []byte) ServerOptions {
	return func(s *Server) {
		s.secret = secret
	}
}
func WithDatabase(repositories *db.Repositories) ServerOptions {
	return func(s *Server) {
		s.db = repositories
	}
}
func WithLanguageDetector(detector languageDetector) ServerOptions {
	return func(s *Server) {
		s.languageDetector = detector
	}
}

func WithEnv(env Env) ServerOptions {
	return func(s *Server) {
		s.env = env
	}
}

func WithReverseProxy(protocol, domain string, port int) ServerOptions {
	return func(s *Server) {
		s.publicProtocol = protocol
		s.publicDomain = domain
		s.publicHttpPort = port
	}
}

type Env int

const (
	Prod Env = iota
	Dev
)
