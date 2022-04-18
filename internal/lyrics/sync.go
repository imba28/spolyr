package lyrics

import (
	"fmt"
	"github.com/imba28/spolyr/internal/db"
	"strings"
)

type tracksSyncFetcherSaver interface {
	Save(track *db.Track) error
	TracksWithoutLyricsError() ([]*db.Track, error)
}

type Syncer struct {
	ready                   chan struct{}
	syncLyricsTracksCurrent int
	syncLyricsTrackTotal    int
	syncLog                 []string

	fetcher Fetcher
	db      tracksSyncFetcherSaver
}

func (s *Syncer) Sync() (<-chan struct{}, error) {
	tracks, err := s.db.TracksWithoutLyricsError()
	if err != nil {
		return nil, err
	}

	finished := make(chan struct{})

	select {
	case s.ready <- struct{}{}:
		s.syncLyricsTracksCurrent = 0
		s.syncLyricsTrackTotal = len(tracks)

		go s.run(tracks, finished)
		return finished, nil
	default:
		return nil, ErrBusy
	}
}

func (s *Syncer) run(tracks []*db.Track, finishedSignal chan<- struct{}) {
	defer func() {
		// Do not block if no one is waiting for us to end.
		select {
		case finishedSignal <- struct{}{}:
		default:
		}
		close(finishedSignal)

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

		if result.Err != nil {
			s.syncLog = append(s.syncLog, fmt.Sprintf("\xE2\x9D\x8C %s - %s: %s", result.Track.Artist, result.Track.Name, result.Err.Error()))
			result.Track.LyricsImportErrorCount++
		} else {
			result.Track.LyricsImportErrorCount = 0
		}

		err = s.db.Save(result.Track)
		if result.Err != nil || err != nil {
			message := result.Err
			if err != nil {
				message = err
			}
			s.syncLog = append(s.syncLog, fmt.Sprintf("\xE2\x9D\x8C %s - %s: %s", result.Track.Name, result.Track.Artist, message.Error()))
		} else {
			s.syncLog = append(s.syncLog, fmt.Sprintf("\xE2\x9C\x85 %s - %s", result.Track.Name, result.Track.Artist))
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
	b := strings.Builder{}
	for i := len(s.syncLog) - 1; i >= 0; i-- {
		b.WriteString(s.syncLog[i] + "<br>")
	}

	return b.String()
}

func NewSyncer(fetcher Fetcher, db tracksSyncFetcherSaver) *Syncer {
	return &Syncer{
		ready:                   make(chan struct{}, 1),
		syncLyricsTracksCurrent: -1,
		fetcher:                 fetcher,
		db:                      db,
	}
}
