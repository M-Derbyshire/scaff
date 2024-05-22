package command_mocks

// GetMockReadFile will create and return a mock function for the command package's ReadFile variable
func GetReadFile(mockScaffoldFileContents []byte) func(string) ([]byte, error) {
	return func(filePath string) ([]byte, error) {
		return mockScaffoldFileContents, nil
	}
}
