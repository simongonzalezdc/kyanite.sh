package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Kyanite/noise/internal/ui/styles"
)

// AudioTool represents different audio tools available
type AudioTool int

const (
	ToolPlayback AudioTool = iota
	ToolMetronome
	ToolRecording
)

// PlaybackState represents the state of audio playback
type PlaybackState int

const (
	StateStopped PlaybackState = iota
	StatePlaying
	StatePaused
)

// MetronomeState represents the state of the metronome
type MetronomeState int

const (
	MetronomeStopped MetronomeState = iota
	MetronomePlaying
)

// RecordingState represents the state of audio recording
type RecordingState int

const (
	RecordingStopped RecordingState = iota
	RecordingRecording
	RecordingPaused
)

// AudioModel handles the audio tools screen
type AudioModel struct {
	// Layout
	width  int
	height int

	// Current active tool
	activeTool AudioTool

	// Playback tool state
	playbackState    PlaybackState
	currentChord     string
	currentScale     string
	playbackProgress int

	// Metronome tool state
	metronomeState  MetronomeState
	metronome       stopwatch.Model
	tempo           int // BPM
	beatsPerMeasure int
	currentBeat     int
	timeSignature   string

	// Recording tool state
	recordingState    RecordingState
	recordingDuration time.Duration

	// UI state
	focusedSection int
	showHelp       bool

	// Styles
	containerStyle    lipgloss.Style
	headerStyle       lipgloss.Style
	sectionStyle      lipgloss.Style
	buttonStyle       lipgloss.Style
	activeButtonStyle lipgloss.Style
	statusStyle       lipgloss.Style
	helpStyle         lipgloss.Style
}

// NewAudioModel creates a new audio model
func NewAudioModel() *AudioModel {
	metronome := stopwatch.NewWithInterval(time.Minute) // Will be controlled manually

	return &AudioModel{
		activeTool:        ToolPlayback,
		playbackState:     StateStopped,
		metronomeState:    MetronomeStopped,
		metronome:         metronome,
		tempo:             120, // Default 120 BPM
		beatsPerMeasure:   4,
		currentBeat:       0,
		timeSignature:     "4/4",
		recordingState:    RecordingStopped,
		focusedSection:    0,
		containerStyle:    styles.Border,
		headerStyle:       styles.H1,
		sectionStyle:      styles.H2,
		buttonStyle:       styles.ButtonSecondary,
		activeButtonStyle: styles.ButtonPrimary,
		statusStyle:       styles.StatusBar,
		helpStyle:         styles.Card,
	}
}

// Init initializes the audio model
func (m *AudioModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the audio tools
func (m *AudioModel) Update(msg tea.Msg) (*AudioModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.nextTool()
		case "shift+tab":
			m.prevTool()
		case "1", "2", "3":
			// Direct tool selection
			toolNum := int(msg.String()[0] - '1')
			if toolNum >= 0 && toolNum <= 2 {
				m.activeTool = AudioTool(toolNum)
			}
		case " ":
			// Spacebar for play/pause actions
			switch m.activeTool {
			case ToolPlayback:
				m.togglePlayback()
			case ToolMetronome:
				m.toggleMetronome()
			case ToolRecording:
				m.toggleRecording()
			}
		case "s":
			// Stop actions
			switch m.activeTool {
			case ToolPlayback:
				m.stopPlayback()
			case ToolMetronome:
				m.stopMetronome()
			case ToolRecording:
				m.stopRecording()
			}
		case "up":
			if m.focusedSection > 0 {
				m.focusedSection--
			}
		case "down":
			maxSections := m.getMaxSections()
			if m.focusedSection < maxSections-1 {
				m.focusedSection++
			}
		case "+", "=":
			// Increase tempo
			if m.activeTool == ToolMetronome && m.tempo < 200 {
				m.tempo += 5
			}
		case "-":
			// Decrease tempo
			if m.activeTool == ToolMetronome && m.tempo > 60 {
				m.tempo -= 5
			}
		case "h", "?":
			m.showHelp = !m.showHelp
		case "esc":
			if m.showHelp {
				m.showHelp = false
			}
		}

		// Tool-specific key handling
		switch m.activeTool {
		case ToolPlayback:
			if cmd := m.updatePlaybackKeys(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ToolMetronome:
			if cmd := m.updateMetronomeKeys(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ToolRecording:
			if cmd := m.updateRecordingKeys(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	// Update metronome if playing
	if m.metronomeState == MetronomePlaying {
		var cmd tea.Cmd
		m.metronome, cmd = m.metronome.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the audio tools screen
func (m *AudioModel) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	var sections []string

	// Header
	header := m.headerStyle.Render("ðŸŽµ Audio Tools")
	sections = append(sections, header)

	// Tool selector
	toolSelector := m.renderToolSelector()
	sections = append(sections, toolSelector)

	// Main content based on active tool
	var content string
	switch m.activeTool {
	case ToolPlayback:
		content = m.renderPlaybackTool()
	case ToolMetronome:
		content = m.renderMetronomeTool()
	case ToolRecording:
		content = m.renderRecordingTool()
	}
	sections = append(sections, content)

	// Status bar
	status := m.renderStatusBar()
	sections = append(sections, status)

	// Combine all sections
	return m.containerStyle.Render(strings.Join(sections, "\n\n"))
}

// Tool navigation methods
func (m *AudioModel) nextTool() {
	tools := []AudioTool{ToolPlayback, ToolMetronome, ToolRecording}
	currentIndex := int(m.activeTool)
	nextIndex := (currentIndex + 1) % len(tools)
	m.activeTool = tools[nextIndex]
}

func (m *AudioModel) prevTool() {
	tools := []AudioTool{ToolPlayback, ToolMetronome, ToolRecording}
	currentIndex := int(m.activeTool)
	prevIndex := (currentIndex - 1 + len(tools)) % len(tools)
	m.activeTool = tools[prevIndex]
}

func (m *AudioModel) getMaxSections() int {
	switch m.activeTool {
	case ToolPlayback:
		return 3
	case ToolMetronome:
		return 4
	case ToolRecording:
		return 3
	default:
		return 3
	}
}

// Playback tool methods
func (m *AudioModel) togglePlayback() {
	if m.playbackState == StatePlaying {
		m.playbackState = StatePaused
	} else {
		m.playbackState = StatePlaying
	}
}

func (m *AudioModel) stopPlayback() {
	m.playbackState = StateStopped
	m.playbackProgress = 0
}

func (m *AudioModel) updatePlaybackKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "c":
		m.currentChord = m.getNextChord()
	case "n":
		// Use 'n' for Next Scale to avoid conflicting with Stop (s)
		m.currentScale = m.getNextScale()
	}
	return nil
}

func (m *AudioModel) getNextChord() string {
	chords := []string{"C", "G", "Am", "F", "Dm", "Em", "Bb", "A"}
	// Simple rotation for demo
	return chords[m.playbackProgress%len(chords)]
}

func (m *AudioModel) getNextScale() string {
	scales := []string{"C Major", "G Major", "A Minor", "E Minor", "F Major"}
	// Simple rotation for demo
	return scales[m.playbackProgress%len(scales)]
}

// Metronome tool methods
func (m *AudioModel) toggleMetronome() {
	if m.metronomeState == MetronomePlaying {
		m.metronomeState = MetronomeStopped
	} else {
		m.metronomeState = MetronomePlaying
	}
}

func (m *AudioModel) stopMetronome() {
	m.metronomeState = MetronomeStopped
	m.currentBeat = 0
}

func (m *AudioModel) updateMetronomeKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "t":
		// Change time signature
		if m.timeSignature == "4/4" {
			m.timeSignature = "3/4"
			m.beatsPerMeasure = 3
		} else {
			m.timeSignature = "4/4"
			m.beatsPerMeasure = 4
		}
	}
	return nil
}

// Recording tool methods
func (m *AudioModel) toggleRecording() {
	if m.recordingState == RecordingRecording {
		m.recordingState = RecordingPaused
	} else {
		m.recordingState = RecordingRecording
	}
}

func (m *AudioModel) stopRecording() {
	m.recordingState = RecordingStopped
	m.recordingDuration = 0
}

func (m *AudioModel) updateRecordingKeys(msg tea.KeyMsg) tea.Cmd {
	// Recording-specific key handling could be added here
	switch msg.String() {
	case "r":
		// Reset recording duration
		m.recordingDuration = 0
	}
	return nil
}

// Rendering methods
func (m *AudioModel) renderToolSelector() string {
	tools := []string{"Playback", "Metronome", "Recording"}
	var buttons []string

	for i, tool := range tools {
		var style lipgloss.Style
		if i == int(m.activeTool) {
			style = m.activeButtonStyle
		} else {
			style = m.buttonStyle
		}
		buttons = append(buttons, style.Render(fmt.Sprintf(" %s ", tool)))
	}

	return m.sectionStyle.Render("Tools: " + strings.Join(buttons, " "))
}

func (m *AudioModel) renderPlaybackTool() string {
	var sections []string

	// Chord/Scale selector
	playbackSection := m.sectionStyle.Render("Chord/Scale Playback")
	sections = append(sections, playbackSection)

	// Current chord/scale display
	chordInfo := fmt.Sprintf("Current Chord: %s | Current Scale: %s",
		m.currentChord, m.currentScale)
	sections = append(sections, chordInfo)

	// Playback controls
	controls := []string{"Play/Pause: Space", "Stop: S", "Next Chord: C", "Next Scale: N"}
	sections = append(sections, "Controls: "+strings.Join(controls, " | "))

	// Progress indicator
	progress := strings.Repeat("â–ˆ", m.playbackProgress%20) +
		strings.Repeat("â–‘", 20-(m.playbackProgress%20))
	playbackState := "Stopped"
	switch m.playbackState {
	case StatePlaying:
		playbackState = "Playing"
	case StatePaused:
		playbackState = "Paused"
	}
	sections = append(sections, fmt.Sprintf("Progress: [%s] %s", progress, playbackState))

	return strings.Join(sections, "\n")
}

func (m *AudioModel) renderMetronomeTool() string {
	var sections []string

	// Metronome settings
	metronomeSection := m.sectionStyle.Render("Metronome Settings")
	sections = append(sections, metronomeSection)

	// Tempo and time signature
	tempoInfo := fmt.Sprintf("Tempo: %d BPM | Time Signature: %s", m.tempo, m.timeSignature)
	sections = append(sections, tempoInfo)

	// Beat indicator
	beatDisplay := strings.Repeat("â—‹", m.currentBeat) +
		"â—" +
		strings.Repeat("â—‹", m.beatsPerMeasure-m.currentBeat-1)
	sections = append(sections, fmt.Sprintf("Beat: %s", beatDisplay))

	// Metronome controls
	controls := []string{"Play/Pause: Space", "Stop: S", "Tempo +/-: +/=", "Time Sig: T"}
	sections = append(sections, "Controls: "+strings.Join(controls, " | "))

	// Metronome state
	state := "Stopped"
	if m.metronomeState == MetronomePlaying {
		state = "Playing"
	}
	sections = append(sections, fmt.Sprintf("Status: %s", state))

	return strings.Join(sections, "\n")
}

func (m *AudioModel) renderRecordingTool() string {
	var sections []string

	// Recording section
	recordingSection := m.sectionStyle.Render("Audio Recording")
	sections = append(sections, recordingSection)

	// Recording duration
	duration := m.recordingDuration.Round(time.Second).String()
	sections = append(sections, fmt.Sprintf("Duration: %s", duration))

	// Recording controls
	controls := []string{"Record/Pause: Space", "Stop: S"}
	sections = append(sections, "Controls: "+strings.Join(controls, " | "))

	// Recording state
	state := "Stopped"
	switch m.recordingState {
	case RecordingRecording:
		state = "Recording"
	case RecordingPaused:
		state = "Paused"
	}
	sections = append(sections, fmt.Sprintf("Status: %s", state))

	return strings.Join(sections, "\n")
}

func (m *AudioModel) renderStatusBar() string {
	var status []string

	// Current tool
	toolNames := []string{"Playback", "Metronome", "Recording"}
	status = append(status, fmt.Sprintf("Tool: %s", toolNames[m.activeTool]))

	// Navigation hints
	status = append(status, "Tab: Next Tool | 1-3: Select Tool | H: Help")

	return m.statusStyle.Render(strings.Join(status, " | "))
}

func (m *AudioModel) renderHelp() string {
	help := m.helpStyle.Render(`Audio Tools Help

Navigation:
  Tab / Shift+Tab    Switch between tools
  1, 2, 3           Jump to specific tool
  â†‘ / â†“             Navigate within tool
  H / ?             Toggle this help
  Esc               Exit help

Playback Tool:
  Space             Play/Pause chord/scale
  S                 Stop playback
  C                 Next chord
  N                 Next scale

Metronome Tool:
  Space             Start/Stop metronome
  S                 Stop metronome
  + / =             Increase tempo
  -                 Decrease tempo
  T                 Toggle time signature

Recording Tool:
  Space             Start/Pause recording
  S                 Stop recording

General:
  All tools support Space for primary action
  and S to stop the current operation`)

	return help
}
