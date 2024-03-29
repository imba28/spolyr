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


import ApiClient from './ApiClient';
import AuthLoginPostRequest from './model/AuthLoginPostRequest';
import Lyrics from './model/Lyrics';
import LyricsImportStatus from './model/LyricsImportStatus';
import Message from './model/Message';
import Model401Unauthorized from './model/Model401Unauthorized';
import Model404NotFound from './model/Model404NotFound';
import Model500InternalError from './model/Model500InternalError';
import OAuthConfiguration from './model/OAuthConfiguration';
import OAuthUserInfo from './model/OAuthUserInfo';
import PaginationMetadata from './model/PaginationMetadata';
import PlaylistInfo from './model/PlaylistInfo';
import PlaylistsGet200Response from './model/PlaylistsGet200Response';
import TrackDetail from './model/TrackDetail';
import TrackDetailAllOf from './model/TrackDetailAllOf';
import TrackInfo from './model/TrackInfo';
import TracksGet200Response from './model/TracksGet200Response';
import TracksStats from './model/TracksStats';
import UserResponse from './model/UserResponse';
import AuthApi from './api/AuthApi';
import ImportApi from './api/ImportApi';
import PlaylistsApi from './api/PlaylistsApi';
import TracksApi from './api/TracksApi';


/**
* JS API client generated by OpenAPI Generator.<br>
* The <code>index</code> module provides access to constructors for all the classes which comprise the public API.
* <p>
* An AMD (recommended!) or CommonJS application will generally do something equivalent to the following:
* <pre>
* var @/openapi = require('index'); // See note below*.
* var xxxSvc = new @/openapi.XxxApi(); // Allocate the API class we're going to use.
* var yyyModel = new @/openapi.Yyy(); // Construct a model instance.
* yyyModel.someProperty = 'someValue';
* ...
* var zzz = xxxSvc.doSomething(yyyModel); // Invoke the service.
* ...
* </pre>
* <em>*NOTE: For a top-level AMD script, use require(['index'], function(){...})
* and put the application logic within the callback function.</em>
* </p>
* <p>
* A non-AMD browser application (discouraged) might do something like this:
* <pre>
* var xxxSvc = new @/openapi.XxxApi(); // Allocate the API class we're going to use.
* var yyy = new @/openapi.Yyy(); // Construct a model instance.
* yyyModel.someProperty = 'someValue';
* ...
* var zzz = xxxSvc.doSomething(yyyModel); // Invoke the service.
* ...
* </pre>
* </p>
* @module index
* @version 1.0.0
*/
export {
    /**
     * The ApiClient constructor.
     * @property {module:ApiClient}
     */
    ApiClient,

    /**
     * The AuthLoginPostRequest model constructor.
     * @property {module:model/AuthLoginPostRequest}
     */
    AuthLoginPostRequest,

    /**
     * The Lyrics model constructor.
     * @property {module:model/Lyrics}
     */
    Lyrics,

    /**
     * The LyricsImportStatus model constructor.
     * @property {module:model/LyricsImportStatus}
     */
    LyricsImportStatus,

    /**
     * The Message model constructor.
     * @property {module:model/Message}
     */
    Message,

    /**
     * The Model401Unauthorized model constructor.
     * @property {module:model/Model401Unauthorized}
     */
    Model401Unauthorized,

    /**
     * The Model404NotFound model constructor.
     * @property {module:model/Model404NotFound}
     */
    Model404NotFound,

    /**
     * The Model500InternalError model constructor.
     * @property {module:model/Model500InternalError}
     */
    Model500InternalError,

    /**
     * The OAuthConfiguration model constructor.
     * @property {module:model/OAuthConfiguration}
     */
    OAuthConfiguration,

    /**
     * The OAuthUserInfo model constructor.
     * @property {module:model/OAuthUserInfo}
     */
    OAuthUserInfo,

    /**
     * The PaginationMetadata model constructor.
     * @property {module:model/PaginationMetadata}
     */
    PaginationMetadata,

    /**
     * The PlaylistInfo model constructor.
     * @property {module:model/PlaylistInfo}
     */
    PlaylistInfo,

    /**
     * The PlaylistsGet200Response model constructor.
     * @property {module:model/PlaylistsGet200Response}
     */
    PlaylistsGet200Response,

    /**
     * The TrackDetail model constructor.
     * @property {module:model/TrackDetail}
     */
    TrackDetail,

    /**
     * The TrackDetailAllOf model constructor.
     * @property {module:model/TrackDetailAllOf}
     */
    TrackDetailAllOf,

    /**
     * The TrackInfo model constructor.
     * @property {module:model/TrackInfo}
     */
    TrackInfo,

    /**
     * The TracksGet200Response model constructor.
     * @property {module:model/TracksGet200Response}
     */
    TracksGet200Response,

    /**
     * The TracksStats model constructor.
     * @property {module:model/TracksStats}
     */
    TracksStats,

    /**
     * The UserResponse model constructor.
     * @property {module:model/UserResponse}
     */
    UserResponse,

    /**
    * The AuthApi service constructor.
    * @property {module:api/AuthApi}
    */
    AuthApi,

    /**
    * The ImportApi service constructor.
    * @property {module:api/ImportApi}
    */
    ImportApi,

    /**
    * The PlaylistsApi service constructor.
    * @property {module:api/PlaylistsApi}
    */
    PlaylistsApi,

    /**
    * The TracksApi service constructor.
    * @property {module:api/TracksApi}
    */
    TracksApi
};
