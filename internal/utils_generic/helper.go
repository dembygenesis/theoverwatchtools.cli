package utils_generic

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

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
