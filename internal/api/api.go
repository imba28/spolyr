package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	template2 "github.com/imba28/spolyr/internal/template"
	"net/http"
)

func New(db *db.Repositories, geniusAPIToken string) *gin.Engine {
	r := gin.Default()

	fetcher := lyrics.New(geniusAPIToken, 3)
	syncer := lyrics.NewSyncer(fetcher, db.Tracks)

	store := cookie.NewStore([]byte("spolyr-cookie-secret"))
	r.Use(sessions.Sessions("session", store))
	r.Use(UserProviderMiddleware)
	r.Use(ErrorHandle)

	authRequired := r.Group("/").Use(AuthRequired)
	{
		authRequired.GET("/sync-tracks", TracksSyncHandler(db.Tracks))
		authRequired.GET("/sync-lyrics", LyricsSyncHandler(syncer))
		authRequired.POST("/sync-lyrics", LyricsSyncHandler(syncer))
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
	r.GET("/search", TrackSearchHandler(db.Tracks))
	r.Static("static", "public")

	r.NoRoute(func(c *gin.Context) {
		p := viewFromContext(c)
		p["Status"] = http.StatusNotFound
		p["Message"] = "Oh no, page not found"
		_ = template2.ErrorPage(c.Writer, p, http.StatusNotFound)
	})

	return r
}
