/*
 * Spolyr
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// TracksApiController binds http requests to an api service and writes the service results to the http response
type TracksApiController struct {
	service      TracksApiServicer
	errorHandler ErrorHandler
}

// TracksApiOption for how the controller is set up.
type TracksApiOption func(*TracksApiController)

// WithTracksApiErrorHandler inject ErrorHandler into controller
func WithTracksApiErrorHandler(h ErrorHandler) TracksApiOption {
	return func(c *TracksApiController) {
		c.errorHandler = h
	}
}

// NewTracksApiController creates a default api controller
func NewTracksApiController(s TracksApiServicer, opts ...TracksApiOption) Router {
	controller := &TracksApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the TracksApiController
func (c *TracksApiController) Routes() Routes {
	return Routes{
		{
			"TracksGet",
			strings.ToUpper("Get"),
			"/api/tracks",
			c.TracksGet,
		},
		{
			"TracksIdGet",
			strings.ToUpper("Get"),
			"/api/tracks/{id}",
			c.TracksIdGet,
		},
		{
			"TracksIdPatch",
			strings.ToUpper("Patch"),
			"/api/tracks/{id}",
			c.TracksIdPatch,
		},
		{
			"TracksStatsGet",
			strings.ToUpper("Get"),
			"/api/tracks-stats",
			c.TracksStatsGet,
		},
	}
}

// TracksGet - Returns a list of tracks
func (c *TracksApiController) TracksGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageParam, err := parseInt32Parameter(query.Get("page"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	limitParam, err := parseInt32Parameter(query.Get("limit"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	queryParam := query.Get("query")
	result, err := c.service.TracksGet(r.Context(), pageParam, limitParam, queryParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}

// TracksIdGet - Returns a track
func (c *TracksApiController) TracksIdGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam := params["id"]

	result, err := c.service.TracksIdGet(r.Context(), idParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}

// TracksIdPatch - Updates lyrics of a track
func (c *TracksApiController) TracksIdPatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idParam := params["id"]

	lyricsParam := Lyrics{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&lyricsParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertLyricsRequired(lyricsParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.TracksIdPatch(r.Context(), idParam, lyricsParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}

// TracksStatsGet - Returns stats about your index
func (c *TracksApiController) TracksStatsGet(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.TracksStatsGet(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)

}