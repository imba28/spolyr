import refreshTokenMiddleware from './superagent';
import {AuthApi} from '@/openapi';

afterEach(() => {
  jest.clearAllMocks();

  // since superagent plugins modify the prototype chain reset the module
  jest.resetModules();
});

describe('refreshTokenMiddleware superagent plugin', () => {
  it('refreshes the token if server returns 401 and retries the original request.', (done) => {
    expect.assertions(3);
    const expectedData = {id: 1, title: 'a protected track'};

    const authMock = jest.spyOn(AuthApi.prototype, 'authRefreshGet')
        .mockImplementation(() => Promise.resolve());
    const endMock = jest.fn()
        .mockImplementationOnce((callback) => {
          callback({status: 401});
        })
        .mockImplementation((callback) => {
          callback({status: 200}, expectedData);
        });

    const superagent = require('superagent');
    const request = superagent('post', '/protected/area');
    request.end = endMock;
    request.use(refreshTokenMiddleware);

    request.end((err, data) => {
      expect(authMock).toHaveBeenCalledTimes(1);
      expect(endMock).toHaveBeenCalledTimes(2);
      expect(data).toEqual(data);

      done();
    });
  });

  it('tries to refresh the token only once', (done) => {
    expect.assertions(2);

    const authMock = jest.spyOn(AuthApi.prototype, 'authRefreshGet')
        .mockImplementation(() => Promise.resolve());
    const endMock = jest.fn()
        .mockImplementation((callback) => {
          callback({status: 401}, null); // endpoint returns an error again
        });

    const superagent = require('superagent');
    const request = superagent('post', '/protected/area');
    request.end = endMock;
    request.use(refreshTokenMiddleware);

    request.end(() => {
      expect(authMock).toHaveBeenCalledTimes(1);
      expect(endMock).toHaveBeenCalledTimes(2);

      done();
    });
  });

  it('does not repeat the request if refreshing the token returns an error', (done) => {
    expect.assertions(3);

    const authMock = jest.spyOn(AuthApi.prototype, 'authRefreshGet')
        .mockImplementation(() => Promise.reject(new Error('something went wrong')));
    const endMock = jest.fn()
        .mockImplementationOnce((callback) => {
          callback({status: 401}, null);
        });

    const superagent = require('superagent');
    const request = superagent('post', '/protected/area');
    request.end = endMock;
    request.use(refreshTokenMiddleware);

    request.end((err) => {
      expect(authMock).toHaveBeenCalledTimes(1);
      expect(endMock).toHaveBeenCalledTimes(1);
      expect(err.message).toEqual('could not refresh token');

      done();
    });
  });
});
