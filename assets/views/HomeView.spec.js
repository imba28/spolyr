import {render} from '@testing-library/vue';
import HomeView from '@/views/HomeView';
import {TracksApi} from '@/openapi';
import {createLocalVue} from '@vue/test-utils';
import {PiniaVuePlugin} from 'pinia';
import {createTestingPinia} from '@pinia/testing';
import {useAuthStore} from '@/stores/auth';

const renderHomeView = (isAuthenticated=false) => {
  const localVue = createLocalVue();
  localVue.use(PiniaVuePlugin);

  const r = render(HomeView, {
    localVue,
    pinia: createTestingPinia(),
    stubs: ['router-link'],
    mocks: {
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

afterEach(() => {
  jest.clearAllMocks();
});

describe('HomeView', () => {
  it('renders a list of the latest songs', async () => {
    const tracks = [
      {spotifyId: '1', title: 'Track 1'},
      {spotifyId: '2', title: 'Track 2'},
    ];
    jest.spyOn(TracksApi.prototype, 'tracksStatsGet').mockImplementation(() => Promise.resolve(null));
    const spy = jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: tracks, meta: {total: 2}}));

    const {findByText} = renderHomeView();

    await findByText('Latest songs');
    await (findByText(tracks[0].title));
    await (findByText(tracks[1].title));
    expect(spy).toHaveBeenCalled();
  });

  it('renders number of indexed tracks', async () => {
    const numberOfTracks = 55;
    const numberOfTracksWithLyrics = 10;
    const tracks = [{title: '1', spotifyId: '1'}, {title: '2', spotifyId: '2'}];
    jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: tracks, meta: {total: 2}}));
    const spy = jest.spyOn(TracksApi.prototype, 'tracksStatsGet')
        .mockImplementation(() => Promise.resolve({
          numberOfTracks,
          numberOfTracksWithLyrics,
        }));

    const {findByText} = renderHomeView();

    await findByText(`Indexed ${numberOfTracks} tracks including ${numberOfTracksWithLyrics} songs with known lyrics.`);
    expect(spy).toHaveBeenCalled();
  });

  it('shows a message if library is empty', async () => {
    jest.spyOn(TracksApi.prototype, 'tracksStatsGet').mockImplementation(() => Promise.resolve(null));
    jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: [], meta: {total: 0}}));

    const {findByText} = renderHomeView();

    await findByText(`Oh no, it seems your library is currently empty.`);
    await findByText(`Maybe you want to sign in and change that?`);
  });
});
