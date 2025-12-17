package models_test

import (
	"testing"

	"github.com/M-Derbyshire/scaff/models"
)

func TestScaffFileGetInvalidJsonErrorReturnsMessage(t *testing.T) {
	scafffile := models.ScaffFile{}

	result := scafffile.GetInvalidJsonError()

	resultMsg := result.Error()
	expectedMsg := "encountered an invalid scaff file. scaff files should contain 2 properties: 'commands' (array of command objects) and 'children' (array of strings)"

	if resultMsg != expectedMsg {
		t.Errorf("expected message to be '%s'. got '%s'", expectedMsg, resultMsg)
	}
}
