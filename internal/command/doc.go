// Package command provides the Command pattern implementation for undo/redo functionality.
//
// The command package implements the Gang of Four Command pattern, allowing
// operations to be encapsulated as objects. This enables:
//   - Undo/Redo functionality for task operations
//   - Command history tracking
//   - Reversible operations
//
// The package includes:
//   - Command interface for all undoable operations
//   - History manager for undo/redo stack management
//   - Concrete commands: AddTask, DeleteTask, CompleteTask, UpdateTask
//
// Example usage:
//
//	history := command.NewHistory()
//	cmd := command.NewAddTaskCommand(engine, parsedTask)
//	history.Execute(cmd)
//	// Later...
//	history.Undo() // Undoes the add
//	history.Redo() // Redoes the add
package command
