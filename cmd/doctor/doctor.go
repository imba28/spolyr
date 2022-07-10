package main

import (
	"fmt"
	"github.com/dmolesUC/go-spinner"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/language"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "spolyr-doctor",
	}

	fixLanguageCmd := &cobra.Command{
		Use: "fix-track-languages",
		Run: func(cmd *cobra.Command, args []string) {
			databaseUsername := getEnv("DATABASE_USER", "root")
			databasePassword := getEnv("DATABASE_PASSWORD", "example")
			databaseHost := getEnv("DATABASE_HOST", "127.0.0.1")
			dbConn, err := db.New(databaseUsername, databasePassword, "spolyr", databaseHost, 3)
			if err != nil {
				log.Fatal(err)
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

			s := spinner.StartNew("Setting language of tracks...")
			defer func() {
				s.Stop()
				fmt.Println("Setting language of tracks: Completed")
			}()

			p := 1
			i := 1
			success := 0
			for {
				tracks, n, err := dbConn.Tracks.AllTracks(p, 25)
				if err != nil {
					log.Fatal(err)
				}
				for _, t := range tracks {
					if t.Loaded && t.Language == "" {
						lang, err := d.Detect(t.Lyrics)
						if err == nil {
							t.Language = lang
							err = dbConn.Tracks.Save(t)
							if err != nil {
								log.Fatal(err)
							}
							success++
						}
					}
					s.Title = fmt.Sprintf("Setting language of tracks (%d/%d)...", i, n)
					i++
				}

				if len(tracks) < 25 {
					break
				}
				p++
			}
		},
	}

	rootCmd.AddCommand(fixLanguageCmd)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
