/**
 * Spolyr
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */


import ApiClient from "../ApiClient";
import PlaylistsGet200Response from '../model/PlaylistsGet200Response';

/**
* Playlists service.
* @module api/PlaylistsApi
* @version 1.0.0
*/
export default class PlaylistsApi {

    /**
    * Constructs a new PlaylistsApi. 
    * @alias module:api/PlaylistsApi
    * @class
    * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
    * default to {@link module:ApiClient#instance} if unspecified.
    */
    constructor(apiClient) {
        this.apiClient = apiClient || ApiClient.instance;
    }



    /**
     * Returns a list of your saved playlists
     * @param {Object} opts Optional parameters
     * @param {Number} opts.page Current page number (default to 1)
     * @param {Number} opts.limit Limits the size of the result size (default to 25)
     * @return {Promise} a {@link https://www.promisejs.org/|Promise}, with an object containing data of type {@link module:model/PlaylistsGet200Response} and HTTP response
     */
    playlistsGetWithHttpInfo(opts) {
      opts = opts || {};
      let postBody = null;

      let pathParams = {
      };
      let queryParams = {
        'page': opts['page'],
        'limit': opts['limit']
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = ['cookieAuth'];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = PlaylistsGet200Response;
      return this.apiClient.callApi(
        '/playlists', 'GET',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null
      );
    }

    /**
     * Returns a list of your saved playlists
     * @param {Object} opts Optional parameters
     * @param {Number} opts.page Current page number (default to 1)
     * @param {Number} opts.limit Limits the size of the result size (default to 25)
     * @return {Promise} a {@link https://www.promisejs.org/|Promise}, with data of type {@link module:model/PlaylistsGet200Response}
     */
    playlistsGet(opts) {
      return this.playlistsGetWithHttpInfo(opts)
        .then(function(response_and_data) {
          return response_and_data.data;
        });
    }


}