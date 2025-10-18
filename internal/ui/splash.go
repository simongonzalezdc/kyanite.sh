package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/lyricforge/internal/ui/styles"
)

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
	s.Style = lipgloss.NewStyle().Foreground(styles.Primary)

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
	// Start fade-in animation for title
	m.animation.FadeTransition("splash_title_fade", 1.0)

	// Start animations for different elements with staggered timing
	return tea.Batch(
		m.spinner.Tick,
		m.startStaggeredAnimations(),
	)
}

// Update handles messages for the splash screen
func (m *SplashModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
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
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
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
		title := "🎵 LyricForge 🎵"
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

// startStaggeredAnimations starts animations for different UI elements with delays
func (m *SplashModel) startStaggeredAnimations() tea.Cmd {
	return func() tea.Msg {
		// Start title fade-in immediately
		m.animation.FadeTransition("splash_title_fade", 1.0)

		// Start subtitle animation after a delay
		time.Sleep(300 * time.Millisecond)
		m.animation.FadeTransition("splash_subtitle_fade", 1.0)

		return nil
	}
}
