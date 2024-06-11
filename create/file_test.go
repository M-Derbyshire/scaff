package create_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/M-Derbyshire/scaff/create"
	"github.com/M-Derbyshire/scaff/mocks"
	"github.com/M-Derbyshire/scaff/models"
)

var mockTemplateData = []byte("my mock template file.")
var mockFileScaffold = models.FileScaffold{
	Name:         "MyNewFile",
	TemplatePath: "MyTemplate.txt",
}

// setup runs any setup code that is generic across all tests for the file func
func findBeforeEach() {
	create.ReadFile = mocks.GetReadFile(mockTemplateData)
	create.WriteFile = mocks.GetWriteFile()
}

func TestWillLoadTheTemplateFileWithTheCorrectFilePath(t *testing.T) {
	findBeforeEach()

	fileScaffold := models.FileScaffold{
		Name:         "MyNewFile",
		TemplatePath: "MyInnerTemplateDir/MyTemplate.txt",
	}
	templatePath := "C:/myTemplatePath"
	expectedFilePath := "C:/myTemplatePath/MyInnerTemplateDir/MyTemplate.txt"

	// We'll capture the path that ReadFile was called with
	var resultFilePath string
	create.ReadFile = func(filePath string) ([]byte, error) {
		resultFilePath = filePath
		return nil, errors.New("test error")
	}

	create.File(fileScaffold, "C:/", templatePath, map[string]string{})

	if resultFilePath != expectedFilePath {
		t.Errorf("file path should have been '%s'. Got '%s'", expectedFilePath, resultFilePath)
	}
}

func TestWillReturnErrorFromReadFileWhenReadingTheTemplateFile(t *testing.T) {
	findBeforeEach()

	expectedErrorText := "my test error 123"
	create.ReadFile = func(filePath string) ([]byte, error) {
		return nil, errors.New(expectedErrorText)
	}

	err := create.File(mockFileScaffold, "", "", map[string]string{})

	if err == nil {
		t.Errorf("expected an error when reading template file, but got nil")
	}

	errorText := err.Error()
	if errorText != expectedErrorText {
		t.Errorf("expected template read error to be '%s'. Got '%s'", expectedErrorText, errorText)
	}
}

func TestWillReturnErrorFromWriteFileWhenCreatingFile(t *testing.T) {
	findBeforeEach()

	expectedErrorText := "my test error 123"
	create.WriteFile = func(name string, data []byte, perm fs.FileMode) error {
		return errors.New(expectedErrorText)
	}

	err := create.File(mockFileScaffold, "", "", map[string]string{})

	if err == nil {
		t.Errorf("expected an error when creating file, but got nil")
	}

	errorText := err.Error()
	if errorText != expectedErrorText {
		t.Errorf("expected file create error to be '%s'. Got '%s'", expectedErrorText, errorText)
	}
}

func TestWillNotReturnErrorIfFileCreated(t *testing.T) {
	findBeforeEach()

	err := create.File(mockFileScaffold, "C:/", "test", map[string]string{})

	if err != nil {
		t.Errorf("expected nil for error when creating file. Got '%s'", err.Error())
	}
}

func TestWillCreateFileWithTheCorrectFilePath(t *testing.T) {
	findBeforeEach()

	var resultFilePath string
	create.WriteFile = func(name string, data []byte, perm fs.FileMode) error {
		resultFilePath = name
		return nil
	}

	fileScaffold := models.FileScaffold{
		Name:         "My-{: var1 :}-New-{: var2 :}-File",
		TemplatePath: "MyTemplate.txt",
	}
	parentDir := "C:/myDir"
	vars := map[string]string{
		"var1": "value1",
		"var2": "value2",
	}

	create.File(fileScaffold, parentDir, "/", vars)

	expectedFilePath := "C:/myDir/My-value1-New-value2-File"

	if resultFilePath != expectedFilePath {
		t.Errorf("expected created file path to be '%s'. Got '%s'", expectedFilePath, resultFilePath)
	}
}

func TestWillCreateFileWithTheCorrectContents(t *testing.T) {
	findBeforeEach()

	fileContents := "{: var1 :} - {: var2 :}"
	create.ReadFile = mocks.GetReadFile([]byte(fileContents)) // Mock the template file

	var resultFileContents string
	create.WriteFile = func(name string, data []byte, perm fs.FileMode) error {
		resultFileContents = string(data)
		return nil
	}

	vars := map[string]string{
		"var1": "value1",
		"var2": "value2",
	}

	create.File(mockFileScaffold, "C:/", "/", vars)

	expectedFileContents := "value1 - value2"

	if resultFileContents != expectedFileContents {
		t.Errorf("expected created file contents to be '%s'. Got '%s'", expectedFileContents, resultFileContents)
	}
}

func TestWillCreateFileWithTheCorrectPermissions(t *testing.T) {
	findBeforeEach()

	var resultFilePerms fs.FileMode
	create.WriteFile = func(name string, data []byte, perm fs.FileMode) error {
		resultFilePerms = perm
		return nil
	}

	create.File(mockFileScaffold, "C:/", "/", map[string]string{})

	var expectedFilePerms fs.FileMode = 0666

	if resultFilePerms != expectedFilePerms {
		t.Errorf("expected created file permissions to be %#o. Got %#o", expectedFilePerms, resultFilePerms)
	}
}
