package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
	"github.com/kyanite/focus/internal/theme"
	"github.com/kyanite/config"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme [theme-name]",
	Short: "🎨 Change the visual theme",
	Long: `Switch between different visual themes.
Available Kyanite Themes:
  - amber-night: Warm amber glow with deep contrast (default)
  - indigo-depths: Deep blue indigo with subtle accents  
  - forest-path: Natural green and brown forest tones
  - cyan-wave: Cool cyan and teal ocean waves
  - electric-rose: Vibrant pink and electric purple
  - monochrome: Clean black and white minimalism`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := strings.ToLower(args[0])

		// Validate theme exists
		availableThemes := design.List()
		validTheme := false
		for _, t := range availableThemes {
			if t == themeName {
				validTheme = true
				break
			}
		}

		if !validTheme {
			fmt.Printf("❌ Unknown theme: %s\n", themeName)
			fmt.Printf("Available themes: %v\n", availableThemes)
			return
		}

		// Set theme globally
		theme.GetManager().SetTheme(themeName)

		// Refresh styles to use new theme colors
		styles.RefreshColors()

		// Update config
		root, err := config.Load()
		if err != nil {
			root = &config.Root{}
		}
		root.Focus.Theme = themeName
		_ = config.Save(root) // Ignore error for UI command

		fmt.Printf("✨ Theme changed to: %s\n", lipgloss.NewStyle().
			Foreground(theme.GetManager().Current().Primary).
			Bold(true).
			Render(themeName))
	},
}

func init() {
	// themeCmd is now registered in root.go to avoid duplication
}
