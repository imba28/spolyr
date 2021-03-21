package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestAuthRequired(t *testing.T) {
	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	r := setUp()
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("userEmail", "test@test.com")
		session.Save()
	})
	r.Use(AuthRequired)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestErrorHandle(t *testing.T) {
	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	r := gin.Default()
	r.Use(ErrorHandle)
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusTeapot)
		c.Error(errors.New("foobar"))
	})
	r.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusTeapot, rr.Code, "does not override an already set response status")
	body, _ := io.ReadAll(rr.Body)
	assert.Contains(t, string(body), strconv.Itoa(http.StatusTeapot), "should contain response status somewhere in the template")
}

func TestErrorHandle_set_an_error_code(t *testing.T) {
	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	r := gin.Default()
	r.Use(ErrorHandle)
	r.GET("/", func(c *gin.Context) {
		c.Error(errors.New("foobar"))
	})
	r.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "should fallback to 500 status if no status was set prior")
}

func TestErrorHandle_does_nothing_if_no_error_occured(t *testing.T) {
	rr := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	r := gin.Default()
	r.Use(ErrorHandle)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusCreated, "Hello world")
	})
	r.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusCreated, rr.Code)
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, string(body), "Hello world")
}
