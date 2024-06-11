package help_test

import (
	"testing"

	"github.com/M-Derbyshire/scaff/help"
)

func TestTextWillReturnPopulatedString(t *testing.T) {
	result := help.Text()

	if len(result) == 0 {
		t.Errorf("expected help text string to be populated. Got an empty string")
	}
}
