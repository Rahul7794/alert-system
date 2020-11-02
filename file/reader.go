package file

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONReader type has json.Decoder as field
type JSONReader struct {
	Parser *json.Decoder
}

// ParseFromJSON decodes json to struct provided
func (reader *JSONReader) ParseFromJSON(c interface{}) error {
	return reader.Parser.Decode(c)
}

// HasNext checks if there is any more elements to be decoded
func (reader *JSONReader) HasNext() bool {
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
	return file, &JSONReader{
		Parser: decoder,
	}, nil
}
