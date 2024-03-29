/*
 * Spolyr
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type PlaylistsGet200Response struct {
	Meta PaginationMetadata `json:"meta,omitempty"`

	Data []PlaylistInfo `json:"data,omitempty"`
}

// AssertPlaylistsGet200ResponseRequired checks if the required fields are not zero-ed
func AssertPlaylistsGet200ResponseRequired(obj PlaylistsGet200Response) error {
	if err := AssertPaginationMetadataRequired(obj.Meta); err != nil {
		return err
	}
	for _, el := range obj.Data {
		if err := AssertPlaylistInfoRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecursePlaylistsGet200ResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of PlaylistsGet200Response (e.g. [][]PlaylistsGet200Response), otherwise ErrTypeAssertionError is thrown.
func AssertRecursePlaylistsGet200ResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aPlaylistsGet200Response, ok := obj.(PlaylistsGet200Response)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertPlaylistsGet200ResponseRequired(aPlaylistsGet200Response)
	})
}
