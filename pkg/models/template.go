package models

import "time"

// TaskTemplate represents a reusable task template
type TaskTemplate struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`                // Template name (e.g., "Daily Standup")
	Description        string            `json:"description"`         // Task description template
	Priority           string            `json:"priority,omitempty"`  // Default priority
	Categories         []string          `json:"categories,omitempty"`// Default categories
	Notes              string            `json:"notes,omitempty"`     // Default notes
	RecurrencePattern  RecurrencePattern `json:"recurrence_pattern,omitempty"`
	RecurrenceInterval int               `json:"recurrence_interval,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// ToTask converts a template to a task
func (t *TaskTemplate) ToTask() ParsedTask {
	return ParsedTask{
		Description:        t.Description,
		Priority:           t.Priority,
		Categories:         t.Categories,
		Notes:              t.Notes,
		RecurrencePattern:  t.RecurrencePattern,
		RecurrenceInterval: t.RecurrenceInterval,
	}
}
