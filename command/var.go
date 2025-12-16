package command

import (
	"io/fs"
	"os"
	"runtime"
)

// These are here to make it easier to mock in tests (default values are in the init() func)
var ReadFile func(filePath string) ([]byte, error)
var FileStat func(filePath string) (fs.FileInfo, error)
var CurrentOS string

func init() {
	ReadFile = os.ReadFile
	FileStat = os.Stat
	CurrentOS = runtime.GOOS
}
