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

/**
 * The OAuthConfiguration model module.
 * @module model/OAuthConfiguration
 * @version 1.0.0
 */
class OAuthConfiguration {
    /**
     * Constructs a new <code>OAuthConfiguration</code>.
     * @alias module:model/OAuthConfiguration
     */
    constructor() { 
        
        OAuthConfiguration.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>OAuthConfiguration</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/OAuthConfiguration} obj Optional instance to populate.
     * @return {module:model/OAuthConfiguration} The populated <code>OAuthConfiguration</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new OAuthConfiguration();

            if (data.hasOwnProperty('redirectUrl')) {
                obj['redirectUrl'] = ApiClient.convertToType(data['redirectUrl'], 'String');
            }
            if (data.hasOwnProperty('clientId')) {
                obj['clientId'] = ApiClient.convertToType(data['clientId'], 'String');
            }
            if (data.hasOwnProperty('scope')) {
                obj['scope'] = ApiClient.convertToType(data['scope'], 'String');
            }
        }
        return obj;
    }


}

/**
 * @member {String} redirectUrl
 */
OAuthConfiguration.prototype['redirectUrl'] = undefined;

/**
 * @member {String} clientId
 */
OAuthConfiguration.prototype['clientId'] = undefined;

/**
 * @member {String} scope
 */
OAuthConfiguration.prototype['scope'] = undefined;






export default OAuthConfiguration;

