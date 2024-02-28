package utility

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"reflect"
)

// GetJSONAndCopyToClipboard generates a JSON string from the input and copies it to the clipboard.
func GetJSONAndCopyToClipboard(i ...interface{}) string {
	if len(i) == 0 || (len(i) == 1 && i[0] == nil) {
		fmt.Println("No input provided or input is nil")
		return ""
	}

	jsonStr := GetJSON(i...) // Simplified call

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

// GetJSON returns the input as a JSON string. It directly returns strings to keep markdown compatibility.
func GetJSON(i ...interface{}) string {
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
