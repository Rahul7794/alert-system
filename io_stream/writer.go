package io_stream

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriterInterface which can be implemented to write any client
type WriterInterface interface {
	ParseToJson(c interface{}) error
}

// JsonWriter
type JsonWriter struct {
	Encoder *json.Encoder
}

func (writer *JsonWriter) ParseToJson(c interface{}) error {
	return writer.Encoder.Encode(c)
}

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
