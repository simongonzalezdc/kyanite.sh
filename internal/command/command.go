package command

import (
	"github.com/kyanite/focus/pkg/models"
)

// Command represents an undoable operation
type Command interface {
	Execute() error
	Undo() error
	Description() string
}

// History manages command history for undo/redo
type History struct {
	commands []Command
	current  int // Points to the last executed command
}

// NewHistory creates a new command history
func NewHistory() *History {
	return &History{
		commands: make([]Command, 0),
		current:  -1,
	}
}

// Execute runs a command and adds it to history
func (h *History) Execute(cmd Command) error {
	// Execute the command
	if err := cmd.Execute(); err != nil {
		return err
	}

	// Remove any commands after current position (they're now invalidated)
	if h.current < len(h.commands)-1 {
		h.commands = h.commands[:h.current+1]
	}

	// Add command to history
	h.commands = append(h.commands, cmd)
	h.current++

	return nil
}

// Undo undoes the last command
func (h *History) Undo() error {
	if !h.CanUndo() {
		return ErrNothingToUndo
	}

	cmd := h.commands[h.current]
	if err := cmd.Undo(); err != nil {
		return err
	}

	h.current--
	return nil
}

// Redo redoes the next command
func (h *History) Redo() error {
	if !h.CanRedo() {
		return ErrNothingToRedo
	}

	h.current++
	cmd := h.commands[h.current]
	if err := cmd.Execute(); err != nil {
		h.current-- // Rollback on error
		return err
	}

	return nil
}

// CanUndo returns true if there are commands to undo
func (h *History) CanUndo() bool {
	return h.current >= 0
}

// CanRedo returns true if there are commands to redo
func (h *History) CanRedo() bool {
	return h.current < len(h.commands)-1
}

// GetLastCommand returns the description of the last executed command
func (h *History) GetLastCommand() string {
	if h.current < 0 || h.current >= len(h.commands) {
		return ""
	}
	return h.commands[h.current].Description()
}

// GetNextCommand returns the description of the next command to redo
func (h *History) GetNextCommand() string {
	if h.current+1 >= len(h.commands) {
		return ""
	}
	return h.commands[h.current+1].Description()
}

// Clear clears the command history
func (h *History) Clear() {
	h.commands = make([]Command, 0)
	h.current = -1
}

// AddTaskCommand represents adding a task
type AddTaskCommand struct {
	engine      TaskEngine
	parsedTask  models.ParsedTask
	createdTask *models.Task
}

// DeleteTaskCommand represents deleting a task
type DeleteTaskCommand struct {
	engine      TaskEngine
	taskID      string
	deletedTask *models.Task
}

// CompleteTaskCommand represents completing a task
type CompleteTaskCommand struct {
	engine TaskEngine
	taskID string
}

// UpdateTaskCommand represents updating a task
type UpdateTaskCommand struct {
	engine  TaskEngine
	newTask models.Task
	oldTask *models.Task
}

// TaskEngine interface defines the operations needed by commands
type TaskEngine interface {
	AddTask(parsedTask models.ParsedTask) (models.Task, error)
	DeleteTask(id string) error
	RestoreTask(task models.Task) error
	GetTask(id string) (models.Task, error)
	UpdateTask(task models.Task) error
	UpdateTaskStatus(id, status string) error
}
