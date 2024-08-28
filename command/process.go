package command

import (
	"github.com/M-Derbyshire/scaff/create"
	"github.com/M-Derbyshire/scaff/models"
)

var (
	CreateFile      func(file models.FileScaffold, parentDirectoryPath, fullTemplatesDirectoryPath string, vars map[string]string) error
	CreateDirectory func(directory models.DirectoryScaffold, parentDirectoryPath, fullTemplatesDirectoryPath string, vars map[string]string) error
)

func init() {
	CreateFile = create.File
	CreateDirectory = create.Directory
}

// Process creates directories/files from the data in the given ScaffoldCommand.
// The workingDirectory is the path to the current working directory
// The fullTemplatesDirectoryPath is the path to the directory that contains templates for files.
// The vars is a map of variables that may be needed to populate the directory/file names, and file contents.
func Process(command models.Command, workingDirectory, fullTemplatesDirectoryPath string, vars map[string]string) error {
	for _, file := range command.Files {
		fileCreateErr := CreateFile(file, workingDirectory, fullTemplatesDirectoryPath, vars)
		if fileCreateErr != nil {
			return fileCreateErr
		}
	}

	for _, directory := range command.Directories {
		dirCreateErr := CreateDirectory(directory, workingDirectory, fullTemplatesDirectoryPath, vars)
		if dirCreateErr != nil {
			return dirCreateErr
		}
	}

	return nil
}
