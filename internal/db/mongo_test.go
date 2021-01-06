package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/imba28/spolyr/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const testDatabaseName = "tests"

func setUp() *Repositories {
	repos, err := New("root", "example", testDatabaseName)
	if err != nil {
		panic(err)
	}
	return repos
}

func tearDown(repositories *Repositories) {
	ctx := context.Background()
	err := repositories.client.Database(testDatabaseName).Collection(TrackCollection).Drop(ctx)
	if err != nil {
		panic(err)
	}
}

func TestMongoTrackStore_Save__inserts_new_document_into_database(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	track := model.Track{
		SpotifyID:  "1",
		Artist:     "Frank Sinatra",
		AlbumName:  "an album",
		ImageURL:   "http://spolyr.app/album.png",
		PreviewURL: "http://spolyr.app/preview.mp3",
		Name:       "Let It Snow",
		Lyrics:     "Oh the weather outside is frightful.",
		Loaded:     true,
	}
	repos := setUp()
	defer tearDown(repos)

	err := repos.Tracks.Save(&track)
	assert.Nil(t, err)

	trackFromDatabase, err := repos.Tracks.FindTrack(track.SpotifyID)
	assert.Nil(t, err)

	testCases := []struct {
		is, want, field string
	}{
		{trackFromDatabase.SpotifyID, track.SpotifyID, "spotifyID"},
		{trackFromDatabase.Artist, track.Artist, "artist"},
		{trackFromDatabase.AlbumName, track.AlbumName, "albumName"},
		{trackFromDatabase.ImageURL, track.ImageURL, "imageURL"},
		{trackFromDatabase.PreviewURL, track.PreviewURL, "previewURL"},
		{trackFromDatabase.Name, track.Name, "name"},
		{trackFromDatabase.Lyrics, track.Lyrics, "lyrics"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("updates field %s", tc.field), func(t *testing.T) {
			assert.Equal(t, tc.is, tc.want, "should save field %q of track in database", tc.field)
		})
	}
}

func TestTrackRepository_Count__count_is_zero_if_database_is_empty(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	repos := setUp()
	defer tearDown(repos)

	count, err := repos.Tracks.Count()
	assert.Equal(t, count, int64(0))
	assert.Nil(t, err)
}

func TestTrackRepository_Count__count_is_equal_to_number_of_documents_in_database(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2"})

	count, err := repos.Tracks.Count()
	assert.Equal(t, count, int64(2))
	assert.Nil(t, err)
}

func TestTrackRepository_CountWithLyrics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Loaded: true, Lyrics: "foobar"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3"})

	count, err := repos.Tracks.CountWithLyrics()
	assert.Equal(t, count, int64(1), "it should only count tracks where lyrics are set")
	assert.Nil(t, err)
}

func TestTrackRepository_CountWithLyrics__is_zero_if_database_is_empty(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	count, err := repos.Tracks.CountWithLyrics()
	assert.Equal(t, count, int64(0), "it should return 0 if database is empty")
	assert.Nil(t, err)
}

func TestTrackRepository_CountWithoutLyrics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Loaded: true, Lyrics: "foobar"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3"})

	count, err := repos.Tracks.CountWithoutLyrics()
	assert.Equal(t, count, int64(2), "it should only count tracks where lyrics are missing")
	assert.Nil(t, err)
}

func TestTrackRepository_CountWithoutLyrics__is_zero_if_database_is_empty(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	count, err := repos.Tracks.CountWithoutLyrics()
	assert.Equal(t, count, int64(0), "it should return 0 if database is empty")
	assert.Nil(t, err)
}

func TestTrackRepository_TracksWithoutLyrics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Loaded: true, Lyrics: "foobar"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2"})

	tracks, err := repos.Tracks.TracksWithoutLyrics()
	assert.Nil(t, err)
	assert.Len(t, tracks, 1, "it should return the tracks where lyrics are missing")
	assert.Equal(t, tracks[0].SpotifyID, "2", "it should return the track spotifyID = '2' that is missing its lyrics")
}

func TestTrackRepository_FindTrack(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Loaded: true, Lyrics: "foobar"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2"})

	track, err := repos.Tracks.FindTrack("1")
	assert.Nil(t, err)
	assert.Equal(t, track.SpotifyID, "1", "it should load the track spotifyID = '1' from database")
}

func TestTrackRepository_FindTrack__track_not_found__empty_database(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	track, err := repos.Tracks.FindTrack("1")
	assert.Error(t, err)
	assert.True(t, track.ID.IsZero(), "if the track does not exist in the database the primary key should be zero")
}

func TestTrackRepository_FindTrack__track_not_found(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "2"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3"})

	track, err := repos.Tracks.FindTrack("1")
	assert.Error(t, err)
	assert.True(t, track.ID.IsZero(), "if the track does not exist in the database the primary key should be zero")
}

func TestMongoTrackStore_FindTrack__database_error(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("FindOne", mock.AnythingOfType("primitive.D"), mock.Anything).
		Return(&model.Track{}, errors.New("db error"))
	repo := NewTrackRepository(mocksStore)
	_, err := repo.FindTrack("foobar")
	assert.Error(t, err)
}

func TestTrackRepository_Search__by_artist_name__partly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Artist: "Frank Sinatra"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Artist: "Dean Martin"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3"})

	tracks, err := repos.Tracks.Search("Frank")

	assert.Nil(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, tracks[0].Artist, "Frank Sinatra")
}

func TestTrackRepository_Search__by_artist_name__full(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Artist: "Frank Sinatra"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Artist: "Dean Martin"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3"})

	tracks, err := repos.Tracks.Search("Frank Sinatra")

	assert.Nil(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, tracks[0].Artist, "Frank Sinatra")
}

func TestTrackRepository_Search__by_album_name__partly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Artist: "Eminem", AlbumName: "The Eminem Show"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Artist: "Eminem", AlbumName: "Encore"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3", Artist: "Eminem", AlbumName: "The Slim Shady LP"})
	repos.Tracks.Save(&model.Track{SpotifyID: "4", Artist: "The Bloodhound Gang", AlbumName: "Show us your hits"})

	tracks, err := repos.Tracks.Search("Show")

	assert.Nil(t, err)
	assert.Len(t, tracks, 2)
	assert.True(t, tracks[0].AlbumName == "The Eminem Show" || tracks[1].AlbumName == "The Eminem Show", "should find a track from Eminem")
	assert.True(t, tracks[0].AlbumName == "Show us your hits" || tracks[1].AlbumName == "Show us your hits", "should find a track from The Bloodhound gang")
}

func TestTrackRepository_Search__by_album_name(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Artist: "Eminem", AlbumName: "The Eminem Show"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Artist: "Eminem", AlbumName: "Encore"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3", Artist: "Eminem", AlbumName: "The Slim Shady LP"})
	repos.Tracks.Save(&model.Track{SpotifyID: "4", Artist: "The Bloodhound Gang", AlbumName: "Show us your hits"})

	tracks, err := repos.Tracks.Search("Encore")

	assert.Nil(t, err)
	assert.Len(t, tracks, 1)
	assert.True(t, tracks[0].AlbumName == "Encore", "should find a track from Eminem's album 'Encore'")
}

func TestMongoTrackStore_Search__database_error(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("Find", mock.Anything, mock.Anything).
		Return([]*model.Track{}, errors.New("db error"))
	repo := NewTrackRepository(mocksStore)

	_, err := repo.Search("car")

	assert.Error(t, err)
}

func TestTrackRepository_Search__by_lyrics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Name: "A", Lyrics: "house mouse money car", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Name: "B", Lyrics: "house sky school", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "3", Name: "C", Lyrics: "fish company tank", Loaded: true})

	tracks, err := repos.Tracks.Search("car")

	assert.Nil(t, err)
	assert.Len(t, tracks, 1)
	assert.True(t, tracks[0].Name == "A", "should find a track whose lyrics contain the term 'car'")
}

func TestTrackRepository_Search__by_lyrics__multiple_results(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Name: "A", Lyrics: "house mouse money car", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Name: "B", Lyrics: "house sky school", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "3", Name: "C", Lyrics: "fish company tank", Loaded: true})

	tracks, err := repos.Tracks.Search("house")

	assert.Nil(t, err)
	assert.Len(t, tracks, 2)
}

func TestTrackRepository_Search__by_lyrics__multiple_query_term(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Name: "A", Lyrics: "house mouse money car", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Name: "B", Lyrics: "house sky school", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "3", Name: "C", Lyrics: "fish company tank", Loaded: true})

	tracks, err := repos.Tracks.Search("house money")

	assert.Nil(t, err)
	assert.Len(t, tracks, 2)
	assert.True(t, tracks[0].Name == "A" || tracks[1].Name == "A", "should find the track A whose lyrics contain the term 'house' or 'money'")
	assert.True(t, tracks[0].Name == "B" || tracks[1].Name == "B", "should find the track B whose lyrics contain the term 'house' or 'money'")
}

func TestTrackRepository_Search__by_lyrics__multiple_query_term__inclusive_search(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Name: "A", Lyrics: "house mouse money car", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Name: "B", Lyrics: "house sky school", Loaded: true})
	repos.Tracks.Save(&model.Track{SpotifyID: "3", Name: "C", Lyrics: "fish company tank", Loaded: true})

	tracks, err := repos.Tracks.Search("house \"money\"")

	assert.Nil(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, tracks[0].Name, "A", "should find a track whose lyrics contain the term 'house' as well as 'money'")
}

func TestTrackRepository_Search__by_name(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1", Name: "Without Me"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2", Name: "Stan"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3", Name: "'Till I Collapse'"})

	tracks, err := repos.Tracks.Search("collapse")

	assert.Nil(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, tracks[0].Name, "'Till I Collapse'", "should find the track A whose title contain the term 'collapse'")
}

func TestTrackRepository_LatestTracks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repos := setUp()
	defer tearDown(repos)

	repos.Tracks.Save(&model.Track{SpotifyID: "1"})
	repos.Tracks.Save(&model.Track{SpotifyID: "2"})
	repos.Tracks.Save(&model.Track{SpotifyID: "3"})
	repos.Tracks.Save(&model.Track{SpotifyID: "4"})

	tracks, err := repos.Tracks.LatestTracks(1)

	assert.Nil(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, tracks[0].SpotifyID, "4", "should return the latest track with regards to the insertion date")
}

func TestTrackRepository_LatestTracks__database_error(t *testing.T) {
	mocksStore := new(TrackStoreMock)
	mocksStore.On("Find", mock.Anything, mock.Anything).
		Return([]*model.Track{}, errors.New("db error"))
	repo := NewTrackRepository(mocksStore)

	_, err := repo.LatestTracks(10)

	assert.Error(t, err)
}
