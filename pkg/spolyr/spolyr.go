package spolyr

import (
	"github.com/imba28/spolyr/internal/api"
	"github.com/imba28/spolyr/internal/db"
	"net/http"
)

func New(dbHost, dbUser, dbPassword, geniusAPIToken, oauthClientId string, secret []byte) (http.Handler, error) {
	dbConn, err := db.New(dbUser, dbPassword, "spolyr", dbHost, 3)
	if err != nil {
		return nil, err
	}

	return api.New(dbConn, geniusAPIToken, oauthClientId, secret), nil
}
