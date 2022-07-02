import TrackDetailView from '@/views/TrackDetailView';
import {fireEvent, render} from '@testing-library/vue';
import {createTestingPinia} from '@pinia/testing';
import {ImportApi, TracksApi} from '@/openapi';
import {createLocalVue} from '@vue/test-utils';
import {PiniaVuePlugin} from 'pinia';
import {useAuthStore} from '@/stores/auth';

beforeEach(() => {
  jest.clearAllMocks();
});

const renderTrackDetailView = (params, isAuthenticated = false) => {
  const localVue = createLocalVue();
  localVue.use(PiniaVuePlugin);

  const r = render(TrackDetailView, {
    localVue,
    pinia: createTestingPinia(),
    stubs: ['router-link'],
    mocks: {
      $route: {
        params,
      },
      $toast: {
        success: jest.fn(),
        error: jest.fn(),
        warning: jest.fn(),
      },
    },
  });

  const authStore = useAuthStore();
  if (isAuthenticated) {
    authStore.avatarUrl = 'https://foobar.com/1.jpg';
    authStore.displayName = 'Test user';
  }

  return r;
};

describe('TrackDetailView', () => {
  it('calls the track info api after mounting the component', async () => {
    const spotifyId = 'a-spotify-id';

    const spy = jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve({}));

    renderTrackDetailView({id: spotifyId});

    expect(spy).toHaveBeenCalledWith(spotifyId);
  });

  it('renders track infos', async () => {
    const trackInfo = {
      spotifyId: '1',
      title: 'a title',
      artists: [
        'Artist A',
        'Artist B',
      ],
      album: 'name of album',
      lyrics: 'His palms are sweaty, knees weak, arms are heavy',
    };
    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));

    const {findByText} = renderTrackDetailView({id: '1'});

    await findByText(trackInfo.title);
    await findByText(trackInfo.artists[0]);
    await findByText(trackInfo.artists[1]);
    await findByText(trackInfo.album);
    await findByText(trackInfo.lyrics);
  });

  it('contains button to edit the lyrics if user is signed in', async () => {
    const trackInfo = {
      lyrics: 'His palms are sweaty, knees weak, arms are heavy',
      hasLyrics: true,
    };

    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));

    const {findByLabelText} = renderTrackDetailView({id: '1'}, true);

    await findByLabelText('edit lyrics');
  });

  it('does not render the edit button if user is not authenticated ', async () => {
    const trackInfo = {title: 'a title'};

    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));

    const {findByText, queryByLabelText} = renderTrackDetailView({id: '1'});
    // wait for dom to update
    await findByText(trackInfo.title);

    expect(queryByLabelText('edit lyrics')).not.toBeInTheDocument();
  });

  it('after clicking the edit button a textarea is visible', async () => {
    const trackInfo = {
      spotifyId: 'a-spotify-id',
      lyrics: 'His palms are sweaty, knees weak, arms are heavy',
      hasLyrics: true,
    };
    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));

    const {findByLabelText, findByRole} = renderTrackDetailView({id: trackInfo.spotifyId}, true);

    const editButton = await findByLabelText('edit lyrics');
    await fireEvent.click(editButton);
    const textarea = await findByRole('textbox');

    expect(textarea).toBeVisible();
  });

  it('shows an audio player if the track contains an audio snippet', async () => {
    const trackInfo = {
      previewURL: 'https://spotify.com/foobar.mp3',
    };
    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));

    const {findByTestId} = renderTrackDetailView({id: trackInfo.spotifyId});

    await findByTestId('audio-player');
  });

  it('does not show an audio player if the track does not contain an audio snippet', async () => {
    const trackInfo = {title: 'foo'};
    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));

    const {queryByTestId, findByText} = renderTrackDetailView({id: trackInfo.spotifyId});
    // wait for dom to update
    await findByText(trackInfo.title);

    expect(await queryByTestId('audio-player')).not.toBeInTheDocument();
  });

  it('clicking the save button calls the tracksIdPatch endpoint and hides the textarea', async () => {
    const updatedLyrics = 'new lyrics';
    const trackInfo = {
      spotifyId: 'a-spotify-id',
      lyrics: 'His palms are sweaty, knees weak, arms are heavy',
      hasLyrics: true,
    };
    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));
    const spy = jest.spyOn(TracksApi.prototype, 'tracksIdPatch')
        .mockImplementation(() => Promise.resolve({lyrics: updatedLyrics}));

    const {findByLabelText, findByRole} = renderTrackDetailView({id: trackInfo.spotifyId}, true);
    await fireEvent.click(await findByLabelText('edit lyrics'));
    const textarea = await findByRole('textbox');
    await fireEvent.update(textarea, updatedLyrics);

    await fireEvent.click(await findByLabelText('save lyrics'));

    expect(spy).toHaveBeenCalledWith(trackInfo.spotifyId, expect.objectContaining({lyrics: updatedLyrics}));
    expect(textarea).not.toBeVisible();
  });

  it('contains button to import the lyrics of a track if user is signed in', async () => {
    const trackInfo = {
      hasLyrics: false,
    };

    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));

    const {findByLabelText} = renderTrackDetailView({id: '1'}, true);

    await findByLabelText('import lyrics');
  });

  it('clicking the import button calls the import endpoint and updates the lyrics', async () => {
    const trackInfo = {
      spotifyId: 'an-id',
      hasLyrics: false,
    };
    const updatedTrackInfo = {
      hasLyrics: true,
      lyrics: 'New lyrics',
    };
    jest.spyOn(TracksApi.prototype, 'tracksIdGet')
        .mockImplementation(() => Promise.resolve(trackInfo));
    const spy = jest.spyOn(ImportApi.prototype, 'importLyricsTrackIdPost')
        .mockImplementation(() => Promise.resolve(updatedTrackInfo));

    const {findByLabelText, findByText} = renderTrackDetailView({id: trackInfo.spotifyId}, true);
    const importButton = await findByLabelText('import lyrics');
    await fireEvent.click(importButton);

    expect(spy).toHaveBeenCalledWith(trackInfo.spotifyId);
    await findByText(updatedTrackInfo.lyrics);
    expect(importButton).not.toBeInTheDocument();
  });
});
