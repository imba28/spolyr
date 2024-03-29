/*
 * Spolyr
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type OAuthConfiguration struct {
	RedirectUrl string `json:"redirectUrl,omitempty"`

	ClientId string `json:"clientId,omitempty"`

	Scope string `json:"scope,omitempty"`
}

// AssertOAuthConfigurationRequired checks if the required fields are not zero-ed
func AssertOAuthConfigurationRequired(obj OAuthConfiguration) error {
	return nil
}

// AssertRecurseOAuthConfigurationRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of OAuthConfiguration (e.g. [][]OAuthConfiguration), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseOAuthConfigurationRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aOAuthConfiguration, ok := obj.(OAuthConfiguration)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertOAuthConfigurationRequired(aOAuthConfiguration)
	})
}
