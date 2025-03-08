package e2e_test

import (
	"testing"
)

func e2eScaffoldBeforeEach(t *testing.T) {
	_, _, err := runShellCmd(".", "sh", []string{}, "reset.sh")

	if err != nil {
		t.Errorf("error while resetting environment: %v", err.Error())
	}
}

func TestWillCreateScaffoldFromCommandWithGivenVarFlags(t *testing.T) {
	e2eScaffoldBeforeEach(t)

	commandName := "command1"

	err := runScaffoldCommand(commandName, []string{}, "var1=val1", "var2=val2", "var3=val3")
	if err != nil {
		t.Errorf("error while running scaff command: %v", err.Error())
		return
	}

	diffs, err := diffScaffoldCommand(commandName)
	if err != nil {
		t.Errorf("error while diffing results of scaff command: %v", err.Error())
		return
	}

	for _, diff := range diffs {
		t.Errorf(diff)
	}
}

func TestScaffoldFromCommandWillPromptForVars(t *testing.T) {
	e2eScaffoldBeforeEach(t)

	commandName := "command2"
	userVariableValues := []string{"val1", "val2", "val3"}

	err := runScaffoldCommand(commandName, userVariableValues)
	if err != nil {
		t.Errorf("error while running scaff command: %v", err.Error())
		return
	}

	diffs, err := diffScaffoldCommand(commandName)
	if err != nil {
		t.Errorf("error while diffing results of scaff command: %v", err.Error())
		return
	}

	for _, diff := range diffs {
		t.Errorf(diff)
	}
}
