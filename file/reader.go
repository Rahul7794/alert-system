package file

import (
	"encoding/json"
	"fmt"
	"os"
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

// Read takes filename as input and creates an Reader object and
// returns *os.File object to close file once finish reading.
func Read(filename string) (ReaderInterface, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open the file %v", err)
	}
	decoder := json.NewDecoder(file)
	return &JSONReader{
		Parser: decoder,
		File:   file,
	}, nil
}
