package alert_processor

import (
	"alert-system/model"
	"alert-system/moving_mean"
	"encoding/json"
	"fmt"
	"math"
)

// 5 minutes rolling window in seconds
const rollingAverageInSec int = 300

type AlertProcessorInterface interface {
	ProcessAlerts(out chan<- *model.AlertFormat, errors chan<- error)
	SendAlert(in *model.AlertFormat) error
}

// AlertProcessor Process the alerts.
type AlertProcessor struct {
	Decoder *json.Decoder // Decoder stream the rates from the input file.
	Encoder *json.Encoder // Encoder stream the alerts to the output file.
}

// checkSpotRateChange checks if there rate has dropped/increased by 10%.
func checkSpotRateChange(currentValue, previousValue float64) bool {
	changePercent := ((currentValue - previousValue) / previousValue) * 100.0
	return math.Abs(changePercent) > 10
}

// ProcessAlerts consumes stream of currency conversion rates and produces alerts for a number of situations.
func (a *AlertProcessor) ProcessAlerts(out chan<- *model.AlertFormat, errors chan<- error) {
	// create a Map to store moving average for each currency pair => Map[CurrencyPair, MovingAverage]
	currencyPairRates := make(map[string]*moving_mean.MovingMean)
	// Decoder.More() streams the json until reaches EOF
	for a.Decoder.More() {
		// Deserialize json to CurrencyConversionRates
		var currentRates model.CurrencyConversionRates
		err := a.Decoder.Decode(&currentRates)
		if err != nil {
			// Outputs error to the out <- channel, if there is an error deserializing json.
			errors <- err
			return
		}
		// Check if currencyPair key exists in the Map
		// if exists add the rates of the currencyPair to the movingMean
		// and calculate the moving mean and check for > 10% change in rate
		if movingMean, ok := currencyPairRates[currentRates.CurrencyPair]; ok {
			// Add rates to the movingMean
			movingMean.Add(currentRates.Rate)
			currencyPairRates[currentRates.CurrencyPair] = movingMean
			// Check for > 10% change in rate
			if checkSpotRateChange(currentRates.Rate, floatRound(movingMean.Average(), 5)) {
				// if there is > 10% change, send the alert in the channel
				out <- &model.AlertFormat{
					Timestamp:    currentRates.Timestamp,
					CurrencyPair: currentRates.CurrencyPair,
					Alert:        "spotChange",
					MovingMean:   movingMean,
				}
			}
			// If the trend crosses 15 minutes
			if movingMean.Trend > 900 {
				out <- &model.AlertFormat{
					Timestamp:    currentRates.Timestamp,
					CurrencyPair: currentRates.CurrencyPair,
					MovingMean:   movingMean,
				}
			}
		} else {
			// Initialize the MovingMean if a new CurrencyPair comes in and adds it to the Map.
			mm := moving_mean.New(rollingAverageInSec)
			mm.Add(currentRates.Rate)
			currencyPairRates[currentRates.CurrencyPair] = mm
		}
	}
	close(out) // close the alert channel once all the currencyPair rates are processed
}

// Writes alerts to an output file
func (a *AlertProcessor) SendAlert(in *model.AlertFormat) error {
	switch trend := in.MovingMean.Trend; {
	case trend > 900: // After 15 minutes continuous rise or fall
		if (trend-900)%60 == 0 { // Throttle down the alert sending to one alert/minute
			// Change alert message for continuous rise or fall
			switch in.MovingMean.Direction { // Change alert type based on Direction
			case moving_mean.Fall:
				in.Alert = "falling"
			case moving_mean.Rise:
				in.Alert = "rising"
			}
			in.Seconds = int32(trend) // Set trend as seconds
		}
	}
	in.MovingMean = nil // make movingMean pointer reference to nil to avoid printing in output file
	// write alert to the output file
	err := a.Encoder.Encode(&in)
	if err != nil {
		return fmt.Errorf("error encoding json to file %v", err)
	}
	return nil
}

// floatRound rounds the precision n times
func floatRound(val float64, n int) float64 {
	floatRound := math.Pow10(n)
	valInt := int64(val * floatRound)
	val = float64(valInt) / floatRound
	return val
}

// NewAlertProcessor initializes AlertProcessor with Encoder and Decoder
func NewAlertProcessor(encoder *json.Encoder, decoder *json.Decoder) AlertProcessorInterface {
	return &AlertProcessor{
		Decoder: decoder,
		Encoder: encoder,
	}
}
