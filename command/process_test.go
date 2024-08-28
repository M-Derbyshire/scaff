package command_test

import (
	"errors"
	"testing"

	"github.com/M-Derbyshire/scaff/command"
	"github.com/M-Derbyshire/scaff/models"
)

// This is used to store calls to the CreateFile or CreateDirectory mocks
type mockCreateCall struct {
	File         models.FileScaffold
	Directory    models.DirectoryScaffold
	ParentDir    string
	TemplatesDir string
	Vars         map[string]string
}

// Command for testing process
var processTestCommand models.Command = models.Command{
	Name: "testId1",
	Files: []models.FileScaffold{
		{
			Name:         "file1.txt",
			TemplatePath: "Template1.txt",
		},
		{
			Name:         "file2.txt",
			TemplatePath: "Template2.txt",
		},
		{
			Name:         "file3.txt",
			TemplatePath: "Template3.txt",
		},
	},
	Directories: []models.DirectoryScaffold{
		{
			Name:        "dir1",
			Files:       []models.FileScaffold{},
			Directories: []models.DirectoryScaffold{},
		},
		{
			Name:        "dir2",
			Files:       []models.FileScaffold{},
			Directories: []models.DirectoryScaffold{},
		},
		{
			Name:        "dir3",
			Files:       []models.FileScaffold{},
			Directories: []models.DirectoryScaffold{},
		},
	},
}

// setup anything that's required before each individual test
func processBeforeEach() {
	command.CreateFile = func(_ models.FileScaffold, _, _ string, _ map[string]string) error {
		return nil
	}

	command.CreateDirectory = func(_ models.DirectoryScaffold, _, _ string, _ map[string]string) error {
		return nil
	}
}

func TestProcessWillCreateAllFilesInTheCommand(t *testing.T) {
	processBeforeEach()

	expectedParentDir := "C:/test123/my-work-dir"
	expectedTemplateDir := "C:/test123/my-work-dir/templates"
	expectedVars := map[string]string{}
	varKey := "testKey"
	expectedVars[varKey] = "testVal123"

	// We're going to mock the CreateFile function, and have it populate a slice of calls with the data it's called with
	mockCalls := make([]mockCreateCall, 0, len(processTestCommand.Files))

	command.CreateFile = func(file models.FileScaffold, parentDirectoryPath, templatesDirectoryPath string, vars map[string]string) error {
		mockCalls = append(mockCalls, mockCreateCall{
			File:         file,
			ParentDir:    parentDirectoryPath,
			TemplatesDir: templatesDirectoryPath,
			Vars:         vars,
		})

		return nil
	}

	// Now run the function, and test the CreateFile is called with the correct parameters
	givenErr := command.Process(processTestCommand, expectedParentDir, expectedTemplateDir, expectedVars)
	if givenErr != nil {
		t.Errorf("expected Process to return no error. Got '%s'", givenErr.Error())
	}

	if len(mockCalls) != len(processTestCommand.Files) {
		t.Errorf(
			"expected Process to make the same amount of calls to CreateFile as the amount of files (%d). Got %d",
			len(processTestCommand.Files),
			len(mockCalls),
		)

		return
	}

	for i, file := range processTestCommand.Files {
		call := mockCalls[i]

		if call.File.Name != file.Name {
			t.Errorf("expected Process to call CreateFile with the given file ('%s'). Got '%s'", file.Name, call.File.Name)
		}

		if call.ParentDir != expectedParentDir {
			t.Errorf(
				"expected Process to call CreateFile with the given parent directory ('%s'). Got '%s'",
				expectedParentDir,
				call.ParentDir,
			)
		}

		if call.TemplatesDir != expectedTemplateDir {
			t.Errorf(
				"expected Process to call CreateFile with the given template directory ('%s'). Got '%s'",
				expectedTemplateDir,
				call.TemplatesDir,
			)
		}

		if call.Vars[varKey] != expectedVars[varKey] {
			t.Errorf(
				"expected Process to call CreateFile with the correct variable map (with a variable value '%s'). Got '%s'",
				expectedVars[varKey],
				call.Vars[varKey],
			)
		}
	}
}

func TestProcessWillReturnAnErrorIfThereWasAnErrorCreatingAFile(t *testing.T) {
	processBeforeEach()

	expectedErrorText := "my test file error"
	command.CreateFile = func(_ models.FileScaffold, _, _ string, _ map[string]string) error {
		return errors.New(expectedErrorText)
	}

	givenErr := command.Process(processTestCommand, "/", "/templates", map[string]string{})

	if givenErr == nil {
		t.Errorf("expected an error from Process. Got nil")
		return
	}

	givenErrText := givenErr.Error()
	if givenErrText != expectedErrorText {
		t.Errorf("Expected error text from Process to be '%s'. Got '%s'", expectedErrorText, givenErrText)
	}
}

func TestProcessWillCreateAllDirectoriesInTheCommand(t *testing.T) {
	processBeforeEach()

	expectedParentDir := "C:/test123/my-work-dir"
	expectedTemplateDir := "C:/test123/my-work-dir/templates"
	expectedVars := map[string]string{}
	varKey := "testKey"
	expectedVars[varKey] = "testVal123"

	// We're going to mock the CreateDirectory function, and have it populate a slice of calls with the data it's called with
	mockCalls := make([]mockCreateCall, 0, len(processTestCommand.Directories))

	command.CreateDirectory = func(directory models.DirectoryScaffold, parentDirectoryPath, templatesDirectoryPath string, vars map[string]string) error {
		mockCalls = append(mockCalls, mockCreateCall{
			Directory:    directory,
			ParentDir:    parentDirectoryPath,
			TemplatesDir: templatesDirectoryPath,
			Vars:         vars,
		})

		return nil
	}

	// Now run the function, and test the CreateDirectory is called with the correct parameters
	givenErr := command.Process(processTestCommand, expectedParentDir, expectedTemplateDir, expectedVars)
	if givenErr != nil {
		t.Errorf("expected Process to return no error. Got '%s'", givenErr.Error())
	}

	if len(mockCalls) != len(processTestCommand.Directories) {
		t.Errorf(
			"expected Process to make the same amount of calls to CreateDirectory as the amount of directories (%d). Got %d",
			len(processTestCommand.Directories),
			len(mockCalls),
		)

		return
	}

	for i, dir := range processTestCommand.Directories {
		call := mockCalls[i]

		if call.Directory.Name != dir.Name {
			t.Errorf(
				"expected Process to call CreateDirectory with the given directory ('%s'). Got '%s'",
				dir.Name,
				call.Directory.Name,
			)
		}

		if call.ParentDir != expectedParentDir {
			t.Errorf(
				"expected Process to call CreateDirectory with the given parent directory ('%s'). Got '%s'",
				expectedParentDir,
				call.ParentDir,
			)
		}

		if call.TemplatesDir != expectedTemplateDir {
			t.Errorf(
				"expected Process to call CreateDirectory with the given template directory ('%s'). Got '%s'",
				expectedTemplateDir,
				call.TemplatesDir,
			)
		}

		if call.Vars[varKey] != expectedVars[varKey] {
			t.Errorf(
				"expected Process to call CreateDirectory with the correct variable map (with a variable value '%s'). Got '%s'",
				expectedVars[varKey],
				call.Vars[varKey],
			)
		}
	}
}

func TestProcessWillReturnAnErrorIfThereWasAnErrorCreatingADirectory(t *testing.T) {
	processBeforeEach()

	expectedErrorText := "my test dir error"
	command.CreateDirectory = func(_ models.DirectoryScaffold, _, _ string, _ map[string]string) error {
		return errors.New(expectedErrorText)
	}

	givenErr := command.Process(processTestCommand, "/", "/templates", map[string]string{})

	if givenErr == nil {
		t.Errorf("expected an error from Process. Got nil")
		return
	}

	givenErrText := givenErr.Error()
	if givenErrText != expectedErrorText {
		t.Errorf("Expected error text from Process to be '%s'. Got '%s'", expectedErrorText, givenErrText)
	}
}
