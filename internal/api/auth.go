package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"net/http"
	"os"
	"strings"
)

var (
	permissions = []string{
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopePlaylistReadPrivate}
	scope = strings.Join(permissions, " ")
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectUrl()), spotifyauth.WithScopes(scope))
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func redirectUrl() string {
	var protocol = getEnv("PROTOCOL", "http")
	var domain = getEnv("DOMAIN", "localhost")

	var httpPort = getEnv("HTTP_PUBLIC_PORT", getEnv("HTTP_PORT", "8080"))
	if httpPort != "80" && httpPort != "443" {
		httpPort = ":" + httpPort
	} else {
		httpPort = ""
	}

	return fmt.Sprintf("%s://%s%s/auth/callback", protocol, domain, httpPort)
}

type AuthApiService struct {
	clientId   string
	hmacSecret []byte
}

func (a AuthApiService) AuthLoginPost(ctx context.Context, request openapi.AuthLoginPostRequest) (openapi.ImplResponse, error) {
	t, err := auth.Exchange(ctx, request.Code)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), errors.New("could not exchange code for token")
	}

	c := spotify.New(auth.Client(ctx, t))
	user, err := c.CurrentUser(ctx)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("could not get user info")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &customClaims{
		Token: *t,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: t.Expiry},
		},
	})
	tokenString, err := token.SignedString(a.hmacSecret)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("could not sign jwt")
	}

	headers := make(map[string][]string)
	headers["Set-Cookie"] = []string{"jwt=" + tokenString + "; HttpOnly; SameSite=None; Secure; Path=/api"}
	body := openapi.OAuthUserInfo{
		DisplayName: user.DisplayName,
	}
	if len(user.Images) > 0 {
		body.AvatarUrl = user.Images[0].URL
	}
	return openapi.ResponseWithHeaders(http.StatusOK, headers, body), nil
}

func (a AuthApiService) AuthConfigurationGet(ctx context.Context) (openapi.ImplResponse, error) {
	res := openapi.OAuthConfiguration{
		RedirectUrl: redirectUrl(),
		ClientId:    a.clientId,
		Scope:       scope,
	}

	return openapi.Response(http.StatusOK, res), nil
}

func (a AuthApiService) AuthRefreshGet(ctx context.Context) (openapi.ImplResponse, error) {
	return openapi.Response(http.StatusNotImplemented, nil), errors.New("AuthRefreshGet method not implemented")
}

func newAuthApiService(clientId string, secret []byte) AuthApiService {
	return AuthApiService{
		clientId:   clientId,
		hmacSecret: secret,
	}
}

var _ openapi.AuthApiServicer = AuthApiService{}
