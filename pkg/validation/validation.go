package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

func ValidateCategories(categories []string) error { // kept for backward compat
	return ValidateTaskCategories(categories)
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

func ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email cannot be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email: %s", email)
	}
	if strings.Contains(email, "..") {
		return fmt.Errorf("invalid email: %s", email)
	}
	return nil
}

func ValidateUrl(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("url cannot be empty")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid url: %s", raw)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid url scheme: %s", u.Scheme)
	}
	return nil
}

func ValidatePhoneNumber(phone string) error {
	if strings.TrimSpace(phone) == "" {
		return fmt.Errorf("phone number cannot be empty")
	}
	re := regexp.MustCompile(`^\+?[0-9\s\-()]{7,}$`)
	if !re.MatchString(phone) {
		return fmt.Errorf("invalid phone number: %s", phone)
	}
	return nil
}

func ValidatePositiveInteger(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("value cannot be empty")
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return fmt.Errorf("must be a positive integer: %s", s)
	}
	return nil
}

func ValidatePositiveNumber(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("value cannot be empty")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f <= 0 {
		return fmt.Errorf("must be a positive number: %s", s)
	}
	return nil
}

func ValidateString(s string, minLen, maxLen int) error {
	if minLen < 0 || maxLen < minLen {
		return fmt.Errorf("invalid length bounds")
	}
	if len(s) < minLen {
		return fmt.Errorf("string too short (min %d)", minLen)
	}
	if len(s) > maxLen {
		return fmt.Errorf("string too long (max %d)", maxLen)
	}
	return nil
}

func ValidateDateTime(dt string) error {
	if strings.TrimSpace(dt) == "" {
		return fmt.Errorf("datetime cannot be empty")
	}
	layouts := []string{time.RFC3339, "2006-01-02T15:04:05"}
	for _, l := range layouts {
		if _, err := time.Parse(l, dt); err == nil {
			return nil
		}
	}
	return fmt.Errorf("invalid datetime: %s", dt)
}

func ValidatePassword(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf("password too short (min 8)")
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pw)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pw)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(pw)
	if !(hasLower && hasUpper && hasDigit) {
		return fmt.Errorf("password must contain upper, lower, and digit")
	}
	return nil
}

func ValidateTags(tags []string) error {
	seen := map[string]bool{}
	for i, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			return fmt.Errorf("tag %d cannot be empty", i+1)
		}
		if strings.Contains(trimmed, " ") {
			return fmt.Errorf("tag cannot contain spaces: %s", trimmed)
		}
		lower := strings.ToLower(trimmed)
		if seen[lower] {
			return fmt.Errorf("duplicate tag: %s", trimmed)
		}
		seen[lower] = true
	}
	return nil
}

func ValidateTask(task map[string]interface{}) error {
	descVal, ok := task["description"]
	if !ok {
		return fmt.Errorf("missing description")
	}
	desc, _ := descVal.(string)
	if err := ValidateTaskDescription(desc); err != nil {
		return err
	}
	if pr, ok := task["priority"].(string); ok && pr != "" {
		if err := ValidateTaskPriority(pr); err != nil { return err }
	}
	if dl, ok := task["deadline"].(string); ok && dl != "" {
		if err := ValidateTaskDeadline(dl); err != nil { return err }
	}
	if cats, ok := task["categories"].([]string); ok {
		if err := ValidateTaskCategories(cats); err != nil { return err }
	}
	return nil
}

func containsSQLInjection(input string) bool {
	sqlPatterns := []string{
		`(?i)(\bunion\b|\bselect\b|\binsert\b|\bupdate\b|\bdelete\b|\bdrop\b|\bcreate\b|\balter\b)`,
		`(?i)(\bexec\b|\bexecute\b|\bsp_|\bxp_)`,
		`(?i)(\bscript\b|\bjavascript\b|\bvbscript\b)`,
		`(?i)(['\"]\s*;\s*(drop|delete|union|select))`,
	}
	for _, pattern := range sqlPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched { return true }
	}
	return false
}

func SanitizeInput(input string) string {
	sanitized := strings.ReplaceAll(input, "\x00", "")
	sanitized = strings.ReplaceAll(sanitized, "\r\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")
	if len(sanitized) > 1000 { sanitized = sanitized[:1000] }
	return strings.TrimSpace(sanitized)
}
