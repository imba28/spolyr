package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"net/http"
)

func New(db *db.Repositories, geniusAPIToken string, sessionKey []byte) *gin.Engine {
	r := gin.Default()

	fetcher := lyrics.New(geniusAPIToken, 3)
	syncer := lyrics.NewSyncer(fetcher, db.Tracks)

	store := cookie.NewStore(sessionKey)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("session", store))
	r.Use(UserProviderMiddleware)
	r.Use(ErrorHandle)

	authRequired := r.Group("/").Use(AuthRequired)
	{
		authRequired.GET("/import", ImportHandler(syncer))
		authRequired.POST("/import/playlist/:ID", ImportPlaylistHandler(db.Tracks))
		authRequired.GET("/sync-tracks", TracksSyncHandler(db.Tracks))
		authRequired.GET("/sync-lyrics", LyricsSyncLogHandler(syncer))
		authRequired.POST("/sync-lyrics", LyricsSyncLogHandler(syncer))
		authRequired.GET("/tracks/id/:spotifyID/edit", TrackEditFormHandler(db.Tracks))
		authRequired.POST("/tracks/id/:spotifyID/edit", TrackUpdateHandler(db.Tracks))
		authRequired.POST("/tracks/id/:spotifyID/sync", LyricsTrackSyncHandler(db.Tracks, fetcher))
	}

	r.GET("/", HomePageHandler(db.Tracks))
	r.GET("/login", LoginHandler)
	r.GET("/logout", LogoutHandler)
	r.GET("/callback", SpotifyAuthCallbackHandler)
	r.GET("/tracks/id/:spotifyID", TrackDetailHandler(db.Tracks))
	r.GET("/tracks/missing-lyrics", TrackMissingLyricsHandler(db.Tracks))
	r.GET("/tracks/no-lyrics-found", TrackNoLyricsFoundHandler(db.Tracks))
	r.GET("/search", TrackSearchHandler(db.Tracks))
	r.Static("static", "public")

	r.NoRoute(NoRouteHandle)

	return r
}
