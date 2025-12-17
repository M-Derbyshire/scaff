package models_test

import (
	"testing"

	"github.com/M-Derbyshire/scaff/models"
)

func TestScaffFileGetInvalidJsonErrorReturnsMessage(t *testing.T) {
	scafffile := models.ScaffFile{}

	result := scafffile.GetInvalidJsonError()
	if result == nil {
		t.Errorf("expected to recieve an error. got nil")
		return
	}

	resultMsg := result.Error()
	expectedMsg := "encountered an invalid scaff file. scaff files should contain 2 properties: 'commands' (array of command objects) and 'children' (array of strings)"

	if resultMsg != expectedMsg {
		t.Errorf("expected message to be '%s'. got '%s'", expectedMsg, resultMsg)
	}
}

func TestScaffFileValidateChildrenArrayReturnsErrorIfEmptyChildStringFound(t *testing.T) {
	scafffile := models.ScaffFile{
		Children: []string{
			"test1",
			"",
			"test2",
		},
	}

	result := scafffile.ValidateChildrenArray()
	if result == nil {
		t.Errorf("expected to recieve an error. got nil")
		return
	}

	resultMsg := result.Error()
	expectedMsg := "encountered an empty file path for a child scaff file"

	if resultMsg != expectedMsg {
		t.Errorf("expected message to be '%s'. got '%s'", expectedMsg, resultMsg)
	}
}

func TestScaffFileValidateChildrenArrayReturnsErrorIfChildStringIsOnlyWhitespace(t *testing.T) {
	scafffile := models.ScaffFile{
		Children: []string{
			"test1",
			"\n\t \r  \n",
			"test2",
		},
	}

	result := scafffile.ValidateChildrenArray()
	if result == nil {
		t.Errorf("expected to recieve an error. got nil")
		return
	}

	resultMsg := result.Error()
	expectedMsg := "encountered an empty file path for a child scaff file"

	if resultMsg != expectedMsg {
		t.Errorf("expected message to be '%s'. got '%s'", expectedMsg, resultMsg)
	}
}

func TestScaffFileValidateChildrenArrayWontReturnErrorIfChildStringsCorrect(t *testing.T) {
	scafffile := models.ScaffFile{
		Children: []string{
			"test1",
			"test2",
			"test3",
		},
	}

	result := scafffile.ValidateChildrenArray()
	if result != nil {
		t.Errorf("expected to recieve nil. got error: '%s'", result.Error())
		return
	}
}
