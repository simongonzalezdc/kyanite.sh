package ai

import (
	"strings"
)

// TaskValidator handles validation of parsed tasks
type TaskValidator struct{}

// NewTaskValidator creates a new task validator
func NewTaskValidator() *TaskValidator {
	return &TaskValidator{}
}

// Validate validates and normalizes a parsed task
func (tv *TaskValidator) Validate(task *ParsedTask) (*ParsedTask, bool) {
	if task == nil {
		return nil, false
	}

	// Validate description
	if task.Description == "" {
		return nil, false
	}

	// Check for low-quality descriptions
	if tv.IsLowQuality(task.Description) {
		return nil, false
	}

	// Normalize priority
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}

	if task.Priority != "" && !validPriorities[strings.ToLower(task.Priority)] {
		task.Priority = "medium" // Default to medium for invalid priorities
	}

	// Validate categories - filter out invalid ones
	validCategories := []string{}
	for _, category := range task.Categories {
		// Category must be between 2 and 20 characters
		if len(category) >= 2 && len(category) <= 20 {
			validCategories = append(validCategories, category)
		}
	}
	task.Categories = validCategories

	// Limit to 10 categories
	if len(task.Categories) > 10 {
		task.Categories = task.Categories[:10]
	}

	return task, true
}

// IsLowQuality checks if a description is too short or vague
func (tv *TaskValidator) IsLowQuality(description string) bool {
	description = strings.TrimSpace(description)

	// Too short
	if len(description) < 3 {
		return true
	}

	// Check for common low-quality patterns
	lowQualityPatterns := []string{
		"do something",
		"do stuff",
		"task",
		"todo",
		"fix",
		"update",
		"change",
	}

	descLower := strings.ToLower(description)
	for _, pattern := range lowQualityPatterns {
		if descLower == pattern {
			return true
		}
	}

	return false
}
