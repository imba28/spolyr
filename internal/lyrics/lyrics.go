package lyrics

import (
	"errors"
	"fmt"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/model"
	lyrics "github.com/rhnvrm/lyric-api-go"
	"log"
	"strings"
)

var (
	ErrBusy = errors.New("sync already started")
)

type Syncer struct {
	ready                   chan struct{}
	syncLyricsTracksCurrent int
	syncLyricsTrackTotal    int
	syncLog                 []string

	lyricsFetcher lyrics.Lyric
	db            *db.Repositories
}

func New(db *db.Repositories, geniusAPIToken string) *Syncer {
	return &Syncer{
		syncLyricsTracksCurrent: -1,
		ready:                   make(chan struct{}, 1),
		db:                      db,

		lyricsFetcher: lyrics.New(
			lyrics.WithGeniusLyrics(geniusAPIToken),
			lyrics.WithSongLyrics(),
			lyrics.WithMusixMatch(),
		),
	}
}

func (s *Syncer) Start(tracks []*model.Track) error {
	select {
	case s.ready <- struct{}{}:
		s.run(tracks)
		return nil
	default:
		return ErrBusy
	}
}

func (s *Syncer) Syncing() bool {
	return s.syncLyricsTracksCurrent > -1
}

func (s *Syncer) SyncedTracks() int {
	return s.syncLyricsTracksCurrent
}

func (s *Syncer) TotalTracks() int {
	return s.syncLyricsTrackTotal
}

func (s *Syncer) Logs() string {
	return strings.Join(s.syncLog, "<br>")
}

func (s *Syncer) run(tracks []*model.Track) {
	s.syncLyricsTracksCurrent = 0
	s.syncLyricsTrackTotal = len(tracks)

	go func() {
		defer func() {
			s.syncLyricsTracksCurrent = -1
			s.syncLog = nil
			<-s.ready
		}()

		for i := range tracks {
			s.downloadLyrics(tracks[i])
		}
	}()
}

func (s *Syncer) downloadLyrics(t *model.Track) {
	s.syncLyricsTracksCurrent++

	artist := t.Artist
	if strings.Index(t.Artist, ", ") > -1 {
		artist = strings.Split(artist, ", ")[0]
	}
	lyric, err := s.lyricsFetcher.Search(artist, t.Name)
	if err != nil {
		s.syncLog = append(s.syncLog, fmt.Sprintf("%s - %s: %s", artist, t.Name, err.Error()))
		log.Println(artist, t.Name, err)
		return
	}

	t.Lyrics = lyric
	t.Loaded = true
	err = s.db.Tracks.Save(t)
	if err != nil {
		s.syncLog = append(s.syncLog, fmt.Sprintf("%s - %s: %s", artist, t.Name, err.Error()))
		log.Println(artist, t.Name, err)
	}
}
