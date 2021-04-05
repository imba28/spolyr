package spotify

import (
	"github.com/imba28/spolyr/internal/model"
	"github.com/zmb3/spotify"
	"log"
)

type userTrackProvider interface {
	Tracks() ([]*model.Track, error)
	Next() error
}

type trackSaver interface {
	Save(track *model.Track) error
}

type UserTrackProvider struct {
	c        spotify.Client
	lastPage *spotify.SavedTrackPage
}

func (p *UserTrackProvider) Tracks() ([]*model.Track, error) {
	if p.lastPage == nil {
		trackPage, err := p.c.CurrentUsersTracks()
		if err != nil {
			return nil, err
		}
		p.lastPage = trackPage
	}

	var tracks []*model.Track
	for i := range p.lastPage.Tracks {
		track := model.NewTrack(p.lastPage.Tracks[i].FullTrack)
		tracks = append(tracks, &track)
	}
	return tracks, nil
}

func (p *UserTrackProvider) Next() error {
	return p.c.NextPage(p.lastPage)
}

func NewSpotifyTrackProvider(client spotify.Client) *UserTrackProvider {
	return &UserTrackProvider{
		c: client,
	}
}

func SyncTracks(client userTrackProvider, store trackSaver) error {
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


type PlaylistProvider struct {
	c spotify.Client
	saver trackSaver
}

func (p PlaylistProvider) Download(ID string) error {
	playlist, err := p.c.GetPlaylistTracks(spotify.ID(ID))
	if err != nil {
		return err
	}

	for {
		for i := range playlist.Tracks {
			track := model.NewTrack(playlist.Tracks[i].Track)
			err = p.saver.Save(&track)
			if err != nil {
				return err
			}
		}

		err = p.c.NextPage(playlist)
		if err != nil {
			if err == spotify.ErrNoMorePages {
				return nil
			}
			return err
		}
	}
}

func NewPlaylistProvider(c spotify.Client, saver trackSaver) PlaylistProvider {
	return PlaylistProvider{
		c : c,
		saver: saver,
	}
}