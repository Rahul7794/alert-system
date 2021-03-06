package alertprocessor

import (
	"math"

	"alert-system/log"
	"alert-system/model"
	"alert-system/movingmean"
)

const rollingAverageInSec int = 300       // 5 minutes rolling window in seconds
const trendIntervalInSec int = 900        // 15 minutes rise/fall trend of CurrencyPair Rates
const throttleAlertIntervalInSec int = 60 // 1 minute window to throttle outgoing alert

// checkSpotRateChange checks if there rate has dropped/increased by 10%.
func checkSpotRateChange(currentValue, previousValue float64) bool {
	changePercent := math.Abs(((currentValue - previousValue) / previousValue) * 100.0)
	return changePercent > 10
}

// ProcessAlerts consumes stream of currency conversion rates and produces alerts for a number of situations.
func (a *InputTypeProcessor) ProcessAlerts() {
	log.Info("processing currency pairs record ...")
	i := 0 // keep count of record processed
	// create a Map to store moving average for each currency pair => Map[CurrencyPair, MovingAverage]
	currencyPairRates := make(map[string]movingmean.MovingMean)
	// Decoder.More() streams the json until reaches EOF
	for a.Reader.HasNext() {
		i++
		// Deserialize json to CurrencyConversionRates
		var currentRates model.CurrencyConversionRates
		err := a.Reader.ParseFromJSON(&currentRates)
		if err != nil {
			// Outputs error to the out <- channel, if there is an error deserializing json.
			a.ErrorChannel <- err
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
			if checkSpotRateChange(currentRates.Rate, round(movingMean.Average(), 5)) {
				// if there is > 10% change, send the alert in the channel
				a.OutChannel <- model.AlertFormat{
					Timestamp:    currentRates.Timestamp,
					CurrencyPair: currentRates.CurrencyPair,
					Alert:        "spotChange",
					MovingMean:   movingMean,
				}
			}
			// If the trend crosses 15 minutes send alerts to out channel
			if movingMean.Trend > trendIntervalInSec {
				a.OutChannel <- model.AlertFormat{
					Timestamp:    currentRates.Timestamp,
					CurrencyPair: currentRates.CurrencyPair,
					MovingMean:   movingMean,
				}
			}
		} else {
			// Initialize the MovingMean if a new CurrencyPair comes in and adds it to the Map.
			mm := movingmean.New(rollingAverageInSec)
			mm.Add(currentRates.Rate)
			currencyPairRates[currentRates.CurrencyPair] = mm
		}
	}
	log.Infof("processed %d currency pairs record", i)
	close(a.OutChannel)
}

// SendAlert listens to OutChannel and writes alerts to an output file
func (a *InputTypeProcessor) SendAlert() {
	for alert := range a.OutChannel {
		trend := alert.MovingMean.Trend
		// After 15 minutes of continuous rise or fall
		if trend > trendIntervalInSec && alert.Alert == "" {
			// Throttle down the alert sending to one alert/minute
			if (trend-trendIntervalInSec)%throttleAlertIntervalInSec == 0 {
				switch alert.MovingMean.Direction {
				case movingmean.Fall:
					alert.Alert = "falling"
				case movingmean.Rise:
					alert.Alert = "rising"
				}
				alert.Seconds = int32(trend)
				err := a.Writer.ParseToJSON(alert)
				if err != nil {
					a.ErrorChannel <- err
				}
			}
		} else {
			err := a.Writer.ParseToJSON(alert)
			if err != nil {
				a.ErrorChannel <- err
			}
		}
	}
	a.IsComplete <- true
}

// round rounds the precision n times
func round(val float64, n int) float64 {
	floatRound := math.Pow10(n)
	valInt := int64(val * floatRound)
	val = float64(valInt) / floatRound
	return val
}

// NewFileTypeInputProcessor initializes alertprocessor object for file as inputType
func NewFileTypeInputProcessor(input *InputTypeProcessor) AlertProcessor {
	return &InputTypeProcessor{
		Reader:       input.Reader,
		Writer:       input.Writer,
		OutChannel:   input.OutChannel,
		ErrorChannel: input.ErrorChannel,
		IsComplete:   input.IsComplete,
	}
}
