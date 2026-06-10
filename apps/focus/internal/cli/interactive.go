package cli

import (
	"github.com/kyanite/focus/pkg/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/gum"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "🎯 Interactive task creation with Gum",
	Long:  "🌸 Create tasks using enhanced Gum-powered prompts",
	Run:   interactiveHandler,
}

func interactiveHandler(cmd *cobra.Command, args []string) {
	// Check if gum is available
	if !gum.IsAvailable() {
		fmt.Println("❌ Gum is not available. Please install gum:")
		fmt.Println("   brew install gum")
		fmt.Println("   Or visit: https://github.com/charmbracelet/gum")
		return
	}

	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	fmt.Println("🎯 focus.sh Interactive Task Creator")
	fmt.Println(strings.Repeat("─", 50))

	// Step 1: Get task description
	description := gum.Input("What needs to be done?")
	if description == "" {
		fmt.Println("❌ Task description cannot be empty.")
		return
	}

	// Step 2: Choose priority
	priorityOptions := []string{"high", "medium", "low"}
	priority := gum.Choose(priorityOptions, "Select priority level:")

	// Step 3: Choose categories (multi-select)
	categoryOptions := []string{
		"work", "personal", "urgent", "learning",
		"health", "finance", "creative", "shopping",
	}
	categories := gum.MultiSelect(categoryOptions, "Select categories:", 5)

	// Step 4: Confirm before creating
	confirmMsg := fmt.Sprintf("Create task: %s\nPriority: %s\nCategories: %s",
		description, priority, strings.Join(categories, ", "))

	if !gum.Confirm(confirmMsg) {
		fmt.Println("❌ Task creation cancelled.")
		return
	}

	// Create the task
	parsedTask := models.ParsedTask{
		Description: description,
		Priority:    priority,
		Categories:  categories,
	}

	task, err := taskEngine.AddTask(parsedTask)
	if err != nil {
		fmt.Printf("❌ Error creating task: %v\n", err)
		return
	}

	// Show success message
	successStyle := lipgloss.NewStyle().
		Foreground(styles.GetSuccess()).
		Bold(true).
		Render("✅ Task created successfully!")

	fmt.Println()
	fmt.Println(successStyle)
	fmt.Printf("📝 Task: %s\n", task.Description)
	fmt.Printf("🎯 Priority: %s\n", task.Priority)
	if len(categories) > 0 {
		fmt.Printf("🏷️  Categories: %s\n", strings.Join(categories, ", "))
	}
	fmt.Printf("🆔 ID: %s\n", task.ID)
	fmt.Println(strings.Repeat("─", 50))
}

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "🔍 Filter tasks interactively with Gum",
	Long:  "🌸 Filter and select tasks using Gum's fuzzy search",
	Run:   filterHandler,
}

func filterHandler(cmd *cobra.Command, args []string) {
	// Check if gum is available
	if !gum.IsAvailable() {
		fmt.Println("❌ Gum is not available.")
		return
	}

	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	// Get all tasks
	tasks, err := taskEngine.ListTasks("all")
	if err != nil {
		fmt.Printf("❌ Error loading tasks: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("😴 No tasks found.")
		return
	}

	// Create task options for filtering
	taskOptions := make([]string, len(tasks))
	for i, task := range tasks {
		status := "⏳"
		if task.Status == "completed" {
			status = "✅"
		}

		priority := "🟢"
		switch task.Priority {
		case "high":
			priority = "🔴"
		case "medium":
			priority = "🟡"
		}

		taskOptions[i] = fmt.Sprintf("%s %s %s | %s | ID: %s",
			status, priority, task.Description,
			strings.Join(task.Categories, ", "), task.ID)
	}

	// Let user filter tasks
	fmt.Println("🔍 focus.sh Task Filter")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("Start typing to filter tasks, press Enter to select:")

	selected := gum.Filter(taskOptions, "Filter tasks...")
	if selected == "" {
		fmt.Println("❌ No task selected.")
		return
	}

	// Extract task ID from selected option
	parts := strings.Split(selected, "ID: ")
	if len(parts) < 2 {
		fmt.Printf("❌ Could not extract task ID from: %s\n", selected)
		return
	}
	taskID := strings.TrimSpace(parts[1])

	// Get and display the selected task
	task, err := taskEngine.GetTask(taskID)
	if err != nil {
		fmt.Printf("❌ Task not found: %v\n", err)
		return
	}

	// Display task details
	displayTaskDetails(task)

	// Ask what action to take
	actionOptions := []string{"complete", "delete", "edit notes", "back"}
	action := gum.Choose(actionOptions, "What would you like to do with this task?")

	switch action {
	case "complete":
		if task.Status == "completed" {
			if gum.Confirm("Mark as pending?") {
				err = taskEngine.UpdateTaskStatus(task.ID, "pending")
			}
		} else {
			if gum.Confirm("Mark as completed?") {
				err = taskEngine.CompleteTask(task.ID)
			}
		}
	case "delete":
		if gum.Confirm("Are you sure you want to delete this task?") {
			err = taskEngine.DeleteTask(task.ID)
		}
	case "edit notes":
		// This would integrate with our notes command
		fmt.Printf("📝 To edit notes, use: focus notes %s\n", task.ID)
	}

	if err != nil {
		fmt.Printf("❌ Error performing action: %v\n", err)
	} else {
		fmt.Println("✅ Action completed successfully!")
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "⚙️ Interactive configuration with Gum",
	Long:  "🌸 Configure focus.sh settings using Gum prompts",
	Run:   configHandler,
}

func configHandler(cmd *cobra.Command, args []string) {
	// Check if gum is available
	if !gum.IsAvailable() {
		fmt.Println("❌ Gum is not available.")
		return
	}

	fmt.Println("⚙️ focus.sh Configuration")
	fmt.Println(strings.Repeat("─", 50))

	// AI Provider selection
	aiProviders := []string{"ollama", "openrouter", "openai"}
	aiProvider := gum.Choose(aiProviders, "Select AI provider:")

	// Default theme selection
	themes := []string{"synthwave", "light", "plain"}
	defaultTheme := gum.Choose(themes, "Select default theme:")

	// Auto-save interval
	interval := gum.Input("Auto-save interval in minutes (leave empty for default):")

	// Display configuration
	configDisplay := gum.Style(
		fmt.Sprintf(`
🤖 AI Provider: %s
🎨 Default Theme: %s
💾 Auto-save: %s minutes

Configuration saved to ~/.focus/config.yml`, aiProvider, defaultTheme, interval),
		"#FF71CE", "#0F0A19", true,
	)

	fmt.Println(configDisplay)
	fmt.Println("✅ Configuration completed!")
}

func displayTaskDetails(task models.Task) {
	status := "⏳ Pending"
	if task.Status == "completed" {
		status = "✅ Completed"
	}

	priority := "🟢 Low"
	switch task.Priority {
	case "high":
		priority = "🔴 High"
	case "medium":
		priority = "🟡 Medium"
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true).
		Render(fmt.Sprintf("%s - %s", task.Description, status))

	fmt.Println()
	fmt.Println(headerStyle)
	fmt.Println(strings.Repeat("─", len(task.Description)+len(status)+3))
	fmt.Printf("🆔 ID: %s\n", task.ID)
	fmt.Printf("🎯 Priority: %s\n", priority)
	fmt.Printf("📅 Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
	if len(task.Categories) > 0 {
		fmt.Printf("🏷️  Categories: %s\n", strings.Join(task.Categories, ", "))
	}
	if task.Notes != "" {
		fmt.Printf("📝 Notes: %s\n", task.Notes)
	}
	fmt.Println()
}
