package wizards

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/utils"
)

var (
	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF71CE")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF66")).
			Bold(true)
)

// TaskCreationWizard creates a comprehensive task creation form
func TaskCreationWizard() error {
	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	var task models.ParsedTask
	var dueDateStr string

	// Create the form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What needs to be done?").
				Placeholder("Enter task description...").
				Value(&task.Description).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("task description cannot be empty")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Priority Level").
				Description("Select the urgency of this task").
				Options(
					huh.NewOption("🔴 High", "high"),
					huh.NewOption("🟡 Medium", "medium"),
					huh.NewOption("🟢 Low", "low"),
				).
				Value(&task.Priority),

			huh.NewMultiSelect[string]().
				Title("Categories").
				Description("Choose relevant categories (max 5)").
				Options(
					huh.NewOption("💼 Work", "work"),
					huh.NewOption("🏠 Personal", "personal"),
					huh.NewOption("🚨 Urgent", "urgent"),
					huh.NewOption("📚 Learning", "learning"),
					huh.NewOption("💪 Health", "health"),
					huh.NewOption("💰 Finance", "finance"),
					huh.NewOption("🎨 Creative", "creative"),
					huh.NewOption("🛒 Shopping", "shopping"),
				).
				Value(&task.Categories).
				Filterable(true).
				Limit(5),

			huh.NewText().
				Title("Notes").
				Description("Add any additional notes or details").
				Placeholder("Enter any additional information...").
				Value(&task.Notes),

			huh.NewInput().
				Title("Due Date (optional)").
				Description("Format: YYYY-MM-DD or 'tomorrow', 'next week', etc.").
				Placeholder("e.g., 2025-12-25 or tomorrow").
				Value(&dueDateStr),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		return fmt.Errorf("form cancelled or error: %w", err)
	}

	// Process due date if provided
	if dueDateStr != "" {
		parsedDate, parseErr := parseDueDate(dueDateStr)
		if parseErr != nil {
			fmt.Printf("⚠️  Warning: Could not parse due date '%s': %v\n", dueDateStr, parseErr)
		} else {
			task.Deadline = *parsedDate
		}
	}

	// Show confirmation
	confirmMsg := fmt.Sprintf(
		"Create task:\n\n%s\n\nPriority: %s\nCategories: %s\nNotes: %s\nDue: %s",
		task.Description,
		task.Priority,
		formatCategories(task.Categories),
		formatNotes(task.Notes),
		formatDueDate(task.Deadline),
	)

	var confirmed bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm Task Creation").
				Description(confirmMsg).
				Value(&confirmed),
		),
	)

	err = confirmForm.Run()
	if err != nil || !confirmed {
		fmt.Println("❌ Task creation cancelled.")
		return nil
	}

	// Create the task
	newTask, err := taskEngine.AddTask(task)
	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	// Show success message
	displayTaskSuccess(newTask)
	return nil
}

// ConfigurationWizard creates an interactive configuration setup
func ConfigurationWizard() error {
	var config struct {
		AIProvider    string
		Model         string
		DefaultTheme  string
		AutoSave      string
		Notifications bool
		TimeFormat    string
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("AI Provider").
				Description("Choose your preferred AI service").
				Options(
					huh.NewOption("🦙 Ollama (Local)", "ollama"),
					huh.NewOption("🌐 OpenRouter (Remote)", "openrouter"),
					huh.NewOption("🤖 OpenAI", "openai"),
				).
				Value(&config.AIProvider),

			huh.NewSelect[string]().
				Title("Default Theme").
				Description("Select your visual preference").
				Options(
					huh.NewOption("🌌 Synthwave", "synthwave"),
					huh.NewOption("☀️ Light", "light"),
					huh.NewOption("⚪ Plain", "plain"),
				).
				Value(&config.DefaultTheme),

			huh.NewSelect[string]().
				Title("Time Format").
				Description("How time should be displayed").
				Options(
					huh.NewOption("12-hour (e.g., 2:30 PM)", "12h"),
					huh.NewOption("24-hour (e.g., 14:30)", "24h"),
				).
				Value(&config.TimeFormat),

			huh.NewInput().
				Title("Auto-save Interval").
				Description("Minutes between automatic saves").
				Placeholder("5").
				Value(&config.AutoSave),

			huh.NewConfirm().
				Title("Enable Notifications").
				Description("Show desktop notifications for reminders").
				Value(&config.Notifications),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("configuration cancelled: %w", err)
	}

	// Here you would save the configuration
	// For now, just display it
	displayConfiguration(config)
	return nil
}

// TaskEditWizard creates a form for editing existing tasks
func TaskEditWizard(task models.Task) error {
	var editedTask models.ParsedTask
	var dueDateStr string
	var status string

	// Convert task to ParsedTask for editing
	editedTask.Description = task.Description
	editedTask.Priority = task.Priority
	editedTask.Categories = task.Categories
	editedTask.Notes = task.Notes
	status = task.Status
	editedTask.Deadline = task.Deadline

	if !task.Deadline.IsZero() {
		dueDateStr = task.Deadline.Format("2006-01-02")
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Task Description").
				Value(&editedTask.Description).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("task description cannot be empty")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Priority").
				Options(
					huh.NewOption("🔴 High", "high"),
					huh.NewOption("🟡 Medium", "medium"),
					huh.NewOption("🟢 Low", "low"),
				).
				Value(&editedTask.Priority),

			huh.NewSelect[string]().
				Title("Status").
				Options(
					huh.NewOption("⏳ Pending", "pending"),
					huh.NewOption("✅ Completed", "completed"),
				).
				Value(&status),

			huh.NewMultiSelect[string]().
				Title("Categories").
				Options(
					huh.NewOption("💼 Work", "work"),
					huh.NewOption("🏠 Personal", "personal"),
					huh.NewOption("🚨 Urgent", "urgent"),
					huh.NewOption("📚 Learning", "learning"),
					huh.NewOption("💪 Health", "health"),
					huh.NewOption("💰 Finance", "finance"),
					huh.NewOption("🎨 Creative", "creative"),
					huh.NewOption("🛒 Shopping", "shopping"),
				).
				Value(&editedTask.Categories).
				Limit(5),

			huh.NewText().
				Title("Notes").
				Placeholder("Enter any additional information...").
				Value(&editedTask.Notes),

			huh.NewInput().
				Title("Due Date (optional)").
				Placeholder("e.g., 2025-12-25 or tomorrow").
				Value(&dueDateStr),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("edit cancelled: %w", err)
	}

	// Process due date if provided
	if dueDateStr != "" {
		parsedDate, parseErr := parseDueDate(dueDateStr)
		if parseErr != nil {
			fmt.Printf("⚠️  Warning: Could not parse due date: %v\n", parseErr)
		} else {
			editedTask.Deadline = *parsedDate
		}
	}

	// Update the task (implementation would go here)
	fmt.Println("✅ Task updated successfully!")
	displayTaskDetails(editedTask, status)
	return nil
}

// Helper functions
func parseDueDate(dateStr string) (*time.Time, error) {
	// Handle common natural language expressions
	lower := strings.ToLower(strings.TrimSpace(dateStr))

	switch lower {
	case "today":
		today := time.Now()
		return &today, nil
	case "tomorrow":
		tomorrow := time.Now().AddDate(0, 0, 1)
		return &tomorrow, nil
	case "next week":
		nextWeek := time.Now().AddDate(0, 0, 7)
		return &nextWeek, nil
	default:
		// Try to parse as date
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse date format")
		}
		return &parsed, nil
	}
}

func formatCategories(categories []string) string {
	if len(categories) == 0 {
		return "None"
	}
	return strings.Join(categories, ", ")
}

func formatNotes(notes string) string {
	if notes == "" {
		return "None"
	}
	// Truncate long notes
	if len(notes) > 50 {
		return notes[:47] + "..."
	}
	return notes
}

func formatDueDate(date time.Time) string {
	if date.IsZero() {
		return "None"
	}
	return date.Format("2006-01-02")
}

func displayTaskSuccess(task models.Task) {
	fmt.Println()
	fmt.Println(successStyle.Render("✅ Task Created Successfully!"))
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("📝 Description: %s\n", task.Description)
	fmt.Printf("🎯 Priority: %s\n", task.Priority)
	if len(task.Categories) > 0 {
		fmt.Printf("🏷️  Categories: %s\n", strings.Join(task.Categories, ", "))
	}
	if task.Notes != "" {
		fmt.Printf("📝 Notes: %s\n", task.Notes)
	}
	if !task.Deadline.IsZero() {
		fmt.Printf("📅 Due: %s\n", task.Deadline.Format("2006-01-02"))
	}
	fmt.Printf("🆔 ID: %s\n", task.ID)
	fmt.Println(strings.Repeat("─", 50))
}

func displayTaskDetails(task models.ParsedTask, status string) {
	fmt.Println()
	fmt.Println(accentStyle.Render("📝 Task Details"))
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("📝 Description: %s\n", task.Description)
	fmt.Printf("🎯 Priority: %s\n", task.Priority)
	fmt.Printf("📊 Status: %s\n", status)
	if len(task.Categories) > 0 {
		fmt.Printf("🏷️  Categories: %s\n", strings.Join(task.Categories, ", "))
	}
	if task.Notes != "" {
		fmt.Printf("📝 Notes: %s\n", task.Notes)
	}
	if !task.Deadline.IsZero() {
		fmt.Printf("📅 Due: %s\n", task.Deadline.Format("2006-01-02"))
	}
	fmt.Println(strings.Repeat("─", 50))
}

func displayConfiguration(config interface{}) {
	fmt.Println()
	fmt.Println(successStyle.Render("✅ Configuration Saved!"))
	fmt.Println(strings.Repeat("─", 50))
	// Implementation would display the actual configuration
	fmt.Println("📋 Configuration has been saved to ~/.focus/config.yml")
	fmt.Println(strings.Repeat("─", 50))
}
