package cliputil

import (
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	log = logger.New(context.TODO())
)

func visit(files *[]string, root string, ignoredPrefixed []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("Encountered an error accessing path %s: %s\n", path, err)
			return err
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			fmt.Printf("Error calculating relative path for %s: %s\n", path, err)
			return err
		}
		for _, prefix := range ignoredPrefixed {
			if strings.HasPrefix(strings.TrimSpace(relativePath), strings.TrimSpace(prefix)) {
				if info.IsDir() {
					return filepath.SkipDir // Skip the entire directory
				}
				return nil // Skip this file
			}
		}

		*files = append(*files, path)

		return nil
	}
}

func CopyRootPathToClipboard(root string, exclude []string) ([]string, error) {

	var files []string

	err := filepath.Walk(root, visit(&files, root, exclude))
	if err != nil {
		log.Warnf("file walk error: %s\n", err)
		return files, fmt.Errorf("file walk: %v", err)
	}

	var contentBuilder strings.Builder
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			if !strings.Contains(err.Error(), "is a directory") {
				log.Warnf("reading file '%s' failed: %s\n", file, err)
			}
		}
		fileContent := string(data)

		contentBuilder.WriteString(fmt.Sprintf("\n\n--- %s ---\n\n", file))
		contentBuilder.WriteString(fileContent)
	}

	if err := clipboard.WriteAll(contentBuilder.String()); err != nil {
		log.Warnf("Clipboard write error: %s\n", err)
		return files, fmt.Errorf("clip: %v", err)
	}

	return files, nil
}

// GetJSONAndCopyToClipboard generates a JSON string from the input and copies it to the clipboard.
func GetJSONAndCopyToClipboard(i ...interface{}) string {
	if len(i) == 0 || (len(i) == 1 && i[0] == nil) {
		fmt.Println("No input provided or input is nil")
		return ""
	}

	jsonStr := strutil.GetAsJson(i...) // Simplified call

	if jsonStr == "" {
		fmt.Println("Generated JSON string is empty")
		return ""
	}

	err := clipboard.WriteAll(jsonStr)
	if err != nil {
		fmt.Printf("Error copying to clipboard: %v\n", err)
		return ""
	}
	return jsonStr
}
