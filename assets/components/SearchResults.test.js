import {fireEvent, render} from '@testing-library/vue';

import SearchResults from '@/components/SearchResults';

const renderSearchResults = (props, routerLinkStub = jest.fn()) => {
  return render(SearchResults, {
    propsData: props,
    components: {
      'router-link': routerLinkStub,
    },
  });
};

describe('SearchResults', () => {
  test('Shows the number of results', async () => {
    const tracks = [
      {title: 'Foo', artists: []},
      {title: 'Bar', artists: []},
      {title: 'Baz', artists: []},
    ];

    const {getByText} = renderSearchResults({
      items: tracks,
    });

    getByText(`${tracks.length} tracks found`);
  });

  test('Renders the list of results', async () => {
    const tracks = [
      {title: 'Foo', artists: ['Artist A', 'Artist B']},
      {title: 'Bar', artists: [], album: 'Album A'},
      {title: 'Baz', artists: []},
    ];

    const {getByText} = renderSearchResults({
      items: tracks,
    });

    getByText('Foo');
    getByText('Artist A, Artist B');

    getByText('Bar');
    getByText('Album A');

    getByText('Baz');
  });

  test('Table uses the router-link component', async () => {
    const tracks = [
      {title: 'Foo', artists: []},
      {title: 'Bar', artists: []},
      {title: 'Baz', artists: []},
    ];

    const routerLinkMock = jest.fn();
    renderSearchResults({
      items: tracks,
    }, routerLinkMock);

    expect(routerLinkMock).toHaveBeenCalled();
  });

  test('Renders a play buttons if track contains an audio snippet', async () => {
    const tracks = [
      {title: 'Foo', artists: [], previewURL: '1'},
      {title: 'Bar', artists: []},
      {title: 'Baz', artists: [], previewURL: null},
    ];
    const {getAllByLabelText} = renderSearchResults({
      items: tracks,
    });

    const playButtons = getAllByLabelText('play fill');
    expect(playButtons).toHaveLength(1);
  });

  test('Renders a stop buttons if track contains an audio snippet and in playing state', async () => {
    const tracks = [
      {title: 'Foo', artists: [], previewURL: 'http://foo.com/foo.mp3'},
      {title: 'Bar', artists: [], previewURL: 'http://bar.com/foo.mp3'},
      {title: 'Baz', artists: [], previewURL: 'http://baz.com/foo.mp3'},
    ];
    const {getAllByLabelText} = renderSearchResults({
      items: tracks,
      playing: tracks[0].previewURL,
    });

    const stopButtons = getAllByLabelText('pause fill');
    expect(stopButtons).toHaveLength(1);
  });

  test('Clicking a play button emits a `play` event', async () => {
    const tracks = [
      {title: 'Foo', artists: [], previewURL: 'http://foo.com/foo.mp3'},
      {title: 'Bar', artists: []},
    ];
    const {getByLabelText, emitted} = renderSearchResults({
      items: tracks,
    });

    await fireEvent.click(getByLabelText('play fill'));

    expect(emitted().play).toHaveLength(1);
    expect(emitted().play[0][0]).toEqual(tracks[0].previewURL);
  });

  test('Clicking a pause button emits a `stop` event', async () => {
    const tracks = [
      {title: 'Foo', artists: [], previewURL: 'http://foo.com/foo.mp3'},
      {title: 'Bar', artists: [], previewURL: 'http://bar.com/foo.mp3'},
      {title: 'Baz', artists: [], previewURL: 'http://baz.com/foo.mp3'},
    ];

    const {getByLabelText, emitted} = renderSearchResults({
      items: tracks,
      playing: tracks[0].previewURL,
    });

    await fireEvent.click(getByLabelText('pause fill'));

    expect(emitted().stop).toHaveLength(1);
    expect(emitted().stop[0][0]).toEqual(tracks[0].previewURL);
  });
});
