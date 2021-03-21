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

func TestGetEnv_returns_value_of_env_var(t *testing.T) {
	key := "foo"
	value := "bar"
	_ = os.Setenv(key, value)
	defer os.Unsetenv(key)

	assert.Equal(t, getEnv(key, ""), value)
}

func TestGetEnv_returns_fallback_value_if_env_var_is_missing(t *testing.T) {
	assert.Equal(t, getEnv("foo", "fallback"), "fallback")
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
