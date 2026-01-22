package noise

import (
	"regexp"
	"strings"
)

// Helper function to check if string contains substring, stripping ANSI escape codes
func contains(s, substr string) bool {
	// Simple ANSI escape code stripper
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cleanS := re.ReplaceAllString(s, "")
	cleanSubstr := re.ReplaceAllString(substr, "")

	return cleanSubstr == "" || strings.Contains(cleanS, cleanSubstr)
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
