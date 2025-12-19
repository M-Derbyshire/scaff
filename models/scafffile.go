package models

import (
	"strings"

	"github.com/M-Derbyshire/scaff/customerrors"
)

// Represents a file that contains a number of user-defined commands
type ScaffFile struct {
	Commands []Command `json:"commands"` // The defined commands
	Children []string  `json:"children"` // A list of filepaths to child scaff-files (each path is relative to this scaff-file)
}

func (sf *ScaffFile) ValidateChildrenArray() error {
	for _, path := range sf.Children {
		trimmedPath := strings.TrimSpace(path)

		if len(trimmedPath) == 0 {
			msg := "encountered an empty file path for a child scaff file"
			return &customerrors.ValidationError{
				Message: msg,
			}
		}
	}

	return nil
}
