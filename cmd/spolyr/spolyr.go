package main

import (
	"fmt"
	"github.com/imba28/spolyr/pkg/api"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/language"
	"log"
	"net/http"
	"os"
	"strings"
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

	env := api.Prod
	if os.Getenv("DEBUG") != "" {
		env = api.Dev
	}

	var d language.Detector
	if os.Getenv("SUPPORTED_LANGUAGES") != "" {
		languages := strings.Split(os.Getenv("SUPPORTED_LANGUAGES"), ",")
		ld, err := language.WithLanguages(languages)
		if err != nil {
			log.Fatal(err)
		}
		d = ld
	} else {
		d = language.New()
	}

	dbConn, err := db.New(databaseUsername, databasePassword, "spolyr", databaseHost, 3)
	if err != nil {
		log.Fatal(err)
	}

	s := api.NewServer(
		api.WithDatabase(dbConn),
		api.WithSecret(secret),
		api.WithLanguageDetector(d),
		api.WithGeniusAPI(geniusAPIToken),
		api.WithOAuth(oauthClientId),
		api.WithEnv(env))

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
