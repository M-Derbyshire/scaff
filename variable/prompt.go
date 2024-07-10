package variable

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	Stdin          = os.Stdin
	PrintFormatted = fmt.Printf
)

// Prompt prompts the user for the value for a user variable (will always be treated as a string).
// Returns the given value (empty strings are considered valid)
func Prompt(varName string) (string, error) {
	PrintFormatted("variable value required for '%s' > ", varName)

	inReader := bufio.NewReader(Stdin)
	input, err := inReader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")
	return input, nil
}
