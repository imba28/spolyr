import {AuthApi} from '@/openapi';

const authApi = new AuthApi();

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
      if (Request.__proto__._refreshed || !err || err.status !== 401) {
        Request.__proto__._refreshed = false;
        callback(err, response);
        return;
      }

      Request.__proto__._refreshed = true;

      try {
        // acquire fresh access token...
        await authApi.authRefreshGet();

        originalEnd.call(Request, callback);
      } catch (e) {
        callback(new Error('could not refresh token'), response);
      }
    });
  };
}

export default refreshTokenMiddleware;
