package ui

import (
	"fmt"
	"strings"

	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/infra/sync"
	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SyncSettingsModel manages PWA sync settings
type SyncSettingsModel struct {
	syncServer *sync.SyncServer
	config     *config.SyncConfig
	width      int
	height     int

	// Selection state
	selectedItem int
	items        []syncSettingItem
	pairingCode  string

	// Styles
	titleStyle       lipgloss.Style
	itemStyle        lipgloss.Style
	selectedStyle    lipgloss.Style
	descriptionStyle lipgloss.Style
	codeStyle        lipgloss.Style
	urlStyle         lipgloss.Style
	errorStyle       lipgloss.Style
	successStyle     lipgloss.Style
}

type syncSettingItem struct {
	label       string
	description string
	itemType    syncItemType
}

type syncItemType int

const (
	syncItemToggle syncItemType = iota
	syncItemPairing
	syncItemDevices
	syncItemInbox
	syncItemBack
)

// NewSyncSettingsModel creates a new sync settings model
func NewSyncSettingsModel(syncServer *sync.SyncServer, cfg *config.SyncConfig) *SyncSettingsModel {
	t := theme.GetManager().Current()

	m := &SyncSettingsModel{
		syncServer:   syncServer,
		config:       cfg,
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
		codeStyle: lipgloss.Style{}.
			Foreground(t.Accent).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary).
			Padding(0, 2),
		urlStyle: lipgloss.Style{}.
			Foreground(t.Primary),
		errorStyle: lipgloss.Style{}.
			Foreground(t.Error),
		successStyle: lipgloss.Style{}.
			Foreground(t.Success),
	}

	m.buildItems()
	return m
}

func (m *SyncSettingsModel) buildItems() {
	m.items = []syncSettingItem{
		{
			label:       "Toggle Sync Server",
			description: "Start or stop the local sync server",
			itemType:    syncItemToggle,
		},
		{
			label:       "Generate Pairing Code",
			description: "Create a code for connecting PWA devices",
			itemType:    syncItemPairing,
		},
		{
			label:       "Paired Devices",
			description: "View and manage connected devices",
			itemType:    syncItemDevices,
		},
		{
			label:       "Idea Inbox",
			description: "View captured ideas from companion app",
			itemType:    syncItemInbox,
		},
		{
			label:       "<- Back",
			description: "",
			itemType:    syncItemBack,
		},
	}
}

// Init initializes the sync settings model
func (m *SyncSettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *SyncSettingsModel) Update(msg tea.Msg) (*SyncSettingsModel, tea.Cmd) {
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

func (m *SyncSettingsModel) handleSelection() tea.Cmd {
	if m.selectedItem >= len(m.items) {
		return nil
	}

	item := m.items[m.selectedItem]

	switch item.itemType {
	case syncItemToggle:
		if m.syncServer == nil {
			return nil
		}
		if m.syncServer.IsRunning() {
			_ = m.syncServer.Stop()
		} else {
			_ = m.syncServer.Start()
		}
		return nil

	case syncItemPairing:
		if m.syncServer != nil && m.syncServer.IsRunning() {
			m.pairingCode = m.syncServer.GeneratePairingCode()
		}
		return nil

	case syncItemDevices:
		// TODO: Show paired devices screen
		return nil

	case syncItemInbox:
		return func() tea.Msg { return ShowIdeaInboxMsg{} }

	case syncItemBack:
		return func() tea.Msg { return BackToSettingsMsg{} }
	}

	return nil
}

// View renders the sync settings screen
func (m *SyncSettingsModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(m.titleStyle.Render("[SYNC] PWA Sync Settings"))
	b.WriteString("\n\n")

	// Status
	if m.syncServer != nil && m.syncServer.IsRunning() {
		url := m.syncServer.GetLocalURL()
		devices := m.syncServer.GetConnectedDevices()
		b.WriteString(m.successStyle.Render("[OK] Sync server running"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("URL: %s\n", m.urlStyle.Render(url)))
		b.WriteString(fmt.Sprintf("Connected devices: %d\n", devices))
	} else {
		b.WriteString(m.errorStyle.Render("[X] Sync server stopped"))
	}
	b.WriteString("\n")

	// Pairing code if generated
	if m.pairingCode != "" {
		b.WriteString("Pairing Code:\n")
		b.WriteString(m.codeStyle.Render(m.pairingCode))
		b.WriteString("\n")
		b.WriteString(m.descriptionStyle.Render("Enter this code in the PWA to connect"))
		b.WriteString("\n\n")
	}

	// Menu items
	for i, item := range m.items {
		var style lipgloss.Style
		if i == m.selectedItem {
			style = m.selectedStyle
		} else {
			style = m.itemStyle
		}

		label := item.label

		// Update toggle label based on state
		if item.itemType == syncItemToggle {
			if m.syncServer != nil && m.syncServer.IsRunning() {
				label = "[*] Stop Sync Server"
			} else {
				label = "[ ] Start Sync Server"
			}
		}

		b.WriteString(style.Render(label))
		b.WriteString("\n")

		if item.description != "" {
			b.WriteString(m.descriptionStyle.Render(item.description))
			b.WriteString("\n")
		}
	}

	// Help
	b.WriteString("\n")
	b.WriteString(m.descriptionStyle.Render("Up/Down: Navigate - Enter: Select - Esc: Back"))

	return b.String()
}

// ShowIdeaInboxMsg requests showing the idea inbox
type ShowIdeaInboxMsg struct{}
