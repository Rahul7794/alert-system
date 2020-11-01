package cmd

import (
	"alert-system/alert_processor"
	"alert-system/io_stream"
	"alert-system/log"
	"alert-system/model"
	"alert-system/version"
	"encoding/json"

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
	// Creates an alert channel which receives alerts from AlertProcessor.
	alertChannel := make(chan *model.AlertFormat)
	// Creates an error channel which receives error from AlertProcessor.
	errorChannel := make(chan error)

	// Creates a reader object for the path provided.
	reader, err := io_stream.Reader(inputPath)
	if err != nil {
		// return error
		return err
	}
	// defer gets executed at the end of the Function.
	defer reader.Close()

	// Create a writer object for the path provided.
	writer, err := io_stream.Writer(outputPath)
	if err != nil {
		return err
	}
	// defer gets executed at the end of the Function and closes the object
	defer writer.Close()

	// Create a new Decoder
	decoder := json.NewDecoder(reader)
	// Create a new Encoder
	encoder := json.NewEncoder(writer)

	// Create a alert_processor object
	processor := alert_processor.NewAlertProcessor(encoder, decoder)
	// Go Routines to Process the alerts
	go processor.ProcessAlerts(alertChannel, errorChannel)

	// Listen to Alert and Error channel
	// If Error channel got something then close the alertChannel and errorChannel
	// and return the error received.
	for true {
		done := false
		select {
		case err := <-errorChannel: // Listen errorChannel
			close(errorChannel)
			return err
		case alerts, ok := <-alertChannel: // Listen AlertChannel
			if ok {
				// if alerts is received, send it
				err := processor.SendAlert(alerts)
				if err != nil {
					return err
				}
			} else {
				// if nothing is being received in the channel, assuming all the currencyPairs rates are processed.
				done = true
			}
		}
		// break the loop if all the currencyPairs rates are processed
		if done {
			break
		}
	}
	return nil
}
