package models

type ScaffoldConfig struct {
	TemplateDirectoryPath string            `json:"templateDirectoryPath"` //Path is relative to this config JSON file
	Commands              []ScaffoldCommand `json:"commands"`
}
