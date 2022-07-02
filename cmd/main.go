package main

import (
	"fmt"
	"github.com/imba28/spolyr/pkg/spolyr"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	databaseUsername := getEnv("DATABASE_USER", "root")
	databasePassword := getEnv("DATABASE_PASSWORD", "example")
	databaseHost := getEnv("DATABASE_HOST", "127.0.0.1")
	httpPort := getEnv("HTTP_PORT", "8080")
	geniusAPIToken := mustGetEnv("GENIUS_API_TOKEN")
	oauthClientId := mustGetEnv("SPOTIFY_ID")
	secret := []byte(mustGetEnv("SESSION_KEY"))

	s, err := spolyr.New(databaseHost, databaseUsername, databasePassword, geniusAPIToken, oauthClientId, secret)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler:      s,
		Addr:         fmt.Sprintf(":%s", httpPort),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
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
