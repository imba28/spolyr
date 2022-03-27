package lyrics

import (
	"errors"
	lyrics "github.com/imba28/lyric-api-go"
	"github.com/imba28/spolyr/internal/db"
	"strings"
	"sync"
)

var (
	ErrBusy = errors.New("sync already started")
)

type Result struct {
	Track *db.Track
	Err   error
}

type Fetcher interface {
	Fetch(*db.Track) error
	FetchAll([]*db.Track) (<-chan Result, error)
}

type provider interface {
	Search(string, string) (string, error)
}

func fetchTrackLyrics(t *db.Track, l provider) error {
	artist := t.Artist
	if strings.Index(t.Artist, ", ") > -1 {
		artist = strings.Split(artist, ", ")[0]
	}
	lyric, err := l.Search(artist, t.Name)
	if err != nil {
		return err
	}

	t.Lyrics = lyric
	t.Loaded = true

	return nil
}

type AsyncFetcher struct {
	concurrency   int
	ready         chan struct{}
	fetchingQueue chan *db.Track
	lyricsFetcher provider
}

func (s AsyncFetcher) Fetch(t *db.Track) error {
	err := fetchTrackLyrics(t, s.lyricsFetcher)
	if err != nil {
		return err
	}
	return nil
}

func (s AsyncFetcher) FetchAll(tracks []*db.Track) (<-chan Result, error) {
	results := make(chan Result)
	var wg sync.WaitGroup

	queue := s.initWorkers(results, &wg)
	go s.run(tracks, queue, &wg)

	return results, nil
}

func New(geniusAPIToken string, concurrencyLevel int) AsyncFetcher {
	provider := lyrics.New(
		lyrics.WithGeniusLyrics(geniusAPIToken),
		lyrics.WithSongLyrics(),
	)
	return AsyncFetcher{
		ready:         make(chan struct{}, 1),
		concurrency:   concurrencyLevel,
		lyricsFetcher: &provider,
	}
}

func (s *AsyncFetcher) initWorkers(results chan<- Result, wg *sync.WaitGroup) chan *db.Track {
	c := make(chan *db.Track, s.concurrency)
	var once sync.Once

	for i := 0; i < s.concurrency; i++ {
		go func() {
			for t := range c {
				err := fetchTrackLyrics(t, s.lyricsFetcher)
				results <- Result{Track: t, Err: err}
				wg.Done()
			}

			once.Do(func() {
				close(results)
			})
		}()
	}

	return c
}

func (s *AsyncFetcher) run(tracks []*db.Track, queue chan<- *db.Track, wg *sync.WaitGroup) {
	defer close(queue)

	for i := range tracks {
		wg.Add(1)
		queue <- tracks[i]
	}

	wg.Wait()
}

var _ Fetcher = AsyncFetcher{}
