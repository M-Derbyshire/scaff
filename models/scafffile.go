package models

// Represents a file that contains a number of user-defined commands
type ScaffFile struct {
	TemplateDirectoryPath string    `json:"templateDirectoryPath"` //Path is relative to this JSON file
	Commands              []Command `json:"commands"`
}
