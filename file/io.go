package file

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
