package models

type ScaffoldCommand struct {
	Name        string              `json:"name"`
	Directories []DirectoryScaffold `json:"directories"`
}
