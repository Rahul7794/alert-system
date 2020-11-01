package model

import "alert-system/moving_mean"

// Incoming Json are decoded to CurrencyConversionRates
type CurrencyConversionRates struct {
	Timestamp    float64 `json:"timestamp"`    // Timestamp of the incoming CurrencyPair
	CurrencyPair string  `json:"currencyPair"` // Incoming CurrencyPair
	Rate         float64 `json:"rate"`         // Spot Rate of a incoming CurrencyPair
}

// Alert Outputs are encoded to AlertFormat
type AlertFormat struct {
	Timestamp    float64                 `json:"timestamp"`    // Timestamp of the outgoing CurrencyPair is Alert
	CurrencyPair string                  `json:"currencyPair"` // Outgoing CurrencyPair
	Alert        string                  `json:"alert"`        // Alert Type
	Seconds      int32                   `json:",omitempty"`   // Number of Seconds the CurrencyPair trend was rising/falling
	MovingMean   *moving_mean.MovingMean `json:",omitempty"`   // MovingMean holding calculated trend and Seconds for that CurrencyPair
}
