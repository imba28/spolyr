import {AuthApi} from '@/openapi';

const authApi = new AuthApi();

/**
 *
 * @param{Request} request
 */
function reset(request) {
  request._endCalled = false;
}

/**
 * This plugin intercepts http responses.
 * If a requests triggers a 401 response, two things happen:
 *   1. The plugin tries to refresh the JWT
 *   2. If the tokens were refreshed, the original request is repeated once again.
 *
 * @param{Request} Request
 */
function refreshTokenMiddleware(Request) {
  const originalEnd = Request.end;

  Request.end = async (callback) => {
    originalEnd.call(Request, async function(err, response) {
      if (!err || err.status !== 401) {
        callback(err, response);
        return;
      }

      try {
        // acquire fresh access token...
        await authApi.authRefreshGet();

        // repeat the original request
        reset(Request);
        originalEnd.call(Request, callback);
      } catch (e) {
        callback(new Error('could not refresh token'), response);
      }
    });
  };
}

export default refreshTokenMiddleware;
