package noise

import "strings"

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return substr == "" || strings.Contains(s, substr)
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
