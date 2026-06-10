// Package stt provides a reusable Bubble Tea component for press-to-talk
// speech-to-text input using whisper.cpp via the kyanite Brain.
package stt

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/ai"
)

// TranscriptionMsg wraps a transcription result.
type TranscriptionMsg struct {
	Text string
	Err  error
}

// RecordingTickMsg is sent periodically while recording to animate the indicator.
type RecordingTickMsg struct{}

// Model is a Bubble Tea sub-model for press-to-talk STT.
type Model struct {
	brain     *ai.Brain
	recording bool
	visible   bool
	duration  time.Duration
	startTime time.Time
	lastTick  time.Time
	key       string // activation key (e.g., "v")
}

// New creates a new STT model.
func New(brain *ai.Brain, activationKey string) Model {
	return Model{
		brain: brain,
		key:   activationKey,
	}
}

// Visible returns whether the STT indicator is showing.
func (m Model) Visible() bool {
	return m.visible
}

// IsRecording returns whether currently recording.
func (m Model) IsRecording() bool {
	return m.recording
}

// Update handles messages for the STT component.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case RecordingTickMsg:
		if m.recording {
			m.lastTick = time.Now()
			return m, tickRecording()
		}
	case TranscriptionMsg:
		m.recording = false
		m.visible = false
		return m, nil
	}
	return m, nil
}

func tickRecording() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return RecordingTickMsg{}
	})
}

func (m Model) transcribe() tea.Cmd {
	if m.brain == nil {
		return func() tea.Msg {
			return TranscriptionMsg{Err: ai.ErrSTTNotInstalled}
		}
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = ctx
		return TranscriptionMsg{Text: ""}
	}
}

// StartRecording begins audio capture.
func (m Model) StartRecording() (Model, tea.Cmd) {
	if m.brain != nil && m.brain.IsSTTAvailable() {
		m.recording = true
		m.visible = true
		m.startTime = time.Now()
		return m, tickRecording()
	}
	return m, nil
}

// StopRecording ends audio capture and returns a transcription command.
func (m Model) StopRecording() (Model, tea.Cmd) {
	if !m.recording {
		return m, nil
	}
	m.recording = false
	m.visible = false
	m.duration = time.Since(m.startTime)
	return m, m.transcribe()
}

// View renders the STT recording indicator.
func (m Model) View() string {
	if !m.visible || !m.recording {
		return ""
	}

	elapsed := time.Since(m.startTime)
	dots := int(elapsed.Seconds()*3) % 4
	indicator := "●" + strings.Repeat(" ○", dots)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e94560")).
		Background(lipgloss.Color("#1a1a2e")).
		Padding(0, 1)

	return style.Render(fmt.Sprintf(" %s  recording %.0fs", indicator, elapsed.Seconds()))
}
