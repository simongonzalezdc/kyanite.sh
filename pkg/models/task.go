package models

import "time"

// RecurrencePattern defines how often a task recurs
type RecurrencePattern string

const (
	RecurrenceNone    RecurrencePattern = ""
	RecurrenceDaily   RecurrencePattern = "daily"
	RecurrenceWeekly  RecurrencePattern = "weekly"
	RecurrenceMonthly RecurrencePattern = "monthly"
	RecurrenceYearly  RecurrencePattern = "yearly"
)

// Task represents a single task in the focus list
type Task struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`   // pending, completed
	Priority    string    `json:"priority"` // low, medium, high
	Deadline    time.Time `json:"deadline,omitempty"`
	Categories  []string  `json:"categories,omitempty"`
	Notes       string    `json:"notes,omitempty"` // NEW: Markdown notes
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Recurrence fields
	RecurrencePattern  RecurrencePattern `json:"recurrence_pattern,omitempty"`  // daily, weekly, monthly, yearly
	RecurrenceInterval int               `json:"recurrence_interval,omitempty"` // e.g., every 2 days, every 3 weeks
	RecurrenceEndDate  time.Time         `json:"recurrence_end_date,omitempty"` // optional end date
	ParentTaskID       string            `json:"parent_task_id,omitempty"`      // for generated recurring instances
}

// IsRecurring returns true if this task has a recurrence pattern
func (t *Task) IsRecurring() bool {
	return t.RecurrencePattern != RecurrenceNone && t.RecurrencePattern != ""
}

// IsRecurringInstance returns true if this task was generated from a recurring parent
func (t *Task) IsRecurringInstance() bool {
	return t.ParentTaskID != ""
}

// ParsedTask represents the output from AI parsing
type ParsedTask struct {
	Description        string            `json:"description"`
	Deadline           time.Time         `json:"deadline,omitempty"`
	Priority           string            `json:"priority"` // low, medium, high
	Categories         []string          `json:"categories,omitempty"`
	Notes              string            `json:"notes,omitempty"`
	RecurrencePattern  RecurrencePattern `json:"recurrence_pattern,omitempty"`
	RecurrenceInterval int               `json:"recurrence_interval,omitempty"`
	RecurrenceEndDate  time.Time         `json:"recurrence_end_date,omitempty"`
}
