import {setActivePinia, createPinia} from 'pinia';
import {useAuthStore} from '@/stores/auth';
import ApiClient from '../openapi/ApiClient';

describe('Authentication Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('sets displayName and avatarUrl after an successful login', async () => {
    const mockApiResponse = {
      avatarUrl: 'https://foobar.com/avatar.png',
      displayName: 'Test',
    };
    const callApiMock = jest.spyOn(ApiClient.prototype, 'callApi')
        .mockImplementation(() => Promise.resolve({data: mockApiResponse}));

    const authStore = useAuthStore();

    await authStore.login('validSpotifyOAuthCode');

    expect(callApiMock).toHaveBeenCalledTimes(1);
    expect(authStore.avatarUrl).toEqual(mockApiResponse.avatarUrl);
    expect(authStore.displayName).toEqual(mockApiResponse.displayName);
  });

  it('throws an error if api returns an error', async () => {
    expect.assertions(3);

    const err = new Error('api error');
    jest.spyOn(ApiClient.prototype, 'callApi')
        .mockImplementation(() => new Promise((_, reject) => reject(err)));

    const authStore = useAuthStore();

    expect(authStore.login('validSpotifyOAuthCode')).rejects.toEqual(err);
    expect(authStore.avatarUrl).toBeNull();
    expect(authStore.displayName).toBeNull();
  });

  it('unsets avatar and display name after signing out', async () => {
    expect.assertions(3);
    const spy = jest.spyOn(ApiClient.prototype, 'callApi')
        .mockImplementation(() => Promise.resolve({}));

    const authStore = useAuthStore();
    authStore.avatarUrl = 'https://spotify.com/1.jpg';
    authStore.displayName = 'Test user';

    await authStore.logout();

    expect(spy).toHaveBeenCalled();
    expect(authStore.avatarUrl).toBeNull();
    expect(authStore.displayName).toBeNull();
  });

  it('throws an error if signing out was unsuccessful', async () => {
    expect.assertions(3);

    const err = new Error('something went wrong');
    jest.spyOn(ApiClient.prototype, 'callApi')
        .mockImplementation(() => new Promise((_, reject) => reject(err)));

    const authStore = useAuthStore();
    authStore.avatarUrl = 'https://spotify.com/1.jpg';
    authStore.displayName = 'Test user';

    expect(authStore.logout()).rejects.toEqual(err);
    expect(authStore.avatarUrl).toEqual('https://spotify.com/1.jpg');
    expect(authStore.displayName).toEqual('Test user');
  });
});
