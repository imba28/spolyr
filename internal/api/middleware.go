package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	template2 "github.com/imba28/spolyr/internal/template"
	"net/http"
)

const userEmailKey = "UserEmail"
const spotifyTokenKey = "SpotifyToken"

func UserProviderMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	c.Set(userEmailKey, session.Get("userEmail"))
	c.Set(spotifyTokenKey, session.Get("token"))

	c.Next()
}

func ErrorHandle(c *gin.Context) {
	c.Next()
	err := c.Errors.Last()
	if err == nil {
		return
	}

	p := gin.H{
		"Status":  http.StatusInternalServerError,
		"Message": "Whoops! Sorry, an error occurred",
	}
	statusCode := http.StatusInternalServerError

	if errors.Is(err.Err, ErrNotFound) {
		p["Message"] = "Page not found"
		p["Status"] = http.StatusNotFound
		statusCode = http.StatusNotFound
	}

	_ = template2.ErrorPage(c.Writer, p, statusCode)
	c.Abort()
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)

	if email := session.Get("userEmail"); email == nil {
		viewData := gin.H{
			"Status":  http.StatusUnauthorized,
			"Message": "You must be logged in to perform this action.",
		}
		_ = template2.ErrorPage(c.Writer, viewData, http.StatusUnauthorized)

		c.Abort()
		return
	}

	c.Next()
}
