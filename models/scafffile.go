package models

// Represents a file that contains a number of user-defined commands
type ScaffFile struct {
	Commands []Command `json:"commands"`
}
