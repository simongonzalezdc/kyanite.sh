package command

import (
	"fmt"
)

// NewCompleteTaskCommand creates a new complete task command
func NewCompleteTaskCommand(engine TaskEngine, taskID string) *CompleteTaskCommand {
	return &CompleteTaskCommand{
		engine: engine,
		taskID: taskID,
	}
}

// Execute marks the task as complete
func (c *CompleteTaskCommand) Execute() error {
	if err := c.engine.UpdateTaskStatus(c.taskID, "completed"); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}
	return nil
}

// Undo marks the task as pending
func (c *CompleteTaskCommand) Undo() error {
	if err := c.engine.UpdateTaskStatus(c.taskID, "pending"); err != nil {
		return fmt.Errorf("failed to uncomplete task: %w", err)
	}
	return nil
}

// Description returns a description of this command
func (c *CompleteTaskCommand) Description() string {
	return fmt.Sprintf("Complete task: %s", c.taskID)
}
