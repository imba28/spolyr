package model

import (
	"github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Track struct {
	ID         primitive.ObjectID `bson:"_id"`
	SpotifyID  string             `bson:"spotify_id"`
	Artist     string             `bson:"artist"`
	AlbumName  string             `bson:"album_name"`
	ImageURL   string             `bson:"image_url"`
	PreviewURL string             `bson:"preview_url"`
	Name       string             `bson:"name"`
	Lyrics     string             `bson:"lyrics"`
	Loaded     bool               `bson:"loaded"`
}

func NewTrack(t spotify.SavedTrack) Track {
	artists := make([]string, len(t.Artists))
	for j := range t.Artists {
		artists[j] = t.Artists[j].Name
	}

	imageUrl := ""
	if len(t.Album.Images) > 0 {
		imageUrl = t.Album.Images[0].URL
	}

	return Track{
		SpotifyID:  t.ID.String(),
		Artist:     strings.Join(artists, ", "),
		AlbumName:  t.Album.Name,
		ImageURL:   imageUrl,
		PreviewURL: t.PreviewURL,
		Name:       t.Name,
	}
}
