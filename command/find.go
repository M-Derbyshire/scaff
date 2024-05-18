package command

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"

	"github.com/M-Derbyshire/scaff/models"
)

// Find moves up the directory tree structure (from the given "currentPath"), searching for a file (with the given "fileNameAndExt"),
// until it finds one that includes the correct command ("commandName").
// The returned "foundCommand" is the ScaffoldCommand that was searched for. If the command isn't found in a file, the "found" return
// value is false.
// The "templatePath" return value is the full template directory path (generated from the info in the found file).
// If there are any errors reading a file, the errors will be printed.
func Find(commandName, fileNameAndExt, currentPath string) (foundCommand models.ScaffoldCommand, foundTemplatePath string, found bool) {

	pathPrefix := "" //Used when constructing file path strings (different depending on OS)
	if runtime.GOOS != "windows" {
		pathPrefix = "/"
	}

	var command models.ScaffoldCommand
	var templatePath string
	commandFound := false

	pathPartsRegex := regexp.MustCompile(`[\\/]`)
	pathParts := pathPartsRegex.Split(currentPath, -1) //Slice of every directory in the current path

	//Keep rebuilding the path, but losing another directory everytime (so we go up the directory structure)
	for i := len(pathParts); i > 0; i-- {

		dirPathToCheck := path.Join(pathParts[0:i]...)
		filePathToCheck := path.Join(pathPrefix, dirPathToCheck, fileNameAndExt)

		//File exists (and can be accessed)
		if _, err := os.Stat(filePathToCheck); err == nil {
			command, templatePath, commandFound = searchFileForCommand(filePathToCheck, commandName)

			if commandFound {
				break
			}
		}
	}

	return command, templatePath, commandFound
}

func searchFileForCommand(filePath, commandName string) (command models.ScaffoldCommand, templatePath string, found bool) {

	emptyCommand := models.ScaffoldCommand{}

	fileBytes, fileErr := os.ReadFile(filePath)
	if fileErr != nil {
		fmt.Println(fileErr.Error())
		return emptyCommand, "", false
	}

	var config models.ScaffoldConfig
	unmarshalErr := json.Unmarshal(fileBytes, &config)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr.Error())
		return emptyCommand, "", false
	}

	for _, command := range config.Commands {
		if command.Name == commandName {
			containingDir, _ := path.Split(filePath)
			return command, path.Join(containingDir, config.TemplateDirectoryPath), true
		}
	}

	return emptyCommand, "", false
}
