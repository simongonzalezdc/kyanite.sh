package command

import (
	"fmt"
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

	// Use RestoreTask to preserve the original task ID and all fields
	if err := c.engine.RestoreTask(*c.deletedTask); err != nil {
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
