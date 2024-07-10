package variable_test

import (
	"testing"

	"github.com/M-Derbyshire/scaff/variable"
)

func TestPopulateWillReturnStringWithNoTagsWithoutModification(t *testing.T) {
	originalText := `Test 1. Test 2.
	Test3. Test4. Test5.
	I'm a little teapot`

	result, err := variable.Populate(originalText, make(map[string]string))
	if err != nil {
		t.Errorf("expected no error. Got %e", err)
	}

	if result != originalText {
		t.Errorf("expected result to be '%s'. Got '%s'", originalText, result)
	}
}

func TestPopulateWillReplaceTagsWithVariableValuesFromGivenMap(t *testing.T) {
	originalText := `My name is {: name :}.
	I am {: age :} years old.
	My hair is {: hair :}`

	vars := map[string]string{
		"name": "Noddy",
		"age":  "12",
		"hair": "brown",
	}

	expectedText := `My name is Noddy.
	I am 12 years old.
	My hair is brown`

	result, err := variable.Populate(originalText, vars)
	if err != nil {
		t.Errorf("expected no error. Got %e", err)
	}

	if result != expectedText {
		t.Errorf("expected result to be '%s'. Got '%s'", expectedText, result)
	}
}

func TestPopulateWillPromptForVariableValueIfNotInMapAndPopulateTheTag(t *testing.T) {
	originalText := `My name is {: name :}.
	I am {: age :} years old.`

	// Doesn't contain the name value
	vars := map[string]string{"age": "12"}

	expectedText := `My name is Noddy.
	I am 12 years old.`

	// Mock the value that would be input by the user
	setupMockStdIn("Noddy\n")

	result, err := variable.Populate(originalText, vars)
	if err != nil {
		t.Errorf("expected no error. Got %e", err)
	}

	if result != expectedText {
		t.Errorf("expected result to be '%s'. Got '%s'", expectedText, result)
	}
}

func TestPopulateWillReplaceEscapedTagsWithTagsAndNotPopulateThem(t *testing.T) {
	originalText := `My name is {\: name :}`
	expectedText := `My name is {: name :}`

	vars := map[string]string{
		"name": "I should never get used",
	}

	result, err := variable.Populate(originalText, vars)
	if err != nil {
		t.Errorf("expected no error. Got %e", err)
	}

	if result != expectedText {
		t.Errorf("expected result to be '%s'. Got '%s'", expectedText, result)
	}
}
