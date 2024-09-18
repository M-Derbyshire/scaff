package models

// Represents a file that contains a number of user-defined commands
type ScaffFile struct {
	Commands []Command `json:"commands"` // The defined commands
	Children []string  `json:"children"` // A list of filepaths to child scaff-files (each path is relative to this scaff-file)
}
