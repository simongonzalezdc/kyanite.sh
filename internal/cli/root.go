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
		// Default to TUI dashboard - TUI-FIRST approach
		fmt.Println("🚀 Initializing focus.sh...")
		fmt.Println("⚡ Loading AI protocols and visual enhancement systems...")
		
		// Start the full-featured dashboard
		err := tui.StartMainDashboard([]tui.DashboardTask{})
		if err != nil {
			fmt.Printf("Error starting mission matrix: %v\n", err)
		}
	},
}

func Execute() error {
	// Setup initial appearance
	pterm.Println(pterm.FgGreen.Sprintf("🎮 focus.sh initialized"))

	// Execute the root command
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(notesCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(filterCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(simpleInteractiveCmd)
	rootCmd.AddCommand(testGumCmd)
	rootCmd.AddCommand(themeCmd)
	rootCmd.AddCommand(wizardCmd)
	rootCmd.AddCommand(configWizardCmd)
	rootCmd.AddCommand(editWizardCmd)
	
	// Add viper config commands
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configPathCmd)
	rootCmd.AddCommand(enhancedConfigCmd)
	rootCmd.AddCommand(unifiedDashboardCmd)
	rootCmd.AddCommand(mcpServerCmd)

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
