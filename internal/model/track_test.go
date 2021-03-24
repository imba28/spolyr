package model

import (
	"fmt"
	"github.com/zmb3/spotify"
	"testing"
)

func TestTrack_sets_SpotifyID(t *testing.T) {
	spotifyTrack := spotify.FullTrack{
		SimpleTrack: spotify.SimpleTrack{
			ID: "some_id",
		},
	}
	track := NewTrack(spotifyTrack)

	if track.SpotifyID != spotifyTrack.ID.String() {
		t.Errorf("expect a track to copy its ID from the spotify track. expected: %q, got %q", track.ID, spotifyTrack.ID)
	}
}

func TestTrack_aggregates_Artist(t *testing.T) {
	var tests = []struct {
		artists []spotify.SimpleArtist
		want    string
	}{
		{
			artists: []spotify.SimpleArtist{
				{Name: "an artist"},
			},
			want: "an artist",
		},
		{
			artists: []spotify.SimpleArtist{
				{Name: "artist A"},
				{Name: "artist B"},
			},
			want: "artist A, artist B",
		},
		{
			artists: []spotify.SimpleArtist{
				{Name: "artist A"},
				{Name: "artist B"},
				{Name: "artist C"},
			},
			want: "artist A, artist B, artist C",
		},
		{
			artists: nil,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v,%s", len(tt.artists), tt.want), func(t *testing.T) {
			spotifyTrack := spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{
					Artists: tt.artists,
				},
			}
			track := NewTrack(spotifyTrack)

			if track.Artist != tt.want {
				t.Errorf("expect a track aggregate all artists to a csv string. expected: %q, got %q", tt.want, track.Artist)
			}
		})
	}
}

func TestTrack_sets_AlbumName(t *testing.T) {
	var tests = []struct {
		album spotify.SimpleAlbum
		want  string
	}{
		{
			album: spotify.SimpleAlbum{
				Name: "album name",
			},
			want: "album name",
		},
		{
			album: spotify.SimpleAlbum{
				Name: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v,%s", tt.album.Name, tt.want), func(t *testing.T) {
			spotifyTrack := spotify.FullTrack{
				Album: tt.album,
			}
			track := NewTrack(spotifyTrack)

			if track.AlbumName != tt.want {
				t.Errorf("expect a track to copy its album name. expected: %q, got %q", tt.want, track.AlbumName)
			}
		})
	}
}

func TestTrack_sets_ImageURL(t *testing.T) {
	var tests = []struct {
		images []spotify.Image
		want   string
	}{
		{
			images: []spotify.Image{},
			want:   "",
		},
		{
			images: nil,
			want:   "",
		},
		{
			images: []spotify.Image{
				{URL: "https://test.com/album_cover.jpg"},
			},
			want: "https://test.com/album_cover.jpg",
		},
		{
			images: []spotify.Image{
				{URL: "https://test.com/album_cover.jpg"},
				{URL: "https://test.com/album_cover_small.jpg"},
			},
			want: "https://test.com/album_cover.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v,%s", tt.images, tt.want), func(t *testing.T) {
			spotifyTrack := spotify.FullTrack{
				Album: spotify.SimpleAlbum{
					Images: tt.images,
				},
			}
			track := NewTrack(spotifyTrack)

			if track.ImageURL != tt.want {
				t.Errorf("expect a track to use the album's first image. expected: %q, got %q", tt.want, track.ImageURL)
			}
		})
	}
}

func TestTrack_sets_PreviewURL(t *testing.T) {
	var tests = []struct {
		previewURL, want string
	}{
		{"https://foo.com/audio.mp3", "https://foo.com/audio.mp3"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v,%s", tt.previewURL, tt.want), func(t *testing.T) {
			spotifyTrack := spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{PreviewURL: tt.previewURL},
			}
			track := NewTrack(spotifyTrack)

			if track.PreviewURL != tt.want {
				t.Errorf("expect a track to use the album's first image. expected: %q, got %q", tt.want, track.ImageURL)
			}
		})
	}
}

func TestTrack_sets_Name(t *testing.T) {
	var tests = []struct {
		name, want string
	}{
		{"album name", "album name"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v,%s", tt.name, tt.want), func(t *testing.T) {
			spotifyTrack := spotify.FullTrack{
				SimpleTrack: spotify.SimpleTrack{Name: tt.name},
			}
			track := NewTrack(spotifyTrack)

			if track.Name != tt.want {
				t.Errorf("expect a track to use the album's first image. expected: %q, got %q", tt.want, track.Name)
			}
		})
	}
}

func TestTrack_does_not_set_ID(t *testing.T) {
	spotifyTrack := spotify.FullTrack{}
	track := NewTrack(spotifyTrack)

	if !track.ID.IsZero() {
		t.Errorf("expect track to not set the ID field, got: %v", track.ID)
	}
}
