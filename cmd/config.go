package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

type config struct {
	databaseUsername     string
	databasePassword     string
	databaseHost         string
	httpPort             int
	geniusAPIToken       string
	spotifyOAuthClientId string
	secret               string

	protocol       string // getEnv("PROTOCOL", "http")
	domain         string // getEnv("DOMAIN", "localhost")
	httpPublicPort int    // = getEnv("HTTP_PUBLIC_PORT", getEnv("HTTP_PORT", "8080"))

	supportedLanguages []string

	debug bool
}

func initConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName("spolyr")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.AutomaticEnv()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = v.BindEnv(f.Name, strings.ToUpper(f.Name))
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}

func initFlags(cmd *cobra.Command, c *config) {
	cmd.Flags().StringVarP(&c.databaseUsername, "database_user", "", "root", "Username of mongodb user")
	cmd.Flags().StringVarP(&c.databasePassword, "database_password", "", "example", "Password of mongodb user")
	cmd.Flags().StringVarP(&c.databaseHost, "database_host", "", "127.0.0.1", "Host of mongodb instance")
	cmd.Flags().IntVarP(&c.httpPort, "http_port", "", 8080, "Port Spolyr should bind to")
	cmd.Flags().StringVarP(&c.spotifyOAuthClientId, "spotify_id", "", "", "Spotify OAuth2 client id")
	cmd.Flags().StringVarP(&c.geniusAPIToken, "genius_api_token", "", "", "Genius.com api token")
	cmd.Flags().StringVarP(&c.secret, "session_key", "", "", "Secret value used for validating session data")
	cmd.Flags().StringSliceVarP(&c.supportedLanguages, "supported_languages", "", []string{}, "List of languages used for language specific database queries")
	cmd.Flags().BoolVarP(&c.debug, "debug", "d", false, "Start api in debug mode. Enables cors for local development.")

	cmd.Flags().StringVarP(&c.protocol, "protocol", "", "http", "Public http protocol. Pick https if Spolyr resides behind a reverse proxy using TLS")
	cmd.Flags().StringVarP(&c.domain, "domain", "", "localhost", "Public hostname")
	cmd.Flags().IntVarP(&c.httpPublicPort, "http_public_port", "", 8080, "Public http port")
}
