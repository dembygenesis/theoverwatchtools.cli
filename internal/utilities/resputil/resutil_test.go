package resputil

import "testing"

func TestIsValidHTTPStatusCode(t *testing.T) {
	testCases := []struct {
		name     string
		code     int
		expected bool
	}{
		{"Valid lower bound", 100, true},
		{"Valid upper bound", 599, true},
		{"Valid middle", 404, true},
		{"Invalid low", 99, false},
		{"Invalid high", 600, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidHTTPStatusCode(tc.code); got != tc.expected {
				t.Errorf("IsValidHTTPStatusCode(%d) = %v, want %v", tc.code, got, tc.expected)
			}
		})
	}
}
