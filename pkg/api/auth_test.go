package api

import (
	"context"
	jwt2 "github.com/imba28/spolyr/pkg/jwt"
	"github.com/imba28/spolyr/pkg/openapi/openapi"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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

			if v := ctx.Value(spotifyTokenKey); v != nil {
				t.Error("oauth access token should not be set, got", v)
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

			if v, ok := ctx.Value(spotifyTokenKey).(oauth2.Token); !ok {
				t.Errorf("oauth access token should be set, got %v", v)
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

			if v := ctx.Value(spotifyTokenKey); v != nil {
				t.Error("oauth access token should not be set, got", v)
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
	assert.Contains(t, res.Headers["Set-Cookie"][0], "jwt=", "should set the cookie `jwt` ")
	assert.Contains(t, res.Headers["Set-Cookie"][1], "jwt-refresh=", "should set the cookie `jwt-refresh` ")
}
