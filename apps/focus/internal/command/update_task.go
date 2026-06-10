package command

import (
	"fmt"

	"github.com/kyanite/focus/pkg/models"
)

// NewUpdateTaskCommand creates a new update task command
func NewUpdateTaskCommand(engine TaskEngine, newTask models.Task) *UpdateTaskCommand {
	return &UpdateTaskCommand{
		engine:  engine,
		newTask: newTask,
	}
}

// Execute updates the task
func (c *UpdateTaskCommand) Execute() error {
	// Store the old task state before updating
	oldTask, err := c.engine.GetTask(c.newTask.ID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	c.oldTask = &oldTask

	if err := c.engine.UpdateTask(c.newTask); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

// Undo restores the old task state
func (c *UpdateTaskCommand) Undo() error {
	if c.oldTask == nil {
		return fmt.Errorf("no old task state to restore")
	}

	if err := c.engine.UpdateTask(*c.oldTask); err != nil {
		return fmt.Errorf("failed to restore task: %w", err)
	}
	return nil
}

// Description returns a description of this command
func (c *UpdateTaskCommand) Description() string {
	return fmt.Sprintf("Update task: %s", c.newTask.Description)
}
