package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateTaskDescription validates task description input
func ValidateTaskDescription(description string) error {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	if len(trimmed) > 1000 {
		return fmt.Errorf("task description too long (max 1000 characters)")
	}
	// Check for potentially dangerous content
	if containsSQLInjection(trimmed) {
		return fmt.Errorf("task description contains invalid characters")
	}
	return nil
}

// ValidateTaskPriority validates task priority
func ValidateTaskPriority(priority string) error {
	validPriorities := []string{"low", "medium", "high", "urgent"}
	for _, valid := range validPriorities {
		if strings.ToLower(priority) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid priority: %s (must be one of: %s)", priority, strings.Join(validPriorities, ", "))
}

// ValidateCategories validates task categories
func ValidateCategories(categories []string) error {
	if len(categories) > 10 {
		return fmt.Errorf("too many categories (max 10)")
	}
	
	validCategories := []string{"work", "personal", "urgent", "learning", "health", "finance", "creative", "shopping", "home", "other"}
	for i, category := range categories {
		trimmed := strings.TrimSpace(category)
		if trimmed == "" {
			return fmt.Errorf("category %d cannot be empty", i+1)
		}
		
		// Check if category is valid or custom
		isValid := false
		for _, valid := range validCategories {
			if strings.ToLower(trimmed) == valid {
				isValid = true
				break
			}
		}
		
		if !isValid {
			// Allow custom categories but validate they're safe
			if len(trimmed) > 50 || containsSQLInjection(trimmed) {
				return fmt.Errorf("invalid category: %s", trimmed)
			}
		}
	}
	return nil
}

// ValidateTheme validates theme selection
func ValidateTheme(theme string) error {
	validThemes := []string{"synthwave", "light", "plain"}
	for _, valid := range validThemes {
		if strings.ToLower(theme) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid theme: %s (must be one of: %s)", theme, strings.Join(validThemes, ", "))
}

// ValidateTimeFormat validates time format
func ValidateTimeFormat(format string) error {
	validFormats := []string{"12h", "24h"}
	for _, valid := range validFormats {
		if strings.ToLower(format) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid time format: %s (must be one of: %s)", format, strings.Join(validFormats, ", "))
}

// containsSQLInjection checks for potential SQL injection patterns
func containsSQLInjection(input string) bool {
	// More targeted regex for SQL injection patterns
	sqlPatterns := []string{
		`(?i)(\bunion\b|\bselect\b|\binsert\b|\bupdate\b|\bdelete\b|\bdrop\b|\bcreate\b|\balter\b)`,
		`(?i)(\bexec\b|\bexecute\b|\bsp_|\bxp_)`,
		`(?i)(\bscript\b|\bjavascript\b|\bvbscript\b)`,
		`(?i)(['\"]\s*;\s*(drop|delete|union|select))`,
	}
	
	for _, pattern := range sqlPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return true
		}
	}
	return false
}

// SanitizeInput removes or sanitizes potentially harmful content
func SanitizeInput(input string) string {
	// Remove potentially dangerous characters
	sanitized := strings.ReplaceAll(input, "\x00", "")
	sanitized = strings.ReplaceAll(sanitized, "\r\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")
	
	// Limit length
	if len(sanitized) > 1000 {
		sanitized = sanitized[:1000]
	}
	
	return strings.TrimSpace(sanitized)
}
