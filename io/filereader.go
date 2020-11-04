package io

import (
	"encoding/json"
	"os"

	"alert-system/log"
)

// JSONReader type has json.Decoder as field
type JSONReader struct {
	Parser *json.Decoder
	File   *os.File
}

// ParseFromJSON decodes json to struct provided
func (reader *JSONReader) ParseFromJSON(c interface{}) error {
	return reader.Parser.Decode(c)
}

// HasNext checks if there is any more elements to be decoded
func (reader *JSONReader) HasNext() bool {
	return reader.Parser.More()
}

// Close closes the File, rendering it unusable for I/O.
func (reader *JSONReader) Close() error {
	return reader.File.Close()
}

// NewFileReader takes filename as input and creates an ReaderInterface object
func NewFileReader(filename string) ReaderInterface {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open the file %v", err)
	}
	decoder := json.NewDecoder(file)
	return &JSONReader{
		Parser: decoder,
		File:   file,
	}
}
