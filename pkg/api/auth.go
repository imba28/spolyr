package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	jwt2 "github.com/imba28/spolyr/pkg/jwt"
	"github.com/imba28/spolyr/pkg/openapi"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

type authContextKey int

const (
	spotifyRefreshTokenKey authContextKey = iota
	jwtRefreshKey
	jwtAccessKey
	spotifyOauthClientKey

	accessTokenExpiry  = time.Minute * 10
	refreshTokenExpiry = time.Hour * 24
)

var (
	permissions = []string{
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopePlaylistReadPrivate}
	scope = strings.Join(permissions, " ")
	auth  = spotifyauth.New(spotifyauth.WithScopes(scope))

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

func oauthClientFromContext(ctx context.Context) *spotify.Client {
	if c, ok := ctx.Value(spotifyOauthClientKey).(*spotify.Client); ok {
		return c
	}
	return nil
}

func oauthRefreshTokenFromContext(ctx context.Context) string {
	if t, ok := ctx.Value(spotifyRefreshTokenKey).(string); ok {
		return t
	}
	return ""
}

func (a AuthApiService) redirectUrl() string {
	port := a.publicHttpPort
	publicPort := ""
	if port != 80 && port != 443 && port != 0 {
		publicPort = fmt.Sprintf(":%d", port)
	}

	return fmt.Sprintf("%s://%s%s/auth/callback", a.publicHttpProtocol, a.publicHostname, publicPort)
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
						ctx = context.WithValue(ctx, jwtAccessKey, c.Value)
						ctx = context.WithValue(ctx, spotifyOauthClientKey, spotify.New(auth.Client(ctx, &claims.Token)))
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

	publicHttpProtocol string
	publicHostname     string
	publicHttpPort     int
}

func (a AuthApiService) cookie(name, value, path string) http.Cookie {
	sameSite := http.SameSiteLaxMode
	if a.publicHttpProtocol == "https" {
		sameSite = http.SameSiteNoneMode
	}

	return http.Cookie{
		Name:     name,
		Path:     path,
		Value:    value,
		Secure:   a.publicHttpProtocol == "https",
		HttpOnly: true,
		SameSite: sameSite,
	}
}

func (a AuthApiService) jwtTokenHeaders(t oauth2.Token, generateRefreshToken bool) (map[string][]string, error) {
	headers := make(map[string][]string)
	var cookies []string

	accessToken, err := a.jwt.NewAccessToken(t, time.Now().Add(accessTokenExpiry))
	if err != nil {
		return nil, errors.New("could not sign access jwt")
	}
	accessTokenCookie := a.cookie("jwt", accessToken, "/api")
	cookies = append(cookies, accessTokenCookie.String())

	if generateRefreshToken {
		// todo: save valid refresh to database
		refreshToken, err := a.jwt.NewRefreshToken(t.RefreshToken, time.Now().Add(refreshTokenExpiry))
		if err != nil {
			return nil, errors.New("could not sign refresh jwt")
		}
		refreshTokenCookie := a.cookie("jwt-refresh", refreshToken, "/api/auth")
		refreshTokenCookie.Expires = time.Now().Add(refreshTokenExpiry)
		cookies = append(cookies, refreshTokenCookie.String())
	}

	headers["Set-Cookie"] = cookies

	return headers, nil
}

func (a AuthApiService) AuthLogoutGet(ctx context.Context) (openapi.ImplResponse, error) {
	accessTokenCookie := a.cookie("jwt", "1", "/api")
	accessTokenCookie.Expires = time.Unix(0, 0)

	// todo: revoke refresh token => delete from database
	refreshTokenCookie := a.cookie("jwt-refresh", "1", "/api/auth")
	refreshTokenCookie.Expires = time.Unix(0, 0)

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
		RedirectUrl: a.redirectUrl(),
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

func newAuthApiService(clientId, clientSecret string, secret []byte, publicProtocol, publicHostname string, publicPort int) AuthApiService {
	a := AuthApiService{
		clientId:           clientId,
		jwt:                jwt2.New(secret),
		publicHttpPort:     publicPort,
		publicHostname:     publicHostname,
		publicHttpProtocol: publicProtocol,
	}

	// todo: this is extremely ugly, but public hostname/port are only available after setting up the service
	// alternatively, we could directly access the environment variables (but that breaks the top-down approach of the configuration flow)
	spotifyauth.WithRedirectURL(a.redirectUrl())(auth)
	spotifyauth.WithClientID(clientId)(auth)
	spotifyauth.WithClientSecret(clientSecret)(auth)

	return a
}

var _ openapi.AuthApiServicer = AuthApiService{}
