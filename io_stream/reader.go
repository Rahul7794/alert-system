package io_stream

import (
	"alert-system/alert_processor"
	"alert-system/model"
	"encoding/json"
	"fmt"
	"os"
)

func ReadJson(filename string) ([]model.AlertFormat, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open the file %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	processor := alert_processor.AlertProcessor{
		Decoder: decoder,
	}
	alerts, err := processor.ProcessAlerts()
	if err != nil {
		return nil, err
	}
	return alerts, nil
}
