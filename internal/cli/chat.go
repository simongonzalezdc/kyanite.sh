package cli

import (
	"bufio"
	"context"
	"fmt"
	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "🤖 Get AI assistance with your tasks and app usage",
	Long:  "Ask questions about your tasks, productivity, or how to use focus.sh",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🤖 Initializing Chat Assistant 🤖")
		fmt.Println("~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~")

		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)
		aiManager := ai.New()

		// Check AI status and show indicator
		status := utils.CheckAIStatus()
		fmt.Print("AI Status: ")
		utils.StreamStatusWithIndicator(status)
		fmt.Println()

		// Get all tasks for context
		tasks, err := engine.ListTasks("all")
		if err != nil {
			fmt.Printf("❌ Error loading tasks: %v\n", err)
			return
		}

		// Create task descriptions for AI context
		taskDescriptions := make([]string, len(tasks))
		for i, task := range tasks {
			status := "pending"
			if task.Status == "completed" {
				status = "completed"
			}
			taskDescriptions[i] = fmt.Sprintf("%s (%s)", task.Description, status)
		}

		// Start chat loop
		fmt.Println("💡 Chat Assistant Ready! Ask me anything about your tasks or app usage.")
		fmt.Println("   Type 'exit' or 'quit' to end the conversation.")
		fmt.Println()

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("You: ")
			if !scanner.Scan() {
				break
			}

			question := strings.TrimSpace(scanner.Text())
			if question == "" {
				continue
			}

			// Check for exit commands
			if question == "exit" || question == "quit" || question == "q" {
				fmt.Println("👋 Goodbye! See you in the grid!")
				break
			}

			// Get AI response with streaming
			aiLabel := lipgloss.NewStyle().
				Foreground(styles.GetAccent()).
				Bold(true).
				Render("AI: ")
			fmt.Print(aiLabel)

			// Show thinking indicator
			thinkingText := "🤔 Processing your query..."
			utils.StreamText(thinkingText, styles.GetWarning(), 30*time.Millisecond)
			fmt.Print("\r")                                      // Clear line
			fmt.Print(strings.Repeat(" ", len(thinkingText)+10)) // Clear with spaces
			fmt.Print("\r")                                      // Return to start
			fmt.Print(aiLabel)

			ctx := context.Background()
			response, err := aiManager.ChatAssistant(ctx, question, taskDescriptions)
			if err != nil {
				errorText := "❌ Sorry, I encountered an error processing your request."
				utils.StreamWithTypingEffect(errorText, styles.GetError())
				fmt.Println()
				continue
			}

			// Stream the AI response with colors
			utils.StreamWithTypingEffect(response, styles.GetForeground())
			fmt.Println()
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("❌ Error reading input: %v\n", err)
		}
	},
}
