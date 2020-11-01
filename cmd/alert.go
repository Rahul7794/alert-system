package cmd

import (
	"alert-system/alert_processor"
	"alert-system/io_stream"
	"alert-system/log"
	"alert-system/model"
	"alert-system/version"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "alert",
	Short: "Alert on currency conversion rate",
	Long:  "Implement an alerting service which will consume a file of currency conversion rates and\nproduce alerts for a number of situations.",
	RunE:  alertCmd,
}

var inputPath string
var outputPath string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&inputPath, "i", "i", "", "input path for the json file")
	err := runCmd.MarkFlagRequired("i")
	runCmd.Flags().StringVarP(&outputPath, "o", "o", "", "output path for the json file")
	err = runCmd.MarkFlagRequired("o")
	handleCommandError(err)
}

func alertCmd(_ *cobra.Command, _ []string) error {
	log.Infof("starting the alerting system with version: %s %s \n on %s => %s", version.Version, version.BuildDate,
		version.OsArch, version.GoVersion)
	// Creates an alert channel which receives alerts from AlertProcessorInterface.
	alertChannel := make(chan *model.AlertFormat)
	// Creates an error channel which receives error from AlertProcessorInterface.
	errorChannel := make(chan error)
	// done resembles completion of process
	done := make(chan bool, 1)

	// Create a reader object for the path provided.
	inFile, decoder, err := io_stream.Read(inputPath)
	if err != nil {
		return err
	}

	// Create a writer object for the path provided.
	outFile, encoder, err := io_stream.Write(outputPath)
	if err != nil {
		return err
	}
	// Defer gets executed at the end of the Function.
	defer inFile.Close()
	defer outFile.Close()

	// Create a alert_processor object
	processor := alert_processor.NewAlertProcessor(decoder, encoder, alertChannel, errorChannel, done)

	// Go Routines to process the alerts
	go processor.ProcessAlerts()
	// Go Routines to send alerts
	go processor.SendAlert()

	// Listen to Alert and Error channel
	// If Error channel got something then close the alertChannel and errorChannel
	// and return the error received.
	for true {
		finish := false
		select {
		case err := <-errorChannel: // Listen errorChannel
			close(errorChannel)
			return err
		case complete := <-done:
			if complete {
				finish = true
			}
		}
		if finish {
			break
		}
	}
	return nil
}
