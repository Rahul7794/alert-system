package io_stream

import (
	"fmt"
	"os"
)

func Reader(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open the file %v", err)
	}
	return file, nil
}
