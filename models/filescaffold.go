package models

import (
	"fmt"
	"path"
	"strings"

	"github.com/M-Derbyshire/scaff/customerrors"
)

// FileScaffold represents a file to be created
type FileScaffold struct {
	Name         string `json:"name"`         // The filename (including extension)
	TemplatePath string `json:"templatePath"` // Path to the file's template (path relative to the template directory)
}

// GetFullTemplatePath returns the full path to the correct template (when given the path to the template directory)
func (fs *FileScaffold) GetFullTemplatePath(templateDirectoryPath string) string {
	return path.Join(templateDirectoryPath, fs.TemplatePath)
}

func (fs *FileScaffold) Validate(templateDirectoryPath string) []customerrors.ValidationError {
	errs := []customerrors.ValidationError{}

	trimmedName := strings.TrimSpace(fs.Name)
	trimmedTemplatePath := strings.TrimSpace(fs.TemplatePath)

	if len(trimmedName) == 0 {
		newErr := customerrors.ValidationError{
			Message: "file scaffold objects should have a 'name' property that is set to a non-empty value",
		}
		errs = append(errs, newErr)
	}

	if len(trimmedTemplatePath) == 0 {
		newErr := customerrors.ValidationError{
			Message: "file scaffold objects should have a 'templatePath' property that is set to a non-empty value",
		}
		errs = append(errs, newErr)
	}

	// We want to confirm that the template file exists
	fullTemplatePath := fs.GetFullTemplatePath(templateDirectoryPath)

	if _, pathErr := FileStat(fullTemplatePath); pathErr != nil {
		newErr := customerrors.ValidationError{
			Message: fmt.Sprintf("unable to locate template file at path: '%s'", fullTemplatePath),
		}

		errs = append(errs, newErr)
	}

	return errs
}
