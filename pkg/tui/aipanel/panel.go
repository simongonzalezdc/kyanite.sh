// Package aipanel provides a reusable Bubble Tea component for displaying
// streaming AI responses in a side panel or overlay.
package aipanel

import (
	"context"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/ai"
)

// StreamChunk is a tea.Msg wrapper for AI stream chunks.
type StreamChunk struct {
	Content string
	Done    bool
	Err     error
}

// ErrorMsg wraps an error for display in the panel.
type ErrorMsg struct {
	Err error
}

// Model is a Bubble Tea sub-model that displays streaming AI responses.
type Model struct {
	width      int
	height     int
	title      string
	content    strings.Builder
	streaming  bool
	visible    bool
	brain      *ai.Brain
	lastUpdate time.Time
	style      panelStyle
}

type panelStyle struct {
	border  lipgloss.Color
	title   lipgloss.Color
	text    lipgloss.Color
	muted   lipgloss.Color
	bg      lipgloss.Color
	accent  lipgloss.Color
}

// New creates a new AI panel model.
func New(brain *ai.Brain, width, height int) Model {
	return Model{
		width:  width,
		height: height,
		brain:  brain,
		style: panelStyle{
			border: "#4a4a6a",
			title:  "#f4a261",
			text:   "#e0e0e0",
			muted:  "#6c6c8a",
			bg:     "#1a1a2e",
			accent: "#e94560",
		},
	}
}

// SetSize updates the panel dimensions.
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// Visible returns whether the panel is showing.
func (m Model) Visible() bool {
	return m.visible
}

// Toggle shows or hides the panel.
func (m Model) Toggle() Model {
	m.visible = !m.visible
	return m
}

// IsStreaming returns whether the panel is currently receiving a stream.
func (m Model) IsStreaming() bool {
	return m.streaming
}

// Generate sends a prompt to the brain and streams the response.
func (m Model) Generate(prompt string) tea.Cmd {
	if m.brain == nil {
		return func() tea.Msg {
			return ErrorMsg{Err: ai.ErrBrainNotInitialized}
		}
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ch, err := m.brain.StreamChat(ctx, []ai.Message{
			{Role: "user", Content: prompt},
		})
		if err != nil {
			return ErrorMsg{Err: err}
		}

		// Collect first chunk synchronously
		select {
		case chunk, ok := <-ch:
			if !ok {
				return StreamChunk{Done: true}
			}
			if chunk.Error != nil {
				return ErrorMsg{Err: chunk.Error}
			}
			return StreamChunk{Content: chunk.Content, Done: chunk.Done}
		case <-ctx.Done():
			return ErrorMsg{Err: ctx.Err()}
		}
	}
}

// ContinueStream reads the next chunk from an active stream.
// This should be called after each StreamChunk that has Done=false.
func (m Model) ContinueStream(ch <-chan ai.StreamChunk) tea.Cmd {
	return func() tea.Msg {
		select {
		case chunk, ok := <-ch:
			if !ok {
				return StreamChunk{Done: true}
			}
			if chunk.Error != nil {
				return ErrorMsg{Err: chunk.Error}
			}
			return StreamChunk{Content: chunk.Content, Done: chunk.Done}
		case <-time.After(30 * time.Second):
			return StreamChunk{Done: true}
		}
	}
}

// Update handles messages for the AI panel.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StreamChunk:
		m.content.WriteString(msg.Content)
		m.lastUpdate = time.Now()
		if msg.Done {
			m.streaming = false
		}
		return m, nil
	case ErrorMsg:
		m.content.WriteString("\n⚠ " + msg.Err.Error())
		m.streaming = false
		return m, nil
	}
	return m, nil
}

// View renders the AI panel.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	panelW := max(m.width/3, 40)
	panelH := max(m.height-4, 10)

	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.style.border)).
		Background(lipgloss.Color(m.style.bg)).
		Foreground(lipgloss.Color(m.style.text)).
		Width(panelW - 2).
		Height(panelH - 2).
		Padding(0, 1)

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.style.title)).
		Bold(true).
		Render("◆ ai")

	streaming := ""
	if m.streaming {
		streaming = lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.style.accent)).
			Render(" ● streaming")
	}

	header := title + streaming

	content := m.content.String()
	if content == "" {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.style.muted)).
			Render("Press Enter or type a prompt to ask the AI...")
	}

	// Truncate content to fit
	lines := strings.Split(content, "\n")
	maxLines := panelH - 4
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	content = strings.Join(lines, "\n")

	body := frame.Render(header + "\n" + strings.Repeat("─", panelW-6) + "\n" + content)
	return body
}

// SetTitle sets the panel header title.
func (m Model) SetTitle(title string) Model {
	m.title = title
	return m
}

// Clear clears the panel content.
func (m Model) Clear() Model {
	m.content.Reset()
	m.streaming = false
	return m
}

// StartStream marks the panel as streaming and clears previous content.
func (m Model) StartStream(title string) Model {
	m.content.Reset()
	m.streaming = true
	m.title = title
	return m
}

// Content returns the current panel content text.
func (m Model) Content() string {
	return m.content.String()
}
