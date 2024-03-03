package file_system

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_deletePaths(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dir)

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

	deletedCount, err := deletePaths(paths)
	if err != nil {
		t.Errorf("deletePaths returned an unexpected error: %v", err)
	}

	expectedDeletedCount := len(paths)
	if deletedCount != expectedDeletedCount {
		t.Errorf("deletePaths deleted %d paths, expected %d", deletedCount, expectedDeletedCount)
	}

	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("deletePaths did not remove path: %s", path)
		}
	}
}

func TestCopyOptions_Validate(t *testing.T) {
	type testCase struct {
		name        string
		inputData   *CopyOptions
		expectedErr string
	}

	testCases := []testCase{
		{
			name: "Valid Input Data",
			inputData: &CopyOptions{
				Source:      "valid_source",
				Destination: "valid_destination",
			},
			expectedErr: "",
		},
		{
			name: "Missing Source",
			inputData: &CopyOptions{
				Source:      "",
				Destination: "valid_destination",
			},
			expectedErr: "'source' must have a value",
		},
		{
			name: "Missing Destination",
			inputData: &CopyOptions{
				Source:      "valid_source",
				Destination: "",
			},
			expectedErr: "'destination' must have a value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.inputData.Validate()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
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
		case fileTypeFolder:
			err := os.MkdirAll(fullFilePath, 0755)
			require.NoError(t, err, "create folder A fileDetail")
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

func Test_deletePaths_OsRemove_fail(t *testing.T) {
	filePath := []string{
		".",
	}

	_, err := deletePaths(filePath)
	require.Error(t, err)
}

// testCreateTestDir creates a folder based on the test name,
// with its main use-case creating an isolated environment for
// fileDetail manipulation.
func testCreateTestDir(t *testing.T) (testDirArr []string, cleanup func()) {
	t.Helper()

	baseDir, err := os.Getwd()
	require.NoError(t, err, "cannot get working dir")

	testDirArr = []string{
		baseDir,
		t.Name(),
	}

	testDir := filepath.Join(testDirArr...)

	err = os.MkdirAll(filepath.Join(testDirArr...), 0755)
	require.NoError(t, err)

	cleanup = func() {
		err = os.RemoveAll(testDir)
		require.NoError(t, err)
	}

	return testDirArr, cleanup
}
