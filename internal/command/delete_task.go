package command

import (
	"fmt"

	"github.com/kyanite/focus/pkg/models"
)

// NewDeleteTaskCommand creates a new delete task command
func NewDeleteTaskCommand(engine TaskEngine, taskID string) *DeleteTaskCommand {
	return &DeleteTaskCommand{
		engine: engine,
		taskID: taskID,
	}
}

// Execute deletes the task
func (c *DeleteTaskCommand) Execute() error {
	// Store the task before deleting so we can restore it
	task, err := c.engine.GetTask(c.taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	c.deletedTask = &task

	if err := c.engine.DeleteTask(c.taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// Undo restores the deleted task
func (c *DeleteTaskCommand) Undo() error {
	if c.deletedTask == nil {
		return fmt.Errorf("no task to restore")
	}

	// Convert Task back to ParsedTask for AddTask
	parsedTask := models.ParsedTask{
		Description:        c.deletedTask.Description,
		Deadline:           c.deletedTask.Deadline,
		Priority:           c.deletedTask.Priority,
		Categories:         c.deletedTask.Categories,
		Notes:              c.deletedTask.Notes,
		RecurrencePattern:  c.deletedTask.RecurrencePattern,
		RecurrenceInterval: c.deletedTask.RecurrenceInterval,
		RecurrenceEndDate:  c.deletedTask.RecurrenceEndDate,
	}

	// This won't preserve the original ID, timestamps, or status,
	// but it's the best we can do without a direct restore operation
	_, err := c.engine.AddTask(parsedTask)
	if err != nil {
		return fmt.Errorf("failed to restore task: %w", err)
	}

	return nil
}

// Description returns a description of this command
func (c *DeleteTaskCommand) Description() string {
	if c.deletedTask != nil {
		return fmt.Sprintf("Delete task: %s", c.deletedTask.Description)
	}
	return fmt.Sprintf("Delete task: %s", c.taskID)
}
