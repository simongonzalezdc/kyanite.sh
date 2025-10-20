package cli

import (
	"context"
	"fmt"
	"strings"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/validation"
	"github.com/kyanite/focus/internal/ai"

	"github.com/spf13/cobra"
	"github.com/charmbracelet/lipgloss"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new mission to the synthwave matrix",
	Long: `Create a new mission with AI-powered parsing and stunning visual effects.
Example: focus add "Complete the synthwave project by Friday"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showAddError()
			return
		}

		description := strings.Join(args, " ")
		
		// Validate input
		if err := validation.ValidateTaskDescription(description); err != nil {
			showValidationError(err.Error())
			return
		}
		
		// Sanitize input
		description = validation.SanitizeInput(description)
		
		// Show processing animation
		showProcessingAnimation(description)
		
		// Initialize components
		store := store.New(utils.GetStoragePath())
		engine := engine.New(store)
		aiManager := ai.New()

		// AI-powered task parsing
		var task models.ParsedTask
		parsedTask, err := aiManager.ParseTask(context.Background(), description)
		if err != nil {
			// Fallback to basic task creation if AI fails
			task = models.ParsedTask{
				Description: description,
				Priority:    "medium",
			}
		} else {
			task = models.ParsedTask{
				Description: parsedTask.Description,
				Priority:    parsedTask.Priority,
				Deadline:    parsedTask.Deadline,
				Categories:  parsedTask.Categories,
			}
		}

		// Add the task
		addedTask, err := engine.AddTask(task)
		if err != nil {
			showAddError()
			return
		}

		// Show success with maximum impact
		showAddSuccess(addedTask)
	},
}

func showProcessingAnimation(_ string) {
	// Loading animation with synthwave styling
	loadingFrames := []string{
		"⚡ INITIALIZING SYNTHWAVE PROTOCOLS",
		"🔮 ANALYZING MISSION PARAMETERS", 
		"💫 CALIBRATING NEURAL INTERFACES",
		"✨ ACTIVATING DIGITAL ENHANCEMENTS",
		"🌌 INTEGRATING WITH THE MATRIX",
	}

	fmt.Println(styles.LoadingAnimation())
	
	for _, frame := range loadingFrames {
		frameStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveCyan).
			Background(styles.DeepSpace).
			Bold(true).
			Render("▶ " + frame)
		fmt.Println(frameStyle)
		fmt.Println()
	}
}

func showValidationError(message string) {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ VALIDATION ERROR\n\n%s", message))
	fmt.Println(errorBox)
}

func showAddError() {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render("❌ MISSION CREATION FAILED\n\nPlease provide a task description:\nfocus add \"Your mission here\"")
	fmt.Println(errorBox)
}

func showAddSuccess(task models.Task) {
	// Epic success message
	successTitle := styles.SynthwaveTitle("🎯 MISSION ACCEPTED")
	fmt.Println(successTitle)
	fmt.Println()

	// Task details with cyber styling
	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"MISSION ID", task.ID, styles.SynthwaveCyan},
		{"DESCRIPTION", task.Description, styles.SynthwavePink},
		{"STATUS", task.Status, styles.SynthwaveGreen},
		{"PRIORITY", task.Priority, styles.SynthwaveYellow},
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

	// Confirmation with holographic effect
	confirmMsg := styles.HolographicText("✨ Mission successfully integrated into the matrix!")
	fmt.Println(confirmMsg)

	// Digital artifact celebration
	artifacts := []string{
		"⚡⚡⚡", "✨✨✨", "🔥🔥🔥", "💫💫💫", "🌟🌟🌟",
	}
	
	artifactLine := ""
	for _, artifact := range artifacts {
		artifactLine += artifact + " "
	}
	
	artifactStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwavePink).
		Background(styles.DeepSpace).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Render(artifactLine)
	fmt.Println(artifactStyle)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
