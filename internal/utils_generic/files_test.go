package utils_generic

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/models"
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

// TestCopyDir_FailInvalidDirectoryA tests that the invalid
func TestCopyDir_FailInvalidDirectoryA(t *testing.T) {
	opts := models.CopyOptions{
		Source:                    "abc",
		SourceExclusions:          nil,
		Destination:               "def",
		WipeDestination:           true,
		WipeDestinationExclusions: nil,
	}

	err := CopyDir(&opts)
	require.Error(t, err, "error expected")
}

// TestCopyDir_Success tests that we can successfully copy one folder
// onto another with respect to the exclusions provided.
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

	opts := models.CopyOptions{
		Source:                    folderA,
		SourceExclusions:          sourceExclusions,
		Destination:               folderB,
		WipeDestination:           true,
		WipeDestinationExclusions: wipeDestinationExclusions,
	}

	err = CopyDir(&opts)
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
