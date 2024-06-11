package mocks

import "io/fs"

// GetWriteFile will create and return a mock function for os.WriteFile
func GetWriteFile() func(string, []byte, fs.FileMode) error {
	return func(name string, data []byte, perm fs.FileMode) error {
		return nil
	}
}
