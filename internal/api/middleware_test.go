package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRequired_returns_401_if_no_auth_is_provided(t *testing.T) {
	store := memstore.NewStore([]byte("secret"))

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))
	r.Use(AuthRequired)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}
