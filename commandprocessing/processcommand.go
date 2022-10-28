package commandprocessing

import "github.com/M-Derbyshire/scaff/models"

// Creates directories/files from the data in the given ScaffoldCommand.
// The workingDirectory is the path to the current working directory
// The templatesDirectoryPath is the path to the directory that contains templates for files.
// The vars is a map of variables that may be needed to populate the directory/file names, and file contents.
func ProcessCommand(command models.ScaffoldCommand, workingDirectory, templatesDirectoryPath string, vars map[string]string) error {

	for _, directory := range command.Directories {
		dirCreateErr := CreateDirectory(directory, workingDirectory, templatesDirectoryPath, vars)
		if dirCreateErr != nil {
			return dirCreateErr
		}
	}

	return nil
}
