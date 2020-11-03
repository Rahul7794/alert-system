package cmd

import (
	"alert-system/alertprocessor"
	"alert-system/io"
	"alert-system/log"
	"alert-system/model"
	"alert-system/version"

	"github.com/spf13/cobra"
)

// Command line Args
var runCmd = &cobra.Command{
	Use:   "alert",                             // SubCommand
	Short: "Alert on currency conversion rate", // Short description of the SubCommand
	Long:  "Implement an alerting service which will consume a file of currency conversion rates and \n produce alerts for a number of situations.",
	RunE:  alertCmd, // alertCmd processing the currency rates and producing alerts.
}

// inputPath for file of currency conversion rates.
var inputPath string

// outputPath for file of alerts produced in number of situations.
var outputPath string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&inputPath, "i", "i", "", "input path for the json file")
	err := runCmd.MarkFlagRequired("i")
	handleCommandError(err)
	runCmd.Flags().StringVarP(&outputPath, "o", "o", "", "output path for the json file")
	err = runCmd.MarkFlagRequired("o")
	handleCommandError(err)
}

func alertCmd(_ *cobra.Command, _ []string) error {
	log.Infof("starting the alerting system with version: %s %s \n on %s => %s", version.Version, version.BuildDate,
		version.OsArch, version.GoVersion)
	// Creates an alert channel which receives alerts
	alertChannel := make(chan model.AlertFormat)
	// Creates an error channel which receives error.
	errorChannel := make(chan error)
	// done resembles completion of process
	done := make(chan bool, 1)

	// Create a reader object for the path provided.
	reader, err := io.Read(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create a writer object for the path provided.
	writer, err := io.Write(outputPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	// Create a alertprocessor object
	processor := alertprocessor.NewAlertProcessor(reader, writer, alertChannel, errorChannel, done)

	// Go routines to process the alerts
	go processor.ProcessAlerts()
	// Go routines to send alerts
	go processor.SendAlert()

	// Listen to error and done channel
	// If error channel receive error, return it
	// if done channel receive signal, end the program
	for {
		select {
		case err := <-errorChannel: // Listen errorChannel
			return err
		case <-done: // Listen done Channel
			log.Info("complete !!!")
			return nil
		}
	}
}
