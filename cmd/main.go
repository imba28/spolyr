package main

import (
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/pkg/spolyr"
	"log"
	"os"
)

func main() {
	databaseUsername := db.GetEnv("DATABASE_USER", "root")
	databasePassword := db.GetEnv("DATABASE_PASSWORD", "example")
	databaseHost := db.GetEnv("DATABASE_HOST", "127.0.0.1")
	geniusAPIToken := os.Getenv("GENIUS_API_TOKEN")

	if len(geniusAPIToken) == 0 {
		log.Fatal("Please specify your genius api token as env var GENIUS_API_TOKEN")
	}

	s, err := spolyr.New(databaseHost, databaseUsername, databasePassword, geniusAPIToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Run(":8080"))
}
