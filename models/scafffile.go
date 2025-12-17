package models

import (
	"errors"
	"strings"
)

// Represents a file that contains a number of user-defined commands
type ScaffFile struct {
	Commands []Command `json:"commands"` // The defined commands
	Children []string  `json:"children"` // A list of filepaths to child scaff-files (each path is relative to this scaff-file)
}

func (sf *ScaffFile) GetInvalidJsonError() error {
	msg := "encountered an invalid scaff file. scaff files should contain 2 properties: 'commands' (array of command objects) and 'children' (array of strings)"
	return errors.New(msg)
}

func (sf *ScaffFile) ValidateChildrenArray() error {
	for _, path := range sf.Children {
		trimmedPath := strings.TrimSpace(path)

		if len(trimmedPath) == 0 {
			msg := "encountered an empty file path for a child scaff file"
			return errors.New(msg)
		}
	}

	return nil
}
