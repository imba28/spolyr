package api

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	template2 "github.com/imba28/spolyr/internal/template"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"html/template"
	"math"
	"net/http"
	"strings"
)

func (co Controller) HomePageHandler(c *gin.Context) {
	trackCount, _ := co.db.Tracks.Count()
	tracksWithLyrics, _ := co.db.Tracks.CountWithLyrics()
	latestTracks, _ := co.db.Tracks.LatestTracks(8)

	viewData := gin.H{
		"TrackCount":            trackCount,
		"TracksWithLyricsCount": tracksWithLyrics,
		"TracksLatest":          latestTracks,
		"User":                  sessions.Default(c).Get("userEmail"),
	}
	_ = template2.HomePage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackDetailHandler(c *gin.Context) {
	track, err := co.db.Tracks.FindTrack(c.Param("spotifyID"))
	if err == mongo.ErrNoDocuments {
		c.String(http.StatusNotFound, "track not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	viewData := gin.H{
		"Track": track,
		"User":  sessions.Default(c).Get("userEmail"),
	}
	_ = template2.TrackPage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackUpdateHandler(c *gin.Context) {
	track, err := co.db.Tracks.FindTrack(c.Param("spotifyID"))
	if err != nil {
		c.String(http.StatusNotFound, "track not found")
		return
	}

	userEmail := sessions.Default(c).Get("userEmail")

	lyrics := strings.TrimSpace(c.PostForm("lyrics"))
	if len(lyrics) == 0 {
		c.HTML(http.StatusBadRequest, "track-edit.html", gin.H{
			"Track":            track,
			"User":             userEmail,
			"TextareaRowCount": 20,
			"Error":            "Please provide some lyrics!",
		})
		return
	}
	track.Lyrics = lyrics
	track.Loaded = true
	err = co.db.Tracks.Save(track)
	if err != nil {
		c.HTML(http.StatusBadRequest, "track-edit.html", gin.H{
			"Track":            track,
			"User":             userEmail,
			"TextareaRowCount": 20,
			"Error":            "Could not update lyrics",
		})
		return
	}

	viewData := gin.H{
		"Track":   track,
		"User":    userEmail,
		"Success": "Lyrics of track updated!",
	}
	_ = template2.TrackPage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackEditFormHandler(c *gin.Context) {
	track, err := co.db.Tracks.FindTrack(c.Param("spotifyID"))
	if err == mongo.ErrNoDocuments {
		c.String(http.StatusNotFound, "track not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	viewData := gin.H{
		"Track":            track,
		"User":             sessions.Default(c).Get("userEmail"),
		"TextareaRowCount": 20,
	}
	_ = template2.TrackEditPage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackMissingLyricsHandler(c *gin.Context) {
	tracks, err := co.db.Tracks.TracksWithoutLyrics()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	viewData := gin.H{
		"Tracks": tracks,
		"User":   sessions.Default(c).Get("userEmail"),
	}
	_ = template2.TracksPage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackSearchHandler(c *gin.Context) {
	query := c.Query("q")

	if strings.Index(query, " ") > -1 && strings.Index(query, "\"") == -1 && strings.Index(query, "-") == -1 {
		qs := strings.Split(query, " ")
		for i := range qs {
			qs[i] = "\"" + qs[i] + "\""
		}
		query = strings.Join(qs, " ")
	}

	tracks, err := co.db.Tracks.Search(query)
	if err != nil && err != mongo.ErrNoDocuments {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	viewData := gin.H{
		"Query":  c.Query("q"),
		"Tracks": tracks,
		"User":   sessions.Default(c).Get("userEmail"),
	}
	_ = template2.SearchPage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackSyncHandler(c *gin.Context) {
	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {
		c.String(http.StatusUnauthorized, "Authentication required")
		return
	}

	var tok oauth2.Token
	err := json.Unmarshal([]byte(token), &tok)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/logout")
		return
	}

	err = SyncTracks(auth.NewClient(&tok), co.db)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (co *Controller) LyricsSyncHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		tracks, err := co.db.Tracks.TracksWithoutLyrics()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		co.startLyricsSync(tracks)
	}

	viewData := gin.H{
		"Syncing":           co.syncLyricsTracksCurrent > -1,
		"SyncedTracks":      co.syncLyricsTracksCurrent,
		"TotalTracksToSync": co.syncLyricsTrackTotal,
		"SyncProgressValue": math.Round(float64(co.syncLyricsTracksCurrent) / float64(co.syncLyricsTrackTotal) * 100),
		"SyncLog":           template.HTML(strings.Join(co.syncLog, "<br>")),
		"User":              sessions.Default(c).Get("userEmail"),
	}
	_ = template2.TrackLyricsSyncPage(c.Writer, viewData, http.StatusOK)
}
