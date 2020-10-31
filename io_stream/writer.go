package io_stream

import (
	"alert-system/model"
	"encoding/json"
	"os"
)

func WriteJson(filename string, alerts []model.AlertFormat) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	for _, alert := range alerts {
		encoder.Encode(&alert)
	}
	return nil
}
