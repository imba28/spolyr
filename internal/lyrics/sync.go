package lyrics

import (
	"fmt"
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/model"
	"strings"
)

type Syncer struct {
	ready                   chan struct{}
	syncLyricsTracksCurrent int
	syncLyricsTrackTotal    int
	syncLog                 []string

	fetcher Fetcher
	db      db.TrackService
}

func (s *Syncer) Sync() error {
	tracks, err := s.db.TracksWithoutLyrics()
	if err != nil {
		return err
	}

	select {
	case s.ready <- struct{}{}:
		s.syncLyricsTracksCurrent = 0
		s.syncLyricsTrackTotal = len(tracks)

		go s.run(tracks)
		return nil
	default:
		return ErrBusy
	}
}

func (s *Syncer) run(tracks []*model.Track) {
	defer func() {
		s.syncLyricsTracksCurrent = -1
		s.syncLog = nil
		<-s.ready
	}()

	c, err := s.fetcher.FetchAll(tracks)
	if err != nil {
		return
	}

	for result := range c {
		s.syncLyricsTracksCurrent++

		if result.err != nil {
			s.syncLog = append(s.syncLog, fmt.Sprintf("\xE2\x9D\x8C %s - %s: %s", result.track.Artist, result.track.Name, result.err.Error()))
		} else {
			err = s.db.Save(result.track)
			if err != nil {
				s.syncLog = append(s.syncLog, fmt.Sprintf("\xE2\x9D\x8C %s - %s: %s", result.track.Name, result.track.Artist, err.Error()))
			} else {
				s.syncLog = append(s.syncLog, fmt.Sprintf("\xE2\x9C\x85 %s - %s", result.track.Name, result.track.Artist))
			}
		}
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

func NewSyncer(fetcher Fetcher, db db.TrackService) *Syncer {
	return &Syncer{
		ready:                   make(chan struct{}, 1),
		syncLyricsTracksCurrent: -1,
		fetcher:                 fetcher,
		db:                      db,
	}
}
