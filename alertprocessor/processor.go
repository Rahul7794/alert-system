package alertprocessor

// ProcessorInterface defines func to be implemented
type ProcessorInterface interface {
	ProcessAlerts() // ProcessAlerts consumes and process incoming CurrencyPairs rates
	SendAlert()     // SendAlert produces alerts for a number of situations.
}
