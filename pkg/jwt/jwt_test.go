package jwt

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestJWT_ValidateAccessToken(t *testing.T) {
	t.Run("Created access tokens should be valid", func(t *testing.T) {
		secret := []byte("something secret")

		oauthToken := oauth2.Token{
			AccessToken:  "accessToken",
			TokenType:    "Bearer",
			RefreshToken: "refreshToken",
			Expiry:       time.Now(),
		}
		j := New(secret)
		token, err := j.NewAccessToken(oauthToken, time.Now().Add(5*time.Minute))

		assert.Nil(t, err)

		// create another token validator using the same secret
		j2 := New(secret)
		claims, valid := j2.ValidateAccessToken(token)

		// we should be able to decode the token
		assert.True(t, valid)
		assert.Equal(t, oauthToken.AccessToken, claims.Token.AccessToken)
		assert.Equal(t, oauthToken.RefreshToken, claims.Token.RefreshToken)
		assert.Equal(t, oauthToken.TokenType, claims.Token.TokenType)
		assert.True(t, oauthToken.Expiry.Equal(claims.Token.Expiry))
	})

	t.Run("Token should be invalid if secret changes", func(t *testing.T) {
		secret := []byte("something secret")

		oauthToken := oauth2.Token{}
		j := New(secret)
		token, err := j.NewAccessToken(oauthToken, time.Now().Add(5*time.Minute))

		assert.Nil(t, err)

		// create another token validator using a different secret
		j2 := New([]byte("secret changed after restarting server"))
		_, valid := j2.ValidateAccessToken(token)

		assert.False(t, valid)
	})
}

func TestJWT_ValidateRefreshToken(t *testing.T) {
	t.Run("Created refresh tokens should be valid", func(t *testing.T) {
		secret := []byte("something secret")
		spotifyOauthRefreshToken := "refreshToken"

		j := New(secret)
		token, err := j.NewRefreshToken(spotifyOauthRefreshToken, time.Now().Add(5*time.Minute))

		assert.Nil(t, err)

		// create another token validator using the same secret
		j2 := New(secret)
		claims, valid := j2.ValidateRefreshToken(token)

		// we should be able to decode the refresh token
		assert.True(t, valid)
		assert.Equal(t, claims.RefreshToken, claims.RefreshToken)
	})

	t.Run("Refresh token should be invalid if secret changes", func(t *testing.T) {
		secret := []byte("something secret")
		spotifyOauthRefreshToken := "refreshToken"

		j := New(secret)
		token, err := j.NewRefreshToken(spotifyOauthRefreshToken, time.Now().Add(5*time.Minute))

		assert.Nil(t, err)

		// create another token validator using a different secret
		j2 := New([]byte("a different secret"))
		_, valid := j2.ValidateRefreshToken(token)

		assert.False(t, valid)
	})
}
