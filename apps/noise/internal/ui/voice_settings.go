package ui

import (
	"fmt"
	"strings"

	"github.com/kyanite/noise/internal/app"
	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VoiceSettingsModel manages voice-to-text settings
type VoiceSettingsModel struct {
	voiceService *app.VoiceService
	width        int
	height       int

	// Selection state
	selectedItem  int
	items         []voiceSettingItem
	testRecording bool

	// Styles
	titleStyle       lipgloss.Style
	itemStyle        lipgloss.Style
	selectedStyle    lipgloss.Style
	descriptionStyle lipgloss.Style
	progressStyle    lipgloss.Style
	errorStyle       lipgloss.Style
	successStyle     lipgloss.Style
}

type voiceSettingItem struct {
	label       string
	description string
	itemType    voiceItemType
	value       string
}

type voiceItemType int

const (
	voiceItemTestMic voiceItemType = iota
	voiceItemBack
)

// NewVoiceSettingsModel creates a new voice settings model
func NewVoiceSettingsModel(voiceService *app.VoiceService) *VoiceSettingsModel {
	t := theme.GetManager().Current()

	m := &VoiceSettingsModel{
		voiceService: voiceService,
		selectedItem: 0,
		titleStyle: lipgloss.Style{}.
			Bold(true).
			Foreground(t.Primary).
			MarginBottom(1),
		itemStyle: lipgloss.Style{}.
			Foreground(t.Text).
			PaddingLeft(2),
		selectedStyle: lipgloss.Style{}.
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true).
			PaddingLeft(2),
		descriptionStyle: lipgloss.Style{}.
			Foreground(t.Secondary).
			PaddingLeft(4),
		progressStyle: lipgloss.Style{}.
			Foreground(t.Accent),
		errorStyle: lipgloss.Style{}.
			Foreground(t.Error),
		successStyle: lipgloss.Style{}.
			Foreground(t.Success),
	}

	m.buildItems()
	return m
}

func (m *VoiceSettingsModel) buildItems() {
	m.items = []voiceSettingItem{
		{
			label:       "Test Microphone",
			description: "Record a short test to verify audio input",
			itemType:    voiceItemTestMic,
		},
		{
			label:       "<- Back",
			description: "",
			itemType:    voiceItemBack,
		},
	}
}

// Init initializes the voice settings model
func (m *VoiceSettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for voice settings
func (m *VoiceSettingsModel) Update(msg tea.Msg) (*VoiceSettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedItem > 0 {
				m.selectedItem--
			}
		case "down", "j":
			if m.selectedItem < len(m.items)-1 {
				m.selectedItem++
			}
		case "enter", " ":
			return m, m.handleSelection()
		case "esc":
			return m, func() tea.Msg { return BackToSettingsMsg{} }
		}
	}

	return m, nil
}

func (m *VoiceSettingsModel) handleSelection() tea.Cmd {
	if m.selectedItem >= len(m.items) {
		return nil
	}

	item := m.items[m.selectedItem]

	switch item.itemType {
	case voiceItemTestMic:
		// Test microphone - quick record and playback
		m.testRecording = !m.testRecording
		if m.testRecording && m.voiceService != nil {
			return func() tea.Msg {
				if err := m.voiceService.StartDictation(); err != nil {
					return VoiceTranscriptionMsg{Error: err}
				}
				return nil
			}
		} else if !m.testRecording && m.voiceService != nil {
			return func() tea.Msg {
				text, err := m.voiceService.StopDictation()
				return VoiceTranscriptionMsg{Text: text, Error: err}
			}
		}
		return nil

	case voiceItemBack:
		return func() tea.Msg { return BackToSettingsMsg{} }
	}

	return nil
}

// View renders the voice settings screen
func (m *VoiceSettingsModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(m.titleStyle.Render("[MIC] Voice-to-Text Settings"))
	b.WriteString("\n\n")

	// Status
	if m.voiceService != nil && m.voiceService.IsAvailable() {
		b.WriteString(m.successStyle.Render("[OK] Voice service ready (brain STT)"))
	} else if m.voiceService != nil {
		b.WriteString(m.errorStyle.Render("[...] Voice service initializing..."))
	} else {
		b.WriteString(m.errorStyle.Render("[X] Voice service not available"))
	}
	b.WriteString("\n\n")

	// Current config
	if m.voiceService != nil {
		cfg := m.voiceService.GetConfig()
		if cfg != nil {
			b.WriteString(fmt.Sprintf("Push-to-talk: %s\n", cfg.PushToTalkKey))
		}
	}
	b.WriteString("\n")

	// Menu items
	for i, item := range m.items {
		var style lipgloss.Style
		if i == m.selectedItem {
			style = m.selectedStyle
		} else {
			style = m.itemStyle
		}

		b.WriteString(style.Render(item.label))
		b.WriteString("\n")

		if item.description != "" {
			b.WriteString(m.descriptionStyle.Render(item.description))
			b.WriteString("\n")
		}
	}

	// Test recording indicator
	if m.testRecording {
		b.WriteString("\n")
		b.WriteString(m.errorStyle.Render("[REC] Recording... Press Enter to stop"))
	}

	// Help
	b.WriteString("\n\n")
	b.WriteString(m.descriptionStyle.Render("Up/Down: Navigate - Enter: Select - Esc: Back"))

	return b.String()
}

// BackToSettingsMsg requests returning to the main settings screen
type BackToSettingsMsg struct{}
