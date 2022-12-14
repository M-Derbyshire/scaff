package commandprocessing

import (
	"os"
	"path"

	"github.com/M-Derbyshire/scaff/models"
	"github.com/M-Derbyshire/scaff/stringprocessing"
)

// CreateFile creates a file, based on the given FileScaffold.
// The parentDirectoryPath is the path to the directory that will contain this file.
// The templatesDirectoryPath is the path to the directory that contains templates (may not be the full path to the specific
// template directory for this file -- it will be joined with the FileScaffold's TemplatePath property).
// The vars is a map of variables to populate the file and filename with.
func CreateFile(file models.FileScaffold, parentDirectoryPath, templatesDirectoryPath string, vars map[string]string) error {

	// Load template
	fullTemplatePath := file.GetFullTemplatePath(templatesDirectoryPath)
	templateBytes, templateErr := os.ReadFile(fullTemplatePath)
	if templateErr != nil {
		return templateErr
	}

	// Populate template with variable values
	templateStr := string(templateBytes)
	populatedTemplate, templatePopulateErr := stringprocessing.PopulateVariablesInString(templateStr, vars)
	if templatePopulateErr != nil {
		return templatePopulateErr
	}

	// Populate file name with variable values, and create the full file path for the new file
	populatedFileName, fileNameErr := stringprocessing.PopulateVariablesInString(file.Name, vars)
	if fileNameErr != nil {
		return fileNameErr
	}
	fullFilePath := path.Join(parentDirectoryPath, populatedFileName)

	// Write the file
	writeErr := os.WriteFile(fullFilePath, []byte(populatedTemplate), 0222)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
