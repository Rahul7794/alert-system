package alert_processor

import (
	"alert-system/io_stream"
	"alert-system/model"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertProcessor_ProcessAlerts(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(filepath string) (*json.Decoder, *os.File)
		filePath        string
		check           func(chan *model.AlertFormat)
		checkIfComplete func(chan bool)
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
			check: func(out chan *model.AlertFormat) {
				i := 0
				for actual := range out {
					if !reflect.DeepEqual(actual.Alert, "spotChange") {
						t.Errorf("alert Expected = %+v , got = %+v", "spotChange", actual.Alert)
					}
					i++
				}
				assert.Equal(t, i, 2)
			},
			checkIfComplete: func(bools chan bool) {
				for actual := range bools {
					if !reflect.DeepEqual(actual, true) {
						t.Errorf("complete = %+v , got = %+v", true, false)
					}
				}
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
			check: func(out chan *model.AlertFormat) {
				i := 0
				for range out {
					i++
				}
				assert.Equal(t, i, 0)

			},
		},
	}
	for _, tt := range tests {
		decoder, file := tt.setup(tt.filePath)
		defer file.Close()
		out := make(chan *model.AlertFormat)
		errors := make(chan error)
		done := make(chan bool)
		go func() {
			<-done
		}()
		alertProcessor := NewAlertProcessor(&io_stream.JsonReader{
			Parser: decoder,
		}, nil, out, errors, done)
		go alertProcessor.ProcessAlerts()
		tt.check(out)
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
