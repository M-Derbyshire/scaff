package models

import (
	"io/fs"
	"os"
)

// These are here to make it easier to mock in tests (default values are in the init() func)
var FileStat func(filePath string) (fs.FileInfo, error)

func init() {
	FileStat = os.Stat
}
