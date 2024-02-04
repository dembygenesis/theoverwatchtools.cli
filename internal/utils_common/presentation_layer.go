package utils_common

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"regexp"
	"strings"
)

func GetJSONAndCopyToClipboard(i ...interface{}) string {
	var jsonStr string
	if getInterfaceArgLength(i) > 1 {
		jsonStr = GetJSON(i)
	} else {
		jsonStr = GetJSON(i[0])
	}

	err := clipboard.WriteAll(jsonStr)
	if err != nil {

	}
	return jsonStr
}

func prettifyJSON(s string) string {
	s = strings.ReplaceAll(s, `\t`, " ")
	s = strings.ReplaceAll(s, `\n`, "")

	reLeadCloseWhiteSpace := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reInsideWhiteSpace := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	formattedStr := reLeadCloseWhiteSpace.ReplaceAllString(s, "")
	formattedStr = reInsideWhiteSpace.ReplaceAllString(s, " ")

	return formattedStr
}

func getInterfaceArgLength(i ...interface{}) int {
	items := make([]interface{}, 0)
	for _, v := range i {
		items = append(items, v)
	}
	return len(items)
}

// GetJSON return an interface as type json
// Returns a generic error formatted string if fails
func GetJSON(i ...interface{}) string {
	if getInterfaceArgLength(i) > 1 {
		d := make([]interface{}, 0)
		for _, v := range i {
			d = append(d, v)
		}
		j, err := json.Marshal(d)
		if err != nil {
			return fmt.Sprintf("Error dumping to json: %v", err)
		} else {
			return string(j)
			return prettifyJSON(string(j))
		}
	} else {
		j, err := json.Marshal(i[0])
		if err != nil {
			return fmt.Sprintf("Error dumping to json: %v", err)
		} else {
			return string(j)
			return prettifyJSON(string(j))
		}
	}
}
