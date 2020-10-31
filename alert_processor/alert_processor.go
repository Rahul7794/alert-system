package alert_processor

import (
	"alert-system/model"
	"alert-system/moving_mean"
	"encoding/json"
	"fmt"
	"math"
)

const rollingAverageInSec int = 300

type AlertProcessorInterface interface {
	ProcessAlerts(decoder *json.Decoder) ([]model.AlertFormat, error)
}

type AlertProcessor struct {
	Decoder *json.Decoder
}

func checkSpotRateChange(currentValue, previousValue float64) bool {
	changePercent := ((currentValue - previousValue) / previousValue) * 100.0
	if math.Abs(changePercent) > 10 {
		return true
	}
	return false
}

func (a *AlertProcessor) ProcessAlerts() ([]model.AlertFormat, error) {
	var alerts []model.AlertFormat
	mapData := make(map[string]*moving_mean.MovingMean)
	for a.Decoder.More() {
		var currencyConRate model.CurrencyConversionRates
		err := a.Decoder.Decode(&currencyConRate)
		if err != nil {
			return nil, fmt.Errorf("error decoding the json %v", err)
		}
		if currencyRates, ok := mapData[currencyConRate.CurrencyPair]; ok {
			currencyRates.Add(floatRound(currencyConRate.Rate, 5))
			mapData[currencyConRate.CurrencyPair] = currencyRates
			if checkSpotRateChange(floatRound(currencyConRate.Rate, 5), floatRound(currencyRates.Average(), 5)) {
				alerts = append(alerts, model.AlertFormat{
					Timestamp:    currencyConRate.Timestamp,
					CurrencyPair: currencyConRate.CurrencyPair,
					Alert:        "spotChange",
				})
			}
		} else {
			mapData[currencyConRate.CurrencyPair] = moving_mean.New(rollingAverageInSec)
		}
	}
	return alerts, nil
}

func floatRound(val float64, round int) float64 {
	floatRound := math.Pow10(round)
	valInt := int64(val * floatRound)
	val = float64(valInt) / floatRound
	return val
}
