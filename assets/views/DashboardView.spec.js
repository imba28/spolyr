import DashboardView from '@/views/DashboardView';
import {render, fireEvent, waitForElementToBeRemoved, getByLabelText} from '@testing-library/vue';

import {PlaylistsApi, ImportApi} from '@/openapi';

beforeEach(() => {
  jest.clearAllMocks();
});

const renderDashboardView = (successToast = jest.fn(), errorToast = jest.fn()) => {
  return render(DashboardView, {
    stubs: ['router-link'],
    mocks: {
      $toast: {
        success: successToast,
        error: errorToast,
      },
    },
  });
};

describe('DashboardView', function() {
  it('initially calls the playlist api endpoint', async () => {
    const spy = jest.spyOn(PlaylistsApi.prototype, 'playlistsGet')
        .mockImplementation(() => Promise.resolve({data: [], meta: {total: 0}}));

    renderDashboardView();

    expect(spy).toHaveBeenCalledWith(expect.objectContaining({
      page: 1,
    }));
  });

  it('renders a list of playlist results', async () => {
    const playlists = [
      {spotifyId: '1', name: 'Playlist A'},
      {spotifyId: '2', name: 'Playlist B'},
      {spotifyId: '3', name: 'Playlist C'},
    ];

    jest.spyOn(PlaylistsApi.prototype, 'playlistsGet')
        .mockImplementation(() => Promise.resolve({data: playlists, meta: {total: playlists.length}}));

    const {findByText} = renderDashboardView();

    await findByText(playlists[0].name);
    await findByText(playlists[1].name);
    await findByText(playlists[2].name);
  });

  it('clicking the import button of a playlists calls the import endpoint', async () => {
    const playlist = {
      spotifyId: '1',
      name: 'Playlist A',
      trackCount: 5,
    };

    jest.spyOn(PlaylistsApi.prototype, 'playlistsGet')
        .mockImplementation(() => Promise.resolve({data: [playlist], meta: {total: 1}}));
    const spy = jest.spyOn(ImportApi.prototype, 'importPlaylistIdPost')
        .mockImplementation(() => Promise.resolve({}));

    const {findByText, getByText} = renderDashboardView();
    await findByText(playlist.name);

    await fireEvent.click(getByText(`Import ${playlist.trackCount} tracks`));

    expect(spy).toHaveBeenCalledWith(playlist.spotifyId);
  });

  it('shows a pagination if more than 8 results exist.', async () => {
    jest.spyOn(PlaylistsApi.prototype, 'playlistsGet').mockImplementationOnce(
        () => Promise.resolve({data: [{spotifyId: '1', name: 'Playlist 1'}], meta: {total: 99}}),
    );

    // wait for page load
    const {findByLabelText} = renderDashboardView();
    await findByLabelText('Pagination');
  });

  it('calls the playlist api when a pagination item is clicked', async () => {
    const spy = jest.spyOn(PlaylistsApi.prototype, 'playlistsGet').mockImplementationOnce(
        () => Promise.resolve({data: [{spotifyId: '1', name: 'Playlist 1'}], meta: {total: 99}}),
    );

    // wait for page load
    const {findByLabelText, getByLabelText} = renderDashboardView();
    await findByLabelText('Pagination');

    await fireEvent.click(getByLabelText('Go to page 2'));

    expect(spy).toHaveBeenCalledWith(expect.objectContaining({
      page: 2,
    }));
  });
});
