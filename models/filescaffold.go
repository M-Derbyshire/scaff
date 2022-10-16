package models

import "path"

type FileScaffold struct {
	Name         string `json:"name"`         // The filename (including extension)
	TemplatePath string `json:"templatePath"` // Path to the file's template (path relative to the command's template directory)
}

func (fs *FileScaffold) GetFullTemplatePath(templateDirectoryPath string) string {
	return path.Join(templateDirectoryPath, fs.TemplatePath)
}
