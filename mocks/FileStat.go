package mocks

import (
	"io/fs"
)

// GetFileStat will create and return a mock function for os.Stat
func GetFileStat(correctFiles []MockFileInfo) func(string) (fs.FileInfo, error) {
	return func(filepath string) (fs.FileInfo, error) {
		for _, fileInfo := range correctFiles {
			if fileInfo.FilePath == filepath {
				return fileInfo, nil
			}
		}

		return nil, fs.ErrNotExist
	}
}

// CreateMockInfo will create a MockFileInfo struct
func CreateMockInfo(path string, isDir bool) MockFileInfo {
	return MockFileInfo{
		FilePath:   path,
		isDirValue: isDir,
	}
}

// Acts as a FileInfo struct, but allows us to modify the result from the IsDir method
type MockFileInfo struct {
	fs.FileInfo
	FilePath   string
	isDirValue bool
}

func (mfi MockFileInfo) IsDir() bool {
	return mfi.isDirValue
}
