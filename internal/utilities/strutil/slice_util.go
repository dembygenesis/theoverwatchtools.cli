package strutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// GetUuidUnderscore generates a uuid but, with underscore ("_"), instead of
// dash ("-").
func GetUuidUnderscore() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "_")
}

// IsValidFilename check if the string is a valid filename.
func IsValidFilename(str string) error {
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

	// Check for empty str
	if str == "" {
		return errors.New("str cannot be empty")
	}

	// Check length
	if len(str) > maxFilenameLength {
		return errors.New("str exceeds maximum length")
	}

	// Check for reserved names
	baseFilename := strings.ToUpper(strings.TrimSuffix(str, filepath.Ext(str)))
	for _, reservedName := range reservedNames {
		if baseFilename == reservedName {
			return errors.New("str uses a reserved name")
		}
	}

	// Check for illegal characters
	for _, char := range illegalCharacters {
		if strings.Contains(str, char) {
			return errors.New("str contains illegal characters")
		}
	}

	// Check for reserved filenames specific to Unix-like systems
	if str == "." || str == ".." {
		return errors.New("str uses a reserved name for Unix-like systems")
	}

	return nil
}

// IsValidEmail check if the string is it is a valid email.
func IsValidEmail(str string) bool {
	var emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")

	if len(str) < 3 && len(str) > 254 {
		return false
	}
	return emailRegex.MatchString(str)
}

// GetAsJson returns the input as a JSON string. It directly returns strings to keep markdown compatibility.
func GetAsJson(i ...interface{}) string {
	if len(i) == 0 || (len(i) == 1 && i[0] == nil) {
		fmt.Println("No input provided or input is nil")
		return ""
	}

	// Check for array of nil values
	allNil := true
	for _, item := range i {
		if item != nil {
			allNil = false
			break
		}
	}
	if allNil {
		fmt.Println("Input is an array of nil values")
		return ""
	}

	// Check if the input is a single string to avoid unnecessary JSON encoding.
	if len(i) == 1 && reflect.TypeOf(i[0]).Kind() == reflect.String {
		return i[0].(string) // Return the string directly.
	}

	// Handle general case where input needs JSON encoding.
	j, err := json.Marshal(i)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return ""
	}
	return string(j)
}

// AppendQueryToURL takes a base URL and a map[string]interface{}, appending the converted URL query parameters to the base URL.
func AppendQueryToURL(baseURL string, params map[string]interface{}) string {
	if len(params) == 0 {
		return baseURL // Return base URL if there are no parameters to append
	}

	values := url.Values{}
	for key, value := range params {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case int:
			values.Add(key, fmt.Sprintf("%d", v))
		case []int:
			for _, iv := range v {
				values.Add(key, fmt.Sprintf("%d", iv))
			}
		case []string:
			for _, sv := range v {
				values.Add(key, sv)
			}
		case bool:
			values.Add(key, fmt.Sprintf("%t", v))
		default:
			// For types not handled, skip or handle accordingly
			fmt.Printf("Unhandled type for key %s: %T\n", key, v)
		}
	}

	encodedParams := values.Encode()
	if strings.Contains(baseURL, "?") {
		return baseURL + "&" + encodedParams
	}
	return baseURL + "?" + encodedParams
}
