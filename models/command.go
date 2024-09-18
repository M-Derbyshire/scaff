package models

// Represents a user-defined command that can be executed
type Command struct {
	Name                  string              `json:"name"`
	TemplateDirectoryPath string              `json:"templateDirectoryPath"` // This path is relative to the containing scaff-file (or child file)
	Files                 []FileScaffold      `json:"files"`
	Directories           []DirectoryScaffold `json:"directories"`
}
