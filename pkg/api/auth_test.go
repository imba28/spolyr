package api

import (
	"context"
	jwt2 "github.com/imba28/spolyr/pkg/jwt"
	"github.com/imba28/spolyr/pkg/openapi"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsAuthenticated(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()

		r := isAuthenticated(ctx)
		if r == true {
			t.Errorf("An empty context should return false, got %v", r)
		}
	})

	t.Run("context containing token", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), jwtAccessKey, "some token")

		r := isAuthenticated(ctx)
		if r == false {
			t.Errorf("A context containing a token should return true, got %v", r)
		}
	})
}

func TestAuthenticationMiddleware(t *testing.T) {
	t.Run("no cookies", func(t *testing.T) {
		jwt := jwt2.New([]byte("secret"))

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if v := ctx.Value(spotifyOauthClientKey); v != nil {
				t.Error("oauth client should not be set, got", v)
			}
			if v := ctx.Value(spotifyRefreshTokenKey); v != nil {
				t.Error("oauth refresh token should not be set, got", v)
			}
			if v := ctx.Value(jwtRefreshKey); v != nil {
				t.Error("jwt refresh token should not be set, got", v)
			}
			if v := ctx.Value(jwtAccessKey); v != nil {
				t.Error("jwt access token should not be set, got", v)
			}
		})
		handlerToTest := AuthenticationMiddleware(jwt)(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})

	t.Run("cookies containing valid jwt", func(t *testing.T) {
		jwt := jwt2.New([]byte("secret"))
		expiry := time.Now().Add(time.Hour)
		accessToken, _ := jwt.NewAccessToken(oauth2.Token{
			AccessToken: "access",
			Expiry:      expiry,
		}, expiry)
		spotifyRefreshTokenValue := "refresh-token"
		refreshToken, _ := jwt.NewRefreshToken(spotifyRefreshTokenValue, expiry)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if v, ok := ctx.Value(spotifyOauthClientKey).(*spotify.Client); !ok {
				t.Errorf("oauth client should be set, got %v", v)
			}
			if v, ok := ctx.Value(spotifyRefreshTokenKey).(string); !ok || v != spotifyRefreshTokenValue {
				t.Errorf("oauth refresh token should be set to %v, got %v", spotifyRefreshTokenValue, v)
			}
			if v, ok := ctx.Value(jwtRefreshKey).(string); !ok || v != refreshToken {
				t.Errorf("jwt refresh token should be set to %v, got %v", refreshToken, v)
			}
			if v, ok := ctx.Value(jwtAccessKey).(string); !ok || v != accessToken {
				t.Errorf("jwt access token should be set to %v, got %v", accessToken, v)
			}
		})
		handlerToTest := AuthenticationMiddleware(jwt)(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)
		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: accessToken,
		})
		req.AddCookie(&http.Cookie{
			Name:  "jwt-refresh",
			Value: refreshToken,
		})
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})

	t.Run("cookies containing expired jwt", func(t *testing.T) {
		jwt := jwt2.New([]byte("secret"))
		expiry := time.Now().Add(-1 * time.Hour)
		accessToken, _ := jwt.NewAccessToken(oauth2.Token{
			AccessToken: "access",
			Expiry:      expiry,
		}, expiry)
		refreshToken, _ := jwt.NewRefreshToken("refresh-token", expiry)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if v := ctx.Value(spotifyOauthClientKey); v != nil {
				t.Error("oauth client should not be set, got", v)
			}
			if v := ctx.Value(spotifyRefreshTokenKey); v != nil {
				t.Error("oauth refresh token should not be set, got", v)
			}
			if v := ctx.Value(jwtRefreshKey); v != nil {
				t.Error("jwt refresh token should not be set, got", v)
			}
			if v := ctx.Value(jwtAccessKey); v != nil {
				t.Error("jwt access token should not be set, got", v)
			}
		})
		handlerToTest := AuthenticationMiddleware(jwt)(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)
		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: accessToken,
		})
		req.AddCookie(&http.Cookie{
			Name:  "jwt-refresh",
			Value: refreshToken,
		})
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})
}

func TestAuthApiService_AuthLoginPost_secure_cookies(t *testing.T) {
	t.Run("if protocol is http do not use secure cookies", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "=~/me$",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"email": "user@test.com",
				})
			},
		)

		httpmock.RegisterResponder("POST", "=~/api/token",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"access_token":  "access-token",
					"token_type":    "Bearer",
					"refresh_token": "refresh_token",
				})
			},
		)

		ctx := context.Background()
		loginRequest := openapi.AuthLoginPostRequest{Code: "oauth-code"}

		auth := AuthApiService{
			publicHttpProtocol: "http",
		}
		res, _ := auth.AuthLoginPost(ctx, loginRequest)

		cookies := parseCookies(res.Headers["Set-Cookie"])
		assert.False(t, cookies[0].Secure)
		assert.False(t, cookies[1].Secure)
	})

	t.Run("if protocol is https define cookies as secure", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "=~/me$",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"email": "user@test.com",
				})
			},
		)

		httpmock.RegisterResponder("POST", "=~/api/token",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"access_token":  "access-token",
					"token_type":    "Bearer",
					"refresh_token": "refresh_token",
				})
			},
		)

		ctx := context.Background()
		loginRequest := openapi.AuthLoginPostRequest{Code: "oauth-code"}

		auth := AuthApiService{
			publicHttpProtocol: "https",
		}
		res, _ := auth.AuthLoginPost(ctx, loginRequest)

		cookies := parseCookies(res.Headers["Set-Cookie"])
		assert.True(t, cookies[0].Secure)
		assert.True(t, cookies[1].Secure)
	})
}

func TestAuthApiService_AuthLoginPost(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "=~/me$",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"email": "user@test.com",
			})
		},
	)

	httpmock.RegisterResponder("POST", "=~/api/token",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"access_token":  "access-token",
				"token_type":    "Bearer",
				"refresh_token": "refresh_token",
			})
		},
	)

	ctx := context.Background()
	loginRequest := openapi.AuthLoginPostRequest{Code: "oauth-code"}

	auth := AuthApiService{}
	res, err := auth.AuthLoginPost(ctx, loginRequest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Len(t, res.Headers["Set-Cookie"], 2, "should set two cookies ")
	c := parseCookies(res.Headers["Set-Cookie"])

	assert.Equal(t, "jwt", c[0].Name, "should set the cookie `jwt`")
	assert.Equal(t, "/api", c[0].Path, "`jwt` should be valid for path /api")
	assert.Equal(t, "jwt-refresh", c[1].Name, "should set the cookie `jwt-refresh`")
	assert.Equal(t, "/api/auth", c[1].Path, "`jwt-refresh` should only be valid for path /api/auth")
}

func TestAuthApiService_AuthLogoutGet(t *testing.T) {
	auth := AuthApiService{}

	res, err := auth.AuthLogoutGet(context.Background())

	assert.Nil(t, err)
	assert.Len(t, res.Headers["Set-Cookie"], 2, "should set two cookies")
	c := parseCookies(res.Headers["Set-Cookie"])
	assert.True(t, c[0].Expires.Before(time.Now()), "jwt cookie should be expired")
	assert.True(t, c[1].Expires.Before(time.Now()), "jwt-refresh cookie should be expired")
}

func parseCookies(cookieHeaders []string) []*http.Cookie {
	header := http.Header{}
	for i := range cookieHeaders {
		header.Add("Set-Cookie", cookieHeaders[i])
	}
	req := http.Response{Header: header}
	return req.Cookies()
}
