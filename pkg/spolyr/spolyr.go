package spolyr

import (
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/api"
	"github.com/imba28/spolyr/internal/db"
)

func New(dbHost, dbUser, dbPassword, geniusAPIToken string, sessionKey []byte) (*gin.Engine, error) {
	dbConn, err := db.New(dbUser, dbPassword, "spolyr", dbHost)
	if err != nil {
		return nil, err
	}

	return api.New(dbConn, geniusAPIToken, sessionKey), nil
}
