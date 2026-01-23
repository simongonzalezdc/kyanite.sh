package ui

import (
	"fmt"
	"time"

	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// splashSubtitleReadyMsg triggers the subtitle animation after a delay
type splashSubtitleReadyMsg struct{}

// SplashModel handles the splash screen
type SplashModel struct {
	spinner      spinner.Model
	width        int
	height       int
	animation    *AnimationManager
	fadeIn       bool
	showTitle    bool
	showSubtitle bool
}

// NewSplashModel creates a new splash screen model
func NewSplashModel() *SplashModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.GetManager().Current().Primary)

	return &SplashModel{
		spinner:      s,
		animation:    NewAnimationManager(),
		fadeIn:       true,
		showTitle:    false,
		showSubtitle: false,
	}
}

// Init initializes the splash model
func (m *SplashModel) Init() tea.Cmd {
	// Start fade-in animation for title immediately
	m.animation.FadeTransition("splash_title_fade", 1.0)

	// Use tea.Tick for staggered subtitle animation instead of blocking sleep
	// This keeps the event loop responsive
	return tea.Batch(
		m.spinner.Tick,
		m.animation.Update(), // Start animation tick
		tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
			return splashSubtitleReadyMsg{}
		}),
	)
}

// Update handles messages for the splash screen
func (m *SplashModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case splashSubtitleReadyMsg:
		// Start subtitle animation after the delay (triggered by tea.Tick)
		m.animation.FadeTransition("splash_subtitle_fade", 1.0)
		cmds = append(cmds, m.animation.Update())

	case AnimationTickMsg:
		// Update animation states based on progress
		titleProgress := m.animation.GetAnimationProgress("splash_title_fade")
		m.showTitle = titleProgress > 0.3
		m.showSubtitle = titleProgress > 0.6

		// Continue animation if still running
		cmd := m.animation.Update()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	// Update spinner
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// View renders the splash screen
func (m *SplashModel) View() string {
	t := theme.GetManager().Current()
	titleStyle := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height)

	// Get animation progress
	titleProgress := m.animation.GetAnimationProgress("splash_title_fade")
	subtitleProgress := m.animation.GetAnimationProgress("splash_subtitle_fade")

	content := ""

	// Apply fade effects based on animation progress
	if m.showTitle || titleProgress > 0 {
		title := "[~] noise.sh [~]"
		if titleProgress < 1.0 && titleProgress > 0 {
			// Apply fade effect to title
			titleOpacity := uint(titleProgress * 255)
			if titleOpacity > 255 {
				titleOpacity = 255
			}
			titleStyle = titleStyle.Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", titleOpacity, titleOpacity/2, titleOpacity)))
		}
		content += title + "\n\n"
	}

	if m.showSubtitle || subtitleProgress > 0 {
		subtitle := "AI-Powered Songwriting Terminal Interface"
		if subtitleProgress < 1.0 && subtitleProgress > 0 {
			// Apply fade effect to subtitle
			subtitleOpacity := uint(subtitleProgress * 255)
			if subtitleOpacity > 255 {
				subtitleOpacity = 255
			}
			subtitleStyle := titleStyle.Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", subtitleOpacity/2, subtitleOpacity/2, subtitleOpacity)))
			content += subtitleStyle.Render(subtitle) + "\n\n"
		} else if subtitleProgress >= 1.0 {
			content += subtitle + "\n\n"
		}
	}

	// Always show spinner and loading text
	content += m.spinner.View() + " Loading..."

	return titleStyle.Render(content)
}

// SetSize sets the dimensions for the splash screen
func (m *SplashModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// startStaggeredAnimations is deprecated - using tea.Tick in Init() instead
// to avoid blocking the event loop with time.Sleep
