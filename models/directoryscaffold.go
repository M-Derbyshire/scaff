package models

// DirectoryScaffold represents a directory to be created
type DirectoryScaffold struct {
	Name        string              `json:"name"`
	Files       []FileScaffold      `json:"files"`
	Directories []DirectoryScaffold `json:"directories"`
}
