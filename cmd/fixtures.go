package cmd

import (
	"github.com/imba28/spolyr/pkg/db"
	"github.com/spf13/cobra"
	"log"
)

var trackFixtures = []db.Track{
	{
		SpotifyID: "1",
		Artist:    "Eminem",
		AlbumName: "Curtain Call",
		ImageURL:  "https://i.scdn.co/image/ab67616d0000b273a8a32ae2a279b9bf03445738",
		Name:      "Lose Yourself",
		Lyrics:    "His palms are sweaty, knees weak, arms are heavy. There's vomit on his sweater already, mom's spaghetti",
		Loaded:    true,
		Language:  "english",
	},
	{
		SpotifyID: "2",
		Artist:    "Smash Mouth",
		AlbumName: "Astro Lounge",
		ImageURL:  "https://i.scdn.co/image/ab67616d0000b2739b997f037ce314dc4d39625f",
		Name:      "All Star",
		Lyrics:    "Somebody once told me the world is gonna roll me. I ain't the sharpest tool in the shed",
		Loaded:    true,
		Language:  "english",
	},
	{
		SpotifyID:  "3",
		Artist:     "Rammstein",
		AlbumName:  "REISE, REISE",
		ImageURL:   "https://i.scdn.co/image/ab67616d0000b273ad873056c55855af6012851c",
		PreviewURL: "",
		Name:       "AMERIKA",
		Lyrics:     "Wir bilden einen lieben Reigen. Die Freiheit spielt auf allen Geigen.",
		Loaded:     true,
		Language:   "german",
	},
	{
		SpotifyID:  "4",
		Artist:     "Eminem, Nate Dogg",
		AlbumName:  "The Eminem Show",
		ImageURL:   "https://i.scdn.co/image/ab67616d0000b2736ca5c90113b30c3c43ffb8f4",
		PreviewURL: "https://p.scdn.co/mp3-preview/5b3f004468ea54351d1a656577faa94a6a193b61?cid=fa338f8a902f43dfb275be2eb27d96e5",
		Name:       "'Till I Collapse",
		Loaded:     false,
	},
}

func NewFixturesCommand() *cobra.Command {
	config := &config{}

	c := &cobra.Command{
		Use:    "fixtures",
		Hidden: true,
		Run:    fixtures(config),
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

func fixtures(c *config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		dbConn, err := db.New(c.databaseUsername, c.databasePassword, "spolyr", c.databaseHost, 3)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Creating fixtures...")

		for i := range trackFixtures {
			if err := dbConn.Tracks.Save(&trackFixtures[i]); err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("Created %d fixtures...", len(trackFixtures))
	}
}
