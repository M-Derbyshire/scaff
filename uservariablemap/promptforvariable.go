package uservariablemap

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptForVariable prompts the user for the value to a variable (will always be treated as a string).
// Returns the given value (empty strings are considered valid)
func PromptForVariable(varName string) (string, error) {

	fmt.Printf("Variable value required for '%s' > ", varName)

	inReader := bufio.NewReader(os.Stdin)
	input, err := inReader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")
	return input, nil
}
