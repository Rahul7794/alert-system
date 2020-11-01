package alert_processor

type AlertProcessorInterface interface {
	ProcessAlerts() // ProcessAlerts consumes and process incoming CurrencyPairs rates
	SendAlert()     // SendAlert produces alerts for a number of situations.
}
