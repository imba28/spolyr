package api

import (
	"encoding/json"
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
		"User":                  c.GetString(userEmailKey),
	}
	_ = template2.HomePage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackDetailHandler(c *gin.Context) {
	track, err := co.db.Tracks.FindTrack(c.Param("spotifyID"))
	if err == mongo.ErrNoDocuments {
		c.Error(ErrNotFound)
		return
	}
	if err != nil {
		c.Error(err)
		return
	}

	viewData := gin.H{
		"Track": track,
		"User":  c.GetString(userEmailKey),
	}
	_ = template2.TrackPage(c.Writer, viewData, http.StatusOK)
}

func (co Controller) TrackUpdateHandler(c *gin.Context) {
	track, err := co.db.Tracks.FindTrack(c.Param("spotifyID"))
	if err != nil {
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
		c.HTML(http.StatusBadRequest, "track-edit.html", view)
		return
	}

	track.Lyrics = lyrics
	track.Loaded = true
	err = co.db.Tracks.Save(track)
	if err != nil {
		view["Error"] = "Could not update lyrics"
		c.HTML(http.StatusBadRequest, "track-edit.html", view)
		return
	}

	view["Success"] = "Lyrics of track updated!"
	_ = template2.TrackPage(c.Writer, view, http.StatusOK)
}

func (co Controller) TrackEditFormHandler(c *gin.Context) {
	track, err := co.db.Tracks.FindTrack(c.Param("spotifyID"))
	if err == mongo.ErrNoDocuments {
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

func (co Controller) TrackMissingLyricsHandler(c *gin.Context) {
	tracks, err := co.db.Tracks.TracksWithoutLyrics()
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

func (co Controller) TrackSyncHandler(c *gin.Context) {
	token := c.GetString(spotifyTokenKey)
	var tok oauth2.Token
	err := json.Unmarshal([]byte(token), &tok)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/logout")
		return
	}

	err = SyncTracks(auth.NewClient(&tok), co.db)
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (co *Controller) LyricsSyncHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		tracks, err := co.db.Tracks.TracksWithoutLyrics()
		if err != nil {
			c.Error(err)
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
		"User":              c.GetString(userEmailKey),
	}
	_ = template2.TrackLyricsSyncPage(c.Writer, viewData, http.StatusOK)
}
