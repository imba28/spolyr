package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/lyrics"
	"github.com/imba28/spolyr/internal/spotify"
	template2 "github.com/imba28/spolyr/internal/template"
	spotify2 "github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"html/template"
	"math"
	"net/http"
	"strings"
)

var (
	ErrNotFound = errors.New("item not found")
)

func HomePageHandler(s db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		trackCount, _ := s.Count()
		tracksWithLyrics, _ := s.CountWithLyrics()
		latestTracks, _ := s.LatestTracks(8)

		viewData := mergeView(gin.H{
			"TrackCount":            trackCount,
			"TracksWithLyricsCount": tracksWithLyrics,
			"TracksLatest":          latestTracks,
		}, viewFromContext(c))
		_ = template2.HomePage(c.Writer, viewData, http.StatusOK)
	}
}

func ImportHandler(s *lyrics.Syncer) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		token := c.GetString(spotifyTokenKey)
		var tok oauth2.Token
		err := json.Unmarshal([]byte(token), &tok)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/logout")
			return
		}
		client := auth.NewClient(&tok)

		var playlists []spotify2.SimplePlaylist
		var savedTracksCount int

		p, err := client.CurrentUsersPlaylists()
		if err == nil {
			playlists = p.Playlists
		} else {
			session.AddFlash("An error occurred while connecting with the Spotify API.", "Error")
		}

		pp, err := client.CurrentUsersTracks()
		if err == nil {
			savedTracksCount = pp.Total
		} else {
			session.AddFlash("An error occurred while connecting with the Spotify API.", "Error")
		}

		viewData := mergeView(gin.H{
			"Syncing":           s.Syncing(),
			"SyncedTracks":      s.SyncedTracks(),
			"TotalTracksToSync": s.TotalTracks(),
			"SyncProgressValue": math.Round(float64(s.SyncedTracks()) / float64(s.TotalTracks()) * 100),
			"LibraryTrackCount": savedTracksCount,
			"Playlists": playlists,
			"Success": session.Flashes("Success"),
			"Error":   session.Flashes("Error"),
		}, viewFromContext(c))

		_ = session.Save()
		_ = template2.ImportPage(c.Writer, viewData, http.StatusOK)
	}
}

func ImportPlaylistHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		defer func() {
			session.Save()
			c.Redirect(http.StatusFound, "/import")
		}()

		token := c.GetString(spotifyTokenKey)
		var tok oauth2.Token
		err := json.Unmarshal([]byte(token), &tok)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/logout")
			return
		}
		client := auth.NewClient(&tok)

		err = spotify.NewPlaylistProvider(client, db).Download(c.Param("ID"))
		if err != nil {
			session.AddFlash(fmt.Sprintf("Could not import playlist: %s", err), "Error")
			return
		}

		session.AddFlash("Successfully imported all tracks of playlist! Next, you can download the lyrics.", "Success")
	}
}

func TrackDetailHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		track, err := db.FindTrack(c.Param("spotifyID"))
		if err == mongo.ErrNoDocuments {
			c.Status(http.StatusNotFound)
			_ = c.Error(ErrNotFound)
			return
		}
		if err != nil {
			_ = c.Error(err)
			return
		}

		session := sessions.Default(c)
		viewData := mergeView(gin.H{
			"Track":   track,
			"Success": session.Flashes("Success"),
			"Error":   session.Flashes("Error"),
		}, viewFromContext(c))
		_ = session.Save()
		_ = template2.TrackPage(c.Writer, viewData, http.StatusOK)
	}
}

func TrackUpdateHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		track, err := db.FindTrack(c.Param("spotifyID"))
		if err != nil {
			c.Status(http.StatusNotFound)
			_ = c.Error(ErrNotFound)
			return
		}

		updatedLyrics := strings.TrimSpace(c.PostForm("lyrics"))
		view := mergeView(gin.H{
			"Track":            track,
			"TextareaRowCount": 20,
		}, viewFromContext(c))

		if len(updatedLyrics) == 0 {
			view["Error"] = "Please provide some lyrics!"
			_ = template2.TrackEditPage(c.Writer, view, http.StatusBadRequest)
			return
		}

		track.Lyrics = updatedLyrics
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
			_ = c.Error(ErrNotFound)
			return
		}
		if err != nil {
			_ = c.Error(err)
			return
		}

		viewData := mergeView(gin.H{
			"Track":            track,
			"TextareaRowCount": 20,
		}, viewFromContext(c))
		_ = template2.TrackEditPage(c.Writer, viewData, http.StatusOK)
	}
}

func TrackMissingLyricsHandler(db db.TrackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracks, err := db.TracksWithoutLyrics()
		if err != nil {
			_ = c.Error(err)
			return
		}

		viewData := mergeView(gin.H{
			"Tracks": tracks,
		}, viewFromContext(c))
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
			_ = c.Error(err)
			return
		}

		viewData := mergeView(gin.H{
			"Query":  c.Query("q"),
			"Tracks": tracks,
		}, viewFromContext(c))
		_ = template2.TracksPage(c.Writer, viewData, http.StatusOK)
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
			_ = c.Error(err)
			return
		}

		session := sessions.Default(c)
		session.AddFlash("Your library has been successfully downloaded.", "Success")
		_ = session.Save()

		c.Redirect(http.StatusTemporaryRedirect, "/import")
	}
}

func LyricsTrackSyncHandler(db db.TrackService, fetcher lyrics.Fetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		track, err := db.FindTrack(c.Param("spotifyID"))
		if err == mongo.ErrNoDocuments {
			_ = c.Error(ErrNotFound)
			return
		}
		if err != nil {
			_ = c.Error(err)
			return
		}

		session := sessions.Default(c)
		defer func() {
			_ = session.Save()
			c.Redirect(http.StatusFound, "/tracks/id/"+track.SpotifyID)
		}()

		err = fetcher.Fetch(track)
		if err != nil {
			session.AddFlash(fmt.Sprintf("An error occurred while trying to download the lyrics of this song: %s", err.Error()), "Error")
			return
		}

		err = db.Save(track)
		if err != nil {
			session.AddFlash(fmt.Sprintf("An error occurred while trying save the song: %s", err.Error()), "Error")
			return
		}

		session.AddFlash("Lyrics have been successfully downloaded!", "Success")
	}
}

func LyricsSyncLogHandler(s *lyrics.Syncer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" {
			_, err := s.Sync()
			if err != nil {
				_ = c.Error(err)
				return
			}
		}

		viewData := mergeView(gin.H{
			"Syncing":           s.Syncing(),
			"SyncedTracks":      s.SyncedTracks(),
			"TotalTracksToSync": s.TotalTracks(),
			"SyncProgressValue": math.Round(float64(s.SyncedTracks()) / float64(s.TotalTracks()) * 100),
			"SyncLog":           template.HTML(s.Logs()),
		}, viewFromContext(c))
		_ = template2.LyricsSyncLogPage(c.Writer, viewData, http.StatusOK)
	}
}

func NoRouteHandle(c *gin.Context) {
	p := viewFromContext(c)
	p["Status"] = http.StatusNotFound
	p["Message"] = "Oh no, page not found"
	_ = template2.ErrorPage(c.Writer, p, http.StatusNotFound)
}
