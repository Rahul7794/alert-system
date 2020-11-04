package io

// WriterInterface provides an abstraction over Writer
type WriterInterface interface {
	ParseToJSON(c interface{}) error
	Close() error
}

// ReaderInterface provides an abstraction over Reader
type ReaderInterface interface {
	ParseFromJSON(c interface{}) error
	HasNext() bool
	Close() error
}

// InputSource store key, value pair, creating object for different input source
var InputSource = map[string]func(input string) ReaderInterface{
	"fileType": NewFileReader,
}

// OutDestination store key, value pair, creating object for different out dest
var OutDestination = map[string]func(output string) WriterInterface{
	"fileType": NewFileWriter,
}

// NewReader returns object of ReaderInterface
func NewReader(source string, input string) ReaderInterface {
	if _, ok := InputSource[source]; ok {
		return InputSource[source](input)
	}
	return nil
}

// NewWriter returns object of WriterInterface
func NewWriter(dest string, input string) WriterInterface {
	if _, ok := OutDestination[dest]; ok {
		return OutDestination[dest](input)
	}
	return nil
}
