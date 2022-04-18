package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	template2 "github.com/imba28/spolyr/internal/template"
	"net/http"
)

const userDisplayNameKey = "UserDisplayName"
const userAvatarKey = "UserAvatar"
const spotifyTokenKey = "SpotifyToken"

func mergeView(v1, v2 gin.H) gin.H {
	for k, v := range v1 {
		v2[k] = v
	}
	return v2
}

func viewFromContext(c *gin.Context) gin.H {
	v, ok := c.Get("view")
	var view gin.H
	if !ok {
		view = gin.H{}
	} else {
		view = v.(gin.H)
	}
	return view
}

func UserProviderMiddleware(c *gin.Context) {
	view := viewFromContext(c)

	session := sessions.Default(c)
	c.Set(userDisplayNameKey, session.Get("userEmail"))
	c.Set(userAvatarKey, session.Get("userAvatar"))
	c.Set(spotifyTokenKey, session.Get("token"))

	view["User"] = session.Get("userEmail")
	view["UserAvatar"] = session.Get("userAvatar")
	c.Set("view", view)

	c.Next()
}

func ErrorHandle(c *gin.Context) {
	c.Next()
	err := c.Errors.Last()
	if err == nil {
		return
	}

	p := viewFromContext(c)
	p["Status"] = http.StatusInternalServerError
	p["Message"] = "Whoops! Sorry, an error occurred"
	statusFallback := http.StatusInternalServerError

	if errors.Is(err.Err, ErrNotFound) {
		p["Message"] = "Page not found"
		statusFallback = http.StatusNotFound
	}

	status := c.Writer.Status()
	// set error status if no status was explicitly set prior
	if status == http.StatusOK {
		status = statusFallback
	}
	p["Status"] = status
	_ = template2.ErrorPage(c.Writer, p, status)
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
