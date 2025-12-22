package models_test

import (
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"testing"

	"github.com/M-Derbyshire/scaff/models"
)

func TestFileScaffoldValidateWillReturnErrorIfEmptyName(t *testing.T) {
	expectedErr := "file scaffold objects should have a 'name' property that is set to a non-empty value"

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name:         "",
		TemplatePath: "/test.txt",
	}

	results := scaffold.Validate("")

	if len(results) != 1 {
		t.Errorf("expected a single error. got %d", len(results))
		return
	}

	resultMsg := results[0].Error()
	if resultMsg != expectedErr {
		t.Errorf("expected error message to be '%s'. got '%s'", expectedErr, resultMsg)
	}
}

func TestFileScaffoldValidateWillReturnErrorIfNameUndefined(t *testing.T) {
	expectedErr := "file scaffold objects should have a 'name' property that is set to a non-empty value"

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		TemplatePath: "/test.txt",
	}

	results := scaffold.Validate("")

	if len(results) != 1 {
		t.Errorf("expected a single error. got %d", len(results))
		return
	}

	resultMsg := results[0].Error()
	if resultMsg != expectedErr {
		t.Errorf("expected error message to be '%s'. got '%s'", expectedErr, resultMsg)
	}
}

func TestFileScaffoldValidateWillReturnErrorIfNameIsOnlyWhitespace(t *testing.T) {
	expectedErr := "file scaffold objects should have a 'name' property that is set to a non-empty value"

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name:         "\t  \n  ",
		TemplatePath: "/test.txt",
	}

	results := scaffold.Validate("")

	if len(results) != 1 {
		t.Errorf("expected a single error. got %d", len(results))
		return
	}

	resultMsg := results[0].Error()
	if resultMsg != expectedErr {
		t.Errorf("expected error message to be '%s'. got '%s'", expectedErr, resultMsg)
	}
}

func TestFileScaffoldValidateWillReturnErrorIfTemplatePathUndefined(t *testing.T) {
	expectedErr := "file scaffold objects should have a 'templatePath' property that is set to a non-empty value"

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name: "test.txt",
	}

	results := scaffold.Validate("")

	// we will not expect the path-not-found error
	if len(results) != 1 {
		t.Errorf("expected a single error. got %d", len(results))
		return
	}

	resultMsg := results[0].Error()
	if resultMsg != expectedErr {
		t.Errorf("expected error message to be '%s'. got '%s'", expectedErr, resultMsg)
	}
}

func TestFileScaffoldValidateWillReturnErrorIfEmptyTemplatePath(t *testing.T) {
	expectedErr := "file scaffold objects should have a 'templatePath' property that is set to a non-empty value"

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name:         "test.txt",
		TemplatePath: "",
	}

	results := scaffold.Validate("")

	// we will not expect the path-not-found error
	if len(results) != 1 {
		t.Errorf("expected a single error. got %d", len(results))
		return
	}

	resultMsg := results[0].Error()
	if resultMsg != expectedErr {
		t.Errorf("expected error message to be '%s'. got '%s'", expectedErr, resultMsg)
	}
}

func TestFileScaffoldValidateWillReturnErrorIfTemplatePathIsOnlyWhitespace(t *testing.T) {
	expectedErr := "file scaffold objects should have a 'templatePath' property that is set to a non-empty value"

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name:         "test.txt",
		TemplatePath: " \t  \n  ",
	}

	results := scaffold.Validate("")

	// we will not expect the path-not-found error
	if len(results) != 1 {
		t.Errorf("expected a single error. got %d", len(results))
		return
	}

	resultMsg := results[0].Error()
	if resultMsg != expectedErr {
		t.Errorf("expected error message to be '%s'. got '%s'", expectedErr, resultMsg)
	}
}

func TestFileScaffoldValidateWillReturnErrorIfTemplatePathDoesntExist(t *testing.T) {
	templatePath := "my_templates/template1.txt"
	templateDirectoryPath := "C:/my_project"
	absTemplatePath := fmt.Sprintf("%s/%s", templateDirectoryPath, templatePath)

	expectedErr := fmt.Sprintf("unable to locate template file at path: '%s'", absTemplatePath)

	models.AbsPath = func(path string) (string, error) {
		if path == templatePath {
			return absTemplatePath, nil
		}

		return path, nil
	}

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		if filepath == absTemplatePath {
			return nil, errors.New("doesn't exist")
		}

		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name:         "test.txt",
		TemplatePath: templatePath,
	}

	results := scaffold.Validate(templateDirectoryPath)

	if len(results) != 1 {
		t.Errorf("expected a single error. got %d", len(results))
		return
	}

	resultMsg := results[0].Error()
	if resultMsg != expectedErr {
		t.Errorf("expected error message to be '%s'. got '%s'", expectedErr, resultMsg)
	}
}

func TestFileScaffoldValidateWillReturnMultipleErrors(t *testing.T) {
	expectedErrs := []string{
		"file scaffold objects should have a 'name' property that is set to a non-empty value",
		"file scaffold objects should have a 'templatePath' property that is set to a non-empty value",
	}

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name:         "",
		TemplatePath: "",
	}

	results := scaffold.Validate("")

	if len(results) != 2 {
		t.Errorf("expected 2 errors. got %d", len(results))
		return
	}

	for _, err := range results {
		msg := err.Error()

		if !slices.Contains(expectedErrs, msg) {
			t.Errorf("recieved an unexpected error message: '%s'", msg)
		}
	}
}

func TestFileScaffoldValidateWillNotReturnErrorsIfValid(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.FileScaffold{
		Name:         "test.txt",
		TemplatePath: "/test.txt",
	}

	results := scaffold.Validate("C:/my_templates")

	if len(results) != 0 {
		for _, err := range results {
			t.Errorf("expected no errors. got '%s'", err.Error())
		}
	}
}
