import PlaylistCard from '@/components/PlaylistCard';
import {render, fireEvent} from '@testing-library/vue';

const renderPlaylistCard = (props) => {
  return render(PlaylistCard, {
    propsData: props,
  });
};

describe('PlaylistCard', () => {
  it('renders the name of a playlist', async () => {
    const playlistData = {
      spotifyId: '1',
      name: 'A simple playlist',
    };
    const {findByText} = renderPlaylistCard(playlistData);
    await findByText(playlistData.name);
  });

  it('renders a button', async () => {
    const playlistData = {
      spotifyId: '1',
      name: 'A simple playlist',
    };
    const {findByRole} = renderPlaylistCard(playlistData);

    await findByRole('button');
  });

  it('emits a `button-click` event when clicking the button', async () => {
    const playlistData = {
      spotifyId: '1',
      name: 'A simple playlist',
    };
    const {getByRole, emitted} = renderPlaylistCard(playlistData);

    await fireEvent.click(getByRole('button'));

    expect(emitted()['button-click']).toHaveLength(1);
    expect(emitted()['button-click'][0][0]).toEqual(playlistData.spotifyId);
  });
});
