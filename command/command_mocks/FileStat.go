package command_mocks

import (
	"io/fs"
)

func GetFileStat(correctFilePath string) func(string) (fs.FileInfo, error) {
	return func(filepath string) (fs.FileInfo, error) {
		if filepath == correctFilePath {
			return nil, nil
		}

		return nil, fs.ErrNotExist
	}
}
