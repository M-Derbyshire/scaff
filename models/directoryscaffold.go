package models

import (
	"strings"

	"github.com/M-Derbyshire/scaff/customerrors"
)

// DirectoryScaffold represents a directory to be created
type DirectoryScaffold struct {
	Name        string              `json:"name"`
	Files       []FileScaffold      `json:"files"`
	Directories []DirectoryScaffold `json:"directories"`
}

// Validates the properties in the DirectoryScaffold, and returns any validation errors
// The templateDirectoryPath is the root template directory for the command
func (ds *DirectoryScaffold) Validate(templateDirectoryPath string) []customerrors.ValidationError {
	errs := []customerrors.ValidationError{}

	trimmedName := strings.TrimSpace(ds.Name)
	if len(trimmedName) == 0 {
		newErr := customerrors.ValidationError{
			Message: "directory scaffold objects should have a 'name' property that is set to a non-empty value",
		}
		errs = append(errs, newErr)
	}

	for _, file := range ds.Files {
		fileErrs := file.Validate(templateDirectoryPath)
		errs = append(errs, fileErrs...)
	}

	for _, directory := range ds.Directories {
		dirErrs := directory.Validate(templateDirectoryPath)
		errs = append(errs, dirErrs...)
	}

	return errs
}
