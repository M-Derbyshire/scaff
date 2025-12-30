package e2e_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func switchToInvalidScaffFile() error {
	envPath, err := filepath.Abs(scaffoldRunPath)
	if err != nil {
		return fmt.Errorf("error while getting absolute path: %v", err)
	}

	err = os.Rename(path.Join(envPath, "scaff.json"), path.Join(envPath, "_scaff.json"))
	if err != nil {
		return fmt.Errorf("error while renaming scaff file: %v", err)
	}

	err = os.Rename(path.Join(envPath, "invalid_scaff.json"), path.Join(envPath, "scaff.json"))
	if err != nil {
		return fmt.Errorf("error while renaming invalid scaff file: %v", err)
	}

	return nil
}

func switchBackToValidScaffFile() error {
	envPath, err := filepath.Abs(scaffoldRunPath)
	if err != nil {
		return fmt.Errorf("error while getting absolute path: %v", err)
	}

	err = os.Rename(path.Join(envPath, "scaff.json"), path.Join(envPath, "invalid_scaff.json"))
	if err != nil {
		return fmt.Errorf("error while reseting name of invalid scaff file: %v", err)
	}

	err = os.Rename(path.Join(envPath, "_scaff.json"), path.Join(envPath, "scaff.json"))
	if err != nil {
		return fmt.Errorf("error while resting name of scaff file: %v", err)
	}

	return nil
}

func TestWillPrintErrorWhenNoCommandOrFlagProvided(t *testing.T) {
	output, errOutput, err := runShellCmd(scaffoldRunPath, "./scaff", []string{})
	if err != nil {
		t.Errorf("error while running command: %v", err.Error())
		return
	}

	if len(output) > 0 {
		t.Errorf("expected no output on stdout. got '%v'", output)
	}

	expectedErrText := "please provide the name of the command to process (or use '--help')"

	if !strings.Contains(errOutput, expectedErrText) {
		t.Errorf("expected error output to contain '%v'. got '%v'", expectedErrText, errOutput)
	}
}

func TestWillPrintErrorIfUnableToFindRequestedCommand(t *testing.T) {
	cmdName := "non-existant-cmd"
	expectedErrText := fmt.Sprintf("unable to find the requested command ('%v')", cmdName)

	output, errOutput, err := runShellCmd(scaffoldRunPath, "./scaff", []string{}, cmdName)

	if err != nil {
		t.Errorf("error while running command: %v", err.Error())
		return
	}

	if len(output) > 0 {
		t.Errorf("expected no output on stdout. got '%v'", output)
	}

	if !strings.Contains(errOutput, expectedErrText) {
		t.Errorf("expected error output to contain '%v'. got '%v'", expectedErrText, errOutput)
	}
}

func TestWillPrintErrorWhenEncounteringInvalidScaffFile(t *testing.T) {
	err := switchToInvalidScaffFile()
	if err != nil {
		panic(err)
	}
	defer switchBackToValidScaffFile()

	cmdName := "command1"
	expectedErrText := "encountered an empty file path for a child scaff file"

	output, errOutput, err := runShellCmd(scaffoldRunPath, "./scaff", []string{}, cmdName)

	if err != nil {
		t.Errorf("error while running command: %v", err.Error())
		return
	}

	if len(output) > 0 {
		t.Errorf("expected no output on stdout. got '%v'", output)
	}

	trimmedErrOutput := strings.TrimSpace(errOutput)
	if trimmedErrOutput != expectedErrText {
		t.Errorf("expected error output to be '%v'. got '%v'", expectedErrText, errOutput)
	}
}

func TestWillPrintValidationErrorsIfInvalidCommand(t *testing.T) {
	cmdName := "invalid_command"

	expectedErrs := []string{
		"command objects should have a 'templateDirectoryPath' property that is set to a non-empty value",
		"file scaffold objects should have a 'name' property that is set to a non-empty value",
		"file scaffold objects should have a 'templatePath' property that is set to a non-empty value",
		"directory scaffold objects should have a 'name' property that is set to a non-empty value",
	}

	output, errOutput, err := runShellCmd(scaffoldRunPath, "./scaff", []string{}, cmdName)

	if err != nil {
		t.Errorf("error while running command: %v", err.Error())
		return
	}

	if len(output) > 0 {
		t.Errorf("expected no output on stdout. got '%v'", output)
	}

	// Confirm the right errors were displayed
	allErrorsText := strings.TrimSuffix(errOutput, "\n")
	resultErrTexts := strings.Split(allErrorsText, "\n")

	if len(resultErrTexts) != len(expectedErrs) {
		t.Errorf("expected %d errors. got %d: %s", len(expectedErrs), len(resultErrTexts), allErrorsText)
	}

	for _, expectedErrText := range expectedErrs {
		if !slices.Contains(resultErrTexts, expectedErrText) {
			t.Errorf("expected error output to contain '%s'. got '%s'", expectedErrText, allErrorsText)
		}
	}
}
