package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"net/http"
	"os"
)

const csrfSessionKey = "spolyrCSRF"

func redirectUrl() string {
	var protocol = getEnv("PROTOCOL", "http")
	var domain = getEnv("DOMAIN", "localhost")

	var httpPort = getEnv("HTTP_PUBLIC_PORT", getEnv("HTTP_PORT", "8080"))
	if httpPort != "80" && httpPort != "443" {
		httpPort = ":" + httpPort
	} else {
		httpPort = ""
	}

	return fmt.Sprintf("%s://%s%s/callback", protocol, domain, httpPort)
}

var auth = spotify.NewAuthenticator(redirectUrl(), spotify.ScopeUserLibraryRead, spotify.ScopeUserReadEmail)

func SpotifyAuthCallbackHandler(c *gin.Context) {
	csrfToken, ok := sessions.Default(c).Get(csrfSessionKey).(string)
	if !ok {
		c.Error(errors.New("could not decode csrf token"))
		return
	}

	tok, err := auth.Token(csrfToken, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func LoginHandler(c *gin.Context) {
	token, err := csrfToken()
	if err != nil {
		c.Error(err)
		return
	}
	session := sessions.Default(c)
	session.Set(csrfSessionKey, token)
	err = session.Save()
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, auth.AuthURL(token))
}

func csrfToken() (string, error) {
	sha := sha256.New()
	b := make([]byte, 256)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	sha.Write(b)

	return base64.URLEncoding.EncodeToString(sha.Sum(nil)), nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
