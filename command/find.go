package command

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"regexp"
	"runtime"

	"github.com/M-Derbyshire/scaff/models"
)

// These are here to make it easier to mock in tests (default values are in the init() func)
var ReadFile func(filePath string) ([]byte, error)
var FileStat func(filePath string) (fs.FileInfo, error)
var CurrentOS string

func init() {
	ReadFile = os.ReadFile
	FileStat = os.Stat
	CurrentOS = runtime.GOOS
}

// Find moves up the directory tree structure (from the given "currentPath"), searching for a file (with the given "fileNameAndExt"),
// until it finds one that includes the correct command ("commandName").
// The returned "foundCommand" is the ScaffoldCommand that was searched for. If the command isn't found in a file, the "isFound" return
// value is false.
// The "templatePath" return value is the full template directory path (generated from the info in the found file).
// If there are any errors reading a file, the errors will be printed.
func Find(commandName, fileNameAndExt, currentPath string) (foundCommand models.ScaffoldCommand, foundTemplatePath string, isFound bool, err error) {
	pathPrefix := "" //Used when constructing file path strings (different depending on OS)
	if CurrentOS != "windows" {
		pathPrefix = "/"
	}

	var command models.ScaffoldCommand
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

func searchFileForCommand(filePath, commandName string) (command models.ScaffoldCommand, templatePath string, isFound bool, err error) {
	emptyCommand := models.ScaffoldCommand{}

	fileBytes, fileErr := ReadFile(filePath)
	if fileErr != nil {
		return emptyCommand, "", false, fileErr
	}

	var config models.ScaffoldConfig
	unmarshalErr := json.Unmarshal(fileBytes, &config)
	if unmarshalErr != nil {
		return emptyCommand, "", false, unmarshalErr
	}

	for _, command := range config.Commands {
		if command.Name == commandName {
			containingDir, _ := path.Split(filePath)
			return command, path.Join(containingDir, config.TemplateDirectoryPath), true, nil
		}
	}

	return emptyCommand, "", false, nil
}
