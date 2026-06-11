package cli

import (
	"context"
	"fmt"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:   "inspire",
	Short: "🔮 Generate AI-inspired mission suggestions",
	Long:  "🌸 Receive AI-powered insights from the digital realm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔮 Tuning into Digital Frequencies 🔮")
		fmt.Println("~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~")

		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)
		aiManager := defaultContainer.GetAIManager()

		// Check AI status and show indicator
		status := utils.CheckAIStatus()
		fmt.Print("AI Status: ")
		utils.StreamStatusWithIndicator(status)
		fmt.Println()

		// Get existing tasks for context
		tasks, err := engine.ListTasks("all")
		if err != nil {
			fmt.Printf("❌ Error accessing mission database: %v\n", err)
			return
		}

		// Create task descriptions for AI
		taskDescriptions := make([]string, len(tasks))
		for i, task := range tasks {
			status := "pending"
			if task.Status == "completed" {
				status = "completed"
			}
			taskDescriptions[i] = fmt.Sprintf("%s (%s)", task.Description, status)
		}

		// Get AI suggestions with streaming
		fmt.Println("🧠 Channeling suggestion matrix...")
		ctx := context.Background()

		// Show streaming AI thinking
		thinkingText := "Processing..."
		utils.StreamWithTypingEffect(thinkingText, styles.GetAccent())

		suggestions, err := aiManager.SuggestTasks(ctx, taskDescriptions)
		if err != nil {
			fmt.Printf("⚠️ Error generating suggestions: %v\n", err)
			return
		}

		if len(suggestions) == 0 {
			emptyText := "No suggestions available."
			utils.StreamWithTypingEffect(emptyText, styles.GetWarning())

			hintText := "Try adding more tasks to unlock suggestions!"
			utils.StreamWithTypingEffect(hintText, styles.GetAccent())
			return
		}

		// Stream suggestions with consistent styling
		titleStyle := lipgloss.NewStyle().
			Foreground(styles.GetSuccess()).
			Bold(true).
			Render("Task Suggestions:")
		fmt.Println(titleStyle)
		fmt.Println()

		for i, suggestion := range suggestions {
			suggestionText := fmt.Sprintf("%d. %s", i+1, suggestion)
			utils.StreamWithTypingEffect(suggestionText, styles.SynthwaveCyan)
		}

		fmt.Println("\nTo upload a suggestion, use: focus add \"<suggestion text>\"")
		fmt.Println("💡 Pro Tip: Suggestions evolve with your mission patterns!")
	},
}
