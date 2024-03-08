package urlutil

import (
	"net/url"
	"strings"
)

// MapToQueryString takes a base URL as a string and a map of string slices representing query parameters,
// and returns the URL with the query parameters appended. If the base URL already contains query parameters,
// additional parameters are appended using '&'; otherwise, '?' is used to start the query string.
// Each key in the map can have multiple values, which are added as separate key-value pairs in the query string.
//
// Examples:
//
//   - Base URL without existing query parameters:
//     MapToQueryString("http://example.com/path", map[string][]string{"key": {"value1", "value2"}})
//     Returns: "http://example.com/path?key=value1&key=value2"
//
//   - Base URL with existing query parameters:
//     MapToQueryString("http://example.com/path?existing=param", map[string][]string{"key": {"value1", "value2"}})
//     Returns: "http://example.com/path?existing=param&key=value1&key=value2"
func MapToQueryString(baseUrl string, params map[string][]string) (string, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	var queryStrings []string
	for key, values := range params {
		for _, value := range values {
			queryStrings = append(queryStrings, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}

	fullQueryString := strings.Join(queryStrings, "&")

	var finalUrl string
	if parsedUrl.RawQuery == "" {
		finalUrl = baseUrl + "?" + fullQueryString
	} else {
		finalUrl = baseUrl + "&" + fullQueryString
	}

	return finalUrl, nil
}
