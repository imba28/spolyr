package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"github.com/imba28/spolyr/internal/model"
	"github.com/imba28/spolyr/internal/spotify"
	template2 "github.com/imba28/spolyr/internal/template"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"html/template"
	"math"
	"net/http"
	"strings"
)

func HomePageHandler(s db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		trackCount, _ := s.Count()
		tracksWithLyrics, _ := s.CountWithLyrics()
		latestTracks, _ := s.LatestTracks(8)

		viewData := gin.H{
			"TrackCount":            trackCount,
			"TracksWithLyricsCount": tracksWithLyrics,
			"TracksLatest":          latestTracks,
			"User":                  c.GetString(userEmailKey),
		}
		_ = template2.HomePage(c.Writer, viewData, http.StatusOK)
	}
}

func TrackDetailHandler(db db.TrackService, s *lyrics.Syncer) gin.HandlerFunc {
	return func(c *gin.Context) {
		track, err := db.FindTrack(c.Param("spotifyID"))
		if err == mongo.ErrNoDocuments {
			c.Status(http.StatusNotFound)
			c.Error(ErrNotFound)
			return
		}
		if err != nil {
			c.Error(err)
			return
		}

		session := sessions.Default(c)
		viewData := gin.H{
			"Track":         track,
			"SyncAvailable": !s.Syncing(),
			"User":          c.GetString(userEmailKey),
			"Success":       session.Flashes("Success"),
			"Error":         session.Flashes("Error"),
		}
		_ = session.Save()
		_ = template2.TrackPage(c.Writer, viewData, http.StatusOK)
	}
}

func TrackUpdateHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		track, err := db.FindTrack(c.Param("spotifyID"))
		if err != nil {
			c.Status(http.StatusNotFound)
			c.Error(ErrNotFound)
			return
		}

		userEmail := c.GetString(userEmailKey)
		lyrics := strings.TrimSpace(c.PostForm("lyrics"))
		view := gin.H{
			"Track":            track,
			"User":             userEmail,
			"TextareaRowCount": 20,
		}

		if len(lyrics) == 0 {
			view["Error"] = "Please provide some lyrics!"
			_ = template2.TrackEditPage(c.Writer, view, http.StatusBadRequest)
			return
		}

		track.Lyrics = lyrics
		track.Loaded = true
		err = db.Save(track)
		if err != nil {
			view["Error"] = "Could not update lyrics"
			_ = template2.TrackEditPage(c.Writer, view, http.StatusInternalServerError)
			return
		}

		session := sessions.Default(c)
		session.AddFlash("Lyrics of track updated!", "Success")
		_ = session.Save()

		c.Redirect(http.StatusFound, "/tracks/id/"+track.SpotifyID)
	}
}

func TrackEditFormHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		track, err := db.FindTrack(c.Param("spotifyID"))
		if err == mongo.ErrNoDocuments {
			c.Status(http.StatusNotFound)
			c.Error(ErrNotFound)
			return
		}
		if err != nil {
			c.Error(err)
			return
		}

		viewData := gin.H{
			"Track":            track,
			"User":             c.GetString(userEmailKey),
			"TextareaRowCount": 20,
		}
		_ = template2.TrackEditPage(c.Writer, viewData, http.StatusOK)
	}
}

func TrackMissingLyricsHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracks, err := db.TracksWithoutLyrics()
		if err != nil {
			c.Error(err)
			return
		}

		viewData := gin.H{
			"Tracks": tracks,
			"User":   c.GetString(userEmailKey),
		}
		_ = template2.TracksPage(c.Writer, viewData, http.StatusOK)
	}
}

func TrackSearchHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")

		if strings.Index(query, " ") > -1 && strings.Index(query, "\"") == -1 && strings.Index(query, "-") == -1 {
			qs := strings.Split(query, " ")
			for i := range qs {
				qs[i] = "\"" + qs[i] + "\""
			}
			query = strings.Join(qs, " ")
		}

		tracks, err := db.Search(query)
		if err != nil && err != mongo.ErrNoDocuments {
			c.Error(err)
			return
		}

		viewData := gin.H{
			"Query":  c.Query("q"),
			"Tracks": tracks,
			"User":   c.GetString(userEmailKey),
		}
		_ = template2.SearchPage(c.Writer, viewData, http.StatusOK)
	}
}

func TracksSyncHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetString(spotifyTokenKey)
		var tok oauth2.Token
		err := json.Unmarshal([]byte(token), &tok)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/logout")
			return
		}

		err = spotify.SyncTracks(spotify.NewSpotifyTrackProvider(auth.NewClient(&tok)), db)
		if err != nil {
			c.Error(err)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func LyricsTrackSyncHandler(db db.TrackService, s *lyrics.Syncer) gin.HandlerFunc {
	return func(c *gin.Context) {
		track, err := db.FindTrack(c.Param("spotifyID"))
		if err == mongo.ErrNoDocuments {
			c.Error(ErrNotFound)
			return
		}
		if err != nil {
			c.Error(err)
			return
		}

		session := sessions.Default(c)
		defer func() {
			_ = session.Save()
			c.Redirect(http.StatusFound, "/tracks/id/"+track.SpotifyID)
		}()

		err = s.Start([]*model.Track{track})
		if err != nil && errors.Is(err, lyrics.ErrBusy) {
			session.AddFlash("Action not available. Please try again later.", "Error")
			return
		}
		if err != nil {
			session.AddFlash("An unknown error occurred", "Error")
			return
		}

		session.AddFlash("Download of song lyrics started! This might take a few seconds.", "Success")
	}
}

func LyricsSyncHandler(db db.TrackService, s *lyrics.Syncer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" {
			tracks, err := db.TracksWithoutLyrics()
			if err != nil {
				c.Error(err)
				return
			}

			s.Start(tracks)
		}

		viewData := gin.H{
			"Syncing":           s.Syncing(),
			"SyncedTracks":      s.SyncedTracks(),
			"TotalTracksToSync": s.TotalTracks(),
			"SyncProgressValue": math.Round(float64(s.SyncedTracks()) / float64(s.TotalTracks()) * 100),
			"SyncLog":           template.HTML(s.Logs()),
			"User":              c.GetString(userEmailKey),
		}
		_ = template2.TrackLyricsSyncPage(c.Writer, viewData, http.StatusOK)
	}
}
