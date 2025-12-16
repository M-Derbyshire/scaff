package e2e_test

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// Copies over the expected files/dirs to the correct environment directory
func setupPreexistingPathsEnvironment(t *testing.T) {
	sourceDirPath := "./expected/preexistingPaths/*"

	_, _, err := runShellCmd(".", "cp", []string{}, "-r", sourceDirPath, scaffoldRunPath)
	if err != nil {
		t.Errorf("error while moving pre-existing paths into environment: %v", err.Error())
	}
}

func TestWillNotMakeChangesIfFilesOrDirectoriesAlreadyExist(t *testing.T) {
	e2eScaffoldBeforeEach(t)
	setupPreexistingPathsEnvironment(t)

	workingDir, err := filepath.Abs(filepath.Join("environment", "child_dir", "grandchild_dir"))
	if err != nil {
		t.Errorf("error while getting environment directory path: %v", err)
	}

	// This command will try to create files/directories that already exist in the environment
	// directory. It will also try to create some that don't (those shouldn't get created, and
	// any existing files should not be modified)
	commandName := "preexistingPaths"

	expectedErrTexts := []string{
		fmt.Sprintf("path already exists: %s", filepath.Join(workingDir, "existing_file_1.txt")),
		fmt.Sprintf("path already exists: %s", filepath.Join(workingDir, "existing_dir_1")),
		fmt.Sprintf("path already exists: %s", filepath.Join(workingDir, "existing_dir_2")),
	}

	err = runScaffoldCommand(commandName, []string{}, "var1=val1", "var2=val2", "var3=val3", "file=file")
	if err == nil {
		t.Errorf("expected error while running command. got none")
		return
	}

	// Confirm the right errors were displayed
	allErrorsText := strings.TrimSuffix(err.Error(), "\n")
	resultErrTexts := strings.Split(allErrorsText, "\n")

	if len(resultErrTexts) != len(expectedErrTexts) {
		t.Errorf("expected %d errors. got %d: %s", len(expectedErrTexts), len(resultErrTexts), allErrorsText)
	}

	for _, expectedErrText := range expectedErrTexts {
		if !slices.Contains(resultErrTexts, expectedErrText) {
			t.Errorf("expected error output to contain '%s'. got '%s'", expectedErrText, allErrorsText)
		}
	}

	// Now confirm nothing was changed in the environment
	diffs, err := diffScaffoldCommand(commandName)
	if err != nil {
		t.Errorf("error while diffing results of scaff command: %v", err.Error())
		return
	}

	for _, diff := range diffs {
		t.Errorf(diff)
	}
}
