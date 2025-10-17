package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SplashModel handles the splash screen
type SplashModel struct {
	spinner spinner.Model
	width   int
	height  int
}

// NewSplashModel creates a new splash screen model
func NewSplashModel() *SplashModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))

	return &SplashModel{
		spinner: s,
	}
}

// Init initializes the splash model
func (m *SplashModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages for the splash screen
func (m *SplashModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return cmd
}

// View renders the splash screen
func (m *SplashModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF69B4")).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height)

	content := "🎵 LyricForge 🎵\n\n"
	content += "AI-Powered Songwriting Terminal Interface\n\n"
	content += m.spinner.View() + " Loading..."

	return titleStyle.Render(content)
}

// SetSize sets the dimensions for the splash screen
func (m *SplashModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
