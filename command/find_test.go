package command_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"testing"

	"github.com/M-Derbyshire/scaff/command"
	"github.com/M-Derbyshire/scaff/customerrors"
	"github.com/M-Derbyshire/scaff/mocks"
	"github.com/M-Derbyshire/scaff/models"
)

var (
	commandFileNameAndExt string         = "scaff.json"
	commandToFind         models.Command = models.Command{
		Name:                  "MyTestingCommand123",
		Files:                 []models.FileScaffold{},
		Directories:           []models.DirectoryScaffold{},
		TemplateDirectoryPath: "/my_templates_1/my_templates_2",
	}
	commandNotToFind models.Command = models.Command{
		Name:                  "MyOtherTestingCommand456",
		Files:                 []models.FileScaffold{},
		Directories:           []models.DirectoryScaffold{},
		TemplateDirectoryPath: "/my_templates_1/",
	}
)

var (
	scaffFile models.ScaffFile = models.ScaffFile{
		Commands: []models.Command{commandNotToFind, commandToFind},
		Children: []string{},
	}
	parentScaffFile models.ScaffFile = models.ScaffFile{
		Commands: []models.Command{},
		Children: []string{
			"/children/child1.json",
			"/children/child2.json",
		},
	}
)

// setup runs any setup code that is generic across all tests for the find func
func findBeforeEach() {
	fileContentsJSON, _ := json.Marshal(scaffFile)
	command.ReadFile = mocks.GetReadFile(fileContentsJSON)

	command.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	command.CurrentOS = "windows"
}

func parentFindBeforeEach() {
	findBeforeEach()

	command.ReadFile = func(filePath string) ([]byte, error) {
		if strings.HasSuffix(filePath, commandFileNameAndExt) { // Parent scaff-file
			parentFileContents, _ := json.Marshal(parentScaffFile)
			return parentFileContents, nil
		} else if strings.HasSuffix(filePath, "children/child1.json") { // Child file without requested command
			wrongChildFileContents, _ := json.Marshal(models.ScaffFile{
				Commands: []models.Command{commandNotToFind},
				Children: []string{},
			})
			return wrongChildFileContents, nil
		} else if strings.HasSuffix(filePath, "children/child2.json") { // Child file with the requested command
			rightChildFileContents, _ := json.Marshal(models.ScaffFile{
				Commands: []models.Command{commandNotToFind, commandToFind},
				Children: []string{},
			})
			return rightChildFileContents, nil
		}

		// Error if we're here
		return nil, fmt.Errorf("An unexpected path was provided to ReadFile: %s", filePath)
	}
}

// ---- Finding in top-level scaff-file ----------------------------------------------

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
		mockFileInfo := mocks.CreateMockInfo(path.Join(dirPath, commandFileNameAndExt), false)
		command.FileStat = mocks.GetFileStat([]mocks.MockFileInfo{mockFileInfo})

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
		mockFileInfo := mocks.CreateMockInfo(path.Join(dirPath, commandFileNameAndExt), false)
		command.FileStat = mocks.GetFileStat([]mocks.MockFileInfo{mockFileInfo})

		_, templateDirPath, _, _ := command.Find(commandName, commandFileNameAndExt, dirPath)

		expectedTemplateDirPath := path.Join(dirPath, commandToFind.TemplateDirectoryPath)

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
		return
	}

	expectedErrText := "encountered a scaff.json file with an invalid structure: 'C:/scaff.json'"
	resultErrText := err.Error()
	if resultErrText != expectedErrText {
		t.Errorf("expected error text to be '%s'. got '%s'", expectedErrText, resultErrText)
	}

	var vErr *customerrors.ValidationError
	if !errors.As(err, &vErr) {
		t.Errorf("expected error type to be ValidationError")
	}
}

func TestFindWillReturnValidationErrorForChildrenArray(t *testing.T) {
	findBeforeEach()

	testFile := models.ScaffFile{
		Commands: []models.Command{commandToFind},
		Children: []string{""},
	}

	fileContentsJSON, _ := json.Marshal(testFile)
	command.ReadFile = mocks.GetReadFile(fileContentsJSON)

	_, _, _, err := command.Find(commandNotToFind.Name, commandFileNameAndExt, "C:/")

	if err == nil {
		t.Error("expected error. got nil")
		return
	}

	expectedErrText := "encountered an empty file path for a child scaff file"
	resultErrText := err.Error()
	if resultErrText != expectedErrText {
		t.Errorf("expected error text to be '%s'. got '%s'", expectedErrText, resultErrText)
	}

	var vErr *customerrors.ValidationError
	if !errors.As(err, &vErr) {
		t.Errorf("expected error type to be ValidationError")
	}
}

func TestFindWillNotReturnValidationErrorIfNoCommandsArray(t *testing.T) {
	findBeforeEach()

	testFile := models.ScaffFile{
		Children: []string{},
	}

	fileContentsJSON, _ := json.Marshal(testFile)
	command.ReadFile = mocks.GetReadFile(fileContentsJSON)

	_, _, _, err := command.Find(commandNotToFind.Name, commandFileNameAndExt, "C:/")

	if err != nil {
		t.Errorf("expected no error. got '%s'", err.Error())
		return
	}
}

func TestFindWillNotReturnValidationErrorIfNoChildrenArray(t *testing.T) {
	findBeforeEach()

	testFile := models.ScaffFile{
		Commands: []models.Command{commandToFind},
	}

	fileContentsJSON, _ := json.Marshal(testFile)
	command.ReadFile = mocks.GetReadFile(fileContentsJSON)

	_, _, _, err := command.Find(commandNotToFind.Name, commandFileNameAndExt, "C:/")

	if err != nil {
		t.Errorf("expected no error. got '%s'", err.Error())
		return
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

// -----------------------------------------------------------------------------------

// ---- Finding in child scaff-files -------------------------------------------------

func TestFindWillSearchThroughChildScaffFilesForCommand(t *testing.T) {
	parentFindBeforeEach()

	foundCommand, _, _, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")
	if err != nil {
		t.Error(err)
		return
	}

	if foundCommand.Name != commandToFind.Name {
		t.Errorf("expected to recieve the correct command from find ('%s'). got '%s'", commandToFind.Name, foundCommand.Name)
	}
}

func TestFindWillReturnErrorIfUnableToFindChildFile(t *testing.T) {
	parentFindBeforeEach()

	childPathToFind := parentScaffFile.Children[0]
	expectedErrorMsg := fmt.Sprintf("unable to locate child scaff file at path: 'C:%s'", childPathToFind)

	command.FileStat = func(filepath string) (fs.FileInfo, error) {
		if strings.HasSuffix(filepath, childPathToFind) {
			return nil, errors.New("my error msg")
		}

		return nil, nil
	}

	_, _, _, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")
	if err == nil {
		t.Error("expected find to return an error, but got none")
		return
	}

	if err.Error() != expectedErrorMsg {
		t.Errorf("expected find to return the correct error message '%s'. got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestFindWillConstructTheCorrectTemplatePathForChildScaffFile(t *testing.T) {
	parentFindBeforeEach()

	expectedTemplatesPath := "C:/my_location/children/my_templates_1/my_templates_2"

	_, resultTemplatesPath, _, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/my_location")
	if err != nil {
		t.Error(err)
		return
	}

	if resultTemplatesPath != expectedTemplatesPath {
		t.Errorf("expected to recieve the correct template path from find ('%s'). got '%s'", expectedTemplatesPath, resultTemplatesPath)
	}
}

func TestFindWillReturnIsFoundAsTrueIfCommandFoundInChildScaffFile(t *testing.T) {
	parentFindBeforeEach()

	_, _, isFound, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")
	if err != nil {
		t.Error(err)
		return
	}

	if !isFound {
		t.Errorf("expected find to return true for isFound value. got false")
	}
}

func TestFindWillReturnIsFoundAsFalseIfCommandNotFoundInChildScaffFile(t *testing.T) {
	parentFindBeforeEach()

	command.ReadFile = func(filePath string) ([]byte, error) {
		if strings.HasSuffix(filePath, commandFileNameAndExt) { // Parent scaff-file
			parentFileContents, _ := json.Marshal(parentScaffFile)
			return parentFileContents, nil
		} else {
			childFileContents, _ := json.Marshal(models.ScaffFile{
				Commands: []models.Command{commandNotToFind},
				Children: []string{},
			})
			return childFileContents, nil
		}
	}

	_, _, isFound, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")
	if err != nil {
		t.Error(err)
		return
	}

	if isFound {
		t.Errorf("expected find to return false for isFound value. got true")
	}
}

func TestFindWillReturnErrorIfUnableToUnmarshalJsonFromChildScaffFile(t *testing.T) {
	parentFindBeforeEach()

	command.ReadFile = func(filePath string) ([]byte, error) {
		if strings.HasSuffix(filePath, commandFileNameAndExt) { // Parent scaff-file
			parentFileContents, _ := json.Marshal(parentScaffFile)
			return parentFileContents, nil
		} else {
			return []byte("not valid json"), nil
		}
	}

	_, _, _, err := command.Find(commandToFind.Name, commandFileNameAndExt, "C:/")
	if err == nil {
		t.Error("expected error when unmarshalling json, but recieved nil")
		return
	}

	expectedErrText := "encountered a scaff.json file with an invalid structure: 'C:/children/child1.json'"
	resultErrText := err.Error()
	if resultErrText != expectedErrText {
		t.Errorf("expected error text to be '%s'. got '%s'", expectedErrText, resultErrText)
	}

	var vErr *customerrors.ValidationError
	if !errors.As(err, &vErr) {
		t.Errorf("expected error type to be ValidationError")
	}
}

func TestFindWillReturnValidationErrorForChildrenArrayInChildScaffFile(t *testing.T) {
	findBeforeEach()

	testFile := models.ScaffFile{
		Commands: []models.Command{commandToFind},
		Children: []string{""},
	}

	fileContentsJSON, _ := json.Marshal(testFile)

	command.ReadFile = func(filePath string) ([]byte, error) {
		if strings.HasSuffix(filePath, commandFileNameAndExt) { // Parent scaff-file
			parentFileContents, _ := json.Marshal(parentScaffFile)
			return parentFileContents, nil
		} else {
			return fileContentsJSON, nil
		}
	}

	_, _, _, err := command.Find(commandNotToFind.Name, commandFileNameAndExt, "C:/")

	if err == nil {
		t.Error("expected error. got nil")
		return
	}

	expectedErrText := "encountered an empty file path for a child scaff file"
	resultErrText := err.Error()
	if resultErrText != expectedErrText {
		t.Errorf("expected error text to be '%s'. got '%s'", expectedErrText, resultErrText)
	}

	var vErr *customerrors.ValidationError
	if !errors.As(err, &vErr) {
		t.Errorf("expected error type to be ValidationError")
	}
}

func TestFindWillNotReturnValidationErrorIfNoCommandsArrayInChildScaffFile(t *testing.T) {
	findBeforeEach()

	testFile := models.ScaffFile{
		Children: []string{},
	}

	fileContentsJSON, _ := json.Marshal(testFile)

	command.ReadFile = func(filePath string) ([]byte, error) {
		if strings.HasSuffix(filePath, commandFileNameAndExt) { // Parent scaff-file
			parentFileContents, _ := json.Marshal(parentScaffFile)
			return parentFileContents, nil
		} else {
			return fileContentsJSON, nil
		}
	}

	_, _, _, err := command.Find(commandNotToFind.Name, commandFileNameAndExt, "C:/")

	if err != nil {
		t.Errorf("expected no error. got '%s'", err.Error())
		return
	}
}

func TestFindWillNotReturnValidationErrorIfNoChildrenArrayInChildScaffFile(t *testing.T) {
	findBeforeEach()

	testFile := models.ScaffFile{
		Commands: []models.Command{commandToFind},
	}

	fileContentsJSON, _ := json.Marshal(testFile)

	command.ReadFile = func(filePath string) ([]byte, error) {
		if strings.HasSuffix(filePath, commandFileNameAndExt) { // Parent scaff-file
			parentFileContents, _ := json.Marshal(parentScaffFile)
			return parentFileContents, nil
		} else {
			return fileContentsJSON, nil
		}
	}

	_, _, _, err := command.Find(commandNotToFind.Name, commandFileNameAndExt, "C:/")

	if err != nil {
		t.Errorf("expected no error. got '%s'", err.Error())
		return
	}
}

// -----------------------------------------------------------------------------------
