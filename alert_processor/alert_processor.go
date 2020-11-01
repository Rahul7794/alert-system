package alert_processor

import (
	"alert-system/io_stream"
	"alert-system/model"
	"alert-system/moving_mean"
	"math"
)

// 5 minutes rolling window in seconds
const rollingAverageInSec int = 300

// AlertProcessor process the alerts.
type AlertProcessor struct {
	OutChannel   chan *model.AlertFormat   // OutChannel contains alerts
	ErrorChannel chan error                // ErrorChannel contains error
	IsComplete   chan bool                 // IsComplete indicates if the processing is complete for gracefully close all the open channels
	Reader       io_stream.ReaderInterface // Decoder stream the rates from the input file.
	Writer       io_stream.WriterInterface // Writer stream the alerts to the output file.
}

// checkSpotRateChange checks if there rate has dropped/increased by 10%.
func checkSpotRateChange(currentValue, previousValue float64) bool {
	changePercent := math.Abs(((currentValue - previousValue) / previousValue) * 100.0)
	return changePercent > 10
}

// ProcessAlerts consumes stream of currency conversion rates and produces alerts for a number of situations.
func (a *AlertProcessor) ProcessAlerts() {
	// create a Map to store moving average for each currency pair => Map[CurrencyPair, MovingAverage]
	currencyPairRates := make(map[string]*moving_mean.MovingMean)
	// Decoder.More() streams the json until reaches EOF
	for a.Reader.HasNext() {
		// Deserialize json to CurrencyConversionRates
		var currentRates model.CurrencyConversionRates
		err := a.Reader.ParseFromJson(&currentRates)
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
			if checkSpotRateChange(currentRates.Rate, floatRound(movingMean.Average(), 5)) {
				// if there is > 10% change, send the alert in the channel
				a.OutChannel <- &model.AlertFormat{
					Timestamp:    currentRates.Timestamp,
					CurrencyPair: currentRates.CurrencyPair,
					Alert:        "spotChange",
					MovingMean:   movingMean,
				}
			}
			// If the trend crosses 15 minutes
			if movingMean.Trend > 900 {
				a.OutChannel <- &model.AlertFormat{
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
	a.IsComplete <- true
	close(a.OutChannel) // close the alert channel once all the currencyPair rates are processed
}

// Writes alerts to an output file
func (a *AlertProcessor) SendAlert() {
	for alert := range a.OutChannel {
		switch trend := alert.MovingMean.Trend; {
		case trend > 900: // After 15 minutes continuous rise or fall
			if (trend-900)%60 == 0 { // Throttle down the alert sending to one alert/minute
				// Change alert message for continuous rise or fall
				switch alert.MovingMean.Direction { // Change alert type based on Direction
				case moving_mean.Fall:
					alert.Alert = "falling"
				case moving_mean.Rise:
					alert.Alert = "rising"
				}
				alert.Seconds = int32(trend) // Set trend as seconds
			}
		}
		alert.MovingMean = nil // make movingMean pointer reference to nil to avoid printing in output file
		// write alert to the output file
		err := a.Writer.ParseToJson(alert)
		if err != nil {
			a.ErrorChannel <- err
		}
	}
}

// floatRound rounds the precision n times
func floatRound(val float64, n int) float64 {
	floatRound := math.Pow10(n)
	valInt := int64(val * floatRound)
	val = float64(valInt) / floatRound
	return val
}

// NewAlertProcessor initializes alert_processor object
func NewAlertProcessor(reader io_stream.ReaderInterface, writer io_stream.WriterInterface,
	out chan *model.AlertFormat, error chan error, done chan bool) AlertProcessorInterface {
	return &AlertProcessor{
		Reader:       reader,
		Writer:       writer,
		OutChannel:   out,
		ErrorChannel: error,
		IsComplete:   done,
	}
}
