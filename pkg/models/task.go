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

	// Hierarchy fields
	ParentID string   `json:"parent_id,omitempty"` // for subtasks
	SubtaskIDs []string `json:"subtask_ids,omitempty"` // IDs of child tasks

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

// IsSubtask returns true if this task has a parent (is a subtask)
func (t *Task) IsSubtask() bool {
	return t.ParentID != ""
}

// HasSubtasks returns true if this task has any subtasks
func (t *Task) HasSubtasks() bool {
	return len(t.SubtaskIDs) > 0
}

// AddSubtask adds a subtask ID to this task
func (t *Task) AddSubtask(subtaskID string) {
	if t.SubtaskIDs == nil {
		t.SubtaskIDs = []string{}
	}
	t.SubtaskIDs = append(t.SubtaskIDs, subtaskID)
}

// RemoveSubtask removes a subtask ID from this task
func (t *Task) RemoveSubtask(subtaskID string) {
	if t.SubtaskIDs == nil {
		return
	}
	for i, id := range t.SubtaskIDs {
		if id == subtaskID {
			t.SubtaskIDs = append(t.SubtaskIDs[:i], t.SubtaskIDs[i+1:]...)
			return
		}
	}
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
