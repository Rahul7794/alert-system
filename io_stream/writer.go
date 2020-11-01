package io_stream

import (
	"fmt"
	"os"
)

func Writer(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("cannot open the file %v", err)
	}
	return file, nil
}
