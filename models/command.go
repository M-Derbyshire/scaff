package models

// Represents a user-defined command that can be executed
type Command struct {
	Name        string              `json:"name"`
	Files       []FileScaffold      `json:"files"`
	Directories []DirectoryScaffold `json:"directories"`
}
