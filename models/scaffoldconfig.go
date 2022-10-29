package models

// ScaffoldConfig represents the config for a number of user-defined commands
type ScaffoldConfig struct {
	TemplateDirectoryPath string            `json:"templateDirectoryPath"` //Path is relative to this config JSON file
	Commands              []ScaffoldCommand `json:"commands"`
}
