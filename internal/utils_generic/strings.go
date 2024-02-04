package utils_generic

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func GetJSONAndCopyToClipboard(i ...interface{}) string {
	_json := GetJSON(i)
	err := clipboard.WriteAll(_json)
	if err != nil {
		fmt.Println(err, "error writing to clipboard")
	}
	return _json
}

// GetJSON return an interface as type json
// Returns a generic error formatted string if fails
func GetJSON(i ...interface{}) string {
	if len(i) > 1 {
		d := make([]interface{}, 0)
		for _, v := range i {
			d = append(d, v)
		}
		j, err := json.Marshal(d)
		if err != nil {
			return "Error dumping to json"
		} else {
			return string(j)
		}
	} else {
		j, err := json.Marshal(i)
		if err != nil {
			return "Error dumping to json"
		} else {
			return string(j)
		}
	}
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
	logger := common.GetLogger(nil)
	var files []string

	err := filepath.Walk(root, visit(&files, root, exclude))
	if err != nil {
		logger.Warnf("File walk error: %s\n", err)
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
