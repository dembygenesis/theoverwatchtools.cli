package fslib

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	log = logger.New(context.TODO())
)

type CopyOptions struct {
	Source                    string   `mapstructure:"source" validate:"required" json:"source"`
	SourceExclusions          []string `mapstructure:"source_exclusions" json:"source_exclusions"`
	Destination               string   `mapstructure:"destination" validate:"required" json:"destination"`
	WipeDestination           bool     `mapstructure:"wipe_destination" json:"wipe_destination"`
	WipeDestinationExclusions []string `mapstructure:"wipe_destination_exclusions" json:"wipe_destination_exclusions"`
}

func (c *CopyOptions) Validate() error {
	return validationutils.Validate(c)
}

func wipeDestinationDir(destination string, exclusions []string) ([]string, error) {
	var pathsToDelete []string

	err := filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %v", path, err)
		}

		skip := false
		for _, exclusion := range exclusions {
			if strings.HasSuffix(path, exclusion) {
				skip = true
				break // Skip the excluded path
			}
		}

		if !skip && path != destination {
			pathsToDelete = append(pathsToDelete, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return pathsToDelete, nil
}

func deletePaths(paths []string) (int, error) {
	deletedCount := 0
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return deletedCount, fmt.Errorf("failed to remove %s: %v", path, err)
		}
		deletedCount++
	}
	return deletedCount, nil
}

func CopyDirToAnother(opts *CopyOptions) error {
	totalDeleted := 0
	totalAdded := 0

	if opts.WipeDestination {
		pathsToDelete, err := wipeDestinationDir(opts.Destination, opts.WipeDestinationExclusions)
		if err != nil {
			return fmt.Errorf("error calculating paths to delete: %w", err)
		}

		deletedCount, err := deletePaths(pathsToDelete)
		if err != nil {
			return fmt.Errorf("error deleting paths: %w", err)
		}
		totalDeleted = deletedCount
	}

	addedCount, err := copyDir(opts.Source, opts.Destination, opts.SourceExclusions)
	if err != nil {
		return fmt.Errorf("copy dir: %w", err)
	}
	totalAdded = addedCount

	log.Info("totalDeleted:", totalDeleted)
	log.Info("totalAdded:", totalAdded)

	return nil
}

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

func copyDir(src, dst string, exclusions []string) (int, error) {
	addedCount := 0
	srcInfo, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !srcInfo.IsDir() {
		return 0, fmt.Errorf("source is not a directory")
	}

	for _, exclusion := range exclusions {
		if strings.HasSuffix(src, exclusion) {
			return 0, nil // Skip the excluded directory
		}
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return addedCount, err
	}
	addedCount++

	entries, err := os.ReadDir(src)
	if err != nil {
		return addedCount, err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			dirAddedCount, err := copyDir(srcPath, dstPath, exclusions)
			if err != nil {
				return addedCount, err
			}
			addedCount += dirAddedCount
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return addedCount, err
			}
			addedCount++
		}
	}
	return addedCount, nil
}

// CreateFileWithDirs creates the file specified by the given path and writes the provided bytes to it,
// ensuring that all parent directories exist. It replaces the file if it already exists.
func CreateFileWithDirs(filePath string, data []byte) error {
	// Extract the directory part of the file path
	dirPath := filepath.Dir(filePath)

	// Ensure all directories in the path exist
	err := os.MkdirAll(dirPath, os.ModePerm) // os.ModePerm is 0777, allowing read, write, and execute
	if err != nil {
		return err
	}

	// Create (or truncate) the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed when this function completes

	// Write the provided bytes to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	// File has been created and data has been written successfully
	return nil
}
