package variable_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/M-Derbyshire/scaff/variable"
)

func setupMockStdIn(inputText string) error {
	inputBytes := []byte(inputText)

	// Setup the file to act as the Stdin
	readToFile, writeToFile, err := os.Pipe()
	if err != nil {
		return err
	}

	// Write the mock input to the file
	_, err = writeToFile.Write(inputBytes)
	if err != nil {
		return err
	}
	writeToFile.Close()

	// Set the Stdin
	variable.Stdin = readToFile
	return nil
}

func TestPromptWillReturnEnteredStringFromStdin(t *testing.T) {
	expectedText := "my test input 123"
	userInput := expectedText + "\n"

	err := setupMockStdIn(userInput)
	if err != nil {
		t.Fatal(err)
	}

	result, err := variable.Prompt("a")
	if err != nil {
		t.Errorf("expected to recieve no error. Got %e", err)
	}

	if result != expectedText {
		t.Errorf("expected result to be '%s'. Got '%s'", expectedText, result)
	}
}

func TestPromptWillPrintPromptTextWithGivenVariableName(t *testing.T) {
	originalPrintF := variable.PrintFormatted
	defer func() { variable.PrintFormatted = originalPrintF }()

	// Mock the PrintFormatted
	var resultFormatStr string
	var resultValues []any
	variable.PrintFormatted = func(format string, a ...any) (n int, err error) {
		resultFormatStr = format
		resultValues = a
		return 0, nil
	}

	// Other mocks
	err := setupMockStdIn("test value\n")
	if err != nil {
		t.Fatal(err)
	}

	// Run the function, and confirm output
	expectedVarName := "MyVarName123"
	expectedFormatStr := "variable value required for '%s' > "

	_, err = variable.Prompt(expectedVarName)
	if err != nil {
		t.Errorf("expected to recieve no error. Got %e", err)
	}

	if resultFormatStr != expectedFormatStr {
		t.Errorf("expected print-format string to be '%s'. Got '%s'", expectedFormatStr, resultFormatStr)
	}

	if len(resultValues) != 1 {
		t.Errorf("expected 1 value to be passed to print-format func. Got %d", len(resultValues))
	}

	if resultValues[0] != expectedVarName {
		t.Errorf("expected value passed to print-format to be '%s'. Got '%s'", expectedVarName, resultValues[0])
	}
}

func TestPromptWillReturnErrorFromReadingInput(t *testing.T) {
	err := setupMockStdIn("test value that's not terminated") // Value not terminated with \n, so should fail
	if err != nil {
		t.Fatal(err)
	}

	result, err := variable.Prompt("a")

	if result != "" {
		t.Errorf("expected result to be empty string. Got '%s'", result)
	}

	if err == nil {
		t.Error("expected to recieve an error. Got nil")
	}

	if err.Error() != "EOF" {
		t.Errorf("expected to recieve EOF error. Got '%e'", err)
	}
}

func TestPromptWillRemoveNewlineFromEndOfInput(t *testing.T) {
	err := setupMockStdIn("test value\n")
	if err != nil {
		t.Fatal(err)
	}

	result, err := variable.Prompt("a")
	if err != nil {
		t.Errorf("expected to recieve no error. Got %e", err)
	}

	if strings.HasSuffix(result, "\n") {
		t.Errorf("expected result to not end with newline. Got string ending with newline")
	}
}

func TestPromptWillRemoveCRLFFromEndOfInput(t *testing.T) {
	err := setupMockStdIn("test value\r\n")
	if err != nil {
		t.Fatal(err)
	}

	result, err := variable.Prompt("a")
	if err != nil {
		t.Errorf("expected to recieve no error. Got %e", err)
	}

	if strings.HasSuffix(result, "\r") {
		t.Errorf("expected result to not end with carriage return. Got string ending with carriage return")
	}

	if strings.HasSuffix(result, "\r\n") {
		t.Errorf("expected result to not end with CRLF. Got string ending with CRLF")
	}
}

func TestPromptWillRemoveSurroundingQuotes(t *testing.T) {
	quoteValue := "\""

	err := setupMockStdIn(fmt.Sprintf("%stest value%s\n", quoteValue, quoteValue))
	if err != nil {
		t.Fatal(err)
	}

	result, err := variable.Prompt("a")
	if err != nil {
		t.Errorf("expected to recieve no error. Got %e", err)
	}

	if strings.HasPrefix(result, quoteValue) {
		t.Errorf("expected result to not start with quote. Got string starting with %s quote", quoteValue)
	}

	if strings.HasSuffix(result, quoteValue) {
		t.Errorf("expected result to not end with quote. Got string ending with %s quote", quoteValue)
	}
}
