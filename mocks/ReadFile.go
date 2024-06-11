package mocks

// GetReadFile will create and return a mock function for os.ReadFile
func GetReadFile(mockFileContents []byte) func(string) ([]byte, error) {
	return func(filePath string) ([]byte, error) {
		return mockFileContents, nil
	}
}
