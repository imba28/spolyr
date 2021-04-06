package main

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/imba28/spolyr/pkg/spolyr"
	"log"
	"os"
)

func main() {
	databaseUsername := getEnv("DATABASE_USER", "root")
	databasePassword := getEnv("DATABASE_PASSWORD", "example")
	databaseHost := getEnv("DATABASE_HOST", "127.0.0.1")
	httpPort := getEnv("HTTP_PORT", "8080")
	sessionKey := getSessionKey()
	geniusAPIToken := mustGetEnv("GENIUS_API_TOKEN")

	s, err := spolyr.New(databaseHost, databaseUsername, databasePassword, geniusAPIToken, sessionKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Run(fmt.Sprintf(":%s", httpPort)))
}

func getSessionKey() []byte {
	key := getEnv("SESSION_KEY", "")

	if len(key) == 0 {
		randomKey := securecookie.GenerateRandomKey(32)
		if randomKey == nil {
			panic("Could not generate random session key!")
		}
		return randomKey
	}

	return []byte(key)
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
