// Package templates provides task template functionality.
//
// The templates package allows users to create reusable task definitions
// that can be instantiated multiple times. This is useful for:
//   - Recurring workflows
//   - Standard operating procedures
//   - Common task patterns
//
// Templates store all task attributes except:
//   - Task ID (generated on creation)
//   - Status (always starts as pending)
//   - Timestamps (set on creation)
//
// Templates are persisted to disk in JSON format.
//
// Example usage:
//
//	store := templates.NewStore("~/.focus/templates.json")
//	template := models.TaskTemplate{
//	    Name: "Daily Standup",
//	    Description: "Attend daily standup meeting",
//	    Priority: "high",
//	    Categories: []string{"meetings", "work"},
//	}
//	store.Add(template)
package templates
