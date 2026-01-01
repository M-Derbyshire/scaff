package command

import (
	"io/fs"
	"os"
	"runtime"
)

// These are here to make it easier to mock in tests (default values are in the init() func)

// ReadFile is used to read files from the filesystem
var ReadFile func(filePath string) ([]byte, error)

// FileStat is used to get details about files in the filesystem (this can also be used to confirm a file exists)
var FileStat func(filePath string) (fs.FileInfo, error)

// CurrentOS identifies the current operating system
var CurrentOS string

func init() {
	ReadFile = os.ReadFile
	FileStat = os.Stat
	CurrentOS = runtime.GOOS
}
