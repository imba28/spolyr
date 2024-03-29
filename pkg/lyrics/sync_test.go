package lyrics

import (
	"errors"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type trackStoreMock struct {
	mock.Mock
}

func (t *trackStoreMock) Save(track *db.Track) error {
	args := t.Called(track)
	return args.Error(0)
}
func (t *trackStoreMock) TracksWithoutLyricsError() ([]*db.Track, error) {
	args := t.Called()
	return args.Get(0).([]*db.Track), args.Error(1)
}

type lyricsFetcherMock struct {
	mock.Mock
}

func (l *lyricsFetcherMock) Fetch(ts *db.Track) error {
	panic("not implemented")
}
func (l *lyricsFetcherMock) FetchAll(ts []*db.Track) (<-chan Result, error) {
	args := l.Called(ts)
	return args.Get(0).(chan Result), args.Error(1)
}
func (l *lyricsFetcherMock) writeFakeResults(tracks []*db.Track, results chan<- Result) {
	for i := range tracks {
		time.Sleep(100 * time.Millisecond)

		tracks[i].Lyrics = "la la la"
		results <- Result{
			Track: tracks[i],
			Err:   nil,
		}
	}
	close(results)
}

func TestSyncer_Sync(t *testing.T) {
	t.Run("downloads and saves lyrics of tracks", func(t *testing.T) {
		withTimeout(func(t *testing.T) {
			tracks := []*db.Track{
				{
					Name: "track A",
				},
				{
					Name: "track B",
				},
				{
					Name: "track C",
				},
			}

			dbMock := trackStoreMock{}
			dbMock.On("Save", mock.AnythingOfType("*db.Track")).Times(len(tracks)).Return(nil)
			dbMock.On("TracksWithoutLyricsError").Times(1).Return(tracks, nil)

			results := make(chan Result)

			fetcherMock := lyricsFetcherMock{}
			fetcherMock.On("FetchAll", mock.AnythingOfType("[]*db.Track")).Times(1).Return(results, nil)

			syncer := NewSyncer(&fetcherMock, &dbMock)
			finished, err := syncer.Sync()

			// simulate fetching of lyrics
			go fetcherMock.writeFakeResults(tracks, results)

			assert.Nil(t, err)
			assert.NotNil(t, finished)

			<-finished

			fetcherMock.AssertExpectations(t)
			dbMock.AssertExpectations(t)
		}, time.Second)(t)
	})

	t.Run("prevents clients from starting multiple syncs", func(t *testing.T) {
		withTimeout(func(t *testing.T) {
			dbMock := trackStoreMock{}
			dbMock.On("TracksWithoutLyricsError").Return([]*db.Track{{}}, nil)

			results := make(chan Result)

			fetcherMock := lyricsFetcherMock{}
			fetcherMock.On("FetchAll", mock.AnythingOfType("[]*db.Track")).Times(1).Return(results, nil)

			syncer := NewSyncer(&fetcherMock, &dbMock)
			finished, err := syncer.Sync()

			assert.Nil(t, err)

			_, err = syncer.Sync()
			assert.ErrorIs(t, err, ErrBusy)

			close(results)
			<-finished

			fetcherMock.AssertExpectations(t)
			dbMock.AssertExpectations(t)
		}, time.Second)(t)
	})
}

func TestSyncer_Syncing(t *testing.T) {
	t.Run("returns correct syncing state", func(t *testing.T) {
		tracks := []*db.Track{
			{
				Name: "track A",
			},
			{
				Name: "track B",
			},
			{
				Name: "track C",
			},
		}

		dbMock := trackStoreMock{}
		dbMock.On("Save", mock.AnythingOfType("*db.Track")).Times(len(tracks)).Return(nil)
		dbMock.On("TracksWithoutLyricsError").Times(1).Return(tracks, nil)

		results := make(chan Result)

		fetcherMock := lyricsFetcherMock{}
		fetcherMock.On("FetchAll", mock.AnythingOfType("[]*db.Track")).Times(1).Return(results, nil)

		syncer := NewSyncer(&fetcherMock, &dbMock)
		finished, _ := syncer.Sync()

		assert.True(t, syncer.Syncing())

		close(results)
		<-finished
		assert.False(t, syncer.Syncing())
	})

	t.Run("returns error if tracks cannot be loaded from database", func(t *testing.T) {
		expectedError := errors.New("something went wrong")
		tracks := []*db.Track{}

		dbMock := trackStoreMock{}
		dbMock.On("TracksWithoutLyricsError").Times(1).Return(tracks, expectedError)

		fetcherMock := lyricsFetcherMock{}

		syncer := NewSyncer(&fetcherMock, &dbMock)
		finished, err := syncer.Sync()

		assert.Nil(t, finished)
		assert.ErrorIs(t, err, expectedError)
		assert.False(t, syncer.Syncing())
	})
}

func TestSyncer_TotalTracks(t *testing.T) {
	tracks := []*db.Track{
		{
			Name: "track A",
		},
		{
			Name: "track B",
		},
	}

	dbMock := trackStoreMock{}
	dbMock.On("TracksWithoutLyricsError").Return(tracks, nil)

	results := make(chan Result)
	defer close(results)

	fetcherMock := lyricsFetcherMock{}
	fetcherMock.On("FetchAll", mock.AnythingOfType("[]*db.Track")).Times(1).Return(results, nil)

	syncer := NewSyncer(&fetcherMock, &dbMock)
	_, _ = syncer.Sync()

	assert.Equal(t, syncer.TotalTracks(), len(tracks))
}

func TestSyncer_SyncedTracks(t *testing.T) {
	t.Run("increases counter with every fetched track", func(t *testing.T) {
		tracks := []*db.Track{
			{
				Name: "track A",
			},
			{
				Name: "track B",
			},
		}

		dbMock := trackStoreMock{}
		dbMock.On("TracksWithoutLyricsError").Return(tracks, nil)

		results := make(chan Result)
		defer close(results)

		fetcherMock := lyricsFetcherMock{}
		dbMock.On("Save", mock.AnythingOfType("*db.Track")).Times(len(tracks)).Return(nil)
		fetcherMock.On("FetchAll", mock.AnythingOfType("[]*db.Track")).Times(1).Return(results, nil)

		syncer := NewSyncer(&fetcherMock, &dbMock)
		_, _ = syncer.Sync()

		assert.Equal(t, syncer.SyncedTracks(), 0)
		results <- Result{
			Track: tracks[0],
		}
		assert.Equal(t, syncer.SyncedTracks(), 1)
	})

	t.Run("resets counter after sync process has finished", func(t *testing.T) {
		withTimeout(func(t *testing.T) {
			tracks := []*db.Track{
				{
					Name: "track A",
				},
				{
					Name: "track B",
				},
			}

			dbMock := trackStoreMock{}
			dbMock.On("TracksWithoutLyricsError").Return(tracks, nil)

			results := make(chan Result)

			fetcherMock := lyricsFetcherMock{}
			dbMock.On("Save", mock.AnythingOfType("*db.Track")).Times(len(tracks)).Return(nil)
			fetcherMock.On("FetchAll", mock.AnythingOfType("[]*db.Track")).Times(1).Return(results, nil)

			syncer := NewSyncer(&fetcherMock, &dbMock)
			finished, _ := syncer.Sync()

			go fetcherMock.writeFakeResults(tracks, results)

			<-finished
			assert.Equal(t, syncer.SyncedTracks(), -1)
		}, time.Second)
	})

	t.Run("marks tracks as failed if an error occurs while importing lyrics", func(t *testing.T) {
		tracks := []*db.Track{
			{
				Name: "track A",
			},
		}

		dbMock := trackStoreMock{}
		dbMock.On("TracksWithoutLyricsError").Return(tracks, nil)

		results := make(chan Result)
		defer close(results)

		fetcherMock := lyricsFetcherMock{}
		dbMock.On("Save", mock.AnythingOfType("*db.Track")).Times(len(tracks)).Return(nil)
		fetcherMock.On("FetchAll", mock.AnythingOfType("[]*db.Track")).Times(1).Return(results, nil)

		syncer := NewSyncer(&fetcherMock, &dbMock)
		_, _ = syncer.Sync()

		results <- Result{
			Track: tracks[0],
			Err:   errors.New("something went wrong during the lyrics import"),
		}

		assert.Equal(t, tracks[0].LyricsImportErrorCount, 1, "should increase error counter if import fails")
		dbMock.AssertExpectations(t)
	})
}
