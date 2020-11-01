package io_stream

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriterInterface provides an abstraction over Writer
type WriterInterface interface {
	ParseToJson(c interface{}) error
}

type JsonWriter struct {
	Encoder *json.Encoder
}

// ParseToJson encodes struct to json
func (writer *JsonWriter) ParseToJson(c interface{}) error {
	return writer.Encoder.Encode(c)
}

// Writes takes filename as input and creates an Writer object and
// returns *os.File object to close the file once finish writing.
func Write(filename string) (*os.File, WriterInterface, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open the file %v", err)
	}
	encoder := json.NewEncoder(file)
	return file, &JsonWriter{
		Encoder: encoder,
	}, nil
}
