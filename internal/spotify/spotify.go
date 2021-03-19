package spotify

import (
	"github.com/imba28/spolyr/internal/model"
	"github.com/zmb3/spotify"
	"log"
)

type UserTrackProvider interface {
	Tracks() ([]*model.Track, error)
	Next() error
}

type TrackSaver interface {
	Save(track *model.Track) error
}

type SpotifyUserTrackProvider struct {
	c        spotify.Client
	lastPage *spotify.SavedTrackPage
}

func (p *SpotifyUserTrackProvider) Tracks() ([]*model.Track, error) {
	if p.lastPage == nil {
		trackPage, err := p.c.CurrentUsersTracks()
		if err != nil {
			return nil, err
		}
		p.lastPage = trackPage
	}

	var tracks []*model.Track
	for i := range p.lastPage.Tracks {
		track := model.NewTrack(p.lastPage.Tracks[i])
		tracks = append(tracks, &track)
	}
	return tracks, nil
}

func (p *SpotifyUserTrackProvider) Next() error {
	return p.c.NextPage(p.lastPage)
}

func NewSpotifyTrackProvider(client spotify.Client) *SpotifyUserTrackProvider {
	return &SpotifyUserTrackProvider{
		c: client,
	}
}

func SyncTracks(client UserTrackProvider, store TrackSaver) error {
	for {
		tracks, err := client.Tracks()
		if err != nil {
			return err
		}

		for i := range tracks {
			err := store.Save(tracks[i])
			log.Println(tracks[i].Name)
			if err != nil {
				return err
			}
		}

		err = client.Next()
		if err != nil {
			if err == spotify.ErrNoMorePages {
				break
			}
			return err
		}
	}

	return nil
}
