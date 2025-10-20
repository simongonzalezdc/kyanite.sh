package cli

import (
	"context"
	"fmt"
	"time"
	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/kyanite/focus/pkg/styles"

	"github.com/spf13/cobra"
	"github.com/charmbracelet/lipgloss"
)

var suggestCmd = &cobra.Command{
	Use:   "inspire",
	Short: "🔮 Generate AI-inspired mission suggestions",
	Long:  "🌸 Receive AI-powered insights from the digital realm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔮 Tuning into Digital Frequencies 🔮")
		fmt.Println("~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~")

		// Initialize components
		store := store.New(utils.GetStoragePath())
		engine := engine.New(store)
		aiManager := ai.New()

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
		thinkingText := "🤖 Processing neural pathways..."
		utils.StreamWithTypingEffect(thinkingText, styles.GetAccent())
		time.Sleep(500 * time.Millisecond)
		
		suggestions, err := aiManager.SuggestTasks(ctx, taskDescriptions)
		if err != nil {
			fmt.Printf("⚠️ Error generating suggestions: %v\n", err)
			return
		}

		if len(suggestions) == 0 {
			emptyText := "😴 No new missions received from the digital realm."
			utils.StreamWithTypingEffect(emptyText, styles.GetWarning())
			
			hintText := "💡 Try adding more missions to unlock suggestions!"
			utils.StreamWithTypingEffect(hintText, styles.GetAccent())
			return
		}

		// Stream suggestions with colors
		titleStyle := lipgloss.NewStyle().
			Foreground(styles.GetSuccess()).
			Bold(true).
			Render("✨ Mission Suggestions from the Grid:")
		fmt.Println(titleStyle)
		fmt.Println()
		
		for i, suggestion := range suggestions {
			// Cycle through synthwave colors for each suggestion
			colors := []lipgloss.Color{
				styles.SynthwavePink,
				styles.SynthwaveCyan,
				styles.SynthwaveGreen,
				styles.SynthwaveYellow,
				styles.SynthwavePurple,
			}
			color := colors[i%len(colors)]
			
			suggestionText := fmt.Sprintf("%d. 🌟 %s", i+1, suggestion)
			utils.StreamWithTypingEffect(suggestionText, color)
			time.Sleep(200 * time.Millisecond)
		}
		
		fmt.Println("\nTo upload a suggestion, use: neon add \"<suggestion text>\"")
		fmt.Println("💡 Pro Tip: Suggestions evolve with your mission patterns!")
	},
}
