package create_test

import (
	"errors"
	"io/fs"
	"strings"
	"testing"

	"github.com/M-Derbyshire/scaff/create"
	"github.com/M-Derbyshire/scaff/mocks"
	"github.com/M-Derbyshire/scaff/models"
)

// setup runs any setup code that is generic across all tests for the directory func
func directoryBeforeEach() {
	create.ReadFile = mocks.GetReadFile([]byte{})
	create.WriteFile = mocks.GetWriteFile()
	create.Mkdir = mocks.GetMkdir()
}

func TestWillCreateTheDirectoryAndFileStructure(t *testing.T) {
	directoryBeforeEach()

	// The directory/file structure we want to create
	structure := models.DirectoryScaffold{
		Name: "mainDir",
		Files: []models.FileScaffold{
			{
				Name:         "file1",
				TemplatePath: "/",
			},
			{
				Name:         "file2",
				TemplatePath: "/",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "dir1",
				Files: []models.FileScaffold{
					{
						Name:         "file3",
						TemplatePath: "/",
					},
					{
						Name:         "file4",
						TemplatePath: "/",
					},
				},
				Directories: []models.DirectoryScaffold{},
			},
			{
				Name:  "dir2",
				Files: []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{
					{
						Name:        "dir3",
						Files:       []models.FileScaffold{},
						Directories: []models.DirectoryScaffold{},
					},
					{
						Name:        "dir4",
						Files:       []models.FileScaffold{},
						Directories: []models.DirectoryScaffold{},
					},
				},
			},
		},
	}

	parentDirPath := "C:/parent"

	// The directory paths we expect to be created, in the order we expect
	expectedDirectoryPaths := []string{
		parentDirPath + "/mainDir",
		parentDirPath + "/mainDir/dir1",
		parentDirPath + "/mainDir/dir2",
		parentDirPath + "/mainDir/dir2/dir3",
		parentDirPath + "/mainDir/dir2/dir4",
	}

	// The directory paths we expect to be created, in the order we expect
	expectedFilePaths := []string{
		parentDirPath + "/mainDir/file1",
		parentDirPath + "/mainDir/file2",
		parentDirPath + "/mainDir/dir1/file3",
		parentDirPath + "/mainDir/dir1/file4",
	}

	// We're going to record the directory paths that Mkdir is called with
	mkdirPathCalls := make([]string, 0, len(expectedDirectoryPaths))
	create.Mkdir = func(s string, _ fs.FileMode) error {
		mkdirPathCalls = append(mkdirPathCalls, s)
		return nil
	}

	// We're going to record the file paths that WriteFile is called with
	writeFilePathCalls := make([]string, 0, len(expectedFilePaths))
	create.WriteFile = func(s string, _ []byte, _ fs.FileMode) error {
		writeFilePathCalls = append(writeFilePathCalls, s)
		return nil
	}

	// Run the function
	create.Directory(structure, parentDirPath, "/", map[string]string{})

	// Check the recorded directory paths match the expected paths
	if len(mkdirPathCalls) != len(expectedDirectoryPaths) {
		t.Errorf("expected %d calls to Mkdir. Got %d", len(expectedDirectoryPaths), len(mkdirPathCalls))
	}

	for i := range mkdirPathCalls {
		if strings.Compare(mkdirPathCalls[i], expectedDirectoryPaths[i]) != 0 {
			t.Errorf("expected Mkdir to be called with '%s'. Got '%s'", expectedDirectoryPaths[i], mkdirPathCalls[i])
		}
	}

	// Check the recorded file paths match the expected paths
	if len(writeFilePathCalls) != len(expectedFilePaths) {
		t.Errorf("expected %d calls to WriteFile. Got %d", len(expectedFilePaths), len(writeFilePathCalls))
	}

	for i := range writeFilePathCalls {
		if strings.Compare(writeFilePathCalls[i], expectedFilePaths[i]) != 0 {
			t.Errorf("expected WriteFile to be called with '%s'. Got '%s'", expectedFilePaths[i], writeFilePathCalls[i])
		}
	}
}

func TestWillPopulateDirectoryNamesWithGivenVars(t *testing.T) {
	directoryBeforeEach()

	vars := map[string]string{
		"var1": "val1",
		"var2": "val2",
	}

	directory := models.DirectoryScaffold{
		Name:  "test {: var1 :}",
		Files: []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{
			{
				Name:        "test {: var2 :}",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	expectedDirPaths := []string{
		"/test val1",
		"/test val1/test val2",
	}

	// We're going to record the directory paths that Mkdir is called with
	mkdirPathCalls := make([]string, 0, 2)
	create.Mkdir = func(s string, _ fs.FileMode) error {
		mkdirPathCalls = append(mkdirPathCalls, s)
		return nil
	}

	create.Directory(directory, "/", "/", vars)

	if len(mkdirPathCalls) != len(expectedDirPaths) {
		t.Errorf("expected Mkdir to have been called %d times. Was called %d times", len(expectedDirPaths), len(mkdirPathCalls))
	}

	for i := range mkdirPathCalls {
		if mkdirPathCalls[i] != expectedDirPaths[i] {
			t.Errorf("expected Mkdir to be called with '%s'. Got '%s'", expectedDirPaths[i], mkdirPathCalls[i])
		}
	}
}

func TestWillPassGivenVarsToFileCreate(t *testing.T) {
	directoryBeforeEach()

	vars := map[string]string{
		"var1": "val1",
		"var2": "val2",
	}

	directory := models.DirectoryScaffold{
		Name: "test",
		Files: []models.FileScaffold{
			{
				Name:         "test {: var1 :}",
				TemplatePath: "/",
			},
			{
				Name:         "test {: var2 :}",
				TemplatePath: "/",
			},
		},
		Directories: []models.DirectoryScaffold{},
	}

	expectedFilePaths := []string{
		"/test/test val1",
		"/test/test val2",
	}

	// We're going to record the file paths that WriteFile is called with
	writeFilePathCalls := make([]string, 0, 2)
	create.WriteFile = func(s string, _ []byte, _ fs.FileMode) error {
		writeFilePathCalls = append(writeFilePathCalls, s)
		return nil
	}

	create.Directory(directory, "/", "/", vars)

	if len(writeFilePathCalls) != len(expectedFilePaths) {
		t.Errorf("expected WriteFile to have been called %d times. Was called %d times", len(expectedFilePaths), len(writeFilePathCalls))
	}

	for i := range writeFilePathCalls {
		if writeFilePathCalls[i] != expectedFilePaths[i] {
			t.Errorf("expected WriteFile to be called with '%s'. Got '%s'", expectedFilePaths[i], writeFilePathCalls[i])
		}
	}
}

func TestWillCreateDirectoryWithCorrectPermissionBits(t *testing.T) {
	directoryBeforeEach()

	directory := models.DirectoryScaffold{
		Name:        "test",
		Files:       []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{},
	}

	var resultPerms fs.FileMode
	create.Mkdir = func(_ string, perms fs.FileMode) error {
		resultPerms = perms
		return nil
	}

	var expectedPerms fs.FileMode = 0777

	create.Directory(directory, "/", "/", map[string]string{})

	if resultPerms != expectedPerms {
		t.Errorf("expected directory to be created with permissions %#o. Got %#o", expectedPerms, resultPerms)
	}
}

func TestWillPassTemplateFilePathToFileCreate(t *testing.T) {
	directoryBeforeEach()

	expectedTemplateDirPath := "C:/my-template-directory"

	directory := models.DirectoryScaffold{
		Name: "test",
		Files: []models.FileScaffold{
			{
				Name:         "test1",
				TemplatePath: "/",
			},
		},
		Directories: []models.DirectoryScaffold{
			{
				Name: "test",
				Files: []models.FileScaffold{
					{
						Name:         "test2",
						TemplatePath: "/",
					},
				},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	// We're going to record the paths that ReadFile is called with
	readFilePathCalls := make([]string, 0, 2)
	create.ReadFile = func(s string) ([]byte, error) {
		readFilePathCalls = append(readFilePathCalls, s)
		return []byte{}, nil
	}

	// Call the function
	create.Directory(directory, "/", expectedTemplateDirPath, map[string]string{})

	if len(readFilePathCalls) != 2 {
		t.Errorf("expected ReadFile to be called 2 times. Got %d", len(readFilePathCalls))
	}

	for _, path := range readFilePathCalls {
		if strings.Compare(expectedTemplateDirPath, path) != 0 {
			t.Errorf("expected ReadFile to be called with '%s'. Got '%s'", expectedTemplateDirPath, path)
		}
	}
}

func TestWillNotReturnErrorOnSuccess(t *testing.T) {
	directoryBeforeEach()

	directory := models.DirectoryScaffold{
		Name:        "test",
		Files:       []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{},
	}

	result := create.Directory(directory, "/", "/", map[string]string{})

	if result != nil {
		t.Errorf("expected returned error to be nil when creating directory. Got '%s'", result.Error())
	}
}

func TestWillReturnErrorFromDirectoryCreate(t *testing.T) {
	directoryBeforeEach()

	errorMessage := "my test error"
	expectedResultErrorMessge := "error while creating directory '/test': " + errorMessage

	create.Mkdir = func(_ string, _ fs.FileMode) error {
		return errors.New(errorMessage)
	}

	directory := models.DirectoryScaffold{
		Name:        "test",
		Files:       []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{},
	}

	result := create.Directory(directory, "/", "/", map[string]string{})

	if result == nil {
		t.Errorf("expected Mkdir error to be returned by directory create. Got nil")
	}

	if strings.Compare(result.Error(), expectedResultErrorMessge) != 0 {
		t.Errorf("expected directory create to return error message from Mkdir ('%s'). Got '%s'", expectedResultErrorMessge, result.Error())
	}
}

func TestWillReturnErrorFromFileCreate(t *testing.T) {
	directoryBeforeEach()

	expectedErrorMessage := "my test error"

	create.WriteFile = func(_ string, _ []byte, _ fs.FileMode) error {
		return errors.New(expectedErrorMessage)
	}

	directory := models.DirectoryScaffold{
		Name: "test",
		Files: []models.FileScaffold{
			{
				Name:         "test",
				TemplatePath: "/",
			},
		},
		Directories: []models.DirectoryScaffold{},
	}

	result := create.Directory(directory, "/", "/", map[string]string{})

	if result == nil {
		t.Errorf("expected WriteFile error to be returned by file create. Got nil")
	}

	if !strings.HasSuffix(result.Error(), expectedErrorMessage) {
		t.Errorf(
			"expected file create to return error message from WriteFile ending with '%s'. Got '%s'",
			expectedErrorMessage,
			result.Error(),
		)
	}
}

func TestWillReturnErrorFromInnerDirectoryCreate(t *testing.T) {
	directoryBeforeEach()

	errorMessage := "my test error"
	expectedResultErrorMessge := "error while creating directory '/test1/test2': " + errorMessage

	isFirstCall := true
	create.Mkdir = func(_ string, _ fs.FileMode) error {
		// The first call is the parent directory
		if isFirstCall {
			isFirstCall = false
			return nil
		}

		return errors.New(errorMessage)
	}

	directory := models.DirectoryScaffold{
		Name:  "test1",
		Files: []models.FileScaffold{},
		Directories: []models.DirectoryScaffold{
			{
				Name:        "test2",
				Files:       []models.FileScaffold{},
				Directories: []models.DirectoryScaffold{},
			},
		},
	}

	result := create.Directory(directory, "/", "/", map[string]string{})

	if result == nil {
		t.Errorf("expected Mkdir error to be returned from inner directory create. Got nil")
	}

	if strings.Compare(result.Error(), expectedResultErrorMessge) != 0 {
		t.Errorf(
			"expected inner directory create error message from Mkdir ('%s'). Got '%s'",
			expectedResultErrorMessge,
			result.Error(),
		)
	}
}
