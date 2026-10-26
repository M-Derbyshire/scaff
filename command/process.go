package command

import (
	"github.com/M-Derbyshire/scaff/command/internal"
	"github.com/M-Derbyshire/scaff/models"
)

// Process creates directories/files from the data in the given ScaffoldCommand.
// The workingDirectory is the path to the current working directory
// The fullTemplatesDirectoryPath is the path to the directory that contains templates for files.
// The vars is a map of variables that may be needed to populate the directory/file names, and file contents.
func Process(command models.Command, workingDirectory, fullTemplatesDirectoryPath string, vars map[string]string) error {
	for _, file := range command.Files {
		fileCreateErr := internal.CreateFile(file, workingDirectory, fullTemplatesDirectoryPath, vars)
		if fileCreateErr != nil {
			return fileCreateErr
		}
	}

	for _, directory := range command.Directories {
		dirCreateErr := internal.CreateDirectory(directory, workingDirectory, fullTemplatesDirectoryPath, vars)
		if dirCreateErr != nil {
			return dirCreateErr
		}
	}

	return nil
}
