package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Alert on currency conversion rate",
	Long:  "Implement an alerting service which will consume a file of currency conversion rates and\nproduce alerts for a number of situations.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize()
}

// handleCommandError will print an error message regarding a command set up before killing the script
func handleCommandError(err error) {
	if err != nil {
		log.Fatalf("could not initialise command %s", err)
	}
}
