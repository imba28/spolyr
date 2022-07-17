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
import TrackDetailAllOf from './TrackDetailAllOf';
import TrackInfo from './TrackInfo';

/**
 * The TrackDetail model module.
 * @module model/TrackDetail
 * @version 1.0.0
 */
class TrackDetail {
    /**
     * Constructs a new <code>TrackDetail</code>.
     * @alias module:model/TrackDetail
     * @implements module:model/TrackInfo
     * @implements module:model/TrackDetailAllOf
     * @param spotifyId {String} 
     * @param title {String} 
     * @param hasLyrics {Boolean} 
     * @param lyrics {String} 
     * @param lyricsImportErrorCount {Number} 
     */
    constructor(spotifyId, title, hasLyrics, lyrics, lyricsImportErrorCount) { 
        TrackInfo.initialize(this, spotifyId, title, hasLyrics);TrackDetailAllOf.initialize(this, lyrics, lyricsImportErrorCount, hasLyrics);
        TrackDetail.initialize(this, spotifyId, title, hasLyrics, lyrics, lyricsImportErrorCount);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, spotifyId, title, hasLyrics, lyrics, lyricsImportErrorCount) { 
        obj['spotifyId'] = spotifyId;
        obj['title'] = title;
        obj['hasLyrics'] = hasLyrics;
        obj['lyrics'] = lyrics;
        obj['lyricsImportErrorCount'] = lyricsImportErrorCount;
    }

    /**
     * Constructs a <code>TrackDetail</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/TrackDetail} obj Optional instance to populate.
     * @return {module:model/TrackDetail} The populated <code>TrackDetail</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new TrackDetail();
            TrackInfo.constructFromObject(data, obj);
            TrackDetailAllOf.constructFromObject(data, obj);

            if (data.hasOwnProperty('spotifyId')) {
                obj['spotifyId'] = ApiClient.convertToType(data['spotifyId'], 'String');
            }
            if (data.hasOwnProperty('title')) {
                obj['title'] = ApiClient.convertToType(data['title'], 'String');
            }
            if (data.hasOwnProperty('album')) {
                obj['album'] = ApiClient.convertToType(data['album'], 'String');
            }
            if (data.hasOwnProperty('coverImage')) {
                obj['coverImage'] = ApiClient.convertToType(data['coverImage'], 'String');
            }
            if (data.hasOwnProperty('previewURL')) {
                obj['previewURL'] = ApiClient.convertToType(data['previewURL'], 'String');
            }
            if (data.hasOwnProperty('artists')) {
                obj['artists'] = ApiClient.convertToType(data['artists'], ['String']);
            }
            if (data.hasOwnProperty('hasLyrics')) {
                obj['hasLyrics'] = ApiClient.convertToType(data['hasLyrics'], 'Boolean');
            }
            if (data.hasOwnProperty('language')) {
                obj['language'] = ApiClient.convertToType(data['language'], 'String');
            }
            if (data.hasOwnProperty('lyrics')) {
                obj['lyrics'] = ApiClient.convertToType(data['lyrics'], 'String');
            }
            if (data.hasOwnProperty('lyricsImportErrorCount')) {
                obj['lyricsImportErrorCount'] = ApiClient.convertToType(data['lyricsImportErrorCount'], 'Number');
            }
        }
        return obj;
    }


}

/**
 * @member {String} spotifyId
 */
TrackDetail.prototype['spotifyId'] = undefined;

/**
 * @member {String} title
 */
TrackDetail.prototype['title'] = undefined;

/**
 * @member {String} album
 */
TrackDetail.prototype['album'] = undefined;

/**
 * @member {String} coverImage
 */
TrackDetail.prototype['coverImage'] = undefined;

/**
 * @member {String} previewURL
 */
TrackDetail.prototype['previewURL'] = undefined;

/**
 * @member {Array.<String>} artists
 */
TrackDetail.prototype['artists'] = undefined;

/**
 * @member {Boolean} hasLyrics
 */
TrackDetail.prototype['hasLyrics'] = undefined;

/**
 * @member {String} language
 */
TrackDetail.prototype['language'] = undefined;

/**
 * @member {String} lyrics
 */
TrackDetail.prototype['lyrics'] = undefined;

/**
 * @member {Number} lyricsImportErrorCount
 */
TrackDetail.prototype['lyricsImportErrorCount'] = undefined;


// Implement TrackInfo interface:
/**
 * @member {String} spotifyId
 */
TrackInfo.prototype['spotifyId'] = undefined;
/**
 * @member {String} title
 */
TrackInfo.prototype['title'] = undefined;
/**
 * @member {String} album
 */
TrackInfo.prototype['album'] = undefined;
/**
 * @member {String} coverImage
 */
TrackInfo.prototype['coverImage'] = undefined;
/**
 * @member {String} previewURL
 */
TrackInfo.prototype['previewURL'] = undefined;
/**
 * @member {Array.<String>} artists
 */
TrackInfo.prototype['artists'] = undefined;
/**
 * @member {Boolean} hasLyrics
 */
TrackInfo.prototype['hasLyrics'] = undefined;
/**
 * @member {String} language
 */
TrackInfo.prototype['language'] = undefined;
// Implement TrackDetailAllOf interface:
/**
 * @member {String} lyrics
 */
TrackDetailAllOf.prototype['lyrics'] = undefined;
/**
 * @member {Number} lyricsImportErrorCount
 */
TrackDetailAllOf.prototype['lyricsImportErrorCount'] = undefined;
/**
 * @member {Boolean} hasLyrics
 */
TrackDetailAllOf.prototype['hasLyrics'] = undefined;




export default TrackDetail;

