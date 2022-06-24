package spotify

import (
	"context"
	"github.com/imba28/spolyr/internal/db"
	"github.com/zmb3/spotify/v2"
)

type userTrackProvider interface {
	Tracks(ctx context.Context) ([]*db.Track, error)
	Next(ctx context.Context) error
}

type trackSaver interface {
	Save(track *db.Track) error
}

type UserTrackProvider struct {
	c        *spotify.Client
	lastPage *spotify.SavedTrackPage
}

func (p *UserTrackProvider) Tracks(ctx context.Context) ([]*db.Track, error) {
	if p.lastPage == nil {
		trackPage, err := p.c.CurrentUsersTracks(ctx)
		if err != nil {
			return nil, err
		}
		p.lastPage = trackPage
	}

	var tracks []*db.Track
	for i := range p.lastPage.Tracks {
		track := db.NewTrack(p.lastPage.Tracks[i].FullTrack)
		tracks = append(tracks, &track)
	}
	return tracks, nil
}

func (p *UserTrackProvider) Next(ctx context.Context) error {
	return p.c.NextPage(ctx, p.lastPage)
}

func NewSpotifyTrackProvider(client *spotify.Client) *UserTrackProvider {
	return &UserTrackProvider{
		c: client,
	}
}

func SyncTracks(ctx context.Context, client userTrackProvider, store trackSaver) error {
	for {
		tracks, err := client.Tracks(ctx)
		if err != nil {
			return err
		}

		for i := range tracks {
			err := store.Save(tracks[i])
			if err != nil {
				return err
			}
		}

		err = client.Next(ctx)
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
	c     *spotify.Client
	saver trackSaver
}

func (p PlaylistProvider) Download(ctx context.Context, ID string) error {
	playlist, err := p.c.GetPlaylistTracks(ctx, spotify.ID(ID))
	if err != nil {
		return err
	}

	for {
		for i := range playlist.Tracks {
			track := db.NewTrack(playlist.Tracks[i].Track)
			err = p.saver.Save(&track)
			if err != nil {
				return err
			}
		}

		err = p.c.NextPage(ctx, playlist)
		if err != nil {
			if err == spotify.ErrNoMorePages {
				return nil
			}
			return err
		}
	}
}

func NewPlaylistProvider(c *spotify.Client, saver trackSaver) PlaylistProvider {
	return PlaylistProvider{
		c:     c,
		saver: saver,
	}
}
