package main

import (
	"github.com/imba28/spolyr/pkg/spolyr"
	"log"
	"os"
)

func main() {
	databaseUsername := getEnv("DATABASE_USER", "root")
	databasePassword := getEnv("DATABASE_PASSWORD", "example")
	databaseHost := getEnv("DATABASE_HOST", "127.0.0.1")
	geniusAPIToken := mustGetEnv("GENIUS_API_TOKEN")

	s, err := spolyr.New(databaseHost, databaseUsername, databasePassword, geniusAPIToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Run(":8080"))
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalf("Please specify the environment variable %s", key)
	}
	return value
}
