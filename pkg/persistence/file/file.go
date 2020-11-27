package file

import (
	"io"
	"os"
)

// File represents the file
type File struct {
	Name string
}

// NewFile returns an initialized *File
func NewFile(name string) *File {
	return &File{
		Name: name,
	}
}

// NewReader returns a new io.ReadCloser for the File
func (f *File) NewReader() (io.ReadCloser, error) {
	return os.Open(f.Name)
}

// NewWriter returns a new io.WriteCloser for the File
func (f *File) NewWriter() (io.WriteCloser, error) {
	return os.Create(f.Name)
}
