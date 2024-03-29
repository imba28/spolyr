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
 * The UserResponse model module.
 * @module model/UserResponse
 * @version 1.0.0
 */
class UserResponse {
    /**
     * Constructs a new <code>UserResponse</code>.
     * @alias module:model/UserResponse
     */
    constructor() { 
        
        UserResponse.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>UserResponse</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/UserResponse} obj Optional instance to populate.
     * @return {module:model/UserResponse} The populated <code>UserResponse</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new UserResponse();

            if (data.hasOwnProperty('username')) {
                obj['username'] = ApiClient.convertToType(data['username'], 'String');
            }
            if (data.hasOwnProperty('avatar')) {
                obj['avatar'] = ApiClient.convertToType(data['avatar'], 'String');
            }
        }
        return obj;
    }


}

/**
 * @member {String} username
 */
UserResponse.prototype['username'] = undefined;

/**
 * @member {String} avatar
 */
UserResponse.prototype['avatar'] = undefined;






export default UserResponse;

