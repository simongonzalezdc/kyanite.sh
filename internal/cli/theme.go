package cli

import (
	"fmt"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/config"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme [synthwave|light|plain]",
	Short: "🎨 Change the visual theme",
	Long: `Switch between different visual themes:
  - synthwave: Cyberpunk aesthetic with neon colors (default)
  - light: Clean, modern light theme
  - plain: Basic terminal appearance for maximum compatibility`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := strings.ToLower(args[0])
		
		var theme styles.ThemeMode
		var description string
		
		switch themeName {
		case "synthwave":
			theme = styles.ThemeSynthwave
			description = "🌌 Cyberpunk synthwave with maximum neon impact"
		case "light":
			theme = styles.ThemeLight
			description = "☀️ Clean and modern light theme"
		case "plain":
			theme = styles.ThemePlain
			description = "📄 Basic terminal appearance"
		default:
			errorStyle := lipgloss.NewStyle().
				Foreground(styles.SynthwaveRed).
				Bold(true).
				Padding(1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.SynthwaveRed).
				Render(fmt.Sprintf("❌ Unknown theme: %s\n\nAvailable themes: synthwave, light, plain", themeName))
			fmt.Println(errorStyle)
			return
		}
		
		// Set the theme for current session
		styles.SetTheme(theme)
		
		// Persist theme to config
		cfg := config.GetConfig()
		if cfg == nil {
			if loaded, err := config.LoadConfig(); err == nil {
				cfg = loaded
			}
		}
		if cfg != nil {
			cfg.Theme = string(theme)
			_ = config.SaveConfig(cfg)
		}
		
		// Show confirmation
		successStyle := lipgloss.NewStyle().
			Foreground(styles.GetSuccess()).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.GetSuccess()).
			Render(fmt.Sprintf("✅ Theme changed to %s", themeName))
		
		fmt.Println(successStyle)
		fmt.Println()
		
		descStyle := lipgloss.NewStyle().
			Foreground(styles.GetForeground()).
			Italic(true).
			Render(description)
		fmt.Println(descStyle)
		fmt.Println()
		
		// Show current settings
		infoStyle := lipgloss.NewStyle().
			Foreground(styles.GetAccent()).
			Render("💡 Theme saved. It will be used across CLI and TUI.")
		fmt.Println(infoStyle)
	},
}

func init() {
	rootCmd.AddCommand(themeCmd)
}
