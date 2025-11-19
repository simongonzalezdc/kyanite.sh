package cli

import (
	"fmt"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [filter]",
	Short: "Display tasks on the task board",
	Long: `Show your missions with beautiful Kyanite theming.
Filters: all, active, completed (default: all)`,
	Run: func(cmd *cobra.Command, args []string) {
		filter := "all"
		if len(args) > 0 {
			filter = args[0]
		}

		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)

		tasks, err := engine.ListTasks(filter)
		if err != nil {
			fmt.Printf("Error loading missions: %v\n", err)
			return
		}

		// Convert engine tasks to models
		modelTasks := make([]models.Task, len(tasks))
		for i, task := range tasks {
			modelTasks[i] = models.Task{
				ID:          task.ID,
				Description: task.Description,
				Status:      task.Status,
				Priority:    task.Priority,
				CreatedAt:   task.CreatedAt,
				UpdatedAt:   task.UpdatedAt,
			}
		}

		// Render with maximum visual impact
		renderSynthwaveList(modelTasks, filter)
	},
}

func renderSynthwaveList(tasks []models.Task, filter string) {
	// Epic header with glitch effects
	header := styles.Header()
	fmt.Println(header)
	fmt.Println()

	// Filter indicator with style
	filterIndicator := lipgloss.NewStyle().
		Foreground(styles.SynthwaveYellow).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(0, 2).
		Render(fmt.Sprintf("🔍 FILTER: %s", strings.ToUpper(filter)))
	fmt.Println(filterIndicator)
	fmt.Println()

	if len(tasks) == 0 {
		// Stunning empty state
		emptyMsg := styles.EmptyStateMessage()
		fmt.Println(emptyMsg)
		return
	}

	// Mission count with cyber styling
	countStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwaveCyan).
		Background(styles.CyberGrid).
		Bold(true).
		Padding(0, 2).
		Render(fmt.Sprintf("📋 %d MISSIONS FOUND", len(tasks)))
	fmt.Println(countStyle)
	fmt.Println()

	// Render each task with maximum impact
	for i, task := range tasks {
		renderTaskWithEffects(i+1, task)
		fmt.Println()
	}

	// Stats footer with holographic text
	active := 0
	completed := 0
	for _, task := range tasks {
		if task.Status == "completed" {
			completed++
		} else {
			active++
		}
	}

	stats := styles.CyberStats(active, completed, len(tasks))
	fmt.Println(stats)
	fmt.Println()

	// Clean footer
	footer := styles.Title("focus.sh Task Management")
	fmt.Println(footer)
}

func renderTaskWithEffects(index int, task models.Task) {
	// Task number with focus effect
	numberStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwavePink).
		Background(styles.DeepSpace).
		Bold(true).
		Render(fmt.Sprintf("%02d", index))

	// Task description with cyber styling
	var descStyle lipgloss.Style
	if task.Status == "completed" {
		descStyle = lipgloss.NewStyle().
			Foreground(styles.SynthwaveGreen).
			Background(styles.DarkVoid).
			Strikethrough(true).
			Bold(true)
	} else {
		descStyle = lipgloss.NewStyle().
			Foreground(styles.SynthwaveCyan).
			Background(styles.DeepSpace).
			Bold(true)
	}

	description := descStyle.Render(task.Description)

	// Priority explosion
	priority := styles.PriorityExplosion(task.Priority)

	// Status indicator
	status := styles.TaskStatus(task.Status)

	// Combine all elements
	taskLine := fmt.Sprintf("%s %s %s %s",
		numberStyle,
		description,
		priority,
		status)

	// Box the entire task
	taskBox := styles.FocusBox(taskLine, styles.SynthwavePurple)
	fmt.Println(taskBox)

	// Add metadata if available
	if len(task.Categories) > 0 {
		tags := make([]string, len(task.Categories))
		for i, cat := range task.Categories {
			tags[i] = styles.CyberTag(cat)
		}
		fmt.Println("   " + strings.Join(tags, " "))
	}

	// Task ID with digital rain effect
	idStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwaveCyan).
		Background(styles.DeepSpace).
		Italic(true).
		Render("🔑 ID: " + task.ID)
	fmt.Println("   " + idStyle)
}

func init() {
	// listCmd is now registered in root.go to avoid duplication
}
