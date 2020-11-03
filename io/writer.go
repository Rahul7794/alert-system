package io

import (
	"encoding/json"
	"fmt"
	"os"
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

// Write takes filename as input and creates an Writer object and
// returns *os.File object to close the file once finish writing.
func Write(filename string) (WriterInterface, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("cannot open the file %v", err)
	}
	encoder := json.NewEncoder(file)
	return &JSONWriter{
		Encoder: encoder,
		File:    file,
	}, nil
}
