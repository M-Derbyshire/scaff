package mocks

import "io/fs"

// GetMkdir will create and return a mock function for os.GetMkdir
func GetMkdir() func(string, fs.FileMode) error {
	return func(path string, perms fs.FileMode) error {
		return nil
	}
}
