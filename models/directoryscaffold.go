package models

type DirectoryScaffold struct {
	Name        string              `json:"name"`
	Files       []FileScaffold      `json:"files"`
	Directories []DirectoryScaffold `json:"directories"`
}
