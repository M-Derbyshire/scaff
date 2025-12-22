package models_test

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"

	"github.com/M-Derbyshire/scaff/models"
)

func TestDirectoryScaffoldValidateShouldReturnErrorIfEmptyName(t *testing.T) {
	expectedErr := "directory scaffold objects should have a 'name' property that is set to a non-empty value"

	scaffold := models.DirectoryScaffold{
		Name:        "",
		Files:       []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{},
	}

	results := scaffold.Validate("test_dir")

	if len(results) != 1 {
		t.Errorf("expected 1 error. got %d", len(results))
		return
	}

	resultErrMsg := results[0].Error()
	if resultErrMsg != expectedErr {
		t.Errorf("expected error to be '%s'. got '%s'", expectedErr, resultErrMsg)
	}
}

func TestDirectoryScaffoldValidateShouldReturnErrorIfNameNotSet(t *testing.T) {
	expectedErr := "directory scaffold objects should have a 'name' property that is set to a non-empty value"

	scaffold := models.DirectoryScaffold{
		Files:       []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{},
	}

	results := scaffold.Validate("test_dir")

	if len(results) != 1 {
		t.Errorf("expected 1 error. got %d", len(results))
		return
	}

	resultErrMsg := results[0].Error()
	if resultErrMsg != expectedErr {
		t.Errorf("expected error to be '%s'. got '%s'", expectedErr, resultErrMsg)
	}
}

func TestDirectoryScaffoldValidateShouldReturnErrorIfNameIsWhitespace(t *testing.T) {
	expectedErr := "directory scaffold objects should have a 'name' property that is set to a non-empty value"

	scaffold := models.DirectoryScaffold{
		Name:        " \t \n  ",
		Files:       []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{},
	}

	results := scaffold.Validate("test_dir")

	if len(results) != 1 {
		t.Errorf("expected 1 error. got %d", len(results))
		return
	}

	resultErrMsg := results[0].Error()
	if resultErrMsg != expectedErr {
		t.Errorf("expected error to be '%s'. got '%s'", expectedErr, resultErrMsg)
	}
}

func TestDirectoryScaffoldValidateShouldNotErrorIfNoFiles(t *testing.T) {
	scaffold := models.DirectoryScaffold{
		Name:        "test",
		Directories: []models.DirectoryScaffold{},
	}

	results := scaffold.Validate("test_dir")

	if len(results) > 0 {
		for _, err := range results {
			t.Errorf("expected no errors. got '%s'", err.Error())
		}
	}
}

func TestDirectoryScaffoldValidateShouldNotErrorIfNoDirectories(t *testing.T) {
	scaffold := models.DirectoryScaffold{
		Name:  "test",
		Files: []models.FileScaffold{},
	}

	results := scaffold.Validate("test_dir")

	if len(results) > 0 {
		for _, err := range results {
			t.Errorf("expected no errors. got '%s'", err.Error())
		}
	}
}

func TestDirectoryScaffoldValidateShouldReturnErrorsFromInnerFilesAndDirectories(t *testing.T) {
	expectedErrs := []string{
		"file scaffold objects should have a 'name' property that is set to a non-empty value",
		"file scaffold objects should have a 'templatePath' property that is set to a non-empty value",
		"directory scaffold objects should have a 'name' property that is set to a non-empty value",
		"file scaffold objects should have a 'name' property that is set to a non-empty value",
		"directory scaffold objects should have a 'name' property that is set to a non-empty value",
	}

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.DirectoryScaffold{
		Name: "test",
		Files: []models.FileScaffold{
			{
				Name:         "",
				TemplatePath: "/test1.txt",
			},
			{
				Name:         "test2.txt",
				TemplatePath: "",
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
						TemplatePath: "test.txt",
					},
				},
				Directories: []models.DirectoryScaffold{
					{
						Name: "",
					},
				},
			},
		},
	}

	results := scaffold.Validate("test_dir")

	if len(results) != len(expectedErrs) {
		t.Errorf("expected %d errors to be returned. got %d", len(expectedErrs), len(results))
		return
	}

	for i, resultErr := range results {
		expectedErr := expectedErrs[i]
		resultMsg := resultErr.Error()

		if expectedErr != resultMsg {
			t.Errorf("expected error to be '%s'. got '%s'", expectedErr, resultMsg)
		}
	}
}

func TestDirectoryScaffoldValidateShouldPassCorrectTemplateDirectoryPathToInnerFilesAndDirectories(t *testing.T) {
	templateDirPath := "C:/project/templates"
	templatePath := "/my_templates/test.txt"
	expectedFullPath := fmt.Sprintf("%s%s", templateDirPath, templatePath)

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		if filepath == expectedFullPath {
			return nil, nil
		}

		return nil, errors.New("not found")
	}

	scaffold := models.DirectoryScaffold{
		Name:  "test_dir1",
		Files: []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir2",
				Files: []models.FileScaffold{
					{
						Name:         "test.txt",
						TemplatePath: templatePath,
					},
				},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	results := scaffold.Validate(templateDirPath)

	for _, resultErr := range results {
		t.Errorf("expected no errors. got '%s'", resultErr.Error())
	}
}

func TestDirectoryScaffoldValidateShouldReturnMultipleErrors(t *testing.T) {
	expectedErrs := []string{
		"directory scaffold objects should have a 'name' property that is set to a non-empty value",
		"file scaffold objects should have a 'templatePath' property that is set to a non-empty value",
		"directory scaffold objects should have a 'name' property that is set to a non-empty value",
	}

	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.DirectoryScaffold{
		Name: "",
		Files: []models.FileScaffold{
			{
				Name:         "test1.txt",
				TemplatePath: "",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "",
			},
		},
	}

	results := scaffold.Validate("test_dir")

	if len(results) != len(expectedErrs) {
		t.Errorf("expected %d errors to be returned. got %d", len(expectedErrs), len(results))
		return
	}

	for i, resultErr := range results {
		expectedErr := expectedErrs[i]
		resultMsg := resultErr.Error()

		if expectedErr != resultMsg {
			t.Errorf("expected error to be '%s'. got '%s'", expectedErr, resultMsg)
		}
	}
}

func TestDirectoryScaffoldValidateShouldNotReturnErrorsWhenValid(t *testing.T) {
	models.FileStat = func(filepath string) (fs.FileInfo, error) {
		return nil, nil
	}

	scaffold := models.DirectoryScaffold{
		Name: "test",
		Files: []models.FileScaffold{
			{
				Name:         "test1.txt",
				TemplatePath: "test1.txt",
			},
			{
				Name:         "test2.txt",
				TemplatePath: "test2.txt",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test_dir1",
			},
			{
				Name: "test_dir2",
				Files: []models.FileScaffold{
					{
						Name:         "test.txt",
						TemplatePath: "test.txt",
					},
				},
				Directories: []models.DirectoryScaffold{
					{
						Name: "test_dir3",
					},
				},
			},
		},
	}

	results := scaffold.Validate("test_dir")

	for _, resultErr := range results {
		t.Errorf("expected no errors. got '%s'", resultErr.Error())
	}
}
