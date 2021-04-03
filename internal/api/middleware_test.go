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

func TestAuthRequired(t *testing.T) {
	t.Run("allows access if request is authenticated", func(t *testing.T) {
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
	})

	t.Run("prevents access to unauthorized requests", func(t *testing.T) {
		r := gin.Default()
		r.Use(sessions.Sessions("session", memstore.NewStore([]byte("secret"))))
		r.Use(AuthRequired)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusUnauthorized)
	})
}

func TestUserProviderMiddleware(t *testing.T) {
	t.Run("saves user and avatar in context", func(t *testing.T) {
		var displayName, avatar string

		expectedDisplayName, expectedAvatar := "test@test.com", "http://foobar.com/avatar.png"

		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		r := gin.Default()
		r.Use(sessions.Sessions("session", memstore.NewStore([]byte("secret"))))
		r.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set("userEmail", expectedDisplayName)
			session.Set("userAvatar", expectedAvatar)

			session.Save()
		})
		r.Use(UserProviderMiddleware)
		r.Use(func(c *gin.Context) {
			view := viewFromContext(c)
			displayName = view["User"].(string)
			avatar = view["UserAvatar"].(string)
		})
		r.ServeHTTP(rr, request)

		assert.Equal(t, expectedDisplayName, displayName)
		assert.Equal(t, expectedAvatar, avatar)
	})
}

func TestErrorHandle(t *testing.T) {
	t.Run("renders custom error template", func(t *testing.T) {
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
	})

	t.Run("falls back to internal error if no error was set prior", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		r := gin.Default()
		r.Use(ErrorHandle)
		r.GET("/", func(c *gin.Context) {
			c.Error(errors.New("foobar"))
		})
		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "should fallback to 500 status if no status was set prior")
	})

	t.Run("Does nothing if no error occurred", func(t *testing.T) {
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
	})
}
