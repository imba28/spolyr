package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"time"
)

type Claims struct {
	Token oauth2.Token `json:"oauth_token"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	RefreshToken string `json:"oauth_refresh_token"`

	jwt.RegisteredClaims
}

type JWT struct {
	secret []byte
}

func (j JWT) validateToken(t *jwt.Token) (interface{}, error) {
	_, ok := t.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}
	return j.secret, nil
}

func (j JWT) ValidateAccessToken(accessToken string) (*Claims, bool) {
	t, err := jwt.ParseWithClaims(accessToken, &Claims{}, j.validateToken)
	if err != nil || !t.Valid {
		return nil, false
	}
	if claims, ok := t.Claims.(*Claims); ok {
		return claims, true
	}
	return nil, false
}

func (j JWT) ValidateRefreshToken(refreshToken string) (*RefreshClaims, bool) {
	t, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, j.validateToken)
	if err != nil || !t.Valid {
		return nil, false
	}
	if claims, ok := t.Claims.(*RefreshClaims); ok {
		return claims, true
	}
	return nil, false
}

func (j JWT) NewAccessToken(oauthToken oauth2.Token, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Token: oauthToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
		},
	})
	return token.SignedString(j.secret)
}

func (j JWT) NewRefreshToken(oauthRefreshToken string, expiresAt time.Time) (string, error) {
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshClaims{
		RefreshToken: oauthRefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})
	return refreshJwt.SignedString(j.secret)
}

func New(secret []byte) JWT {
	return JWT{
		secret: secret,
	}
}
