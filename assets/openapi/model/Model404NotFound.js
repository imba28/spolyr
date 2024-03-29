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
import Message from './Message';

/**
 * The Model404NotFound model module.
 * @module model/Model404NotFound
 * @version 1.0.0
 */
class Model404NotFound {
    /**
     * Constructs a new <code>Model404NotFound</code>.
     * @alias module:model/Model404NotFound
     * @implements module:model/Message
     * @param code {Number} 
     * @param message {String} 
     */
    constructor(code, message) { 
        Message.initialize(this, code, message);
        Model404NotFound.initialize(this, code, message);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, code, message) { 
        obj['code'] = code;
        obj['message'] = message;
    }

    /**
     * Constructs a <code>Model404NotFound</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/Model404NotFound} obj Optional instance to populate.
     * @return {module:model/Model404NotFound} The populated <code>Model404NotFound</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new Model404NotFound();
            Message.constructFromObject(data, obj);

            if (data.hasOwnProperty('code')) {
                obj['code'] = ApiClient.convertToType(data['code'], 'Number');
            }
            if (data.hasOwnProperty('message')) {
                obj['message'] = ApiClient.convertToType(data['message'], 'String');
            }
        }
        return obj;
    }


}

/**
 * @member {Number} code
 */
Model404NotFound.prototype['code'] = undefined;

/**
 * @member {String} message
 */
Model404NotFound.prototype['message'] = undefined;


// Implement Message interface:
/**
 * @member {Number} code
 */
Message.prototype['code'] = undefined;
/**
 * @member {String} message
 */
Message.prototype['message'] = undefined;




export default Model404NotFound;

