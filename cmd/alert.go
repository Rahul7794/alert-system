package cmd

import (
	"alert-system/io_stream"
	"alert-system/log"
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
	runCmd.Flags().StringVarP(&inputPath, "ipath", "i", "", "input path for the json file")
	err := runCmd.MarkFlagRequired("ipath")
	runCmd.Flags().StringVarP(&outputPath, "opath", "o", "", "output path for the json file")
	err = runCmd.MarkFlagRequired("opath")
	handleCommandError(err)
}

func alertCmd(_ *cobra.Command, _ []string) error {
	log.Infof("starting the alerting system with version: %s %s \n on %s => %s", version.Version, version.BuildDate,
		version.OsArch, version.GoVersion)
	alerts, err := io_stream.ReadJson(inputPath)
	if err != nil {
		return err
	}
	err = io_stream.WriteJson(outputPath, alerts)
	if err != nil {
		return err
	}
	return nil
}
