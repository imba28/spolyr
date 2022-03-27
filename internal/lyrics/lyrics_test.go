package lyrics

import (
	"errors"
	"fmt"
	"github.com/imba28/spolyr/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type providerMock struct {
	mock.Mock
}

func (p providerMock) Search(a string, b string) (string, error) {
	args := p.Called(a, b)
	return args.Get(0).(string), args.Error(1)
}

var _ provider = providerMock{}

func TestAsyncFetcher_Fetch(t *testing.T) {
	t.Run("fetches lyrics", func(t *testing.T) {
		artist, song := "artist", "a song"
		expectedLyrics := "la la la la"
		track := db.Track{
			Artist: artist,
			Name:   song,
		}
		providerMock := providerMock{}
		providerMock.On("Search", artist, song).Return(expectedLyrics, nil)

		fetcher := AsyncFetcher{lyricsFetcher: providerMock}

		err := fetcher.Fetch(&track)

		assert.Nil(t, err)
		assert.Equal(t, track.Lyrics, expectedLyrics)
		assert.True(t, track.Loaded)
		providerMock.AssertExpectations(t)
	})

	t.Run("picks the first artist if Track has multiple artists", func(t *testing.T) {
		artist, song := "Eminem, Nate Dog", "'Till I Collapse"
		expectedLyrics := "la la la la"
		track := db.Track{
			Artist: artist,
			Name:   song,
		}
		providerMock := providerMock{}
		providerMock.On("Search", "Eminem", song).Return(expectedLyrics, nil)

		fetcher := AsyncFetcher{lyricsFetcher: providerMock}

		err := fetcher.Fetch(&track)

		assert.Nil(t, err)
		providerMock.AssertExpectations(t)
	})

	t.Run("returns error if provider return error", func(t *testing.T) {
		track := db.Track{}
		expectedErr := errors.New("something went wrong")

		providerMock := providerMock{}
		providerMock.On("Search", track.Artist, track.Name).Return("", expectedErr)
		fetcher := AsyncFetcher{lyricsFetcher: providerMock}

		err := fetcher.Fetch(&track)

		assert.Errorf(t, err, expectedErr.Error())
		providerMock.AssertExpectations(t)
	})
}

func TestAsyncFetcher_FetchAll(t *testing.T) {
	t.Run("it works with different concurrency levels", withTimeout(func(t *testing.T) {
		tests := []int{1, 2, 5, 10}

		tracks := []*db.Track{
			{Artist: "a", Name: "a"},
			{Artist: "b", Name: "b"},
			{Artist: "c", Name: "c"},
			{Artist: "d", Name: "d"},
			{Artist: "e", Name: "e"},
			{Artist: "f", Name: "f"},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("with %d workers", tt), func(t *testing.T) {
				expectedLyrics := "la la la la"
				providerMock := providerMock{}
				providerMock.
					On("Search", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Times(len(tracks)).
					Return(expectedLyrics, nil)
				fetcher := AsyncFetcher{lyricsFetcher: providerMock, concurrency: tt}

				c, err := fetcher.FetchAll(tracks)

				for r := range c {
					assert.Nil(t, r.Err)
					assert.Equal(t, expectedLyrics, r.Track.Lyrics)
					assert.True(t, r.Track.Loaded)
				}

				assert.Nil(t, err)
				providerMock.AssertExpectations(t)
			})
		}
	}, 2*time.Second))

	t.Run("it writes errors to the result struct", withTimeout(func(t *testing.T) {
		tracks := []*db.Track{
			{Artist: "a", Name: "a"},
			{Artist: "b", Name: "b"},
		}
		expectedLyrics := ""
		expectedError := errors.New("something went wrong")

		providerMock := providerMock{}
		providerMock.
			On("Search", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Times(len(tracks)).
			Return(expectedLyrics, expectedError)
		fetcher := AsyncFetcher{lyricsFetcher: providerMock, concurrency: 2}

		c, err := fetcher.FetchAll(tracks)
		assert.Nil(t, err)

		for r := range c {
			assert.EqualError(t, r.Err, expectedError.Error())
			assert.Equal(t, r.Track.Lyrics, expectedLyrics)
		}
		for i := range tracks {
			t.Run(tracks[i].Artist, func(t *testing.T) {
				assert.False(t, tracks[i].Loaded)
			})
		}
		providerMock.AssertExpectations(t)
	}, 2*time.Second))
}

type testerFunc func(t *testing.T)

func withTimeout(f testerFunc, d time.Duration) testerFunc {
	return func(t *testing.T) {
		timeout := time.After(d)
		r := make(chan struct{})
		go func() {
			f(t)
			r <- struct{}{}
		}()

		select {
		case <-timeout:
			t.Fatalf("operation should complete in less than %v!", d)
		case <-r:
			return
		}
	}
}
