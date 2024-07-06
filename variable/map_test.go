package variable_test

import (
	"fmt"
	"testing"

	"github.com/M-Derbyshire/scaff/variable"
)

func TestMapWillSplitVariableArgsIntoMapEntries(t *testing.T) {
	args := []string{
		"key1=val1",
		"key2=val2",
		"key3=val3",
	}

	result := variable.Map(args)

	if len(result) != len(args) {
		t.Errorf("expected %d values in result map. Got %d", len(args), len(result))
	}

	for i := 1; i <= 3; i++ {
		testKey := fmt.Sprintf("key%d", i)
		expectedValue := fmt.Sprintf("val%d", i)

		value, ok := result[testKey]
		if !ok {
			t.Errorf("expected %s value to exist, but it does not.", testKey)
		}
		if value != expectedValue {
			t.Errorf("expected %s value to be '%s'. Got '%s'", testKey, expectedValue, value)
		}
	}
}

func TestMapWillIgnoreArgumentsThatAreNotVariables(t *testing.T) {
	args := []string{
		"test1",
		"key1=val1",
		"test2",
		"test3",
	}

	result := variable.Map(args)

	if len(result) != 1 {
		t.Errorf("expected 1 value in result map. Got %d", len(result))
	}

	if _, ok := result["key1"]; !ok {
		t.Errorf("expected key1 value to exist, but it does not.")
	}
}
