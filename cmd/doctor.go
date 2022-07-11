package cmd

import (
	"fmt"
	"github.com/dmolesUC/go-spinner"
	"github.com/imba28/spolyr/pkg/db"
	"github.com/imba28/spolyr/pkg/language"
	"github.com/spf13/cobra"
	"log"
)

func NewDoctorCommand() *cobra.Command {
	config := &config{}

	c := &cobra.Command{
		Use: "doctor",
		Run: doctor(config),
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

func doctor(c *config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		dbConn, err := db.New(c.databaseUsername, c.databasePassword, "spolyr", c.databaseHost, 3)
		if err != nil {
			log.Fatal(err)
		}

		var d language.Detector
		if len(c.supportedLanguages) > 0 {
			ld, err := language.WithLanguages(c.supportedLanguages)
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
	}
}
