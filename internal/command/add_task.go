package command

import (
	"fmt"

	"github.com/kyanite/focus/pkg/models"
)

// NewAddTaskCommand creates a new add task command
func NewAddTaskCommand(engine TaskEngine, parsedTask models.ParsedTask) *AddTaskCommand {
	return &AddTaskCommand{
		engine:     engine,
		parsedTask: parsedTask,
	}
}

// Execute adds the task
func (c *AddTaskCommand) Execute() error {
	task, err := c.engine.AddTask(c.parsedTask)
	if err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}
	c.createdTask = &task
	return nil
}

// Undo deletes the added task
func (c *AddTaskCommand) Undo() error {
	if c.createdTask == nil {
		return fmt.Errorf("no task to undo")
	}
	if err := c.engine.DeleteTask(c.createdTask.ID); err != nil {
		return fmt.Errorf("failed to undo add task: %w", err)
	}
	return nil
}

// Description returns a description of this command
func (c *AddTaskCommand) Description() string {
	if c.createdTask != nil {
		return fmt.Sprintf("Add task: %s", c.createdTask.Description)
	}
	return fmt.Sprintf("Add task: %s", c.parsedTask.Description)
}
