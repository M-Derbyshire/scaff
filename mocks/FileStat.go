package mocks

import (
	"io/fs"
)

// GetFileStat will create and return a mock function for os.Stat
func GetFileStat(correctFilePath string) func(string) (fs.FileInfo, error) {
	return func(filepath string) (fs.FileInfo, error) {
		if filepath == correctFilePath {
			return nil, nil
		}

		return nil, fs.ErrNotExist
	}
}
