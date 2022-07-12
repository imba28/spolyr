package cmd

import (
	"fmt"
	"github.com/imba28/spolyr/pkg/api"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/language"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

func NewWebCommand() *cobra.Command {
	config := &config{}

	c := &cobra.Command{
		Use: "web",
		Run: web(config),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := initConfig(cmd)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	initFlags(c, config)

	return c
}

func web(c *config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		env := api.Prod
		if c.debug {
			env = api.Dev
		}

		var languageDetector language.Detector
		if len(c.supportedLanguages) > 0 {
			ld, err := language.WithLanguages(c.supportedLanguages)
			if err != nil {
				log.Fatal(err)
			}
			languageDetector = ld
		} else {
			languageDetector = language.New()
		}

		dbConn, err := db.New(
			c.databaseUsername,
			c.databasePassword,
			"spolyr",
			c.databaseHost,
			3)
		if err != nil {
			log.Fatal(err)
		}

		s := api.NewServer(
			api.WithDatabase(dbConn),
			api.WithSecret([]byte(c.secret)),
			api.WithLanguageDetector(languageDetector),
			api.WithGeniusAPI(c.geniusAPIToken),
			api.WithOAuth(c.spotifyOAuthClientId),
			api.WithEnv(env),
			api.WithReverseProxy(c.protocol, c.domain, c.httpPublicPort))

		srv := &http.Server{
			Handler:      s,
			Addr:         fmt.Sprintf(":%d", c.httpPort),
			WriteTimeout: 1 * time.Minute,
			ReadTimeout:  10 * time.Second,
		}

		log.Printf("Starting web server http://127.0.0.1:%d", c.httpPort)

		log.Fatal(srv.ListenAndServe())
	}
}
