package utils_common

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"reflect"
)

// GetJSONAndCopyToClipboard generates a JSON string from the input and copies it to the clipboard.
func GetJSONAndCopyToClipboard(i ...interface{}) string {
	var jsonStr string
	jsonStr = GetJSON(i...) // Simplified call

	err := clipboard.WriteAll(jsonStr)
	if err != nil {
		// Consider logging the error or handling it as needed.
		fmt.Printf("Error copying to clipboard: %v\n", err)
	}
	return jsonStr
}

// GetJSON returns the input as a JSON string. It directly returns strings to keep markdown compatibility.
func GetJSON(i ...interface{}) string {
	// Check if the input is a single string to avoid unnecessary JSON encoding.
	if len(i) == 1 && reflect.TypeOf(i[0]).Kind() == reflect.String {
		return i[0].(string) // Return the string directly.
	}

	// Handle general case where input needs JSON encoding.
	j, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("Error converting to JSON: %v", err)
	}
	return string(j)
}
