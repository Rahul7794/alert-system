package alert_processor

type AlertProcessorInterface interface {
	ProcessAlerts()
	SendAlert()
}
