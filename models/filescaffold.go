package models

import "path"

// FileScaffold represents a file to be created
type FileScaffold struct {
	Name         string `json:"name"`         // The filename (including extension)
	TemplatePath string `json:"templatePath"` // Path to the file's template (path relative to the template directory)
}

// GetFullTemplatePath returns the full path to the correct template (when given the path to the template directory)
func (fs *FileScaffold) GetFullTemplatePath(templateDirectoryPath string) string {
	return path.Join(templateDirectoryPath, fs.TemplatePath)
}
