package ui

import (
	"fmt"
	"strings"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/infra/voice"
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VoiceSettingsModel manages voice-to-text settings
type VoiceSettingsModel struct {
	voiceService *app.VoiceService
	width        int
	height       int

	// Selection state
	selectedItem     int
	items            []voiceSettingItem
	downloadingModel string
	downloadProgress float64
	testRecording    bool

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
	voiceItemModelSelect voiceItemType = iota
	voiceItemDownload
	voiceItemTestMic
	voiceItemBack
)

// ModelDownloadProgressMsg reports model download progress
type ModelDownloadProgressMsg struct {
	ModelName  string
	Downloaded int64
	Total      int64
	Complete   bool
	Error      error
}

// NewVoiceSettingsModel creates a new voice settings model
func NewVoiceSettingsModel(voiceService *app.VoiceService) *VoiceSettingsModel {
	t := theme.GetManager().Current()

	m := &VoiceSettingsModel{
		voiceService: voiceService,
		selectedItem: 0,
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			MarginBottom(1),
		itemStyle: lipgloss.NewStyle().
			Foreground(t.Text).
			PaddingLeft(2),
		selectedStyle: lipgloss.NewStyle().
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true).
			PaddingLeft(2),
		descriptionStyle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			PaddingLeft(4),
		progressStyle: lipgloss.NewStyle().
			Foreground(t.Accent),
		errorStyle: lipgloss.NewStyle().
			Foreground(t.Error),
		successStyle: lipgloss.NewStyle().
			Foreground(t.Success),
	}

	m.buildItems()
	return m
}

func (m *VoiceSettingsModel) buildItems() {
	m.items = []voiceSettingItem{
		{
			label:       "Model: Base English (~142MB)",
			description: "Recommended balance of speed and accuracy",
			itemType:    voiceItemModelSelect,
			value:       voice.ModelBaseEN,
		},
		{
			label:       "Model: Tiny English (~75MB)",
			description: "Fastest, lower accuracy",
			itemType:    voiceItemModelSelect,
			value:       voice.ModelTinyEN,
		},
		{
			label:       "Model: Small English (~466MB)",
			description: "Best accuracy, slower",
			itemType:    voiceItemModelSelect,
			value:       voice.ModelSmallEN,
		},
		{
			label:       "Download Selected Model",
			description: "Download the whisper model if not available",
			itemType:    voiceItemDownload,
		},
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

	case ModelDownloadProgressMsg:
		if msg.Error != nil {
			m.downloadingModel = ""
			// Could show error notification
		} else if msg.Complete {
			m.downloadingModel = ""
			m.downloadProgress = 0
		} else if msg.Total > 0 {
			m.downloadProgress = float64(msg.Downloaded) / float64(msg.Total)
		} else {
			m.downloadProgress = 0
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
	case voiceItemModelSelect:
		// Select this model
		if m.voiceService != nil {
			cfg := m.voiceService.GetConfig()
			if cfg != nil {
				cfg.Model = item.value
			}
		}
		return nil

	case voiceItemDownload:
		if m.voiceService == nil {
			return nil
		}
		cfg := m.voiceService.GetConfig()
		if cfg == nil {
			return nil
		}
		modelName := cfg.Model
		if modelName == "" {
			modelName = voice.ModelBaseEN
		}
		m.downloadingModel = modelName
		m.downloadProgress = 0

		// Start async download
		return func() tea.Msg {
			manager := m.voiceService.GetModelManager()
			if manager == nil {
				return ModelDownloadProgressMsg{Error: fmt.Errorf("no model manager")}
			}

			_, err := manager.EnsureModel(modelName, func(downloaded, total int64) {
				// Progress callback - we can't easily send tea.Msg from here
				// so progress updates are limited
			})
			if err != nil {
				return ModelDownloadProgressMsg{ModelName: modelName, Error: err}
			}
			return ModelDownloadProgressMsg{ModelName: modelName, Complete: true}
		}

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
		b.WriteString(m.successStyle.Render("[OK] Voice service ready"))
	} else if m.voiceService != nil {
		b.WriteString(m.errorStyle.Render("[...] Voice service initializing..."))
	} else {
		b.WriteString(m.errorStyle.Render("[X] Voice service not available"))
	}
	b.WriteString("\n\n")

	// Current model
	if m.voiceService != nil {
		cfg := m.voiceService.GetConfig()
		if cfg != nil {
			b.WriteString(fmt.Sprintf("Current model: %s\n", cfg.Model))
			b.WriteString(fmt.Sprintf("Push-to-talk: %s\n", cfg.PushToTalkKey))
		}
	}
	b.WriteString("\n")

	// Download progress
	if m.downloadingModel != "" {
		progress := int(m.downloadProgress * 40)
		bar := strings.Repeat("#", progress) + strings.Repeat("-", 40-progress)
		b.WriteString(m.progressStyle.Render(fmt.Sprintf("Downloading %s: [%s] %.0f%%\n\n",
			m.downloadingModel, bar, m.downloadProgress*100)))
	}

	// Menu items
	for i, item := range m.items {
		var style lipgloss.Style
		if i == m.selectedItem {
			style = m.selectedStyle
		} else {
			style = m.itemStyle
		}

		// Mark current model
		label := item.label
		if item.itemType == voiceItemModelSelect && m.voiceService != nil {
			cfg := m.voiceService.GetConfig()
			if cfg != nil && cfg.Model == item.value {
				label = "[*] " + label
			}
		}

		// Mark downloaded models
		if item.itemType == voiceItemModelSelect && m.voiceService != nil {
			manager := m.voiceService.GetModelManager()
			if manager != nil && manager.IsModelAvailable(item.value) {
				label += " [downloaded]"
			}
		}

		b.WriteString(style.Render(label))
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
