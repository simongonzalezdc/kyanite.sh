package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AnimatedLoadingSpinner creates an enhanced loading spinner with animations
type AnimatedLoadingSpinner struct {
	spinner   spinner.Model
	animation *AnimationManager
	message   string
	width     int
	height    int
	fadeIn    bool
}

// NewAnimatedLoadingSpinner creates a new animated loading spinner
func NewAnimatedLoadingSpinner(message string) *AnimatedLoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	t := theme.GetManager().Current()
	s.Style = lipgloss.NewStyle().Foreground(t.Primary)

	return &AnimatedLoadingSpinner{
		spinner:   s,
		animation: NewAnimationManager(),
		message:   message,
		fadeIn:    true,
	}
}

// Init initializes the animated spinner
func (als *AnimatedLoadingSpinner) Init() tea.Cmd {
	als.animation.FadeTransition("spinner_fade", 1.0)
	return tea.Batch(
		als.spinner.Tick,
		als.startPulseAnimation(),
	)
}

// Update handles messages for the animated spinner
func (als *AnimatedLoadingSpinner) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case pulseStartMsg:
		// Handle pulse animation start
		cmd := als.HandlePulseStart()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case AnimationTickMsg:
		cmd := als.animation.Update()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		als.width = msg.Width
		als.height = msg.Height
	}

	var cmd tea.Cmd
	als.spinner, cmd = als.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// View renders the animated spinner
func (als *AnimatedLoadingSpinner) View() string {
	fadeProgress := als.animation.GetAnimationProgress("spinner_fade")
	pulseProgress := als.animation.GetAnimationProgress("spinner_pulse")
	t := theme.GetManager().Current()

	// Create animated spinner style
	spinnerStyle := lipgloss.NewStyle().
		Foreground(t.Primary).
		Align(lipgloss.Center).
		Width(als.width).
		Height(als.height)

	// Apply fade effect if animation is in progress
	if fadeProgress < 1.0 && fadeProgress > 0 {
		opacity := uint(fadeProgress * 255)
		if opacity > 255 {
			opacity = 255
		}
		spinnerStyle = spinnerStyle.Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", opacity, opacity/2, opacity)))
	}

	// Apply pulse effect for active state
	if pulseProgress > 0 && pulseProgress < 1 {
		// Pulsing effect - modify the spinner color slightly
		pulseOpacity := uint((0.7 + 0.3*pulseProgress) * 255)
		if pulseOpacity > 255 {
			pulseOpacity = 255
		}
		spinnerStyle = spinnerStyle.Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", pulseOpacity, pulseOpacity/2, pulseOpacity)))
	}

	content := als.spinner.View()
	if als.message != "" {
		content += " " + als.message
	}

	return spinnerStyle.Render(content)
}

// SetSize sets the dimensions for the spinner
func (als *AnimatedLoadingSpinner) SetSize(width, height int) {
	als.width = width
	als.height = height
}

// pulseStartMsg triggers the pulse animation after a delay
type pulseStartMsg struct{}

// startPulseAnimation starts a pulsing animation for the spinner using tea.Tick
func (als *AnimatedLoadingSpinner) startPulseAnimation() tea.Cmd {
	// Use tea.Tick for non-blocking delay instead of time.Sleep
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return pulseStartMsg{}
	})
}

// HandlePulseStart processes the pulse start message
func (als *AnimatedLoadingSpinner) HandlePulseStart() tea.Cmd {
	als.animation.PulseAnimation("spinner_pulse", 1.0)
	return als.animation.Update()
}

// AnimatedStatusBar creates an animated status bar for operations
type AnimatedStatusBar struct {
	message   string
	progress  float64
	animation *AnimationManager
	width     int
	status    StatusType
}

// StatusType represents different status types
type StatusType int

const (
	StatusNormal StatusType = iota
	StatusSuccess
	StatusError
	StatusWarning
	StatusLoading
)

// NewAnimatedStatusBar creates a new animated status bar
func NewAnimatedStatusBar(message string) *AnimatedStatusBar {
	return &AnimatedStatusBar{
		message:   message,
		progress:  0.0,
		animation: NewAnimationManager(),
		status:    StatusNormal,
	}
}

// SetProgress sets the progress and updates animations
func (asb *AnimatedStatusBar) SetProgress(progress float64) {
	asb.progress = progress
	if progress > 0 && progress < 1.0 {
		asb.animation.FadeTransition("progress_fade", progress)
	}
}

// SetStatus sets the status and triggers appropriate animations
func (asb *AnimatedStatusBar) SetStatus(status StatusType, message string) {
	asb.status = status
	asb.message = message

	switch status {
	case StatusSuccess:
		asb.animation.PulseAnimation("status_success", 1.0)
	case StatusError:
		asb.animation.FadeTransition("status_error", 1.0)
	case StatusWarning:
		asb.animation.SlideTransition("status_warning", 1.0)
	case StatusLoading:
		asb.animation.FadeTransition("status_loading", 1.0)
	}
}

// Update handles messages for the animated status bar
func (asb *AnimatedStatusBar) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case AnimationTickMsg:
		cmd := asb.animation.Update()
		return cmd
	}

	return nil
}

// View renders the animated status bar
func (asb *AnimatedStatusBar) View() string {
	// Base style
	var baseStyle lipgloss.Style
	t := theme.GetManager().Current()

	switch asb.status {
	case StatusSuccess:
		baseStyle = lipgloss.NewStyle().Foreground(t.Success)
	case StatusError:
		baseStyle = lipgloss.NewStyle().Foreground(t.Error)
	case StatusWarning:
		baseStyle = lipgloss.NewStyle().Foreground(t.Warning)
	case StatusLoading:
		baseStyle = lipgloss.NewStyle().Foreground(t.Primary)
	default:
		baseStyle = lipgloss.NewStyle().Foreground(t.Text)
	}

	// Apply animation effects
	successProgress := asb.animation.GetAnimationProgress("status_success")
	errorProgress := asb.animation.GetAnimationProgress("status_error")
	_ = asb.animation.GetAnimationProgress("status_warning") // For future use
	_ = asb.animation.GetAnimationProgress("status_loading") // For future use

	if successProgress > 0 && successProgress < 1 {
		// Pulsing success animation
		intensity := 0.8 + 0.2*successProgress
		green := uint(intensity * 255)
		baseStyle = baseStyle.Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", 0, green, 0)))
	}

	if errorProgress > 0 && errorProgress < 1 {
		// Fading error animation
		intensity := errorProgress
		red := uint(intensity * 255)
		baseStyle = baseStyle.Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", red, 0, 0)))
	}

	// Create progress bar if progress > 0
	if asb.progress > 0 {
		barWidth := int(float64(asb.width) * asb.progress)
		if barWidth > asb.width {
			barWidth = asb.width
		}

		progressBar := strings.Repeat("â–ˆ", barWidth)
		if barWidth < asb.width {
			progressBar += strings.Repeat("â–‘", asb.width-barWidth)
		}

		return baseStyle.Render(fmt.Sprintf("%s [%s]", asb.message, progressBar))
	}

	return baseStyle.Render(asb.message)
}

// AnimatedNotification creates animated notifications
type AnimatedNotification struct {
	message   string
	notifType string
	animation *AnimationManager
	duration  time.Duration
	startTime time.Time
	active    bool
	// Add categorization for better notification management
	category string
	priority int
}

// NewAnimatedNotification creates a new animated notification
func NewAnimatedNotification(message, notifType string, duration time.Duration) *AnimatedNotification {
	return &AnimatedNotification{
		message:   message,
		notifType: notifType,
		animation: NewAnimationManager(),
		duration:  duration,
		active:    true,
		category:  "general", // Default category
		priority:  0,         // Default priority
	}
}

// NewCategorizedNotification creates a notification with category and priority
func NewCategorizedNotification(message, notifType, category string, priority int, duration time.Duration) *AnimatedNotification {
	return &AnimatedNotification{
		message:   message,
		notifType: notifType,
		animation: NewAnimationManager(),
		duration:  duration,
		active:    true,
		category:  category,
		priority:  priority,
	}
}

// Init initializes the notification animation
func (an *AnimatedNotification) Init() tea.Cmd {
	an.startTime = time.Now()
	an.animation.SlideTransition("notification_slide", 1.0)
	return an.animation.Update()
}

// Update handles notification messages
func (an *AnimatedNotification) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case AnimationTickMsg:
		// Check if notification should be dismissed
		if time.Since(an.startTime) > an.duration {
			an.active = false
			an.animation.FadeTransition("notification_fade_out", 1.0)
		}

		cmd := an.animation.Update()
		return cmd
	}

	return nil
}

// View renders the animated notification
func (an *AnimatedNotification) View() string {
	if !an.active {
		return ""
	}

	slideProgress := an.animation.GetAnimationProgress("notification_slide")
	fadeProgress := an.animation.GetAnimationProgress("notification_fade_out")
	t := theme.GetManager().Current()

	// Apply slide-in effect
	var style lipgloss.Style
	if slideProgress < 1.0 {
		// Sliding in from top
		marginTop := int((1 - slideProgress) * 5)
		style = lipgloss.NewStyle().MarginTop(marginTop)
	} else if fadeProgress > 0 {
		// Fading out
		opacity := uint((1 - fadeProgress) * 255)
		style = lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", opacity, opacity, opacity)))
	}

	// Style based on notification type
	switch an.notifType {
	case "success":
		style = style.Background(t.Success).Foreground(t.Background).Bold(true).Padding(0, 1)
	case "error":
		style = style.Background(t.Error).Foreground(t.Background).Bold(true).Padding(0, 1)
	case "warning":
		style = style.Background(t.Warning).Foreground(t.Background).Bold(true).Padding(0, 1)
	default:
		style = style.Background(t.Primary).Foreground(t.Background).Bold(true).Padding(0, 1)
	}

	return style.Render(an.message)
}

// IsActive returns whether the notification is still active
func (an *AnimatedNotification) IsActive() bool {
	return an.active && an.animation.GetAnimationProgress("notification_fade_out") < 1.0
}

// GetCategory returns the notification category
func (an *AnimatedNotification) GetCategory() string {
	return an.category
}

// GetPriority returns the notification priority
func (an *AnimatedNotification) GetPriority() int {
	return an.priority
}

// SetCategory sets the notification category
func (an *AnimatedNotification) SetCategory(category string) {
	an.category = category
}

// SetPriority sets the notification priority
func (an *AnimatedNotification) SetPriority(priority int) {
	an.priority = priority
}
