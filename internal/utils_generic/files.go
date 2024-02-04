package utils_generic

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/models"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CopyDir initiates the directory copying process based on the provided CopyOptions.
func CopyDir(opts *models.CopyOptions) error {
	if opts.WipeDestination {
		if err := wipeDestinationDir(opts.Destination, opts.WipeDestinationExclusions); err != nil {
			return fmt.Errorf("error wiping destination directory: %w", err)
		}
	}
	return copyDir(opts.Source, opts.Destination, opts.SourceExclusions)
}

// wipeDestinationDir wipes the destination directory, excluding the specified paths.
func wipeDestinationDir(destination string, exclusions []string) error {
	return filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, exclusion := range exclusions {
			if strings.HasSuffix(path, exclusion) {
				return nil // Skip the excluded path
			}
		}

		if path != destination { // Avoid removing the root destination directory itself
			return os.RemoveAll(path)
		}

		return nil
	})
}

// copyFile copies a single fileDetail from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// copyDir recursively copies a directory from src to dst, ignoring paths in the exclusions array.
func copyDir(src, dst string, exclusions []string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	for _, exclusion := range exclusions {
		if strings.HasSuffix(src, exclusion) {
			return nil // Skip the excluded directory
		}
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath, exclusions); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
