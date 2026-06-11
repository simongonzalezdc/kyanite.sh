package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Pre-compiled regular expressions for performance
var sqlInjectionRegexes = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(\bunion\b|\bselect\b|\binsert\b|\bupdate\b|\bdelete\b|\bdrop\b|\bcreate\b|\balter\b)`),
	regexp.MustCompile(`(?i)(\bexec\b|\bexecute\b|\bsp_|\bxp_)`),
	regexp.MustCompile(`(?i)(\bscript\b|\bjavascript\b|\bvbscript\b)`),
	regexp.MustCompile(`(?i)(['\"]\s*;\s*(drop|delete|union|select))`),
}

func ValidateTaskDescription(description string) error {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	if len(trimmed) < 2 {
		return fmt.Errorf("task description too short (min 2 characters)")
	}
	if len(trimmed) > 200 {
		return fmt.Errorf("task description too long (max 200 characters)")
	}
	if containsSQLInjection(trimmed) {
		return fmt.Errorf("task description contains invalid characters")
	}
	return nil
}

func ValidateTaskPriority(priority string) error {
	validPriorities := []string{"low", "medium", "high", "urgent"}
	for _, valid := range validPriorities {
		if priority == valid { // case-sensitive: must be exact lower-case
			return nil
		}
	}
	return fmt.Errorf("invalid priority: %s (must be one of: %s)", priority, strings.Join(validPriorities, ", "))
}


func ValidateTaskCategories(categories []string) error {
	if len(categories) > 10 {
		return fmt.Errorf("too many categories (max 10)")
	}
	seen := map[string]bool{}
	for i, category := range categories {
		trimmed := strings.TrimSpace(category)
		if trimmed == "" {
			return fmt.Errorf("category %d cannot be empty", i+1)
		}
		lower := strings.ToLower(trimmed)
		if seen[lower] {
			return fmt.Errorf("duplicate category: %s", trimmed)
		}
		seen[lower] = true
		if len(trimmed) > 50 || containsSQLInjection(trimmed) {
			return fmt.Errorf("invalid category: %s", trimmed)
		}
	}
	return nil
}

func ValidateTheme(theme string) error {
	validThemes := []string{"synthwave", "light", "plain"}
	for _, valid := range validThemes {
		if strings.ToLower(theme) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid theme: %s (must be one of: %s)", theme, strings.Join(validThemes, ", "))
}

func ValidateTimeFormat(format string) error {
	validFormats := []string{"12h", "24h"}
	for _, valid := range validFormats {
		if strings.ToLower(format) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid time format: %s (must be one of: %s)", format, strings.Join(validFormats, ", "))
}

func ValidateTaskDeadline(deadline string) error {
	if strings.TrimSpace(deadline) == "" {
		return errors.New("deadline cannot be empty")
	}
	layouts := []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05"}
	for _, layout := range layouts {
		if _, err := time.Parse(layout, deadline); err == nil {
			return nil
		}
	}
	return fmt.Errorf("invalid deadline format: %s", deadline)
}
func ValidateTask(task map[string]interface{}) error {
	descVal, ok := task["description"]
	if !ok {
		return fmt.Errorf("missing description")
	}
	desc, ok := descVal.(string)
	if !ok {
		return fmt.Errorf("description must be a string")
	}
	if err := ValidateTaskDescription(desc); err != nil {
		return err
	}
	if pr, ok := task["priority"].(string); ok && pr != "" {
		if err := ValidateTaskPriority(pr); err != nil {
			return err
		}
	}
	if dl, ok := task["deadline"].(string); ok && dl != "" {
		if err := ValidateTaskDeadline(dl); err != nil {
			return err
		}
	}
	if cats, ok := task["categories"].([]string); ok {
		if err := ValidateTaskCategories(cats); err != nil {
			return err
		}
	}
	return nil
}

func containsSQLInjection(input string) bool {
	for _, regex := range sqlInjectionRegexes {
		if regex.MatchString(input) {
			return true
		}
	}
	return false
}

func SanitizeInput(input string) string {
	sanitized := strings.ReplaceAll(input, "\x00", "")
	sanitized = strings.ReplaceAll(sanitized, "\r\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")
	// Use rune length to properly handle multi-byte UTF-8 characters
	runes := []rune(sanitized)
	if len(runes) > 1000 {
		sanitized = string(runes[:1000])
	}
	return strings.TrimSpace(sanitized)
}
