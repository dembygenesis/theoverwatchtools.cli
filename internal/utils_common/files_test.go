package utils_common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	fileTypeFile = iota
	fileTypeFolder
)

type fileDetail struct {
	Type int
	Name string
}

// FileRemover defines the interface for removing files.
type FileRemover interface {
	RemoveAll(path string) error
}

// deletePaths removes paths using the provided FileRemover.
func deletePaths2(paths []string, remover FileRemover) (int, error) {
	deletedCount := 0
	for _, path := range paths {
		if err := remover.RemoveAll(path); err != nil {
			return deletedCount, fmt.Errorf("failed to remove %s: %v", path, err)
		}
		deletedCount++
	}
	return deletedCount, nil
}

// MockFileRemover is a mock implementation of FileRemover.
type MockFileRemover struct {
	RemovedPaths []string
}

// RemoveAll mocks the RemoveAll method.
func (m *MockFileRemover) RemoveAll(path string) error {
	m.RemovedPaths = append(m.RemovedPaths, path)
	return nil
}

func Test_deletePaths(t *testing.T) {
	// Create temporary files and directories for testing
	dir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir) // Clean up

	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")

	_, err = os.Create(file1)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	_, err = os.Create(file2)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	paths := []string{file1, file2}

	// Call the function under test
	deletedCount, err := deletePaths(paths)
	if err != nil {
		t.Errorf("deletePaths returned an unexpected error: %v", err)
	}

	// Assert that the number of deleted paths matches the expected count
	expectedDeletedCount := len(paths)
	if deletedCount != expectedDeletedCount {
		t.Errorf("deletePaths deleted %d paths, expected %d", deletedCount, expectedDeletedCount)
	}

	// Assert that the paths were deleted
	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("deletePaths did not remove path: %s", path)
		}
	}
}

func TestCopyOptions_Validate(t *testing.T) {
	// Set up the test cases with input data and expected validation results
	testCases := []struct {
		InputData     *CopyOptions
		ExpectedError string
		TestCaseID    int
	}{
		// Test case 1: Valid input data
		{
			InputData: &CopyOptions{
				Source:      "valid_source",
				Destination: "valid_destination",
			},
			ExpectedError: "", // Expect no validation errors
			TestCaseID:    1,
		},
		// Test case 2: Missing source
		{
			InputData: &CopyOptions{
				Source:      "",
				Destination: "valid_destination",
			},
			ExpectedError: "'source' must have a value",
			TestCaseID:    2,
		},
		// Test case 3: Missing destination
		{
			InputData: &CopyOptions{
				Source:      "valid_source",
				Destination: "",
			},
			ExpectedError: "'destination' must have a value",
			TestCaseID:    3,
		},
		// Add more test cases as needed
	}

	// Loop through the test cases
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestCaseID=%d", tc.TestCaseID), func(t *testing.T) {
			// Call the Validate method with the test data
			err := tc.InputData.Validate()

			// Assert that the error message matches the expected error message
			if tc.ExpectedError != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func testCreateFilesAndFolders(t *testing.T, outFolder []string, files []fileDetail) {
	t.Helper()

	for _, fileDetails := range files {
		arrFolderDir := append(outFolder, fileDetails.Name)

		fullFilePath := filepath.Join(arrFolderDir...)

		switch fileDetails.Type {
		case fileTypeFile:
			f, err := os.Create(fullFilePath)
			require.NoError(t, err, "create fileDetail")

			err = f.Close()
			require.NoError(t, err, "close fileDetail")
			break
		case fileTypeFolder:
			err := os.MkdirAll(fullFilePath, 0755)
			require.NoError(t, err, "create folder A fileDetail")
			break
		default:
			require.Error(t, fmt.Errorf("fileType: %v is undefined", fileDetails.Type))
		}
	}
}

func TestCopyDir_FailInvalidDirectoryA(t *testing.T) {
	opts := CopyOptions{
		Source:                    "abc",
		SourceExclusions:          nil,
		Destination:               "def",
		WipeDestination:           true,
		WipeDestinationExclusions: nil,
	}

	err := CopyDirToAnother(&opts)
	require.Error(t, err, "error expected")
}

func TestCopyDir_Success(t *testing.T) {
	var (
		err error
	)

	arrTestDir, cleanup := testCreateTestDir(t)
	defer cleanup()

	arrFolderA := append(arrTestDir, "A")
	arrFolderB := append(arrTestDir, "B")

	// Create folders "A", and "B".
	err = os.MkdirAll(filepath.Join(arrFolderA...), 0755)
	require.NoError(t, err, "create folder A")

	err = os.MkdirAll(filepath.Join(arrFolderB...), 0755)
	require.NoError(t, err, "create folder B")

	// Populate folders "A", and "B"
	folderAFiles := []fileDetail{
		{
			Name: ".git",
			Type: fileTypeFolder,
		},
		{
			Name: ".idea",
			Type: fileTypeFolder,
		},
		{
			Name: "hello world.txt",
			Type: fileTypeFile,
		},
		{
			Name: "pic.jpg",
			Type: fileTypeFile,
		},
	}

	folderBFiles := []fileDetail{
		{
			Name: ".git2",
			Type: fileTypeFolder,
		},
		{
			Name: ".idea2",
			Type: fileTypeFolder,
		},
		{
			Name: "sample.txt",
			Type: fileTypeFile,
		},
	}

	testCreateFilesAndFolders(t, arrFolderA, folderAFiles)
	testCreateFilesAndFolders(t, arrFolderB, folderBFiles)

	// Run copy
	folderA := filepath.Join(arrFolderA...)
	folderB := filepath.Join(arrFolderB...)

	sourceExclusions := []string{".git", ".idea"}
	wipeDestinationExclusions := []string{".git2"}

	opts := CopyOptions{
		Source:                    folderA,
		SourceExclusions:          sourceExclusions,
		Destination:               folderB,
		WipeDestination:           true,
		WipeDestinationExclusions: wipeDestinationExclusions,
	}

	err = CopyDirToAnother(&opts)
	require.NoError(t, err, "copy files")

	// Do assertions
	expectedFilesInFolderB := []string{
		"hello world.txt",
		"pic.jpg",
		".git2",
	}

	files, err := os.ReadDir(folderB)
	require.NoError(t, err, "read test dir")

	assert.Equal(t, len(expectedFilesInFolderB), len(files))

	// Make a map of the actual existing files
	hashFilesInFolderB := make(map[string]bool)
	for _, file := range files {
		hashFilesInFolderB[file.Name()] = true
	}

	for _, fileName := range expectedFilesInFolderB {
		assert.True(t, true, hashFilesInFolderB[fileName])
	}
}
