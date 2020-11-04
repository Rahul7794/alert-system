package alertprocessor

import (
	"alert-system/io"
	"alert-system/model"
)

// AlertProcessor defines func to be implemented
type AlertProcessor interface {
	ProcessAlerts() // ProcessAlerts consumes and process incoming CurrencyPairs rates
	SendAlert()     // SendAlert produces alerts for a number of situations.
}

// InputTypeProcessor process the currency pairs from file as Input.
type InputTypeProcessor struct {
	OutChannel   chan model.AlertFormat // OutChannel contains alerts
	ErrorChannel chan<- error           // ErrorChannel contains error
	IsComplete   chan<- bool            // IsComplete indicates if the processing is complete for gracefully close all the open channels
	Reader       io.ReaderInterface     // Reader stream the rates from the input file.
	Writer       io.WriterInterface     // Writer stream the alerts to the output file.
}

// Sources store key, value pair of inputType and function returning object of that inputType processor
var Sources = map[string]func(input *InputTypeProcessor) AlertProcessor{
	"fileType": NewFileTypeInputProcessor,
}

// NewProcessorObject return the object of the different processor type based on the source provided.
func NewProcessorObject(source string, processor InputTypeProcessor) AlertProcessor {
	inputType := &InputTypeProcessor{
		OutChannel:   processor.OutChannel,
		ErrorChannel: processor.ErrorChannel,
		IsComplete:   processor.IsComplete,
		Reader:       processor.Reader,
		Writer:       processor.Writer,
	}
	if _, ok := Sources[source]; ok {
		return Sources[source](inputType)
	}
	return nil
}
