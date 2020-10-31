package alert_processor

import (
	"alert-system/model"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestAlertProcessor_ProcessAlerts(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(filepath string) (*json.Decoder, *os.File)
		filePath string
		output   []model.AlertFormat
	}{
		{
			name: "detect spot rate changes with the stream of events coming in",
			setup: func(filepath string) (*json.Decoder, *os.File) {
				file, err := os.Open(filepath)
				if err != nil {
					return nil, nil
				}
				return json.NewDecoder(file), file
			},
			filePath: "../test_files/input_alert.json",
			output: []model.AlertFormat{
				{
					Timestamp:    1554933794.023,
					CurrencyPair: "CNYAUD",
					Alert:        "spotChange",
					Seconds:      0,
				},
				{
					Timestamp:    1554933794.023,
					CurrencyPair: "ABCDEF",
					Alert:        "spotChange",
					Seconds:      0,
				},
			},
		},
		{
			name: "no alert as the rate did not increase or decrease by 10%",
			setup: func(filepath string) (*json.Decoder, *os.File) {
				file, err := os.Open(filepath)
				if err != nil {
					return nil, nil
				}
				return json.NewDecoder(file), file
			},
			filePath: "../test_files/input_no_alert.json",
			output:   nil,
		},
	}
	for _, tt := range tests {
		decoder, file := tt.setup(tt.filePath)
		defer file.Close()
		alertProcessor := AlertProcessor{
			Decoder: decoder,
		}
		alerts, err := alertProcessor.ProcessAlerts()
		if err != nil {
			fmt.Println(err)
		}
		if !reflect.DeepEqual(alerts, tt.output) {
			t.Errorf("ProcessAlert=%v, wanted=%v", alerts, tt.output)
		}
	}

}

func TestAlertProcessor_CheckSpotRateChange(t *testing.T) {
	tests := []struct {
		name              string
		currentValue      float64
		previousValue     float64
		is10PercentChange bool
	}{
		{
			name:              "spot rate increase more than 10%",
			currentValue:      0.51234,
			previousValue:     0.32244,
			is10PercentChange: true,
		},
		{
			name:              "spot rate decrease more than 10%",
			currentValue:      0.11111,
			previousValue:     0.32244,
			is10PercentChange: true,
		},
		{
			name:              "spot rate increase less than 10%",
			currentValue:      0.29999,
			previousValue:     0.32244,
			is10PercentChange: false,
		},
	}
	for _, tt := range tests {
		actualValue := checkSpotRateChange(tt.currentValue, tt.previousValue)
		if !reflect.DeepEqual(actualValue, tt.is10PercentChange) {
			t.Errorf("FloatRound(%v, %v)=%v, wanted=%v", tt.currentValue, tt.previousValue, actualValue, tt.is10PercentChange)
		}

	}
}

func TestAlertProcessor_FloatRoundFive(t *testing.T) {
	tests := []struct {
		name        string
		input       float64
		output      float64
		roundNumber int
	}{
		{
			name:        "Successfully round",
			input:       5.34444444444,
			output:      5.34444,
			roundNumber: 5,
		}, {
			name:        "Successfully round another round value",
			input:       5.34444444444,
			output:      5.34,
			roundNumber: 2,
		},
	}

	for _, tt := range tests {
		actualValue := floatRound(tt.input, tt.roundNumber)
		if !reflect.DeepEqual(actualValue, tt.output) {
			t.Errorf("FloatRound(%v, %v)=%v, wanted=%v", tt.input, tt.roundNumber, actualValue, tt.output)
		}

	}
}
