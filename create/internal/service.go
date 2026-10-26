package internal

import (
	"io/fs"
	"os"
)

var (
	// Mkdir is used to create a directory in the filesystem
	Mkdir func(string, fs.FileMode) error
	// ReadFile is used to read files from the filesystem
	ReadFile func(string) ([]byte, error)
	// WriteFile is used to create files in the filesystem
	WriteFile func(string, []byte, fs.FileMode) error
)

func init() {
	ReadFile = os.ReadFile
	WriteFile = os.WriteFile
	Mkdir = os.Mkdir
}
