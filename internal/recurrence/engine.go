package recurrence

import (
	"time"

	"github.com/kyanite/focus/pkg/models"
)

// Engine handles recurring task generation
type Engine struct{}

// New creates a new recurrence engine
func New() *Engine {
	return &Engine{}
}

// GenerateInstances generates task instances for a recurring task up to a given date
func (e *Engine) GenerateInstances(parentTask models.Task, until time.Time) []models.Task {
	if !parentTask.IsRecurring() {
		return nil
	}

	instances := make([]models.Task, 0)
	current := parentTask.Deadline

	if current.IsZero() {
		current = parentTask.CreatedAt
	}

	// Set default interval if not specified
	interval := parentTask.RecurrenceInterval
	if interval <= 0 {
		interval = 1
	}

	// Generate instances until the specified date or recurrence end date
	endDate := until
	if !parentTask.RecurrenceEndDate.IsZero() && parentTask.RecurrenceEndDate.Before(until) {
		endDate = parentTask.RecurrenceEndDate
	}

	for current.Before(endDate) || current.Equal(endDate) {
		// Skip the first occurrence if it's the parent task itself
		if !current.Equal(parentTask.Deadline) && !current.Equal(parentTask.CreatedAt) {
			instance := models.Task{
				Description:       parentTask.Description,
				Status:            "pending",
				Priority:          parentTask.Priority,
				Deadline:          current,
				Categories:        parentTask.Categories,
				Notes:             parentTask.Notes,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				ParentTaskID:      parentTask.ID,
				RecurrencePattern: models.RecurrenceNone, // Instances don't recur themselves
			}
			instances = append(instances, instance)
		}

		// Calculate next occurrence
		current = e.nextOccurrence(current, parentTask.RecurrencePattern, interval)
	}

	return instances
}

// nextOccurrence calculates the next occurrence date based on the pattern and interval
func (e *Engine) nextOccurrence(current time.Time, pattern models.RecurrencePattern, interval int) time.Time {
	switch pattern {
	case models.RecurrenceDaily:
		return current.AddDate(0, 0, interval)
	case models.RecurrenceWeekly:
		return current.AddDate(0, 0, interval*7)
	case models.RecurrenceMonthly:
		return current.AddDate(0, interval, 0)
	case models.RecurrenceYearly:
		return current.AddDate(interval, 0, 0)
	default:
		return current
	}
}

// GetNextInstance returns the next instance that should be created for a recurring task
func (e *Engine) GetNextInstance(parentTask models.Task, lastInstance time.Time) *models.Task {
	if !parentTask.IsRecurring() {
		return nil
	}

	interval := parentTask.RecurrenceInterval
	if interval <= 0 {
		interval = 1
	}

	nextDate := e.nextOccurrence(lastInstance, parentTask.RecurrencePattern, interval)

	// Check if we've passed the end date
	if !parentTask.RecurrenceEndDate.IsZero() && nextDate.After(parentTask.RecurrenceEndDate) {
		return nil
	}

	instance := &models.Task{
		Description:       parentTask.Description,
		Status:            "pending",
		Priority:          parentTask.Priority,
		Deadline:          nextDate,
		Categories:        parentTask.Categories,
		Notes:             parentTask.Notes,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		ParentTaskID:      parentTask.ID,
		RecurrencePattern: models.RecurrenceNone,
	}

	return instance
}

// ShouldGenerateNext checks if a new instance should be generated based on the current date
func (e *Engine) ShouldGenerateNext(parentTask models.Task, lastInstanceDate time.Time) bool {
	if !parentTask.IsRecurring() {
		return false
	}

	interval := parentTask.RecurrenceInterval
	if interval <= 0 {
		interval = 1
	}

	nextDate := e.nextOccurrence(lastInstanceDate, parentTask.RecurrencePattern, interval)

	// Check if the next occurrence is in the past or today
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if nextDate.Before(today) || nextDate.Equal(today) {
		// Check if we haven't passed the end date
		if parentTask.RecurrenceEndDate.IsZero() || nextDate.Before(parentTask.RecurrenceEndDate) || nextDate.Equal(parentTask.RecurrenceEndDate) {
			return true
		}
	}

	return false
}
