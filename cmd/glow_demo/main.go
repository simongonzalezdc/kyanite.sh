package main

import (
	"fmt"

	"github.com/kyanite/focus/pkg/glow"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/utils"
	"strings"
)

func main() {
	fmt.Println("🌌 NEON Focus - Glow Enhancement Demo")
	fmt.Println(strings.Repeat("─", 50))
	
	// Initialize components
	storage := store.New(utils.GetStoragePath())
	taskEngine := engine.New(storage)
	
	// Create Glow styler with synthwave theme
	styler := glow.NewGlowStyler("synthwave")
	
	// Render header with Glow styling
	fmt.Println(styler.RenderHeaderWithGlow(
		"🌌 NEON ENHANCED TASK MANAGER",
		"Featuring Glow & Charm Libraries Integration",
	))
	
	// Get some tasks
	tasks, err := taskEngine.ListTasks("all")
	if err != nil {
		fmt.Println(styler.RenderErrorWithGlow(fmt.Sprintf("Error loading tasks: %v", err)))
		return
	}
	
	if len(tasks) == 0 {
		fmt.Println(styler.RenderInfoWithGlow("No tasks found. Adding demo tasks..."))
		
		// Add demo tasks
		demoTasks := []struct {
			desc string
			prio string
			status string
		}{
			{"Complete Glow integration", "high", "completed"},
			{"Test new visual enhancements", "medium", "pending"},
			{"Implement audio notifications", "low", "pending"},
		}
		
		for _, dt := range demoTasks {
			parsedTask := models.ParsedTask{
				Description: dt.desc,
				Priority:    dt.prio,
				Categories:  []string{"demo"},
				Notes:       fmt.Sprintf("This is a demo task for testing Glow enhancements. Status: %s", dt.status),
			}
			
			_, err := taskEngine.AddTask(parsedTask)
			if err != nil {
				fmt.Println(styler.RenderErrorWithGlow(fmt.Sprintf("Error adding demo task: %v", err)))
			}
		}
		
		// Reload tasks
		tasks, _ = taskEngine.ListTasks("all")
	}
	
	// Render tasks with Glow styling
	fmt.Println(styler.RenderSectionWithGlow(
		"📋 Enhanced Task Display",
		fmt.Sprintf("Found %d tasks with Glow-enhanced visual styling", len(tasks)),
		"#00FFF0",
	))
	
	// Render each task
	for i, task := range tasks {
		fmt.Println(styler.RenderTaskWithGlow(
			task.ID,
			task.Description,
			task.Status,
			task.Priority,
			i+1,
		))
	}
	
	// Show progress bar with Glow styling
	completedCount := 0
	for _, task := range tasks {
		if task.Status == "completed" {
			completedCount++
		}
	}
	
	fmt.Println(styler.RenderProgressBarWithGlow(
		completedCount,
		len(tasks),
		"Task Completion",
	))
	
	// Show statistics section
	stats := fmt.Sprintf(`📊 TASK STATISTICS
• Total Tasks: %d
• Completed: %d (%.1f%%)
• Pending: %d (%.1f%%)
• In Progress: %d (%.1f%%)
• High Priority: %d
• Medium Priority: %d
• Low Priority: %d`,
		len(tasks),
		completedCount,
		float64(completedCount)/float64(len(tasks))*100,
		len(tasks)-completedCount,
		float64(len(tasks)-completedCount)/float64(len(tasks))*100,
		0,
		0.0,
		getPriorityCount(tasks, "high"),
		getPriorityCount(tasks, "medium"),
		getPriorityCount(tasks, "low"),
	)
	
	fmt.Println(styler.RenderSectionWithGlow("📈 Performance Metrics", stats, "#FF71CE"))
	
	// Success message
	fmt.Println(styler.RenderSuccessWithGlow("Glow integration demonstration completed successfully!"))
	fmt.Println(styler.RenderSeparatorWithGlow("✨"))
		fmt.Println("🌈 Charm libraries provide excellent visual enhancements for NEON Focus!")
}

func getPriorityCount(tasks []models.Task, priority string) int {
	count := 0
	for _, task := range tasks {
		if task.Priority == priority {
			count++
		}
	}
	return count
}
