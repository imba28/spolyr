package main

import (
	"github.com/imba28/spolyr/cmd"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "spolyr",
	}

	rootCmd.AddCommand(cmd.NewDoctorCommand())
	rootCmd.AddCommand(cmd.NewWebCommand())
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
