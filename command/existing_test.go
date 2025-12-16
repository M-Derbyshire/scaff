package command_test

import (
	"path/filepath"
	"testing"

	"slices"

	"github.com/M-Derbyshire/scaff/command"
	"github.com/M-Derbyshire/scaff/mocks"
	"github.com/M-Derbyshire/scaff/models"
)

func TestIdentifyExistingPathsWillReturnExistingPaths(t *testing.T) {
	workingDirectory := "C:/project"

	files := []mocks.MockFileInfo{
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_file.txt"), false),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_file2.txt"), false),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dir"), true),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dir2"), true),
	}

	command.FileStat = mocks.GetFileStat(files)

	testCommand := models.Command{
		Name:                  "test",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "my_file.txt",
				TemplatePath: "",
			},
			{
				Name:         "my_file2.txt",
				TemplatePath: "",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name:        "my_dir",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
			{
				Name:        "my_dir2",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	results, _ := command.IdentifyExistingPaths(testCommand, workingDirectory, map[string]string{})

	if len(results) != len(files) {
		t.Errorf("expected length of results from IdentifyExistingPaths to be %d. got %d", len(files), len(results))
		return
	}

	for idx, file := range files {
		if results[idx] != file.FilePath {
			t.Errorf("expected results from IdentifyExistingPaths to contain \"%s\", but it did not", file.FilePath)
		}
	}
}

func TestIdentifyExistingPathsWillNotReturnNonExistantPaths(t *testing.T) {
	workingDirectory := "C:/project"

	existingFiles := []mocks.MockFileInfo{
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_file.txt"), false),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dir"), true),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_file2.txt"), false),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dir2"), true),
	}
	nonExistingPaths := []string{
		filepath.Join(workingDirectory, "my_file3.txt"),
		filepath.Join(workingDirectory, "my_dir3"),
		filepath.Join(workingDirectory, "my_file4.txt"),
		filepath.Join(workingDirectory, "my_dir4"),
	}

	command.FileStat = mocks.GetFileStat(existingFiles)

	// Has a mixture of existing and non-existing files/directories
	testCommand := models.Command{
		Name:                  "test",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "my_file.txt",
				TemplatePath: "",
			},
			{
				Name:         "my_file2.txt",
				TemplatePath: "",
			},
			{
				Name:         "my_file3.txt",
				TemplatePath: "",
			},
			{
				Name:         "my_file4.txt",
				TemplatePath: "",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name:        "my_dir",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
			{
				Name:        "my_dir2",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
			{
				Name:        "my_dir3",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
			{
				Name:        "my_dir4",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	results, _ := command.IdentifyExistingPaths(testCommand, workingDirectory, map[string]string{})

	for _, path := range results {
		if slices.Contains(nonExistingPaths, path) {
			t.Errorf("expected results from IdentifyExistingPaths not to contain non-existant paths, but got \"%s\"", path)
		}
	}
}

func TestIdentifyExistingPathsWillReturnEmptySliceIfNoExistingPaths(t *testing.T) {
	workingDirectory := "C:/project"

	command.FileStat = mocks.GetFileStat([]mocks.MockFileInfo{})

	testCommand := models.Command{
		Name:                  "test",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "my_file.txt",
				TemplatePath: "",
			},
			{
				Name:         "my_file2.txt",
				TemplatePath: "",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name:        "my_dir",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
			{
				Name:        "my_dir2",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	results, _ := command.IdentifyExistingPaths(testCommand, workingDirectory, map[string]string{})

	if len(results) > 0 {
		t.Errorf("expected IdentifyExistingPaths to return empty slice. got length of %d", len(results))
	}
}

func TestIdentifyExistingPathsWillPopulatePathNamesUsingGivenVars(t *testing.T) {
	workingDirectory := "C:/project"

	vars := map[string]string{
		"files_var": "files_text",
		"dirs_var":  "dirs_text",
	}

	files := []mocks.MockFileInfo{
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_files_text_file.txt"), false),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dirs_text_dir"), true),
	}

	command.FileStat = mocks.GetFileStat(files)

	testCommand := models.Command{
		Name:                  "test",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "my_{: files_var :}_file.txt",
				TemplatePath: "",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name:        "my_{: dirs_var :}_dir",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	results, _ := command.IdentifyExistingPaths(testCommand, workingDirectory, vars)

	if len(results) != len(files) {
		t.Errorf("expected length of results from IdentifyExistingPaths to be %d. got %d", len(files), len(results))
		return
	}

	for idx, result := range results {
		if result != files[idx].FilePath {
			t.Errorf("expected path from IdentifyExistingPaths to be \"%s\". got \"%s\"", files[idx].FilePath, result)
		}
	}
}

func TestIdentifyExistingPathsWillNotCheckPathsWithinInnerDirectories(t *testing.T) {
	workingDirectory := "C:/project"

	rootFiles := []mocks.MockFileInfo{
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_file.txt"), false),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dir"), true),
	}
	innerFiles := []mocks.MockFileInfo{
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dir", "my_inner_dir"), true),
		mocks.CreateMockInfo(filepath.Join(workingDirectory, "my_dir", "my_inner_file.txt"), false),
	}

	command.FileStat = mocks.GetFileStat(append(rootFiles, innerFiles...))

	testCommand := models.Command{
		Name:                  "test",
		TemplateDirectoryPath: "/test",
		Files: []models.FileScaffold{
			{
				Name:         "my_file.txt",
				TemplatePath: "",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "my_dir",
				Files: []models.FileScaffold{
					{
						Name:         "my_inner_file.txt",
						TemplatePath: "",
					},
				},
				Directories: []models.DirectoryScaffold{
					{
						Name:        "my_inner_dir",
						Files:       []models.FileScaffold{},
						Directories: []models.DirectoryScaffold{},
					},
				},
			},
		},
	}

	results, _ := command.IdentifyExistingPaths(testCommand, workingDirectory, map[string]string{})

	for _, innerFile := range innerFiles {
		if slices.Contains(results, innerFile.FilePath) {
			t.Errorf("expected IdentifyExistingPaths not to check for any paths within inner directories, but it checked \"%s\"", innerFile.FilePath)
		}
	}
}
