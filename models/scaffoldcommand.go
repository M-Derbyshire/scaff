package models

// ScaffoldCommand represents a user-defined command that can be executed
type ScaffoldCommand struct {
	Name        string              `json:"name"`
	Files       []FileScaffold      `json:"files"`
	Directories []DirectoryScaffold `json:"directories"`
}
