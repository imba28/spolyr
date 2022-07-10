package spolyr

import (
	"github.com/imba28/spolyr/pkg/api"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/language"
	"net/http"
)

func New(dbHost, dbUser, dbPassword, geniusAPIToken, oauthClientId string, secret []byte, d language.Detector) (http.Handler, error) {
	dbConn, err := db.New(dbUser, dbPassword, "spolyr", dbHost, 3)
	if err != nil {
		return nil, err
	}

	return api.New(dbConn, geniusAPIToken, oauthClientId, secret, d), nil
}
