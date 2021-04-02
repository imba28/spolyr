package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Run("returns value of env var", func(t *testing.T) {
		key := "foo"
		value := "bar"
		_ = os.Setenv(key, value)
		defer os.Unsetenv(key)

		assert.Equal(t, getEnv(key, ""), value)
	})

	t.Run("returns fallback value if env var is missing", func(t *testing.T) {
		assert.Equal(t, getEnv("foo", "fallback"), "fallback")
	})
}

func TestLoginHandler(t *testing.T) {
	var token string
	var ok bool

	r := gin.Default()
	r.Use(sessions.Sessions("session", memstore.NewStore([]byte("secret"))))
	r.Use(LoginHandler)
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		token, ok = session.Get(csrfSessionKey).(string)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusTemporaryRedirect)
	assert.True(t, ok)
	assert.Greater(t, len(token), 10, "should save a csrf token to the session")
}

func TestLogoutHandler(t *testing.T) {
	var hasSessionKey bool
	var flashes []interface{}

	r := gin.Default()
	r.Use(sessions.Sessions("session", memstore.NewStore([]byte("secret"))))
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", "value")
		session.AddFlash("a flash")
		session.Save()
	})
	r.Use(LogoutHandler)
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		_, hasSessionKey = session.Get("key").(string)
		flashes = session.Flashes()
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusFound)
	assert.False(t, hasSessionKey, "should delete all session keys")
	assert.Nil(t, flashes, 0, "should delete all flash messages")
}

func TestRedirectUrl(t *testing.T) {
	t.Run("reads http protocol from environment variable", func(t *testing.T) {
		_ = os.Setenv("PROTOCOL", "https")
		defer os.Unsetenv("PROTOCOL")

		assert.Equal(t, redirectUrl(), "https://localhost:8080/callback")
	})

	t.Run("reads hostname from DOMAIN environment variable", func(t *testing.T) {
		_ = os.Setenv("DOMAIN", "spolyr.com")
		defer os.Unsetenv("DOMAIN")

		assert.Equal(t, redirectUrl(), "http://spolyr.com:8080/callback")
	})

	t.Run("defaults to http port if no custom port is set", func(t *testing.T) {
		_ = os.Setenv("HTTP_PORT", "1234")
		defer os.Unsetenv("HTTP_PORT")

		assert.Equal(t, redirectUrl(), "http://localhost:1234/callback")
	})

	t.Run("adds http port to url if a non default port is used", func(t *testing.T) {
		_ = os.Setenv("HTTP_PUBLIC_PORT", "4444")
		defer os.Unsetenv("HTTP_PUBLIC_PORT")

		assert.Equal(t, redirectUrl(), "http://localhost:4444/callback")
	})

	t.Run("does not append port to callback url if port is either 443 or 80", func(t *testing.T) {
		t.Run("port 443", func(t *testing.T) {
			_ = os.Setenv("HTTP_PUBLIC_PORT", "443")
			defer os.Unsetenv("HTTP_PUBLIC_PORT")

			assert.Equal(t, redirectUrl(), "http://localhost/callback")
		})

		t.Run("port 80", func(t *testing.T) {
			_ = os.Setenv("HTTP_PUBLIC_PORT", "80")
			defer os.Unsetenv("HTTP_PUBLIC_PORT")

			assert.Equal(t, redirectUrl(), "http://localhost/callback")
		})
	})
}
