package utility

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

func GetUuidUnderscore() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "_")
}

// DecodeToStruct decodes the "in", into the "out".
// Basically this is the equivalent, or marshalling something to JSON,
// and unmarshalling it to another struct.
func DecodeToStruct(in interface{}, out interface{}) error {
	var jsonVal []byte
	var err error

	switch reflect.TypeOf(in).Kind() {
	case reflect.String:
		jsonVal = []byte(in.(string))
	case reflect.Slice:
		// We need to check if the slice actually contains bytes.
		s, ok := in.([]byte)
		if !ok {
			return errors.New("in was a slice but not of bytes")
		}
		jsonVal = s
	default:
		jsonVal, err = json.Marshal(in)
		if err != nil {
			return errors.Wrap(err, "error marshalling 'in'")
		}
	}

	err = json.Unmarshal(jsonVal, out)
	if err != nil {

		return errors.Wrap(err, "error unmarshalling 'out'")
	}
	return nil
}

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
	logger := common.GetLogger(context.TODO())
	var files []string

	err := filepath.Walk(root, visit(&files, root, exclude))
	if err != nil {
		logger.Warnf("file walk error: %s\n", err)
		return files, fmt.Errorf("file walk: %v", err)
	}

	var contentBuilder strings.Builder
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			if !strings.Contains(err.Error(), "is a directory") {
				logger.Warnf("reading file '%s' failed: %s\n", file, err)
			}
		}
		fileContent := string(data)

		contentBuilder.WriteString(fmt.Sprintf("\n\n--- %s ---\n\n", file))
		contentBuilder.WriteString(fileContent)
	}

	if err := clipboard.WriteAll(contentBuilder.String()); err != nil {
		logger.Warnf("Clipboard write error: %s\n", err)
		return files, fmt.Errorf("clip: %v", err)
	}

	return files, nil
}

func IsValidFilename(filename string) error {
	// Define constraints
	const maxFilenameLength = 255
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
	illegalCharacters := []string{
		"<", ">", ":", "\"", "/", "\\", "|", "?", "*", string(rune(0)),
	}

	// Check for empty filename
	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	// Check length
	if len(filename) > maxFilenameLength {
		return errors.New("filename exceeds maximum length")
	}

	// Check for reserved names
	baseFilename := strings.ToUpper(strings.TrimSuffix(filename, filepath.Ext(filename)))
	for _, reservedName := range reservedNames {
		if baseFilename == reservedName {
			return errors.New("filename uses a reserved name")
		}
	}

	// Check for illegal characters
	for _, char := range illegalCharacters {
		if strings.Contains(filename, char) {
			return errors.New("filename contains illegal characters")
		}
	}

	// Check for reserved filenames specific to Unix-like systems
	if filename == "." || filename == ".." {
		return errors.New("filename uses a reserved name for Unix-like systems")
	}

	return nil
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")

	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}
