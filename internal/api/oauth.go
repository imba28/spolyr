package api

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"net/http"
)

const state = "spolyrCSRF"

var redirectURI = getenv("HOSTNAME", "http://localhost:8080") + "/callback"
var auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserLibraryRead, spotify.ScopeUserReadEmail)

func (co Controller) SpotifyAuthCallbackHandler(c *gin.Context) {
	tok, err := auth.Token(state, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if st := c.Request.FormValue("state"); st != state {
		c.String(http.StatusNotFound, "Invalid csrf token")
		return
	}

	client := auth.NewClient(tok)
	user, err := client.CurrentUser()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	token, _ := json.Marshal(tok)
	session := sessions.Default(c)
	session.Set("token", string(token))
	session.Set("userEmail", user.Email)
	err = session.Save()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (co Controller) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (co Controller) LoginHandler(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, auth.AuthURL(state))
}
