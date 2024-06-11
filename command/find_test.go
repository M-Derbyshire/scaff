package command_test

import (
	"encoding/json"
	"errors"
	"io/fs"
	"path"
	"strings"
	"testing"

	"github.com/M-Derbyshire/scaff/command"
	"github.com/M-Derbyshire/scaff/mocks"
	"github.com/M-Derbyshire/scaff/models"
)

var (
	commandFileNameAndExt string                 = "scaffconfig.json"
	commandToFind         models.ScaffoldCommand = models.ScaffoldCommand{
		Name:        "MyTestingCommand123",
		Files:       []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{},
	}
)

var (
	scaffoldConfig models.ScaffoldConfig = models.ScaffoldConfig{
		Commands:              []models.ScaffoldCommand{commandToFind},
		TemplateDirectoryPath: "/my_templates_1/my_templates_2",
	}
	mockScaffoldFileContents, _ = json.Marshal(scaffoldConfig)
)

// setup runs any setup code that is generic across all tests for the find func
func findBeforeEach() {
	command.CurrentOS = "windows"
	command.ReadFile = mocks.GetReadFile(mockScaffoldFileContents)
	command.FileStat = mocks.GetFileStat(path.Join("C:/", commandFileNameAndExt))
}

func TestFindWillFindACommandInAnyDirectoryInThePath(t *testing.T) {
	findBeforeEach()

	commandName := commandToFind.Name
	fullDirPath := "C:/test1/test2/test3/test4" // The path we'll start at
	dirPathsToTest := []string{
		"C:/test1/test2/test3/test4",
		"C:/test1/test2/test3",
		"C:/test1/test2",
		"C:/test1",
		"C:/",
	}

	// Find should traverse up the given path (checking each parent directory).
	// This ensures it will find a file in any directory in the given path
	for _, dirPath := range dirPathsToTest {
		command.FileStat = mocks.GetFileStat(path.Join(dirPath, commandFileNameAndExt))

		foundCommand, _, _, err := command.Find(commandName, commandFileNameAndExt, fullDirPath)

		if err != nil {
			t.Errorf("err should have been nil. Got %s", err.Error())
		}

		if foundCommand.Name != commandToFind.Name {
			t.Errorf("foundCommand.Name should have been '%s'. Got '%s'", commandToFind.Name, foundCommand.Name)
		}
	}
}

func TestFindWillConstructTheCorrectTemplatePath(t *testing.T) {
	findBeforeEach()

	commandName := commandToFind.Name
	dirPathsToTest := []string{
		"C:/test1/test2/test3/test4",
		"C:/test1/test2/test3",
		"C:/test1/test2",
		"C:/test1",
		"C:/",
	}

	for _, dirPath := range dirPathsToTest {
		command.FileStat = mocks.GetFileStat(path.Join(dirPath, commandFileNameAndExt))

		_, templateDirPath, _, _ := command.Find(commandName, commandFileNameAndExt, dirPath)

		expectedTemplateDirPath := path.Join(dirPath, scaffoldConfig.TemplateDirectoryPath)

		if templateDirPath != expectedTemplateDirPath {
			t.Errorf("expected template directory path to be '%s'. Got '%s'", expectedTemplateDirPath, templateDirPath)
		}
	}
}

func TestFindWillReturnIsFoundAsTrueIfCommandFound(t *testing.T) {
	findBeforeEach()

	_, _, isFound, _ := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")

	if !isFound {
		t.Error("expected true for isFound, but recieved false")
	}
}

func TestFindWillReturnIsFoundAsFalseIfCommandNotFound(t *testing.T) {
	findBeforeEach()

	_, _, isFound, _ := command.Find("nonExistantCommand", commandFileNameAndExt, "C:/")

	if isFound {
		t.Error("expected false for isFound, but recieved true")
	}
}

func TestFindWillNotReturnErrorIfCommandNotFound(t *testing.T) {
	findBeforeEach()

	_, _, _, err := command.Find("nonExistantCommand", commandFileNameAndExt, "C:/")

	if err != nil {
		t.Errorf("expected nil for error, but recieved '%s'", err.Error())
	}
}

func TestFindWillReturnErrorIfUnableToReadFile(t *testing.T) {
	findBeforeEach()

	expectedErrText := "my test error 123"

	command.ReadFile = func(filePath string) ([]byte, error) {
		return nil, errors.New(expectedErrText)
	}

	_, _, _, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")

	if err == nil {
		t.Error("expected error when reading file, but recieved nil")
	}

	if err.Error() != expectedErrText {
		t.Errorf("expected error text to be '%s'. Got '%s'", expectedErrText, err.Error())
	}
}

func TestFindWillReturnErrorIfUnableToUnmarshalJson(t *testing.T) {
	findBeforeEach()

	command.ReadFile = func(filePath string) ([]byte, error) {
		return []byte("not valid json"), nil
	}

	_, _, _, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")

	if err == nil {
		t.Error("expected error when unmarshalling json, but recieved nil")
	}
}

func TestFindWontAddSlashToStartOfPathWhenOnWindows(t *testing.T) {
	findBeforeEach()

	givenFilePath := ""
	command.FileStat = func(filepath string) (fs.FileInfo, error) {
		givenFilePath = filepath
		return nil, nil
	}

	originalPath := "C:/test"
	command.Find(commandToFind.Name, commandFileNameAndExt, originalPath)

	if strings.HasPrefix(givenFilePath, "/") {
		t.Errorf("expected file path not to begin with '/', but got '%s'", givenFilePath)
	}
}

func TestFindWillAddSlashToStartOfPathWhenNotOnWindows(t *testing.T) {
	findBeforeEach()

	givenFilePath := ""
	command.FileStat = func(filepath string) (fs.FileInfo, error) {
		givenFilePath = filepath
		return nil, nil
	}

	command.CurrentOS = "not-windows"

	originalPath := "test/test2"
	command.Find(commandToFind.Name, commandFileNameAndExt, originalPath)

	if !strings.HasPrefix(givenFilePath, "/") {
		t.Errorf("expected file path to begin with '/', but got '%s'", givenFilePath)
	}
}
