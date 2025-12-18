package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"regexp"

	"github.com/M-Derbyshire/scaff/models"
)

// Find moves up the directory tree structure (from the given "currentPath"), searching for a file (with the given "fileNameAndExt"),
// until it finds one that includes the correct command ("commandName").
// The returned "foundCommand" is the ScaffoldCommand that was searched for. If the command isn't found in a file, the "isFound" return
// value is false.
// The "templatePath" return value is the full template directory path (generated from the info in the found file).
// If there are any errors reading a file, the errors will be printed.
func Find(commandName, fileNameAndExt, currentPath string) (foundCommand models.Command, fullTemplatePath string, isFound bool, err error) {
	pathPrefix := "" //Used when constructing file path strings (different depending on OS)
	if CurrentOS != "windows" {
		pathPrefix = "/"
	}

	var command models.Command
	var templatePath string
	commandFound := false
	var searchErr error

	pathPartsRegex := regexp.MustCompile(`[\\/]`)
	pathParts := pathPartsRegex.Split(currentPath, -1) //Slice of every directory in the current path

	//Keep rebuilding the path, but losing another directory everytime (so we go up the directory structure)
	for i := len(pathParts); i > 0; i-- {
		dirPathToCheck := path.Join(pathParts[0:i]...)
		filePathToCheck := path.Join(pathPrefix, dirPathToCheck, fileNameAndExt)

		//File exists (and can be accessed)
		if _, statErr := FileStat(filePathToCheck); statErr == nil {
			command, templatePath, commandFound, searchErr = searchFileForCommand(filePathToCheck, commandName)

			if searchErr != nil {
				return command, templatePath, commandFound, searchErr
			}

			if commandFound {
				break
			}
		}
	}

	return command, templatePath, commandFound, searchErr
}

func searchFileForCommand(filePath, commandName string) (command models.Command, fullTemplatePath string, isFound bool, err error) {
	emptyCommand := models.Command{}
	containingDir, _ := path.Split(filePath)

	fileBytes, fileErr := ReadFile(filePath)
	if fileErr != nil {
		return emptyCommand, "", false, fileErr
	}

	var scaffFile models.ScaffFile
	unmarshalErr := json.Unmarshal(fileBytes, &scaffFile)
	if unmarshalErr != nil {
		invalidJsonMsg := fmt.Sprintf("encountered a scaff.json file with an invalid structure: '%s'", filePath)
		return emptyCommand, "", false, errors.New(invalidJsonMsg)
	}

	// Search through the commands array
	for _, command := range scaffFile.Commands {
		if command.Name == commandName {
			return command, path.Join(containingDir, command.TemplateDirectoryPath), true, nil
		}
	}

	// Search through any child files
	if validationErr := scaffFile.ValidateChildrenArray(); validationErr != nil {
		return emptyCommand, "", false, validationErr
	}

	for _, childPath := range scaffFile.Children {
		fullChildPath := path.Join(containingDir, childPath)

		if _, childPathErr := FileStat(fullChildPath); childPathErr != nil {
			childFindErrMsg := fmt.Sprintf("unable to locate child scaff file at path: '%s'", fullChildPath)
			return emptyCommand, "", false, errors.New(childFindErrMsg)
		}

		childCommand, childTemplatePath, foundInChild, childErr := searchFileForCommand(fullChildPath, commandName)
		if childErr != nil {
			return emptyCommand, "", false, childErr
		}

		if foundInChild {
			return childCommand, childTemplatePath, true, nil
		}
	}

	return emptyCommand, "", false, nil
}
