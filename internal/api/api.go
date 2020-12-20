package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/db"
	"github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"log"
	"strings"

	"net/http"
)

type Controller struct {
	db *db.Access
}

func (co Controller) HomePageHandler(c *gin.Context) {
	trackCount, _ := co.db.TrackCount()
	tracksWithLyrics, _ := co.db.TracksWithLyricsCount()

	latestTracks, _ := co.db.LatestTracks(8)

	session := sessions.Default(c)
	userEmail := session.Get("userEmail")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"TrackCount":            trackCount,
		"TracksWithLyricsCount": tracksWithLyrics,
		"TracksLatest":          latestTracks,
		"User":                  userEmail,
	})
}

func (co Controller) TrackDetailHandler(c *gin.Context) {
	track, err := co.db.FindTrack(c.Param("spotifyID"))
	if err == mongo.ErrNoDocuments {
		c.String(http.StatusNotFound, "track not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(c)
	userEmail := session.Get("userEmail")

	c.HTML(http.StatusOK, "track.html", gin.H{
		"Track": track,
		"User":  userEmail,
	})
}

func (co Controller) TrackUpdateHandler(c *gin.Context) {
	track, err := co.db.FindTrack(c.Param("spotifyID"))
	if err != nil {
		c.String(http.StatusNotFound, "track not found")
		return
	}

	session := sessions.Default(c)
	userEmail := session.Get("userEmail")

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
	err = co.db.SaveTrack(track)
	if err != nil {
		c.HTML(http.StatusBadRequest, "track-edit.html", gin.H{
			"Track":            track,
			"User":             userEmail,
			"TextareaRowCount": 20,
			"Error":            "Could not update lyrics",
		})
		return
	}

	c.HTML(http.StatusOK, "track.html", gin.H{
		"Track":   track,
		"User":    userEmail,
		"Success": "Lyrics of track updated!",
	})
}

func (co Controller) TrackEditFormHandler(c *gin.Context) {
	track, err := co.db.FindTrack(c.Param("spotifyID"))
	if err == mongo.ErrNoDocuments {
		c.String(http.StatusNotFound, "track not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(c)
	userEmail := session.Get("userEmail")

	c.HTML(http.StatusOK, "track-edit.html", gin.H{
		"Track":            track,
		"User":             userEmail,
		"TextareaRowCount": 20,
	})
}

func (co Controller) TrackMissingLyricsHandler(c *gin.Context) {
	tracks, err := co.db.TracksWithoutLyrics()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(c)
	userEmail := session.Get("userEmail")

	c.HTML(http.StatusOK, "tracks.html", gin.H{
		"Tracks": tracks,
		"User":   userEmail,
	})
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

	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{
			"$text": bson.M{
				"$search": query,
			},
		}
	}

	tracks, err := co.db.FindTracks(filter)
	if err != nil && err != mongo.ErrNoDocuments {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(c)
	userEmail := session.Get("userEmail")

	c.HTML(http.StatusOK, "search.html", gin.H{
		"Query":  c.Query("q"),
		"Tracks": tracks,
		"User":   userEmail,
	})
}

const redirectURI = "http://localhost:8080/callback"
const state = "spolyrCSRF"

var auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserLibraryRead, spotify.ScopeUserReadEmail)

func (co Controller) SpotifyAuthCallbackHandler(c *gin.Context) {
	tok, err := auth.Token(state, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if st := c.Request.FormValue("state"); st != state {
		c.String(http.StatusNotFound, "Invalid csrf token")
		return
	}

	client := auth.NewClient(tok)
	user, err := client.CurrentUser()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	token, _ := json.Marshal(tok)
	session := sessions.Default(c)
	session.Set("token", string(token))
	session.Set("userEmail", user.Email)
	err = session.Save()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (co Controller) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (co Controller) LoginHandler(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, auth.AuthURL(state))
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

func (co Controller) LyricsSyncHandler(c *gin.Context) {
	tracks, err := co.db.FindTracks(nil)
	if err != nil {
		c.Abort()
		c.Error(err)
		return
	}

	l := lyrics.New(
		lyrics.WithGeniusLyrics("DBGzQI4tQoQ3sBTXbHXI1Yxa1GKWbOIJc3u84VNYQxJLqNXfDXX3p88_Ix7xAwbi"),
		lyrics.WithMusixMatch(),
		lyrics.WithSongLyrics(),
	)

	for i := range tracks {
		if len(tracks[i].Lyrics) > 0 {
			continue
		}

		artist := tracks[i].Artist
		if strings.Index(tracks[i].Artist, ", ") > -1 {
			artist = strings.Split(artist, ", ")[0]
		}
		lyric, err := l.Search(artist, tracks[i].Name)
		if err != nil {
			log.Println(artist, tracks[i].Name, err)
			continue
		}
		tracks[i].Lyrics = lyric
		tracks[i].Loaded = true
		err = co.db.SaveTrack(tracks[i])
		if err != nil {
			c.Error(err)
			return
		}
		fmt.Println(artist, tracks[i].Name, "SAVED")
	}
	c.String(http.StatusOK, "OK")
}

func New(db *db.Access) Controller {
	return Controller{db}
}
