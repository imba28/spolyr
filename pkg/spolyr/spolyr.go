package spolyr

import (
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/api"
	"github.com/imba28/spolyr/internal/db"
)

func New(dbHost, dbUser, dbPassword string) (*gin.Engine, error) {
	dbConn, err := db.New(dbUser, dbPassword, "spolyr", dbHost)
	if err != nil {
		return nil, err
	}

	controller := api.NewController(dbConn)

	return api.New(controller), nil
}
