package persistence

import (
	"io"
)

// Persistence defines the interface
type Persistence interface {
	NewReader() (io.ReadCloser, error)
	NewWriter() io.WriteCloser
}
