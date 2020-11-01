package io_stream

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReaderInterface provides an abstraction over Reader
type ReaderInterface interface {
	ParseFromJson(c interface{}) error
	HasNext() bool
}

type JsonReader struct {
	Parser *json.Decoder
}

// ParseFromJson decodes json to struct provided
func (reader *JsonReader) ParseFromJson(c interface{}) error {
	return reader.Parser.Decode(c)
}

// HasNext checks if there is any more elements to be decoded
func (reader *JsonReader) HasNext() bool {
	return reader.Parser.More()
}

// Read takes filename as input and creates an Reader object and
// returns *os.File object to close file once finish reading.
func Read(filename string) (*os.File, ReaderInterface, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open the file %v", err)
	}
	decoder := json.NewDecoder(file)
	return file, &JsonReader{
		Parser: decoder,
	}, nil
}
