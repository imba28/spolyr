package api

import (
	"fmt"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/model"
	lyrics "github.com/rhnvrm/lyric-api-go"
	"log"
	"strings"
)

type Controller struct {
	syncLyrics              chan bool
	syncLyricsTracksCurrent int
	syncLyricsTrackTotal    int
	syncLog                 []string

	db *db.Repositories
}

func (co *Controller) startLyricsSync(tracks []*model.Track) {
	select {
	case co.syncLyrics <- true:
		co.syncLyricsTracksCurrent = 0
		co.syncLyricsTrackTotal = len(tracks)

		go func() {
			defer func() {
				<-co.syncLyrics
				co.syncLyricsTracksCurrent = -1
				co.syncLog = nil
			}()

			l := lyrics.New(
				lyrics.WithGeniusLyrics("DBGzQI4tQoQ3sBTXbHXI1Yxa1GKWbOIJc3u84VNYQxJLqNXfDXX3p88_Ix7xAwbi"),
				lyrics.WithSongLyrics(),
				lyrics.WithMusixMatch(),
			)

			for i := range tracks {
				co.syncLyricsTracksCurrent++

				artist := tracks[i].Artist
				if strings.Index(tracks[i].Artist, ", ") > -1 {
					artist = strings.Split(artist, ", ")[0]
				}
				lyric, err := l.Search(artist, tracks[i].Name)
				if err != nil {
					co.syncLog = append(co.syncLog, fmt.Sprintf("%s - %s: %s", artist, tracks[i].Name, err.Error()))
					log.Println(artist, tracks[i].Name, err)
					continue
				}

				tracks[i].Lyrics = lyric
				tracks[i].Loaded = true
				err = co.db.Tracks.Save(tracks[i])
				if err != nil {
					co.syncLog = append(co.syncLog, fmt.Sprintf("%s - %s: %s", artist, tracks[i].Name, err.Error()))
					log.Println(artist, tracks[i].Name, err)
				}
			}
		}()
	default:
		//
	}
}

func NewController(db *db.Repositories) Controller {
	return Controller{
		db:                      db,
		syncLyricsTracksCurrent: -1,
		syncLyrics:              make(chan bool, 1),
	}
}
