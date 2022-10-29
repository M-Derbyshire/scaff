package models

// ScaffoldCommand represents a user-defined command that can be executed
type ScaffoldCommand struct {
	Name        string              `json:"name"`
	Directories []DirectoryScaffold `json:"directories"`
}
