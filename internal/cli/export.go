package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	exportFormat   string
	exportOutput   string
	exportFilter   string
	exportPriority string
	exportStatus   string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export tasks to various formats",
	Long: `Export tasks to CSV, Markdown, or JSON format.
Examples:
  focus export --format=csv --output=tasks.csv
  focus export --format=markdown --output=tasks.md
  focus export --format=json --output=tasks.json
  focus export --format=csv --filter=active --priority=high`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate format
		if exportFormat == "" {
			exportFormat = "json" // Default to JSON
		}

		format := strings.ToLower(exportFormat)
		if format != "csv" && format != "markdown" && format != "json" {
			showExportError("Invalid format. Supported formats: csv, markdown, json")
			return
		}

		// Initialize components
		store := store.New(utils.GetStoragePath())
		engine := engine.New(store)

		// Get tasks based on filter
		var tasks []models.Task
		var err error

		if exportFilter != "" {
			tasks, err = engine.ListTasks(exportFilter)
		} else {
			tasks, err = engine.ListTasks("all")
		}

		if err != nil {
			showExportError(fmt.Sprintf("Failed to load tasks: %s", err.Error()))
			return
		}

		// Apply additional filters
		tasks = filterExportTasks(tasks)

		if len(tasks) == 0 {
			showExportWarning("No tasks to export")
			return
		}

		// Generate output filename if not specified
		if exportOutput == "" {
			timestamp := time.Now().Format("2006-01-02")
			exportOutput = fmt.Sprintf("focus-tasks-%s.%s", timestamp, format)
		}

		// Export based on format
		var exportErr error
		switch format {
		case "csv":
			exportErr = exportToCSV(tasks, exportOutput)
		case "markdown":
			exportErr = exportToMarkdown(tasks, exportOutput)
		case "json":
			exportErr = exportToJSON(tasks, exportOutput)
		}

		if exportErr != nil {
			showExportError(fmt.Sprintf("Export failed: %s", exportErr.Error()))
			return
		}

		showExportSuccess(exportOutput, len(tasks), format)
	},
}

func filterExportTasks(tasks []models.Task) []models.Task {
	filtered := make([]models.Task, 0, len(tasks))

	for _, task := range tasks {
		// Priority filter
		if exportPriority != "" {
			if !strings.EqualFold(task.Priority, exportPriority) {
				continue
			}
		}

		// Status filter
		if exportStatus != "" {
			if !strings.EqualFold(task.Status, exportStatus) {
				continue
			}
		}

		filtered = append(filtered, task)
	}

	return filtered
}

func exportToCSV(tasks []models.Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Description", "Status", "Priority", "Deadline", "Categories", "Created", "Updated"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write tasks
	for _, task := range tasks {
		deadline := ""
		if !task.Deadline.IsZero() {
			deadline = task.Deadline.Format("2006-01-02 15:04:05")
		}

		categories := strings.Join(task.Categories, "; ")

		record := []string{
			task.ID,
			task.Description,
			task.Status,
			task.Priority,
			deadline,
			categories,
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			task.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func exportToMarkdown(tasks []models.Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var content strings.Builder

	// Write header
	content.WriteString("# Focus Tasks Export\n\n")
	content.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("**Total Tasks:** %d\n\n", len(tasks)))
	content.WriteString("---\n\n")

	// Group tasks by status
	pendingTasks := make([]models.Task, 0)
	completedTasks := make([]models.Task, 0)

	for _, task := range tasks {
		if task.Status == "completed" {
			completedTasks = append(completedTasks, task)
		} else {
			pendingTasks = append(pendingTasks, task)
		}
	}

	// Write pending tasks
	if len(pendingTasks) > 0 {
		content.WriteString("## Pending Tasks\n\n")
		for _, task := range pendingTasks {
			writeMarkdownTask(&content, task, false)
		}
	}

	// Write completed tasks
	if len(completedTasks) > 0 {
		content.WriteString("## Completed Tasks\n\n")
		for _, task := range completedTasks {
			writeMarkdownTask(&content, task, true)
		}
	}

	_, err = file.WriteString(content.String())
	return err
}

func writeMarkdownTask(content *strings.Builder, task models.Task, completed bool) {
	checkbox := "- [ ]"
	if completed {
		checkbox = "- [x]"
	}

	// Priority emoji
	priorityEmoji := ""
	switch task.Priority {
	case "high":
		priorityEmoji = "🔴"
	case "medium":
		priorityEmoji = "🟡"
	case "low":
		priorityEmoji = "🟢"
	}

	content.WriteString(fmt.Sprintf("%s %s **%s**\n", checkbox, priorityEmoji, task.Description))
	content.WriteString(fmt.Sprintf("  - **ID:** `%s`\n", task.ID))
	content.WriteString(fmt.Sprintf("  - **Priority:** %s\n", task.Priority))

	if !task.Deadline.IsZero() {
		content.WriteString(fmt.Sprintf("  - **Deadline:** %s\n", task.Deadline.Format("2006-01-02 15:04")))
	}

	if len(task.Categories) > 0 {
		content.WriteString(fmt.Sprintf("  - **Categories:** %s\n", strings.Join(task.Categories, ", ")))
	}

	content.WriteString(fmt.Sprintf("  - **Created:** %s\n", task.CreatedAt.Format("2006-01-02 15:04")))

	if completed {
		content.WriteString(fmt.Sprintf("  - **Completed:** %s\n", task.UpdatedAt.Format("2006-01-02 15:04")))
	}

	content.WriteString("\n")
}

func exportToJSON(tasks []models.Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create export structure
	export := struct {
		ExportedAt time.Time      `json:"exported_at"`
		TotalTasks int            `json:"total_tasks"`
		Tasks      []models.Task  `json:"tasks"`
	}{
		ExportedAt: time.Now(),
		TotalTasks: len(tasks),
		Tasks:      tasks,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(export)
}

func showExportSuccess(filename string, count int, format string) {
	successTitle := styles.SynthwaveTitle("📤 EXPORT COMPLETE")
	fmt.Println(successTitle)
	fmt.Println()

	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"FILE", filename, styles.SynthwaveCyan},
		{"TASKS", fmt.Sprintf("%d", count), styles.SynthwaveGreen},
		{"FORMAT", strings.ToUpper(format), styles.SynthwaveYellow},
	}

	for _, detail := range details {
		detailLine := lipgloss.NewStyle().
			Foreground(detail.color).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 1).
			Render(fmt.Sprintf("▸ %s: %s", detail.label, detail.value))
		fmt.Println(detailLine)
	}
	fmt.Println()
}

func showExportWarning(message string) {
	warningStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwaveYellow).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Render(fmt.Sprintf("⚠️  %s", message))
	fmt.Println(warningStyle)
}

func showExportError(message string) {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ EXPORT FAILED\n\n%s", message))
	fmt.Println(errorBox)
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Add flags
	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "Export format (csv, markdown, json)")
	exportCmd.Flags().StringVar(&exportOutput, "output", "", "Output filename (auto-generated if not specified)")
	exportCmd.Flags().StringVar(&exportFilter, "filter", "all", "Filter tasks (all, active, completed)")
	exportCmd.Flags().StringVar(&exportPriority, "priority", "", "Filter by priority (low, medium, high)")
	exportCmd.Flags().StringVar(&exportStatus, "status", "", "Filter by status (pending, completed)")
}
