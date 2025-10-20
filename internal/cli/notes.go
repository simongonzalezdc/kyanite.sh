package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var notesCmd = &cobra.Command{
	Use:   "notes [task_id]",
	Short: "📝 Edit task notes with markdown editor",
	Long:  "🌸 Open markdown editor to add/edit notes for tasks",
	Args:  cobra.MaximumNArgs(1),
	Run:   notesHandler,
}

func notesHandler(cmd *cobra.Command, args []string) {
	// Initialize components
	storage := store.New(utils.GetStoragePath())
	taskEngine := engine.New(storage)

	if len(args) == 0 {
		// Show all notes browser
		showAllNotes(taskEngine)
	} else {
		// Edit specific task notes
		editTaskNotes(args[0], taskEngine)
	}
}

func showAllNotes(engine *engine.Engine) {
	fmt.Println("📝 Loading all task notes...")
	fmt.Println(strings.Repeat("─", 50))

	// Get all tasks with notes
	tasks, err := engine.ListTasks("all")
	if err != nil {
		fmt.Printf("❌ Error loading tasks: %v\n", err)
		return
	}

	// Filter tasks with notes
	var tasksWithNotes []models.Task
	for _, task := range tasks {
		if task.Notes != "" {
			tasksWithNotes = append(tasksWithNotes, task)
		}
	}

	if len(tasksWithNotes) == 0 {
		fmt.Println("😴 No notes found.")
		fmt.Println("💡 Add notes with: focus notes <task_id>")
		return
	}

	// Display notes
	for _, task := range tasksWithNotes {
		statusIcon := "⏳"
		if task.Status == "completed" {
			statusIcon = "✅"
		}

		priorityIcon := "🟢"
		switch task.Priority {
		case "high":
			priorityIcon = "🔴"
		case "medium":
			priorityIcon = "🟡"
		}

		// Task header
		headerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF71CE")).
			Bold(true)

		fmt.Println(headerStyle.Render(fmt.Sprintf("%s %s", statusIcon, task.Description)))
		fmt.Printf("ID: %s | Priority: %s | Status: %s\n",
			task.ID, priorityIcon, task.Status)
		if len(task.Categories) > 0 {
			fmt.Printf("Categories: %s\n", strings.Join(task.Categories, ", "))
		}
		fmt.Println(strings.Repeat("─", 40))

		// Notes content
		notesStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			PaddingLeft(2)

		fmt.Println(notesStyle.Render(task.Notes))
		fmt.Println()
		fmt.Println(strings.Repeat("═", 50))
		fmt.Println()
	}

	// Option to launch editor
	fmt.Println("💡 Tip: Use 'focus notes <task_id>' to edit specific task notes")
}

func editTaskNotes(taskID string, engine *engine.Engine) {
	// Get the task
	task, err := engine.GetTask(taskID)
	if err != nil {
		fmt.Printf("❌ Task not found: %v\n", err)
		return
	}

	// Prepare initial content
	statusIcon := "⏳"
	if task.Status == "completed" {
		statusIcon = "✅"
	}

	initialContent := fmt.Sprintf(`# %s %s

**Task ID:** %s  
**Priority:** %s  
**Status:** %s  
**Created:** %s

---

## 📝 Notes

%s

---

*Edit this markdown file to update your task notes. Save and exit to apply changes.*`,
		statusIcon,
		task.Description,
		task.ID,
		task.Priority,
		task.Status,
		task.CreatedAt.Format("2006-01-02 15:04"),
		task.Notes,
	)

	// Create temporary file
	tempFile, err := os.CreateTemp("", fmt.Sprintf("focus-notes-%s-*.md", taskID))
	if err != nil {
		fmt.Printf("❌ Error creating temp file: %v\n", err)
		return
	}
	tempFileName := tempFile.Name()
	defer func() { _ = os.Remove(tempFileName) }()

	// Write initial content to temp file
	if _, err := tempFile.WriteString(initialContent); err != nil {
		fmt.Printf("❌ Error writing to temp file: %v\n", err)
		return
	}
	_ = tempFile.Close()

	// Try to find an editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Try common editors
		editors := []string{"code", "nano", "vim", "notepad"}
		for _, e := range editors {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}

	if editor == "" {
		fmt.Println("❌ No suitable editor found. Please set EDITOR environment variable.")
		fmt.Println("Example: export EDITOR=code")
		return
	}

	// Launch editor
	fmt.Printf("🚀 Opening notes editor for task: %s\n", task.Description)
	fmt.Printf("📝 Using editor: %s\n", editor)
	fmt.Println("💾 Save and exit to apply changes")

	cmd := exec.Command(editor, tempFileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Error launching editor: %v\n", err)
		return
	}

	// Read the updated content
	updatedContent, err := os.ReadFile(tempFileName)
	if err != nil {
		fmt.Printf("❌ Error reading updated notes: %v\n", err)
		return
	}

	// Extract notes from the markdown (everything after "## 📝 Notes")
	content := string(updatedContent)
	notesStart := strings.Index(content, "## 📝 Notes")
	if notesStart != -1 {
		notesStart += len("## 📝 Notes\n")

		// Find the start of actual notes content
		lines := strings.Split(content[notesStart:], "\n")
		var notesLines []string
		foundContent := false

		for i, line := range lines {
			// Skip empty lines and the line that says "---"
			if line == "" || (i == 0 && strings.TrimSpace(line) == "") {
				continue
			}
			if strings.TrimSpace(line) == "---" {
				if foundContent {
					break // Stop at the closing ---
				}
				continue
			}
			if !foundContent && strings.TrimSpace(line) == "" {
				continue
			}
			foundContent = true
			notesLines = append(notesLines, line)
		}

		task.Notes = strings.Join(notesLines, "\n")
	}

	// Update the task
	if err := engine.UpdateTask(task); err != nil {
		fmt.Printf("❌ Error updating task: %v\n", err)
		return
	}

	// Show success message
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF66")).
		Bold(true).
		Render("✅ Notes updated successfully!")

	fmt.Println()
	fmt.Println(successStyle)
	fmt.Printf("📝 Task: %s\n", task.Description)
	if len(task.Notes) > 100 {
		fmt.Printf("📄 Preview: %s...\n", task.Notes[:100])
	} else if task.Notes != "" {
		fmt.Printf("📄 Preview: %s\n", task.Notes)
	}
}
