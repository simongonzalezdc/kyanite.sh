package cli

import (
	"fmt"
	"github.com/kyanite/focus/internal/tui"
	"github.com/spf13/cobra"
	pterm "github.com/pterm/pterm"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusBlue    = lipgloss.Color("#00FFFF")
	focusGreen   = lipgloss.Color("#00FF66")
	focusPurple  = lipgloss.Color("#B967C7")
	accentColor  = lipgloss.Color("#FFC0CB")

	headerStyle = lipgloss.NewStyle().
		Foreground(focusBlue).
		Bold(true).
		AlignHorizontal(lipgloss.Center)





	boxStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(focusBlue).
		Padding(1).
		Width(80)
)

var rootCmd = &cobra.Command{
	Use:   "focus",
	Short: "🌟 focus.sh - Kyanite Suite Task Manager",
	Long: boxStyle.Render(
		headerStyle.Render("🌟 W E L C O M E   T O   F O C U S . S H 🌟") + "\n\n" +
		lipgloss.NewStyle().Foreground(focusGreen).Render(
			"🚀 Part of the Kyanite Creative Suite\n" +
			"✨ Professional task management with AI assistance\n" +
			"🌈 Clean, focused productivity experience\n\n" +
			"Controls: A(dd) | C(omplete) | L(ist) | J(ournal) | T(heme) | ?(help)") + "\n" +
		lipgloss.NewStyle().Foreground(focusBlue).Render(
			"────────────────────────────────────────────────────────────────────────────────") + "\n" +
		lipgloss.NewStyle().Foreground(focusPurple).Render(
			"💡 Part of Kyanite Suite: noise.sh (music) | focus.sh (tasks) | syntax.sh (writing) | prism.sh (visual) | wave.sh (audio)")),
	Run: func(cmd *cobra.Command, args []string) {
		// Launch directly into TUI dashboard
		err := tui.StartMainDashboard([]tui.DashboardTask{})
		if err != nil {
			fmt.Printf("Error starting dashboard: %v\n", err)
		}
	},
}

func Execute() error {
	// Kyanite Suite loading sequence
	pterm.Println(pterm.FgGreen.Sprintf("🎮 focus.sh initialized"))
	
	// Show Kyanite suite loading animation
	pterm.Info.Println("🌌 Loading Kyanite Suite components...")
	pterm.Success.Println("✨ focus.sh - Task Management System Online")
	pterm.Info.Println("🔗 Connecting to Kyanite Creative Suite...")
	
	// Execute the root command
	return rootCmd.Execute()
}

func init() {
	// Core commands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(notesCmd)
	rootCmd.AddCommand(journalCmd)
	rootCmd.AddCommand(calendarCmd)
	
	// Interactive commands  
	rootCmd.AddCommand(wizardCmd)
	rootCmd.AddCommand(interactiveCmd)
	
	// Configuration commands
	rootCmd.AddCommand(enhancedConfigCmd)
	rootCmd.AddCommand(themeCmd)
	
	// Add viper config subcommands
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configPathCmd)
	rootCmd.AddCommand(configCmd)
	
	// Dashboard commands
	rootCmd.AddCommand(unifiedDashboardCmd)
	
	// Additional commands
	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(editWizardCmd)
	rootCmd.AddCommand(simpleInteractiveCmd)
	
	// Hidden/technical commands (not for regular users)
	// These can be added but hidden from help if needed
	rootCmd.AddCommand(mcpServerCmd)
	rootCmd.AddCommand(filterCmd)
	rootCmd.AddCommand(testGumCmd)

	// Setup pterm styles
	pterm.Info.Prefix = pterm.Prefix{
		Text:  "💡",
		Style: pterm.Info.Prefix.Style,
	}
	pterm.Error.Prefix = pterm.Prefix{
		Text:  "❌",
		Style: pterm.Error.Prefix.Style,
	}
}
