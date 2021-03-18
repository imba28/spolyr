package main

import (
	"github.com/imba28/spolyr/internal/db"
	"github.com/imba28/spolyr/pkg/spolyr"
	"log"
)

func main() {
	databaseUsername := db.GetEnv("DATABASE_USER", "root")
	databasePassword := db.GetEnv("DATABASE_PASSWORD", "example")
	databaseHost := db.GetEnv("DATABASE_HOST", "127.0.0.1")

	s, err := spolyr.New(databaseHost, databaseUsername, databasePassword)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Run(":8080"))
}
