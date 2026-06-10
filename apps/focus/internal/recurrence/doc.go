// Package recurrence provides task recurrence functionality.
//
// The recurrence package handles recurring tasks with support for:
//   - Multiple recurrence patterns: daily, weekly, monthly, yearly
//   - Custom intervals (e.g., every 2 days, every 3 weeks)
//   - End dates for recurring tasks
//   - Automatic generation of task instances
//
// The Engine type is the main interface for working with recurring tasks.
// It generates task instances based on recurrence rules and manages the
// lifecycle of recurring tasks.
//
// Example usage:
//
//	engine := recurrence.New()
//	instances := engine.GenerateInstances(parentTask, endDate)
//	for _, instance := range instances {
//	    // Process each generated instance
//	}
package recurrence
