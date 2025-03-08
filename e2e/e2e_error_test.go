package e2e_test

import (
	"fmt"
	"strings"
	"testing"
)

func TestWillPrintErrorWhenNoCommandOrFlagProvided(t *testing.T) {
	output, errOutput, err := runShellCmd("..", "go", []string{}, "run", ".")
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

	output, errOutput, err := runShellCmd("..", "go", []string{}, "run", ".", cmdName)

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
