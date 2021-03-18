package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/api"
	"github.com/imba28/spolyr/internal/db"
	"log"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	databaseUsername := db.GetEnv("DATABASE_USER", "root")
	databasePassword := db.GetEnv("DATABASE_PASSWORD", "example")
	databaseHost := db.GetEnv("DATABASE_HOST", "127.0.0.1")

	dbConn, err := db.New(databaseUsername, databasePassword, "spolyr", databaseHost)
	if err != nil {
		log.Fatal(err)
	}

	store := cookie.NewStore([]byte("spolyr-cookie-secret"))
	r.Use(sessions.Sessions("session", store))

	controller := api.New(dbConn)

	r.GET("/", controller.HomePageHandler)
	r.GET("/login", controller.LoginHandler)
	r.GET("/logout", controller.LogoutHandler)
	r.GET("/callback", controller.SpotifyAuthCallbackHandler)
	r.GET("/sync-tracks", controller.TrackSyncHandler)
	r.GET("/sync-lyrics", controller.LyricsSyncHandler)
	r.POST("/sync-lyrics", controller.LyricsSyncHandler)
	r.GET("/tracks/id/:spotifyID", controller.TrackDetailHandler)
	r.POST("/tracks/id/:spotifyID/edit", controller.TrackUpdateHandler)
	r.GET("/tracks/id/:spotifyID/edit", controller.TrackEditFormHandler)
	r.GET("/tracks/missing-lyrics", controller.TrackMissingLyricsHandler)
	r.GET("/search", controller.TrackSearchHandler)
	r.Static("static", "public")

	return r
}

func main() {
	r := setupRouter()
	log.Fatal(r.Run(":8080"))
}
