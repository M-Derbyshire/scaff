package models_test

import (
	"io/fs"
	"slices"
	"testing"

	"github.com/M-Derbyshire/scaff/customerrors"
	"github.com/M-Derbyshire/scaff/models"
)

func TestCommandValidateShouldReturnNoErrorsIfValid(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "test",
				TemplatePath: "template1.txt",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir",
			},
		},
	}

	results := command.Validate("C:/test")

	resultLength := len(results)

	if resultLength > 0 {
		t.Errorf("expected no errors. got %d", resultLength)
	}
}

func TestCommandValidateShouldReturnErrorIfEmptyTemplateDirectoryPath(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	expectedErr := "command objects should have a 'templateDirectoryPath' property that is set to a non-empty value"

	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "",
		Files: []models.FileScaffold{
			{
				Name:         "test",
				TemplatePath: "template1.txt",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir",
			},
		},
	}

	results := command.Validate("C:/test")

	if len(results) != 1 {
		t.Errorf("expected 1 error. got %d", len(results))
		return
	}

	errMsg := results[0].Message
	if errMsg != expectedErr {
		t.Errorf("expected error to be '%s'. got '%s'", expectedErr, errMsg)
	}
}

func TestCommandValidateShouldReturnErrorIfTemplateDirectoryPathNotSet(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	expectedErr := "command objects should have a 'templateDirectoryPath' property that is set to a non-empty value"

	command := models.Command{
		Name: "test1",
		Files: []models.FileScaffold{
			{
				Name:         "test",
				TemplatePath: "template1.txt",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir",
			},
		},
	}

	results := command.Validate("C:/test")

	if len(results) != 1 {
		t.Errorf("expected 1 error. got %d", len(results))
		return
	}

	errMsg := results[0].Message
	if errMsg != expectedErr {
		t.Errorf("expected error to be '%s'. got '%s'", expectedErr, errMsg)
	}
}

func TestCommandValidateShouldReturnErrorIfTemplateDirectoryPathSetToOnlyWhitespace(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	expectedErr := "command objects should have a 'templateDirectoryPath' property that is set to a non-empty value"

	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "  \t   \n    ",
		Files: []models.FileScaffold{
			{
				Name:         "test",
				TemplatePath: "template1.txt",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir",
			},
		},
	}

	results := command.Validate("C:/test")

	if len(results) != 1 {
		t.Errorf("expected 1 error. got %d", len(results))
		return
	}

	errMsg := results[0].Message
	if errMsg != expectedErr {
		t.Errorf("expected error to be '%s'. got '%s'", expectedErr, errMsg)
	}
}

func TestCommandValidateShouldNotErrorIfNoFiles(t *testing.T) {
	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "/test",
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir",
			},
		},
	}

	results := command.Validate("C:/test")

	resultLength := len(results)

	if resultLength > 0 {
		t.Errorf("expected no errors. got %d", resultLength)
	}
}

func TestCommandValidateShouldNotErrorIfNoDirectories(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "test",
				TemplatePath: "template1.txt",
			},
		},
	}

	results := command.Validate("C:/test")

	resultLength := len(results)

	if resultLength > 0 {
		t.Errorf("expected no errors. got %d", resultLength)
	}
}

func TestCommandValidateShouldCallValidateOnAllFiles(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	expectedErrs := []customerrors.ValidationError{
		{
			Message: "file scaffold objects should have a 'name' property that is set to a non-empty value",
		},
		{
			Message: "file scaffold objects should have a 'templatePath' property that is set to a non-empty value",
		},
	}

	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "",
				TemplatePath: "template1.txt",
			},
			{
				Name:         "test",
				TemplatePath: "",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir",
			},
		},
	}

	results := command.Validate("C:/test")

	if len(results) != 2 {
		t.Errorf("expected 2 errors. got %d", len(results))
		return
	}

	for _, expectedErr := range expectedErrs {
		if !slices.Contains(results, expectedErr) {
			t.Errorf("expected result errors to include '%s'. it did not", expectedErr)
		}
	}
}

func TestCommandValidateShouldCallValidateOnAllDirectories(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	expectedErrs := []customerrors.ValidationError{
		{
			Message: "directory scaffold objects should have a 'name' property that is set to a non-empty value",
		},
		{
			Message: "file scaffold objects should have a 'name' property that is set to a non-empty value",
		},
	}

	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "test1.txt",
				TemplatePath: "template1.txt",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "",
			},
			{
				Name: "test_dir",
				Files: []models.FileScaffold{
					{
						Name:         "",
						TemplatePath: "/test",
					},
				},
			},
		},
	}

	results := command.Validate("C:/test")

	if len(results) != 2 {
		t.Errorf("expected 2 errors. got %d", len(results))
		return
	}

	for _, expectedErr := range expectedErrs {
		if !slices.Contains(results, expectedErr) {
			t.Errorf("expected result errors to include '%s'. it did not", expectedErr)
		}
	}
}

func TestCommandValidateShouldReturnMultipleErrors(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	expectedErrs := []customerrors.ValidationError{
		{
			Message: "command objects should have a 'templateDirectoryPath' property that is set to a non-empty value",
		},
		{
			Message: "file scaffold objects should have a 'name' property that is set to a non-empty value",
		},
		{
			Message: "directory scaffold objects should have a 'name' property that is set to a non-empty value",
		},
	}

	command := models.Command{
		Name:                  "test1",
		TemplateDirectoryPath: "",
		Files: []models.FileScaffold{
			{
				Name:         "",
				TemplatePath: "template1.txt",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "",
			},
		},
	}

	results := command.Validate("C:/test")

	if len(results) != 3 {
		t.Errorf("expected 3 errors. got %d", len(results))
		return
	}

	for _, expectedErr := range expectedErrs {
		if !slices.Contains(results, expectedErr) {
			t.Errorf("expected result errors to include '%s'. it did not", expectedErr)
		}
	}
}
