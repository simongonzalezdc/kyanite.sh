package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kyanite/noise/internal/logging"
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// =============================================================================
// Constants
// =============================================================================

const onboardingFileName = ".noise_onboarding_complete"

// OnboardingStep represents a step in the onboarding flow
type OnboardingStep struct {
	Title       string
	Description string
	Icon        string
	Tips        []string
}

// =============================================================================
// Onboarding Model
// =============================================================================

// OnboardingModel displays a first-run welcome experience
type OnboardingModel struct {
	visible     bool
	currentStep int
	steps       []OnboardingStep
	width       int
	height      int

	// Styles
	containerStyle    lipgloss.Style
	titleStyle        lipgloss.Style
	descriptionStyle  lipgloss.Style
	tipStyle          lipgloss.Style
	hintStyle         lipgloss.Style
	progressStyle     lipgloss.Style
	progressDotActive lipgloss.Style
	progressDotDone   lipgloss.Style
}

// NewOnboardingModel creates a new onboarding model
func NewOnboardingModel() *OnboardingModel {
	t := theme.GetManager().Current()

	m := &OnboardingModel{
		visible:     false,
		currentStep: 0,
		width:       80,
		height:      24,
		steps:       defaultOnboardingSteps(),

		containerStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Padding(2, 4),

		titleStyle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginBottom(1),

		descriptionStyle: lipgloss.NewStyle().
			Foreground(t.Text),

		tipStyle: lipgloss.NewStyle().
			Foreground(t.Accent).
			MarginTop(1).
			MarginLeft(2),

		hintStyle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Italic(true).
			MarginTop(2),

		progressStyle: lipgloss.NewStyle().
			Foreground(t.Secondary),

		progressDotActive: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		progressDotDone: lipgloss.NewStyle().
			Foreground(t.Success),
	}

	return m
}

// defaultOnboardingSteps returns the default onboarding steps
func defaultOnboardingSteps() []OnboardingStep {
	return []OnboardingStep{
		{
			Title:       "Welcome to noise.sh",
			Description: "Your personal music production notebook.\nCapture ideas, write lyrics, and organize your creative workflow.",
			Icon:        "♪",
			Tips: []string{
				"Press Tab to switch between panes",
				"Press ? to see all keyboard shortcuts",
			},
		},
		{
			Title:       "The Editor",
			Description: "Write lyrics, chord progressions, and production notes.\nReal-time preview shows formatted output as you type.",
			Icon:        "📝",
			Tips: []string{
				"Use markdown syntax for formatting",
				"Ctrl+S saves your work automatically",
			},
		},
		{
			Title:       "Quick Tools",
			Description: "Access built-in tools to speed up your workflow.",
			Icon:        "🛠",
			Tips: []string{
				"Ctrl+F opens the chord picker",
				"Ctrl+Shift+B opens the BPM tapper",
			},
		},
		{
			Title:       "Mobile Companion",
			Description: "Capture ideas on the go with the noise.sh PWA.\nVoice memos, photos, and notes sync to your desktop.",
			Icon:        "📱",
			Tips: []string{
				"Enable sync in Settings > Sync",
				"Scan the QR code to pair your phone",
			},
		},
		{
			Title:       "Themes",
			Description: "Customize your workspace with beautiful themes.\nFind a color scheme that inspires your creativity.",
			Icon:        "🎨",
			Tips: []string{
				"Ctrl+T cycles through themes",
				"10 themes available out of the box",
			},
		},
		{
			Title:       "You're All Set!",
			Description: "Start creating your next masterpiece.\nPress Enter to begin or ? anytime to see shortcuts.",
			Icon:        "🚀",
			Tips: []string{
				"Happy writing!",
			},
		},
	}
}

// Init initializes the model
func (m *OnboardingModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *OnboardingModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if !m.visible {
			return nil
		}

		switch msg.String() {
		case "enter", "right", "l", " ":
			if m.currentStep < len(m.steps)-1 {
				m.currentStep++
			} else {
				m.Complete()
			}
		case "left", "h":
			if m.currentStep > 0 {
				m.currentStep--
			}
		case "esc":
			m.Complete()
		case "home":
			m.currentStep = 0
		case "end":
			m.currentStep = len(m.steps) - 1
		}
	}

	return nil
}

// View renders the onboarding screen
func (m *OnboardingModel) View() string {
	if !m.visible || len(m.steps) == 0 {
		return ""
	}

	step := m.steps[m.currentStep]

	var content strings.Builder

	// Progress indicator
	var dots []string
	for i := range m.steps {
		if i < m.currentStep {
			dots = append(dots, m.progressDotDone.Render("[x]"))
		} else if i == m.currentStep {
			dots = append(dots, m.progressDotActive.Render("[*]"))
		} else {
			dots = append(dots, m.progressStyle.Render("[ ]"))
		}
	}
	progress := strings.Join(dots, " ")
	content.WriteString(lipgloss.NewStyle().MarginBottom(2).Render(progress))
	content.WriteString("\n\n")

	// Icon and title
	icon := lipgloss.NewStyle().
		Foreground(theme.GetManager().Current().Primary).
		Bold(true).
		Render(step.Icon)
	content.WriteString(icon + "  " + m.titleStyle.Render(step.Title))
	content.WriteString("\n\n")

	// Description
	content.WriteString(m.descriptionStyle.Render(step.Description))
	content.WriteString("\n")

	// Tips
	if len(step.Tips) > 0 {
		content.WriteString("\n")
		for _, tip := range step.Tips {
			content.WriteString(m.tipStyle.Render("- " + tip))
			content.WriteString("\n")
		}
	}

	// Navigation hint
	isFirst := m.currentStep == 0
	isLast := m.currentStep == len(m.steps)-1
	var hints []string
	if !isFirst {
		hints = append(hints, "<- Back")
	}
	if isLast {
		hints = append(hints, "Enter: Get Started")
	} else {
		hints = append(hints, "Enter/->: Next")
	}
	hints = append(hints, "Esc: Skip")
	content.WriteString(m.hintStyle.Render(strings.Join(hints, "  -  ")))

	// Step counter
	stepCount := fmt.Sprintf("\n\nStep %d of %d", m.currentStep+1, len(m.steps))
	content.WriteString(m.progressStyle.Render(stepCount))

	// Calculate container size
	containerWidth := min(70, m.width-8)

	// Render container
	container := m.containerStyle.
		Width(containerWidth).
		Render(content.String())

	// Center on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		container,
	)
}

// =============================================================================
// Public Methods
// =============================================================================

// IsVisible returns whether the onboarding is visible
func (m *OnboardingModel) IsVisible() bool {
	return m.visible
}

// Show displays the onboarding flow
func (m *OnboardingModel) Show() {
	m.visible = true
	m.currentStep = 0
}

// Hide hides the onboarding without marking complete
func (m *OnboardingModel) Hide() {
	m.visible = false
}

// Complete marks onboarding as complete and hides it
func (m *OnboardingModel) Complete() {
	m.visible = false
	m.markComplete()
}

// ShouldShow returns true if onboarding should be shown (first run)
func (m *OnboardingModel) ShouldShow() bool {
	return !isOnboardingComplete()
}

// SetDimensions sets the display dimensions
func (m *OnboardingModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// UpdateTheme refreshes styles when theme changes
func (m *OnboardingModel) UpdateTheme() {
	t := theme.GetManager().Current()

	m.containerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Padding(2, 4)

	m.titleStyle = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		MarginBottom(1)

	m.descriptionStyle = lipgloss.NewStyle().
		Foreground(t.Text)

	m.tipStyle = lipgloss.NewStyle().
		Foreground(t.Accent).
		MarginTop(1).
		MarginLeft(2)

	m.hintStyle = lipgloss.NewStyle().
		Foreground(t.Secondary).
		Italic(true).
		MarginTop(2)

	m.progressStyle = lipgloss.NewStyle().
		Foreground(t.Secondary)

	m.progressDotActive = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)

	m.progressDotDone = lipgloss.NewStyle().
		Foreground(t.Success)
}

// =============================================================================
// Persistence
// =============================================================================

// isOnboardingComplete checks if onboarding was already completed
func isOnboardingComplete() bool {
	configDir, err := os.UserConfigDir()
	if err != nil {
		logging.Errorf("Failed to get user config dir: %v", err)
		return false
	}

	markerPath := filepath.Join(configDir, "noise", onboardingFileName)
	_, err = os.Stat(markerPath)
	return err == nil
}

// markComplete creates a marker file to indicate onboarding is complete
func (m *OnboardingModel) markComplete() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		logging.Errorf("Failed to get user config dir for marking completion: %v", err)
		return
	}

	noiseDir := filepath.Join(configDir, "noise")
	if err := os.MkdirAll(noiseDir, 0755); err != nil {
		logging.Errorf("Failed to create noise config directory: %v", err)
		return
	}

	markerPath := filepath.Join(noiseDir, onboardingFileName)
	f, err := os.Create(markerPath)
	if err != nil {
		logging.Errorf("Failed to create onboarding completion marker: %v", err)
		return
	}
	defer f.Close()
}

// ResetOnboarding removes the completion marker (for testing)
func ResetOnboarding() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	markerPath := filepath.Join(configDir, "noise", onboardingFileName)
	return os.Remove(markerPath)
}
