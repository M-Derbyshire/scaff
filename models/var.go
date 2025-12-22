package models

import (
	"io/fs"
	"os"
	"path/filepath"
)

// These are here to make it easier to mock in tests (default values are in the init() func)
var FileStat func(filePath string) (fs.FileInfo, error)
var AbsPath func(path string) (string, error)

func init() {
	FileStat = os.Stat
	AbsPath = filepath.Abs
}
