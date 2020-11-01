package io_stream

import (
	"encoding/json"
	"fmt"
	"os"
)

//
type ReaderInterface interface {
	ParseFromJson(c interface{}) error
	HasNext() bool
}

type JsonReader struct {
	Parser *json.Decoder
}

func (reader *JsonReader) ParseFromJson(c interface{}) error {
	return reader.Parser.Decode(c)
}

func (reader *JsonReader) HasNext() bool {
	return reader.Parser.More()
}

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
