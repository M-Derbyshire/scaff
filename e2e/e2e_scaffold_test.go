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

func TestWillCreateScaffoldFromCommandInChildScaffFile(t *testing.T) {
	e2eScaffoldBeforeEach(t)

	commandName := "childCommand1"

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

// The rest of the test-commands do not require any values from the user, so can be ran by the same code
var scaffoldTestsWithoutUserVarsTable = []struct {
	Name    string
	Command string
}{
	{"Second Child Command", "childCommand2"},
	{"Command in Parent Directory Scaff File", "commandInParentDir"},
	{"Command in Parent Directory Child Scaff File", "childCommandInParentDir"},
	{"Command in Grandparent Directory Scaff File", "commandInGrandparentDir"},
	{"Command in Grandparent Directory Child Scaff File", "childCommandInGrandparentDir"},
}

func TestScaffoldCommandsWithoutUserVars(t *testing.T) {
	for _, tt := range scaffoldTestsWithoutUserVarsTable {
		t.Run(tt.Name, func(t *testing.T) {
			e2eScaffoldBeforeEach(t)

			err := runScaffoldCommand(tt.Command, []string{})
			if err != nil {
				t.Errorf("error while running scaff command: %v", err.Error())
				return
			}

			diffs, err := diffScaffoldCommand(tt.Command)
			if err != nil {
				t.Errorf("error while diffing results of scaff command: %v", err.Error())
				return
			}

			for _, diff := range diffs {
				t.Errorf(diff)
			}
		})
	}
}
