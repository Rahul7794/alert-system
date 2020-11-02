package model

import "alert-system/movingmean"

// CurrencyConversionRates are decoded from Incoming json
type CurrencyConversionRates struct {
	Timestamp    float64 `json:"timestamp"`    // Timestamp of the incoming CurrencyPair
	CurrencyPair string  `json:"currencyPair"` // Incoming CurrencyPair
	Rate         float64 `json:"rate"`         // Spot Rate of a incoming CurrencyPair
}

// AlertFormat encoded from Alert Outputs
type AlertFormat struct {
	Timestamp    float64               `json:"timestamp"`    // Timestamp of the outgoing CurrencyPair is Alert
	CurrencyPair string                `json:"currencyPair"` // Outgoing CurrencyPair
	Alert        string                `json:"alert"`        // Alert Type
	Seconds      int32                 `json:",omitempty"`   // Number of Seconds the CurrencyPair trend was rising/falling
	MovingMean   movingmean.MovingMean `json:"-"`            // MovingMean holding calculated trend and Seconds for that CurrencyPair
}
