package e2e_test

import (
	"strings"
	"testing"
)

func TestWillDisplayHelpText(t *testing.T) {
	expectedOutText := `Creates directories and files in the current working directory, based on the structures defined in a scaff.json file (using the given variables).

SCAFF [commandname] [variablename]=[variablevalue]

SCAFF will work its way up the directory-tree, from the current working directory, searching for a scaff.json file that contains the requested command (if multiple commands are found with the same name, the first one in the array is used).
For full instructions on how to structure commands in a scaff.json file, visit https://github.com/M-Derbyshire/scaff

[commandname] - The name of the command (in a scaff.json file) that defines the files/directories to create.
[variablename]=[variablevalue] - Variables that are needed by the requested command can be defined in this format. See the below examples:

var1=myValue
var2="my longer value"

You can provide multiple variables in this way. If a variable is needed, but not provided, SCAFF will prompt you to provide it.

For full instructions on the use of SCAFF, visit https://github.com/M-Derbyshire/scaff`

	output, errOutput, err := runShellCmd(scaffoldRunPath, "./scaff", []string{}, "--help")

	if err != nil {
		t.Errorf("error while running command: %v", err.Error())
	}

	if len(errOutput) > 0 {
		t.Errorf("expected nothing to be output on Stderr. got '%v'", errOutput)
	}

	preppedOutStr := strings.ReplaceAll(output, "\n", "")
	preppedExpectedStr := strings.ReplaceAll(expectedOutText, "\n", "")

	if preppedOutStr != preppedExpectedStr {
		t.Errorf(
			"the help text did not match the expected help text. expected:\n\n*********\n%v\n*********\n\ngot:\n\n*********\n%v\n*********",
			expectedOutText,
			output,
		)
	}
}
