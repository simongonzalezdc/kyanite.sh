package errors

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ValidationRule represents a validation rule for input data
type ValidationRule struct {
	Field   string
	Value   interface{}
	Rule    string
	Message string
}

// ValidationResult contains the result of validation
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Validator provides input validation functionality
type Validator struct {
	rules map[string][]ValidationRule
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string][]ValidationRule),
	}
}

// AddRule adds a validation rule for a field
func (v *Validator) AddRule(field string, rule ValidationRule) {
	v.rules[field] = append(v.rules[field], rule)
}

// Validate validates a value against all rules for its field
func (v *Validator) Validate(field string, value interface{}) []ValidationError {
	var errors []ValidationError

	rules, exists := v.rules[field]
	if !exists {
		return errors
	}

	for _, rule := range rules {
		if !v.validateRule(rule, value) {
			errors = append(errors, ValidationError{
				Field:   field,
				Message: rule.Message,
				Code:    rule.Rule,
			})
		}
	}

	return errors
}

// ValidateStruct validates a struct with tagged fields
func (v *Validator) ValidateStruct(s interface{}) ValidationResult {
	// This is a simplified implementation
	// In a full implementation, this would use reflection to validate struct fields
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	// For now, return valid - this would be expanded with proper struct validation
	return result
}

// validateRule validates a single rule against a value
func (v *Validator) validateRule(rule ValidationRule, value interface{}) bool {
	switch rule.Rule {
	case "required":
		return v.validateRequired(value)
	case "min_length":
		return v.validateMinLength(value, rule.Message)
	case "max_length":
		return v.validateMaxLength(value, rule.Message)
	case "email":
		return v.validateEmail(value)
	case "numeric":
		return v.validateNumeric(value)
	case "min":
		return v.validateMin(value, rule.Message)
	case "max":
		return v.validateMax(value, rule.Message)
	case "pattern":
		return v.validatePattern(value, rule.Message)
	case "not_empty":
		return v.validateNotEmpty(value)
	case "valid_filepath":
		return v.validateFilepath(value)
	case "valid_song_title":
		return v.validateSongTitle(value)
	case "valid_artist_name":
		return v.validateArtistName(value)
	case "valid_key":
		return v.validateMusicKey(value)
	case "valid_tempo":
		return v.validateTempo(value)
	case "valid_time_signature":
		return v.validateTimeSignature(value)
	default:
		return true // Unknown rule passes by default
	}
}

// Individual validation functions

func (v *Validator) validateRequired(value interface{}) bool {
	if value == nil {
		return false
	}

	switch val := value.(type) {
	case string:
		return strings.TrimSpace(val) != ""
	case int, int32, int64:
		return true // Numbers are considered present if not nil
	default:
		return true // Other types are considered present if not nil
	}
}

func (v *Validator) validateMinLength(value interface{}, param string) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	minLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	return utf8.RuneCountInString(str) >= minLen
}

func (v *Validator) validateMaxLength(value interface{}, param string) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	maxLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	return utf8.RuneCountInString(str) <= maxLen
}

func (v *Validator) validateEmail(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	// Simple email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(str)
}

func (v *Validator) validateNumeric(value interface{}) bool {
	switch val := value.(type) {
	case string:
		_, err := strconv.ParseFloat(val, 64)
		return err == nil
	case int, int32, int64, float32, float64:
		return true
	default:
		return false
	}
}

func (v *Validator) validateMin(value interface{}, param string) bool {
	minVal, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return false
	}

	switch val := value.(type) {
	case string:
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return false
		}
		return num >= minVal
	case int:
		return float64(val) >= minVal
	case int32:
		return float64(val) >= minVal
	case int64:
		return float64(val) >= minVal
	case float32:
		return float64(val) >= minVal
	case float64:
		return val >= minVal
	default:
		return false
	}
}

func (v *Validator) validateMax(value interface{}, param string) bool {
	maxVal, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return false
	}

	switch val := value.(type) {
	case string:
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return false
		}
		return num <= maxVal
	case int:
		return float64(val) <= maxVal
	case int32:
		return float64(val) <= maxVal
	case int64:
		return float64(val) <= maxVal
	case float32:
		return float64(val) <= maxVal
	case float64:
		return val <= maxVal
	default:
		return false
	}
}

func (v *Validator) validatePattern(value interface{}, pattern string) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return regex.MatchString(str)
}

func (v *Validator) validateNotEmpty(value interface{}) bool {
	return v.validateRequired(value)
}

func (v *Validator) validateFilepath(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	// Basic filepath validation - not empty and doesn't contain invalid characters
	if strings.TrimSpace(str) == "" {
		return false
	}

	// Check for invalid characters (platform specific)
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(str, char) {
			return false
		}
	}

	return true
}

func (v *Validator) validateSongTitle(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	// Song title should not be empty and should be reasonable length
	title := strings.TrimSpace(str)
	if title == "" {
		return false
	}

	// Reasonable length limits for song titles
	return utf8.RuneCountInString(title) <= 200
}

func (v *Validator) validateArtistName(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	// Artist name should not be empty and should be reasonable length
	name := strings.TrimSpace(str)
	if name == "" {
		return false
	}

	// Reasonable length limits for artist names
	return utf8.RuneCountInString(name) <= 100
}

func (v *Validator) validateMusicKey(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	key := strings.TrimSpace(strings.ToUpper(str))
	if key == "" {
		return true // Empty key is allowed
	}

	// Valid musical keys
	validKeys := []string{
		"C", "C#", "DB", "D", "D#", "EB", "E", "F", "F#", "GB", "G", "G#", "AB", "A", "A#", "BB", "B",
		"CM", "C#M", "DBM", "DM", "D#M", "EBM", "EM", "FM", "F#M", "GBM", "GM", "G#M", "ABM", "AM", "A#M", "BBM", "BM",
	}

	for _, validKey := range validKeys {
		if key == validKey {
			return true
		}
	}

	return false
}

func (v *Validator) validateTempo(value interface{}) bool {
	switch val := value.(type) {
	case string:
		if strings.TrimSpace(val) == "" {
			return true // Empty tempo is allowed
		}
		tempo, err := strconv.Atoi(val)
		if err != nil {
			return false
		}
		return tempo > 0 && tempo <= 300 // Reasonable tempo range
	case int:
		return val > 0 && val <= 300
	case nil:
		return true // Nil tempo is allowed
	default:
		return false
	}
}

func (v *Validator) validateTimeSignature(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	timeSig := strings.TrimSpace(str)
	if timeSig == "" {
		return true // Empty time signature is allowed
	}

	// Common time signatures (numerator/denominator)
	commonTimeSigs := []string{
		"2/4", "3/4", "4/4", "5/4", "6/4", "7/4", "2/2", "3/2", "4/2",
		"3/8", "5/8", "6/8", "7/8", "9/8", "12/8",
	}

	for _, validSig := range commonTimeSigs {
		if timeSig == validSig {
			return true
		}
	}

	return false
}

// Validation helper functions

// ValidateSongMetadata validates song metadata fields
func ValidateSongMetadata(title, artist, key string, tempo int, timeSignature string) ValidationResult {
	validator := NewValidator()
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	// Add validation rules
	validator.AddRule("title", ValidationRule{
		Field:   "title",
		Value:   title,
		Rule:    "valid_song_title",
		Message: "Song title is required and must be less than 200 characters",
	})

	validator.AddRule("artist", ValidationRule{
		Field:   "artist",
		Value:   artist,
		Rule:    "valid_artist_name",
		Message: "Artist name must be less than 100 characters",
	})

	if key != "" {
		validator.AddRule("key", ValidationRule{
			Field:   "key",
			Value:   key,
			Rule:    "valid_key",
			Message: "Invalid musical key",
		})
	}

	if tempo > 0 {
		validator.AddRule("tempo", ValidationRule{
			Field:   "tempo",
			Value:   tempo,
			Rule:    "valid_tempo",
			Message: "Tempo must be between 1 and 300 BPM",
		})
	}

	if timeSignature != "" {
		validator.AddRule("time_signature", ValidationRule{
			Field:   "time_signature",
			Value:   timeSignature,
			Rule:    "valid_time_signature",
			Message: "Invalid time signature",
		})
	}

	// Validate each field
	fields := []string{"title", "artist", "key", "tempo", "time_signature"}
	for _, field := range fields {
		var value interface{}
		switch field {
		case "title":
			value = title
		case "artist":
			value = artist
		case "key":
			value = key
		case "tempo":
			value = tempo
		case "time_signature":
			value = timeSignature
		}

		errors := validator.Validate(field, value)
		result.Errors = append(result.Errors, errors...)
	}

	result.Valid = len(result.Errors) == 0
	return result
}

// ValidateFilepath validates a file path
func ValidateFilepath(filepath string) ValidationResult {
	validator := NewValidator()
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	validator.AddRule("filepath", ValidationRule{
		Field:   "filepath",
		Value:   filepath,
		Rule:    "required",
		Message: "File path is required",
	})

	validator.AddRule("filepath", ValidationRule{
		Field:   "filepath",
		Value:   filepath,
		Rule:    "valid_filepath",
		Message: "Invalid file path",
	})

	errors := validator.Validate("filepath", filepath)
	result.Errors = errors
	result.Valid = len(errors) == 0

	return result
}

// ValidateSearchQuery validates a search query
func ValidateSearchQuery(query string) ValidationResult {
	validator := NewValidator()
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	validator.AddRule("query", ValidationRule{
		Field:   "query",
		Value:   query,
		Rule:    "max_length:500",
		Message: "Search query must be less than 500 characters",
	})

	errors := validator.Validate("query", query)
	result.Errors = errors
	result.Valid = len(errors) == 0

	return result
}

// ValidateCommand validates a command string
func ValidateCommand(command string) ValidationResult {
	validator := NewValidator()
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	validator.AddRule("command", ValidationRule{
		Field:   "command",
		Value:   command,
		Rule:    "required",
		Message: "Command is required",
	})

	validator.AddRule("command", ValidationRule{
		Field:   "command",
		Value:   command,
		Rule:    "max_length:100",
		Message: "Command must be less than 100 characters",
	})

	errors := validator.Validate("command", command)
	result.Errors = errors
	result.Valid = len(errors) == 0

	return result
}

// ValidateSection validates a song section
func ValidateSection(sectionType string, content string) ValidationResult {
	validator := NewValidator()
	result := ValidationResult{Valid: true, Errors: []ValidationError{}}

	validator.AddRule("section_type", ValidationRule{
		Field:   "section_type",
		Value:   sectionType,
		Rule:    "required",
		Message: "Section type is required",
	})

	validator.AddRule("content", ValidationRule{
		Field:   "content",
		Value:   content,
		Rule:    "max_length:10000",
		Message: "Section content must be less than 10000 characters",
	})

	errors := validator.Validate("section_type", sectionType)
	result.Errors = append(result.Errors, errors...)

	errors = validator.Validate("content", content)
	result.Errors = append(result.Errors, errors...)

	result.Valid = len(result.Errors) == 0
	return result
}

// IsValid returns true if the validation result is valid
func (vr ValidationResult) IsValid() bool {
	return vr.Valid
}

// Error returns a formatted error message from validation errors
func (vr ValidationResult) Error() string {
	if vr.Valid {
		return ""
	}

	var messages []string
	for _, err := range vr.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (vr ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// GetErrors returns all validation errors
func (vr ValidationResult) GetErrors() []ValidationError {
	return vr.Errors
}
