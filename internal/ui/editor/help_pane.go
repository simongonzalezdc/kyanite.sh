package editor

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/lyricforge/internal/ui/styles"
)

// HelpPaneModel displays keyboard shortcuts and help information
type HelpPaneModel struct {
	width   int
	height  int
	focused bool

	// Keyboard shortcuts
	shortcutManager *ShortcutManager

	// Styles
	containerStyle lipgloss.Style
	titleStyle     lipgloss.Style
	categoryStyle  lipgloss.Style
	shortcutStyle  lipgloss.Style
	descStyle      lipgloss.Style
	borderStyle    lipgloss.Style

	// Responsive layout
	compactMode     bool
	showMinimalHelp bool
	showShortKeys   bool
}

// NewHelpPaneModel creates a new help pane model
func NewHelpPaneModel(shortcutManager *ShortcutManager) *HelpPaneModel {
	model := &HelpPaneModel{
		shortcutManager: shortcutManager,
		containerStyle:  styles.BorderActive,
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Accent).
			MarginBottom(1),
		categoryStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Primary).
			MarginTop(1).
			MarginBottom(1),
		shortcutStyle: lipgloss.NewStyle().
			Foreground(styles.Info).
			Bold(true),
		descStyle: lipgloss.NewStyle().
			Foreground(styles.TextSecondary),
		borderStyle:     styles.Border,
		compactMode:     false,
		showMinimalHelp: false,
		showShortKeys:   false,
	}

	return model
}

// Init initializes the help pane
func (m *HelpPaneModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the help pane
func (m *HelpPaneModel) Update(msg tea.Msg) (*HelpPaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			// Exit help mode
			if m.shortcutManager != nil {
				m.shortcutManager.SetHelpMode(false)
			}
		}
	}

	return m, nil
}

// View renders the help pane
func (m *HelpPaneModel) View() string {
	// Update responsive mode based on current dimensions
	m.updateResponsiveMode()

	// Render content based on responsive mode
	if m.showMinimalHelp {
		return m.renderMinimalHelp()
	}

	if m.compactMode {
		return m.renderCompactHelp()
	}

	return m.renderFullHelp()
}

// renderFullHelp renders the complete help content
func (m *HelpPaneModel) renderFullHelp() string {
	title := m.titleStyle.Render("🎹 Keyboard Shortcuts Reference")
	title = lipgloss.NewStyle().Width(m.width - 4).Align(lipgloss.Center).Render(title)

	content := m.renderShortcutsHelp()

	fullContent := lipgloss.JoinVertical(lipgloss.Left, title, content)

	// Add footer with navigation hint
	footer := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Align(lipgloss.Center).
		Width(m.width - 4).
		Render("\nPress ESC, Q, or Enter to return to editor")

	fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, footer)

	if m.width > 0 && m.height > 0 {
		return m.containerStyle.Width(m.width).Height(m.height).Render(fullContent)
	}

	return m.containerStyle.Render(fullContent)
}

// renderCompactHelp renders compact help content
func (m *HelpPaneModel) renderCompactHelp() string {
	title := m.titleStyle.Render("🎹 Shortcuts")
	title = lipgloss.NewStyle().Width(m.width - 4).Align(lipgloss.Center).Render(title)

	content := m.renderCompactShortcutsHelp()

	fullContent := lipgloss.JoinVertical(lipgloss.Left, title, content)

	// Add compact footer
	footer := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Align(lipgloss.Center).
		Width(m.width - 4).
		Render("\nESC/Q/Enter: back")

	fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, footer)

	if m.width > 0 && m.height > 0 {
		return m.containerStyle.Width(m.width).Height(m.height).Render(fullContent)
	}

	return m.containerStyle.Render(fullContent)
}

// renderMinimalHelp renders minimal help content for very small terminals
func (m *HelpPaneModel) renderMinimalHelp() string {
	title := m.titleStyle.Render("❓ Help")
	title = lipgloss.NewStyle().Width(m.width - 4).Align(lipgloss.Center).Render(title)

	content := m.renderMinimalShortcutsHelp()

	fullContent := lipgloss.JoinVertical(lipgloss.Left, title, content)

	// Add minimal footer
	footer := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Align(lipgloss.Center).
		Width(m.width - 4).
		Render("\nESC: back")

	fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, footer)

	if m.width > 0 && m.height > 0 {
		return m.containerStyle.Width(m.width).Height(m.height).Render(fullContent)
	}

	return m.containerStyle.Render(fullContent)
}

// renderShortcutsHelp renders the shortcuts help content
func (m *HelpPaneModel) renderShortcutsHelp() string {
	var sections []string

	// Get current context
	context := ContextGlobal
	if m.shortcutManager != nil {
		context = m.shortcutManager.GetContext()
	}

	// Define categories in display order
	categories := []string{"Navigation", "Edit", "Search", "File", "View", "Application", "Tools"}

	for _, category := range categories {
		bindings := m.getBindingsByCategory(category, context)
		if len(bindings) > 0 {
			sections = append(sections, m.renderCategorySection(category, bindings))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// getBindingsByCategory returns bindings for a category in the current context
func (m *HelpPaneModel) getBindingsByCategory(category string, context KeyContext) []*KeyBinding {
	if m.shortcutManager == nil {
		return nil
	}

	var bindings []*KeyBinding
	allBindings := m.shortcutManager.GetAllBindings()

	for _, binding := range allBindings {
		if binding.Category == category && (binding.Context == context || binding.Context == ContextGlobal) {
			bindings = append(bindings, binding)
		}
	}

	return bindings
}

// renderCategorySection renders a category section with its bindings
func (m *HelpPaneModel) renderCategorySection(category string, bindings []*KeyBinding) string {
	var lines []string

	// Category header
	lines = append(lines, m.categoryStyle.Render(fmt.Sprintf("📂 %s", category)))

	// Bindings
	for _, binding := range bindings {
		keyStr := binding.Key.Help().Key
		desc := binding.Description

		// Format: "Ctrl+C    Copy text"
		line := fmt.Sprintf("  %-12s %s", m.shortcutStyle.Render(keyStr), m.descStyle.Render(desc))
		lines = append(lines, line)
	}

	// Add spacing between categories
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

// SetDimensions sets the pane dimensions
func (m *HelpPaneModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// Focus focuses the help pane
func (m *HelpPaneModel) Focus() {
	m.focused = true
}

// Blur blurs the help pane
func (m *HelpPaneModel) Blur() {
	m.focused = false
}

// SetShortcutManager sets the shortcut manager
func (m *HelpPaneModel) SetShortcutManager(shortcutManager *ShortcutManager) {
	m.shortcutManager = shortcutManager
}

// GetShortcutManager returns the shortcut manager
func (m *HelpPaneModel) GetShortcutManager() *ShortcutManager {
	return m.shortcutManager
}

// updateResponsiveMode updates the responsive display mode based on terminal width
func (m *HelpPaneModel) updateResponsiveMode() {
	// Enable compact mode for smaller terminals
	compactMode := m.width < 100
	minimalHelp := m.width < 80
	shortKeys := m.width < 90

	// Only update if mode has actually changed
	if m.compactMode != compactMode || m.showMinimalHelp != minimalHelp || m.showShortKeys != shortKeys {
		m.compactMode = compactMode
		m.showMinimalHelp = minimalHelp
		m.showShortKeys = shortKeys
	}
}

// renderCompactShortcutsHelp renders compact shortcuts help content
func (m *HelpPaneModel) renderCompactShortcutsHelp() string {
	var sections []string

	// Get current context
	context := ContextGlobal
	if m.shortcutManager != nil {
		context = m.shortcutManager.GetContext()
	}

	// Define essential categories only for compact mode
	categories := []string{"Navigation", "Edit", "File"}

	for _, category := range categories {
		bindings := m.getBindingsByCategory(category, context)
		if len(bindings) > 0 {
			sections = append(sections, m.renderCompactCategorySection(category, bindings))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderMinimalShortcutsHelp renders minimal shortcuts help content
func (m *HelpPaneModel) renderMinimalShortcutsHelp() string {
	var sections []string

	// Get current context
	context := ContextGlobal
	if m.shortcutManager != nil {
		context = m.shortcutManager.GetContext()
	}

	// Show only the most essential shortcuts for minimal mode
	essentialKeys := []string{"Navigation", "Edit"}
	essentialBindings := []*KeyBinding{}

	for _, category := range essentialKeys {
		bindings := m.getBindingsByCategory(category, context)
		essentialBindings = append(essentialBindings, bindings...)
	}

	// Show only first 8 most important bindings
	maxBindings := len(essentialBindings)
	if maxBindings > 8 {
		maxBindings = 8
	}

	if len(essentialBindings) > 0 {
		sections = append(sections, m.renderMinimalBindings(essentialBindings[:maxBindings]))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderCompactCategorySection renders a compact category section
func (m *HelpPaneModel) renderCompactCategorySection(category string, bindings []*KeyBinding) string {
	var lines []string

	// Compact category header
	lines = append(lines, m.categoryStyle.Render(fmt.Sprintf("📂 %s", category)))

	// Compact bindings (key only, no description for very small spaces)
	for _, binding := range bindings {
		keyStr := binding.Key.Help().Key

		// Format: "Ctrl+C"
		line := fmt.Sprintf("  %s", m.shortcutStyle.Render(keyStr))
		lines = append(lines, line)
	}

	// Add spacing between categories
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

// renderMinimalBindings renders minimal binding list
func (m *HelpPaneModel) renderMinimalBindings(bindings []*KeyBinding) string {
	var lines []string

	// Simple list of key combinations
	for _, binding := range bindings {
		keyStr := binding.Key.Help().Key
		line := fmt.Sprintf("  %s", m.shortcutStyle.Render(keyStr))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
