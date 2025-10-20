package validation

import (
	"strings"
	"testing"
)

func TestValidateTask(t *testing.T) {
	// Test valid task
	validTask := map[string]interface{}{
		"description": "Test task description",
		"priority":     "medium",
		"deadline":     "2025-01-01",
		"categories":   []string{"work", "personal"},
		"notes":        "Task notes content",
	}
	
	err := ValidateTask(validTask)
	if err != nil {
		t.Errorf("Valid task should pass validation: %v", err)
	}
	
	// Test invalid task (missing description)
	invalidTask := map[string]interface{}{
		"priority":   "medium",
		"deadline":   "2025-01-01",
	}
	
	err = ValidateTask(invalidTask)
	if err == nil {
		t.Error("Invalid task should fail validation")
	}
}

func TestValidateTaskDescription(t *testing.T) {
	// Test valid descriptions
	validDescriptions := []string{
		"Complete project documentation",
		"Fix bug in user authentication",
		"Design new landing page",
		"Write unit tests for API endpoints",
		"Review pull requests",
	}
	
	for _, desc := range validDescriptions {
		err := ValidateTaskDescription(desc)
		if err != nil {
			t.Errorf("Valid description should pass validation: '%s' - %v", desc, err)
		}
	}
	
	// Test invalid descriptions
	invalidDescriptions := []string{
		"",                    // Empty
		"   ",                 // Spaces only
		"a",                   // Too short
		strings.Repeat("a", 300), // Too long
	}
	
	for _, desc := range invalidDescriptions {
		err := ValidateTaskDescription(desc)
		if err == nil {
			t.Errorf("Invalid description should fail validation: '%s'", desc)
		}
	}
}

func TestValidateTaskPriority(t *testing.T) {
	// Test valid priorities
	validPriorities := []string{"low", "medium", "high", "urgent"}
	
	for _, priority := range validPriorities {
		err := ValidateTaskPriority(priority)
		if err != nil {
			t.Errorf("Valid priority should pass validation: '%s' - %v", priority, err)
		}
	}
	
	// Test invalid priorities
	invalidPriorities := []string{
		"",           // Empty
		"invalid",    // Not in valid list
		"LOW",        // Wrong case
		"medium ",    // Trailing space
	}
	
	for _, priority := range invalidPriorities {
		err := ValidateTaskPriority(priority)
		if err == nil {
			t.Errorf("Invalid priority should fail validation: '%s'", priority)
		}
	}
}

func TestValidateTaskDeadline(t *testing.T) {
	// Test valid deadlines
	validDeadlines := []string{
		"2025-01-01",
		"2025-12-31",
		"2024-02-29", // Leap year
	}
	
	for _, deadline := range validDeadlines {
		err := ValidateTaskDeadline(deadline)
		if err != nil {
			t.Errorf("Valid deadline should pass validation: '%s' - %v", deadline, err)
		}
	}
	
	// Test invalid deadlines
	invalidDeadlines := []string{
		"",                // Empty
		"invalid-date",    // Invalid format
		"2025-13-01",      // Invalid month
		"2025-02-30",      // Invalid day (non-leap year)
		"2025-01-32",      // Invalid day
	}
	
	for _, deadline := range invalidDeadlines {
		err := ValidateTaskDeadline(deadline)
		if err == nil {
			t.Errorf("Invalid deadline should fail validation: '%s'", deadline)
		}
	}
}

func TestValidateTaskCategories(t *testing.T) {
	// Test valid categories
	validCategories := [][]string{
		{"work", "personal"},
		{"development", "design"},
		{"urgent", "project"},
		{}, // Empty categories
	}
	
	for _, categories := range validCategories {
		err := ValidateTaskCategories(categories)
		if err != nil {
			t.Errorf("Valid categories should pass validation: %v - %v", categories, err)
		}
	}
	
	// Test invalid categories
	invalidCategories := [][]string{
		{"", "work"},        // Empty string
		{"work", "personal", "work"}, // Duplicate
	}
	
	for _, categories := range invalidCategories {
		err := ValidateTaskCategories(categories)
		if err == nil {
			t.Errorf("Invalid categories should fail validation: %v", categories)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	// Test valid emails
	validEmails := []string{
		"user@example.com",
		"test.email+tag@domain.co.uk",
		"first.last@company.org",
		"user123@test-domain.com",
	}
	
	for _, email := range validEmails {
		err := ValidateEmail(email)
		if err != nil {
			t.Errorf("Valid email should pass validation: '%s' - %v", email, err)
		}
	}
	
	// Test invalid emails
	invalidEmails := []string{
		"",                    // Empty
		"invalid-email",       // No @
		"@domain.com",         // No local part
		"user@",               // No domain
		"user@.com",           // Invalid domain
		"user..name@domain.com", // Double dot
	}
	
	for _, email := range invalidEmails {
		err := ValidateEmail(email)
		if err == nil {
			t.Errorf("Invalid email should fail validation: '%s'", email)
		}
	}
}

func TestValidateUrl(t *testing.T) {
	// Test valid URLs
	validUrls := []string{
		"https://example.com",
		"http://localhost:8080",
		"https://api.example.com/v1/users",
		"https://sub.domain.co.uk/path?query=value",
	}
	
	for _, url := range validUrls {
		err := ValidateUrl(url)
		if err != nil {
			t.Errorf("Valid URL should pass validation: '%s' - %v", url, err)
		}
	}
	
	// Test invalid URLs
	invalidUrls := []string{
		"",                     // Empty
		"not-a-url",            // Invalid format
		"http://",              // No domain
		"https://",             // No domain
		"ftp://example.com",    // Invalid protocol
	}
	
	for _, url := range invalidUrls {
		err := ValidateUrl(url)
		if err == nil {
			t.Errorf("Invalid URL should fail validation: '%s'", url)
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	// Test valid phone numbers
	validPhoneNumbers := []string{
		"+1234567890",
		"+1 (555) 123-4567",
		"555-123-4567",
		"(555) 123-4567",
		"5551234567",
	}
	
	for _, phone := range validPhoneNumbers {
		err := ValidatePhoneNumber(phone)
		if err != nil {
			t.Errorf("Valid phone number should pass validation: '%s' - %v", phone, err)
		}
	}
	
	// Test invalid phone numbers
	invalidPhoneNumbers := []string{
		"",               // Empty
		"123",           // Too short
		"invalid-phone", // Invalid format
		"phone",         // Contains letters
	}
	
	for _, phone := range invalidPhoneNumbers {
		err := ValidatePhoneNumber(phone)
		if err == nil {
			t.Errorf("Invalid phone number should fail validation: '%s'", phone)
		}
	}
}

func TestValidatePositiveInteger(t *testing.T) {
	// Test valid positive integers
	validValues := []string{
		"1",
		"10",
		"100",
		"999999",
	}
	
	for _, value := range validValues {
		err := ValidatePositiveInteger(value)
		if err != nil {
			t.Errorf("Valid positive integer should pass validation: '%s' - %v", value, err)
		}
	}
	
	// Test invalid positive integers
	invalidValues := []string{
		"",        // Empty
		"0",       // Zero
		"-1",      // Negative
		"invalid", // Not a number
		"1.5",     // Not integer
	}
	
	for _, value := range invalidValues {
		err := ValidatePositiveInteger(value)
		if err == nil {
			t.Errorf("Invalid positive integer should fail validation: '%s'", value)
		}
	}
}

func TestValidatePositiveNumber(t *testing.T) {
	// Test valid positive numbers
	validValues := []string{
		"1",
		"10",
		"100",
		"999999",
		"1.5",
		"100.99",
	}
	
	for _, value := range validValues {
		err := ValidatePositiveNumber(value)
		if err != nil {
			t.Errorf("Valid positive number should pass validation: '%s' - %v", value, err)
		}
	}
	
	// Test invalid positive numbers
	invalidValues := []string{
		"",        // Empty
		"0",       // Zero
		"-1",      // Negative
		"-1.5",    // Negative decimal
		"invalid", // Not a number
	}
	
	for _, value := range invalidValues {
		err := ValidatePositiveNumber(value)
		if err == nil {
			t.Errorf("Invalid positive number should fail validation: '%s'", value)
		}
	}
}

func TestValidateString(t *testing.T) {
	// Test string validation with different constraints
	
	// Test non-empty string
	testCases := []struct {
		name      string
		input     string
		minLength int
		maxLength int
		shouldErr bool
	}{
		{"Non-empty string", "hello", 1, 100, false},
		{"Empty string", "", 1, 100, true},
		{"Too short", "ab", 3, 100, true},
		{"Too long", strings.Repeat("a", 200), 1, 100, true},
		{"Just right", "test", 1, 100, false},
	}
	
	for _, tc := range testCases {
		err := ValidateString(tc.input, tc.minLength, tc.maxLength)
		
		if tc.shouldErr && err == nil {
			t.Errorf("Test case '%s' should fail but passed: '%s'", tc.name, tc.input)
		}
		
		if !tc.shouldErr && err != nil {
			t.Errorf("Test case '%s' should pass but failed: '%s' - %v", tc.name, tc.input, err)
		}
	}
}

func TestValidateDateTime(t *testing.T) {
	// Test valid date-time strings
	validDateTimes := []string{
		"2025-01-01T12:00:00Z",
		"2025-12-31T23:59:59Z",
		"2024-02-29T12:00:00Z", // Leap year
		"2025-01-01T00:00:00",   // No timezone
	}
	
	for _, dt := range validDateTimes {
		err := ValidateDateTime(dt)
		if err != nil {
			t.Errorf("Valid datetime should pass validation: '%s' - %v", dt, err)
		}
	}
	
	// Test invalid date-time strings
	invalidDateTimes := []string{
		"",                     // Empty
		"invalid-datetime",    // Invalid format
		"2025-13-01T12:00:00Z", // Invalid month
		"2025-01-32T12:00:00Z", // Invalid day
		"2025-01-01T24:00:00Z", // Invalid hour
	}
	
	for _, dt := range invalidDateTimes {
		err := ValidateDateTime(dt)
		if err == nil {
			t.Errorf("Invalid datetime should fail validation: '%s'", dt)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	// Test valid passwords
	validPasswords := []string{
		"SecureP@ssw0rd!",
		"MyStr0ngPass123",
		"ComplexP@ss#123",
		"ValidPass2025!",
	}
	
	for _, password := range validPasswords {
		err := ValidatePassword(password)
		if err != nil {
			t.Errorf("Valid password should pass validation: '%s' - %v", password, err)
		}
	}
	
	// Test invalid passwords
	invalidPasswords := []string{
		"",               // Empty
		"123456",         // Only numbers
		"password",       // Only lowercase
		"PASSWORD",       // Only uppercase
		"pass",           // Too short
	}
	
	for _, password := range invalidPasswords {
		err := ValidatePassword(password)
		if err == nil {
			t.Errorf("Invalid password should fail validation: '%s'", password)
		}
	}
}

func TestValidateTags(t *testing.T) {
	// Test valid tags
	validTags := [][]string{
		{"work", "personal"},
		{"development", "design", "testing"},
		{"urgent", "important"},
		{}, // Empty tags
	}
	
	for _, tags := range validTags {
		err := ValidateTags(tags)
		if err != nil {
			t.Errorf("Valid tags should pass validation: %v - %v", tags, err)
		}
	}
	
	// Test invalid tags
	invalidTags := [][]string{
		{"", "work"},        // Empty tag
		{"work", "work"},    // Duplicate tag
		{"work", "work work"}, // Tag with spaces
	}
	
	for _, tags := range invalidTags {
		err := ValidateTags(tags)
		if err == nil {
			t.Errorf("Invalid tags should fail validation: %v", tags)
		}
	}
}

func TestValidationErrors(t *testing.T) {
	// Test that validation errors provide helpful messages
	type testCase struct {
		name  string
		input string
		test  func(string) error
	}
	
	testCases := []testCase{
		{"Empty email", "", ValidateEmail},
		{"Invalid priority", "invalid", ValidateTaskPriority},
		{"Empty description", "", ValidateTaskDescription},
		{"Too short", "ab", func(s string) error { return ValidateString(s, 3, 100) }},
	}
	
	for _, tc := range testCases {
		err := tc.test(tc.input)
		if err == nil {
			t.Errorf("Test case '%s' should return error for input '%s'", tc.name, tc.input)
		}
		
		// Check that error message is not empty
		if err != nil && err.Error() == "" {
			t.Errorf("Test case '%s' should return non-empty error message", tc.name)
		}
	}
}