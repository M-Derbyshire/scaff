package create

import (
	"fmt"
	"os"
	"path"

	"github.com/M-Derbyshire/scaff/models"
	"github.com/M-Derbyshire/scaff/variable"
)

// Directory creates a directory (and its inner directories and files), based on the given DirectoryScaffold.
// The parentDirectoryPath is the path to the directory that will contain this directory.
// The templatesDirectoryPath is the path to the directory that contains templates for files.
// The vars is a map of variables that may be needed to populate the directory name.
func Directory(directory models.DirectoryScaffold, parentDirectoryPath, templatesDirectoryPath string, vars map[string]string) error {

	//Generate the full path to this directory
	populatedDirectoryName, populateNameErr := variable.Populate(directory.Name, vars)
	if populateNameErr != nil {
		return populateNameErr
	}
	fullDirPath := path.Join(parentDirectoryPath, populatedDirectoryName)

	//Create this directory
	dirCreateErr := os.Mkdir(fullDirPath, 0777)
	if dirCreateErr != nil {
		return fmt.Errorf("error while creating directory '%s': %v", populatedDirectoryName, dirCreateErr.Error())
	}

	// Create the files within this directory
	for _, file := range directory.Files {
		fileCreateErr := File(file, fullDirPath, templatesDirectoryPath, vars)
		if fileCreateErr != nil {
			return fileCreateErr
		}
	}

	//Create the directories within this directory
	for _, innerDirectory := range directory.Directories {
		innerDirCreateErr := Directory(innerDirectory, fullDirPath, templatesDirectoryPath, vars)
		if innerDirCreateErr != nil {
			return innerDirCreateErr
		}
	}

	return nil
}
