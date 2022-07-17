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

import ApiClient from '../ApiClient';
import PaginationMetadata from './PaginationMetadata';
import TrackInfo from './TrackInfo';

/**
 * The TracksGet200Response model module.
 * @module model/TracksGet200Response
 * @version 1.0.0
 */
class TracksGet200Response {
    /**
     * Constructs a new <code>TracksGet200Response</code>.
     * @alias module:model/TracksGet200Response
     * @param meta {module:model/PaginationMetadata} 
     * @param data {Array.<module:model/TrackInfo>} 
     */
    constructor(meta, data) { 
        
        TracksGet200Response.initialize(this, meta, data);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, meta, data) { 
        obj['meta'] = meta;
        obj['data'] = data;
    }

    /**
     * Constructs a <code>TracksGet200Response</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/TracksGet200Response} obj Optional instance to populate.
     * @return {module:model/TracksGet200Response} The populated <code>TracksGet200Response</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new TracksGet200Response();

            if (data.hasOwnProperty('meta')) {
                obj['meta'] = PaginationMetadata.constructFromObject(data['meta']);
            }
            if (data.hasOwnProperty('data')) {
                obj['data'] = ApiClient.convertToType(data['data'], [TrackInfo]);
            }
        }
        return obj;
    }


}

/**
 * @member {module:model/PaginationMetadata} meta
 */
TracksGet200Response.prototype['meta'] = undefined;

/**
 * @member {Array.<module:model/TrackInfo>} data
 */
TracksGet200Response.prototype['data'] = undefined;






export default TracksGet200Response;

