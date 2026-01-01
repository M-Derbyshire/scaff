package models

import (
	"strings"

	"github.com/M-Derbyshire/scaff/customerrors"
)

// Command represents a user-defined command that can be executed
type Command struct {
	Name                  string              `json:"name"`
	TemplateDirectoryPath string              `json:"templateDirectoryPath"` // This path is relative to the containing scaff-file (or child file)
	Files                 []FileScaffold      `json:"files"`
	Directories           []DirectoryScaffold `json:"directories"`
}

// Validate validates the properties in the Command, and returns any validation errors
// The absoluteTemplateDirPath is the root template directory for the command
func (c *Command) Validate(absoluteTemplateDirPath string) []customerrors.ValidationError {
	errs := []customerrors.ValidationError{}

	trimmedTemplatePath := strings.TrimSpace(c.TemplateDirectoryPath)

	if len(trimmedTemplatePath) == 0 {
		newErr := customerrors.ValidationError{
			Message: "command objects should have a 'templateDirectoryPath' property that is set to a non-empty value",
		}

		errs = append(errs, newErr)
	}

	for _, file := range c.Files {
		fileErrs := file.Validate(absoluteTemplateDirPath)
		errs = append(errs, fileErrs...)
	}

	for _, directory := range c.Directories {
		dirErrs := directory.Validate(absoluteTemplateDirPath)
		errs = append(errs, dirErrs...)
	}

	return errs
}
