package variable

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	// Stdin is an open file, pointing to the standard input
	Stdin = os.Stdin
	// PrintFormatted is used to print a formatted string to standard output
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

	// Remove CRLF (or might just be line feed) from end of value
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	// Remove surrounding quotes
	input = strings.TrimPrefix(input, "\"")
	input = strings.TrimSuffix(input, "\"")

	return input, nil
}
