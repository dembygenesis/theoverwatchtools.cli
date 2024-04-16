package resputil

// IsValidHTTPStatusCode checks if the given status code is a valid HTTP status code.
func IsValidHTTPStatusCode(code int) bool {
	return code >= 100 && code <= 599
}
