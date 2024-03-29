/*
 * Spolyr
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type PlaylistInfo struct {
	SpotifyId string `json:"spotifyId,omitempty"`

	Name string `json:"name,omitempty"`

	CoverImage string `json:"coverImage,omitempty"`

	TrackCount int32 `json:"trackCount,omitempty"`

	Owner string `json:"owner,omitempty"`

	IsPublic bool `json:"isPublic,omitempty"`

	IsCollaborative bool `json:"isCollaborative,omitempty"`
}

// AssertPlaylistInfoRequired checks if the required fields are not zero-ed
func AssertPlaylistInfoRequired(obj PlaylistInfo) error {
	return nil
}

// AssertRecursePlaylistInfoRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of PlaylistInfo (e.g. [][]PlaylistInfo), otherwise ErrTypeAssertionError is thrown.
func AssertRecursePlaylistInfoRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aPlaylistInfo, ok := obj.(PlaylistInfo)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertPlaylistInfoRequired(aPlaylistInfo)
	})
}
