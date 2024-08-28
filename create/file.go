package create

import (
	"io/fs"
	"os"
	"path"

	"github.com/M-Derbyshire/scaff/models"
	"github.com/M-Derbyshire/scaff/variable"
)

// These are here to make it easier to mock in tests (default values are in the init() func)
var ReadFile func(string) ([]byte, error)
var WriteFile func(string, []byte, fs.FileMode) error

func init() {
	ReadFile = os.ReadFile
	WriteFile = os.WriteFile
}

// File creates a file, based on the given FileScaffold.
// The parentDirectoryPath is the path to the directory that will contain this file.
// The fullTemplatesDirectoryPath is the path to the directory that contains templates (may not be the full path to the specific
// template directory for this file -- it will be joined with the FileScaffold's TemplatePath property).
// The vars is a map of variables to populate the file and filename with.
func File(file models.FileScaffold, parentDirectoryPath, fullTemplatesDirectoryPath string, vars map[string]string) error {
	// Load template
	fullTemplatePath := file.GetFullTemplatePath(fullTemplatesDirectoryPath)
	templateBytes, templateErr := ReadFile(fullTemplatePath)
	if templateErr != nil {
		return templateErr
	}

	// Populate template with variable values
	templateStr := string(templateBytes)
	populatedTemplate, templatePopulateErr := variable.Populate(templateStr, vars)
	if templatePopulateErr != nil {
		return templatePopulateErr
	}

	// Populate file name with variable values, and create the full file path for the new file
	populatedFileName, fileNameErr := variable.Populate(file.Name, vars)
	if fileNameErr != nil {
		return fileNameErr
	}
	fullFilePath := path.Join(parentDirectoryPath, populatedFileName)

	// Write the file
	writeErr := WriteFile(fullFilePath, []byte(populatedTemplate), 0666)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
