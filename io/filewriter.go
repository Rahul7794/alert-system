package io

import (
	"encoding/json"
	"os"

	"alert-system/log"
)

// JSONWriter type has json.Encoder as field
type JSONWriter struct {
	Encoder *json.Encoder
	File    *os.File
}

// Close closes the File, rendering it unusable for I/O.
func (writer *JSONWriter) Close() error {
	return writer.File.Close()
}

// ParseToJSON encodes struct to json
func (writer *JSONWriter) ParseToJSON(c interface{}) error {
	return writer.Encoder.Encode(c)
}

// NewFileWriter takes string as input and creates an WriterInterface object
func NewFileWriter(filename string) WriterInterface {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot open the file %v", err)
	}
	encoder := json.NewEncoder(file)
	return &JSONWriter{
		Encoder: encoder,
		File:    file,
	}
}
