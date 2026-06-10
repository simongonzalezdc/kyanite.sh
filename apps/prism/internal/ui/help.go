package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/theme"
)

// RenderHelp renders the help screen
func RenderHelp(tm *theme.Manager, width, height int) string {
	styles := NewStyles(tm.CurrentTheme())

	var b strings.Builder

	// Title
	title := styles.Title.Render("PRISM.SH - Keyboard Shortcuts")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Global shortcuts
	b.WriteString(styles.Primary.Render("Global Shortcuts:"))
	b.WriteString("\n")
	b.WriteString(renderShortcut(styles, "Ctrl+Q", "Quit application"))
	b.WriteString(renderShortcut(styles, "Ctrl+H", "Toggle this help screen"))
	b.WriteString(renderShortcut(styles, "Ctrl+Shift+T", "Cycle theme"))
	b.WriteString(renderShortcut(styles, "Esc", "Back to menu / Quit"))
	b.WriteString("\n")

	// Navigation
	b.WriteString(styles.Primary.Render("Navigation:"))
	b.WriteString("\n")
	b.WriteString(renderShortcut(styles, "↑/↓ or k/j", "Navigate up/down"))
	b.WriteString(renderShortcut(styles, "←/→ or h/l", "Navigate left/right"))
	b.WriteString(renderShortcut(styles, "Enter or Space", "Select/Confirm"))
	b.WriteString("\n")

	// Screen-specific
	b.WriteString(styles.Primary.Render("Color Wheel:"))
	b.WriteString("\n")
	b.WriteString(renderShortcut(styles, "←/→", "Change hue"))
	b.WriteString(renderShortcut(styles, "↑/↓", "Change lightness"))
	b.WriteString(renderShortcut(styles, "+/-", "Change saturation"))
	b.WriteString(renderShortcut(styles, "C", "Copy color to clipboard"))
	b.WriteString("\n")

	b.WriteString(styles.Primary.Render("Palette Generator:"))
	b.WriteString("\n")
	b.WriteString(renderShortcut(styles, "↑/↓", "Select harmony rule"))
	b.WriteString(renderShortcut(styles, "Enter", "Generate palette"))
	b.WriteString(renderShortcut(styles, "C", "Copy palette colors"))
	b.WriteString(renderShortcut(styles, "S", "Save palette"))
	b.WriteString(renderShortcut(styles, "E", "Export palette"))
	b.WriteString("\n")

	b.WriteString(styles.Primary.Render("WCAG Checker:"))
	b.WriteString("\n")
	b.WriteString(renderShortcut(styles, "Enter", "Check contrast"))
	b.WriteString(renderShortcut(styles, "C", "Copy result to clipboard"))
	b.WriteString("\n")

	b.WriteString(styles.Primary.Render("Color Search:"))
	b.WriteString("\n")
	b.WriteString(renderShortcut(styles, "Type", "Search colors by name"))
	b.WriteString(renderShortcut(styles, "↑/↓", "Navigate results"))
	b.WriteString(renderShortcut(styles, "C", "Copy selected color"))
	b.WriteString("\n")

	b.WriteString(styles.Primary.Render("Palette Manager:"))
	b.WriteString("\n")
	b.WriteString(renderShortcut(styles, "D", "Delete selected palette"))
	b.WriteString(renderShortcut(styles, "R", "Refresh list"))
	b.WriteString("\n")

	// Footer
	footer := styles.Muted.Render("Press Ctrl+H or Esc to close this help")
	b.WriteString("\n")
	b.WriteString(footer)

	// Wrap in border
	content := styles.Border.Width(ContentWidthNarrow).Render(b.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// renderShortcut renders a keyboard shortcut line
func renderShortcut(styles Styles, key, desc string) string {
	keyStyle := styles.HelpKey.Render(key)
	descStyle := styles.HelpDesc.Render(desc)
	return lipgloss.JoinHorizontal(lipgloss.Left,
		keyStyle,
		strings.Repeat(" ", 20-lipgloss.Width(key)),
		descStyle,
	) + "\n"
}
