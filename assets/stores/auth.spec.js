import {setActivePinia, createPinia} from 'pinia';
import {useAuthStore} from '@/stores/auth';
import ApiClient from '../openapi/ApiClient';

describe('Counter Store', () => {
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
    const err = new Error('api error');
    jest.spyOn(ApiClient.prototype, 'callApi')
        .mockImplementation(() => new Promise((_, reject) => reject(err)));

    const authStore = useAuthStore();

    expect.assertions(3);
    expect(authStore.login('validSpotifyOAuthCode')).rejects.toEqual(err);

    expect(authStore.avatarUrl).toBeNull();
    expect(authStore.displayName).toBeNull();
  });
});
