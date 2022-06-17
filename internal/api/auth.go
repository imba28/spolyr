package api

import (
	"context"
	"errors"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"net/http"
)

type AuthApiService struct {
}

func (a AuthApiService) AuthCallbackGet(ctx context.Context) (openapi.ImplResponse, error) {
	// TODO - update AuthCallbackGet with the required logic for this service method.
	// Add api_auth_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return openapi.Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return openapi.Response(400, nil),nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("AuthCallbackGet method not implemented")
}

func (a AuthApiService) AuthLoginGet(ctx context.Context) (openapi.ImplResponse, error) {
	// TODO - update AuthLoginGet with the required logic for this service method.
	// Add api_auth_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(302, {}) or use other options such as http.Ok ...
	//return openapi.Response(302, nil),nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("AuthLoginGet method not implemented")
}

func (a AuthApiService) AuthRefreshGet(ctx context.Context) (openapi.ImplResponse, error) {
	// TODO - update AuthRefreshGet with the required logic for this service method.
	// Add api_auth_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, {}) or use other options such as http.Ok ...
	//return openapi.Response(201, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return openapi.Response(401, nil),nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("AuthRefreshGet method not implemented")
}

func NewAuthApiService() AuthApiService {
	return AuthApiService{}
}

var _ openapi.AuthApiServicer = AuthApiService{}
