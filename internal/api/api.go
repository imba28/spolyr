package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func New(c Controller) *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("spolyr-cookie-secret"))

	r.Use(sessions.Sessions("session", store))
	r.Use(UserProviderMiddleware)
	r.Use(ErrorHandle)

	authRequired := r.Group("/").Use(AuthRequired)
	{
		authRequired.GET("/sync-tracks", c.TrackSyncHandler)
		authRequired.GET("/sync-lyrics", c.LyricsSyncHandler)
		authRequired.POST("/sync-lyrics", c.LyricsSyncHandler)
		authRequired.POST("/tracks/id/:spotifyID/edit", c.TrackUpdateHandler)
		authRequired.GET("/tracks/id/:spotifyID/edit", c.TrackEditFormHandler)
	}

	r.GET("/", c.HomePageHandler)
	r.GET("/login", c.LoginHandler)
	r.GET("/logout", c.LogoutHandler)
	r.GET("/callback", c.SpotifyAuthCallbackHandler)
	r.GET("/tracks/id/:spotifyID", c.TrackDetailHandler)
	r.GET("/tracks/missing-lyrics", c.TrackMissingLyricsHandler)
	r.GET("/search", c.TrackSearchHandler)
	r.Static("static", "public")

	return r
}
