package internal

import (
	"io/fs"
	"os"
	"runtime"

	"github.com/M-Derbyshire/scaff/create"
	"github.com/M-Derbyshire/scaff/models"
)

var (
	// CreateFile is used to create files in the filesystem
	CreateFile func(file models.FileScaffold, parentDirectoryPath, fullTemplatesDirectoryPath string, vars map[string]string) error
	// CreateDirectory is used to create directories in the filesystem
	CreateDirectory func(directory models.DirectoryScaffold, parentDirectoryPath, fullTemplatesDirectoryPath string, vars map[string]string) error
	// ReadFile is used to read files from the filesystem
	ReadFile func(filePath string) ([]byte, error)
	// FileStat is used to get details about files in the filesystem (this can also be used to confirm a file exists)
	FileStat func(filePath string) (fs.FileInfo, error)
	// CurrentOS identifies the current operating system
	CurrentOS string
)

func init() {
	CreateFile = create.File
	CreateDirectory = create.Directory
	ReadFile = os.ReadFile
	FileStat = os.Stat
	CurrentOS = runtime.GOOS
}
