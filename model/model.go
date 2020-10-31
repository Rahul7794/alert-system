package model

type CurrencyConversionRates struct {
	Timestamp    float64 `json:"timestamp"`
	CurrencyPair string  `json:"currencyPair"`
	Rate         float64 `json:"rate"`
}

type AlertFormat struct {
	Timestamp    float64 `json:"timestamp"`
	CurrencyPair string  `json:"currencyPair"`
	Alert        string  `json:"alert"`
	Seconds      int32   `json:",omitempty"`
}
