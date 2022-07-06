import {AuthApi} from '@/openapi';

const authApi = new AuthApi();

/**
 *
 * @param{any} request
 */
function reset(request) {
  const headers = request.req._headers;
  const path = request.req.path;

  request.req.abort();
  request.called = false;
  delete request.req;

  for (const k in headers) {
    if (Object.prototype.hasOwnProperty.call(headers, k)) {
      request.set(k, headers[k]);
    }
  }

  if (!request.qs) {
    request.req.path = path;
  }
}

reset;

/**
 * @param{Number} Request
 */
function refreshTokenMiddleware(Request) {
  const originalEnd = Request.end;

  Request.end = async (callback) => {
    originalEnd.call(Request, async function(err, response) {
      if (err && err.status === 401) {
        debugger;

        try {
          await authApi.authRefreshGet();

          originalEnd.call(Request, callback);
        } catch (e) {
          err = new Error('Could not refresh token');
        }
      }

      callback(err, response);
    });
  };
}

export default refreshTokenMiddleware;
