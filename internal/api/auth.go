package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	jwt2 "github.com/imba28/spolyr/internal/jwt"
	"github.com/imba28/spolyr/internal/openapi/openapi"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
	"time"
)

type authContextKey int

const (
	spotifyTokenKey authContextKey = iota
	spotifyRefreshTokenKey
	jwtRefreshKey
	jwtAccessKey
)

var (
	permissions = []string{
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopePlaylistReadPrivate}
	scope = strings.Join(permissions, " ")
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectUrl()), spotifyauth.WithScopes(scope))

	accessTokenExpiry  = time.Minute * 10
	refreshTokenExpiry = time.Hour * 24
)

var (
	ErrNotAuthenticated = errors.New("no authentication provided")
)

func refreshTokenFromContext(ctx context.Context) *string {
	if t, ok := ctx.Value(jwtRefreshKey).(string); ok {
		return &t
	}
	return nil
}

func accessTokenFromContext(ctx context.Context) *string {
	if t, ok := ctx.Value(jwtAccessKey).(string); ok {
		return &t
	}
	return nil
}

func oauthTokenFromContext(ctx context.Context) *oauth2.Token {
	if t, ok := ctx.Value(spotifyTokenKey).(oauth2.Token); ok {
		return &t
	}
	return nil
}

func oauthClientFromContext(ctx context.Context) *spotify.Client {
	t := oauthTokenFromContext(ctx)
	if t == nil {
		return nil
	}
	return spotify.New(auth.Client(ctx, t))
}

func oauthRefreshTokenFromContext(ctx context.Context) string {
	if t, ok := ctx.Value(spotifyRefreshTokenKey).(string); ok {
		return t
	}
	return ""
}

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

func isAuthenticated(ctx context.Context) bool {
	t := accessTokenFromContext(ctx)
	return t != nil
}

func hasValidRefreshToken(ctx context.Context) bool {
	t := refreshTokenFromContext(ctx)
	// todo: check database if token has been revoked
	return t != nil
}

func AuthenticationMiddleware(j jwt2.JWT) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if r.Method != http.MethodOptions {
				if c, err := r.Cookie("jwt"); err == nil {
					claims, valid := j.ValidateAccessToken(c.Value)
					if valid {
						ctx = context.WithValue(ctx, spotifyTokenKey, claims.Token)
						ctx = context.WithValue(ctx, jwtAccessKey, c.Value)
					}
				}
				if c, err := r.Cookie("jwt-refresh"); err == nil {
					if claims, valid := j.ValidateRefreshToken(c.Value); valid {
						ctx = context.WithValue(ctx, spotifyRefreshTokenKey, claims.RefreshToken)
						ctx = context.WithValue(ctx, jwtRefreshKey, c.Value)
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type AuthApiService struct {
	clientId string
	jwt      jwt2.JWT
}

func (a AuthApiService) jwtTokenHeaders(t oauth2.Token, generateRefreshToken bool) (map[string][]string, error) {
	headers := make(map[string][]string)
	var cookies []string

	accessToken, err := a.jwt.NewAccessToken(t, time.Now().Add(accessTokenExpiry))
	if err != nil {
		return nil, errors.New("could not sign access jwt")
	}
	accessTokenCookie := http.Cookie{
		Name:     "jwt",
		Path:     "/api",
		Value:    accessToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	cookies = append(cookies, accessTokenCookie.String())

	if generateRefreshToken {
		// todo: save valid refresh to database
		refreshToken, err := a.jwt.NewRefreshToken(t.RefreshToken, time.Now().Add(refreshTokenExpiry))
		if err != nil {
			return nil, errors.New("could not sign refresh jwt")
		}
		refreshTokenCookie := http.Cookie{
			Name:     "jwt-refresh",
			Path:     "/api/auth",
			Value:    refreshToken,
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		}
		cookies = append(cookies, refreshTokenCookie.String())
	}

	headers["Set-Cookie"] = cookies

	return headers, nil
}

func (a AuthApiService) AuthLogoutGet(ctx context.Context) (openapi.ImplResponse, error) {
	accessTokenCookie := http.Cookie{
		Name:     "jwt",
		Path:     "/api",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	// todo: revoke refresh token => delete from database
	refreshTokenCookie := http.Cookie{
		Name:     "jwt-refresh",
		Path:     "/api/auth/refresh",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	headers := make(map[string][]string)
	headers["Set-Cookie"] = []string{accessTokenCookie.String(), refreshTokenCookie.String()}

	return openapi.ResponseWithHeaders(http.StatusOK, headers, nil), nil
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

	headers, err := a.jwtTokenHeaders(*t, true)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

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
	if !hasValidRefreshToken(ctx) {
		return openapi.Response(http.StatusUnauthorized, nil), ErrNotAuthenticated
	}

	spotifyToken := &oauth2.Token{
		TokenType:    "Bearer",
		RefreshToken: oauthRefreshTokenFromContext(ctx),
	}
	c := spotify.New(auth.Client(ctx, spotifyToken))
	newToken, err := c.Token()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("refreshing oauth2 access token failed")
	}

	headers, err := a.jwtTokenHeaders(*newToken, false)

	return openapi.ResponseWithHeaders(http.StatusOK, headers, nil), nil
}

func newAuthApiService(clientId string, secret []byte) AuthApiService {
	return AuthApiService{
		clientId: clientId,
		jwt:      jwt2.New(secret),
	}
}

var _ openapi.AuthApiServicer = AuthApiService{}
