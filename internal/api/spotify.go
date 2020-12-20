package api

import (
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/internal/model"
	"github.com/zmb3/spotify"
	"log"
)

func SyncTracks(client spotify.Client, access *db.Access) error {
	tracks, err := client.CurrentUsersTracks()
	if err != nil {
		return err
	}
	for {
		for i := range tracks.Tracks {
			track := model.NewTrack(tracks.Tracks[i])
			err := access.SaveTrack(&track)
			if err != nil {
				return err
			}
			log.Println("saved track", track)
		}

		err = client.NextPage(tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
