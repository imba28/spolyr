import {waitFor, render, fireEvent} from '@testing-library/vue';

import SearchView from './SearchView';
import {TracksApi} from '@/openapi';
import Player from '@/track-page/player';

// Set all module functions to jest.fn
jest.mock('@/track-page/player');

beforeEach(() => {
  jest.clearAllMocks();
});

const renderSearchView = (params = {q: ''}, searchResultsStub = jest.fn()) => {
  return render(SearchView, {
    mocks: {
      $route: {
        params,
      },
    },
    stubs: ['router-link'],
  });
};

describe('SearchResults', () => {
  test('calls the search function of the api client', async () => {
    const spy = jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: [], meta: {total: 0}}));

    const query = 'a search query';
    renderSearchView({q: query});

    expect(spy).toHaveBeenCalledWith(expect.objectContaining({
      query,
      page: 1,
    }));
  });

  test('shows the number of total results', async () => {
    const track = {artists: [], spotifyId: 1};
    jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: [track, track], meta: {total: 2}}));

    const {getByText} = renderSearchView();

    await waitFor(() => getByText('2 tracks found'));
  });

  test('shows error page if no results are found', async () => {
    jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: [], meta: {total: 0}}));

    const {getByText} = renderSearchView();

    await waitFor(() => getByText(/ no results found/));
  });

  test('shows error page if api returns error', async () => {
    jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => new Promise((_, reject) => reject(new Error())));

    const {getByText} = renderSearchView();

    await waitFor(() => getByText(/ no results found/));
  });

  test('hides the pagination if total number of tracks is smaller than the page size', async () => {
    const track = {artists: [], spotifyId: 1, title: 'Test'};
    jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: [track], meta: {total: 1}}));

    const {queryByRole, getByText} = renderSearchView();
    await waitFor(() => getByText('Test'));

    const pagination = queryByRole('Pagination');
    expect(pagination).toBeNull();
  });

  test('shows the pagination if total number of tracks is greater than the page size', async () => {
    const track = {artists: [], spotifyId: 1, title: 'Test'};
    jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: [track], meta: {total: 999}}));

    const {getAllByText, getByText} = renderSearchView();
    await waitFor(() => getByText('Test'));

    getAllByText('4');
  });

  test('clicking a pagination item triggers an api call', async () => {
    const spy = jest.spyOn(TracksApi.prototype, 'tracksGet')
        .mockImplementation(() => Promise.resolve({data: [
          {artists: [], spotifyId: 1, title: 'Test'},
        ],
        meta: {total: 999}}),
        );

    const {getAllByText, getByText} = renderSearchView();
    await waitFor(() => getByText('Test'));

    const pageToClick = 4;
    const paginationItem = getAllByText(pageToClick)[0];
    await fireEvent.click(paginationItem);

    expect(spy).toHaveBeenCalledWith(expect.objectContaining({page: pageToClick}));
  });
});

