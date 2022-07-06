package db

import (
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

const TrackCollection = "tracks"

//go:embed migrations
var migrationFiles embed.FS

func createIndices(db *mongo.Database) error {
	driver, err := mongodb.WithInstance(db.Client(), &mongodb.Config{
		DatabaseName: db.Name(),
	})
	if err != nil {
		return err
	}

	source, err := httpfs.New(http.FS(migrationFiles), "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("httpfs", source, db.Name(), driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
