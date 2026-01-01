package models

import (
	"io/fs"
	"os"
)

// FileStat is used to get details about files in the filesystem (this can also be used to confirm a file exists)
var FileStat func(filePath string) (fs.FileInfo, error)

func init() {
	FileStat = os.Stat
}
