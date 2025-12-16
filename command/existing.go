package command

import (
	"path/filepath"

	"github.com/M-Derbyshire/scaff/models"
	"github.com/M-Derbyshire/scaff/variable"
)

// IdentifyExistingPaths identifies file/directory paths in the given command that already exist
// The workingDirectory is the path to the current working directory
// The vars is a map of variables that may be needed to populate the directory/file names
func IdentifyExistingPaths(command models.Command, workingDirectory string, vars map[string]string) ([]string, error) {
	results := []string{}

	for _, file := range command.Files {
		populatedName, err := variable.Populate(file.Name, vars)
		if err != nil {
			return results, err
		}

		filePathToCheck := filepath.Join(workingDirectory, populatedName)

		if fileInfo, _ := FileStat(filePathToCheck); fileInfo != nil && !fileInfo.IsDir() {
			results = append(results, filePathToCheck)
		}
	}

	// We only need to check the first layer of directories (the inner files/directories don't need to be checked)
	for _, directory := range command.Directories {
		populatedName, err := variable.Populate(directory.Name, vars)
		if err != nil {
			return results, err
		}

		directoryPathToCheck := filepath.Join(workingDirectory, populatedName)

		if dirInfo, _ := FileStat(directoryPathToCheck); dirInfo != nil && dirInfo.IsDir() {
			results = append(results, directoryPathToCheck)
		}
	}

	return results, nil
}
