package mocks

import (
	"io/fs"
	"slices"
)

// GetFileStat will create and return a mock function for os.Stat
func GetFileStat(correctFilePaths []string) func(string) (fs.FileInfo, error) {
	return func(filepath string) (fs.FileInfo, error) {
		if slices.Contains(correctFilePaths, filepath) {
			return nil, nil
		}

		return nil, fs.ErrNotExist
	}
}
