package ui

import (
	"fmt"

	"github.com/kyanite/noise/internal/infra/sync"
	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SyncStatusModel displays the current PWA sync status in the status bar
type SyncStatusModel struct {
	syncServer *sync.SyncServer
	width      int

	// Styles
	runningStyle lipgloss.Style
	stoppedStyle lipgloss.Style
	deviceStyle  lipgloss.Style
}

// SyncStatusUpdateMsg updates the sync status
type SyncStatusUpdateMsg struct {
	Running  bool
	Devices  int
	LocalURL string
}

// NewSyncStatusModel creates a new sync status model
func NewSyncStatusModel(syncServer *sync.SyncServer) *SyncStatusModel {
	t := theme.GetManager().Current()

	return &SyncStatusModel{
		syncServer: syncServer,
		runningStyle: lipgloss.Style{}.
			Foreground(t.Success),
		stoppedStyle: lipgloss.Style{}.
			Foreground(t.Secondary),
		deviceStyle: lipgloss.Style{}.
			Foreground(t.Accent),
	}
}

// Init initializes the sync status model
func (m *SyncStatusModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *SyncStatusModel) Update(msg tea.Msg) (*SyncStatusModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

// View renders the sync status
func (m *SyncStatusModel) View() string {
	if m.syncServer == nil {
		return m.stoppedStyle.Render("⊘ Sync: Off")
	}

	if m.syncServer.IsRunning() {
		devices := m.syncServer.GetConnectedDevices()
		url := m.syncServer.GetLocalURL()

		deviceStr := ""
		if devices > 0 {
			deviceStr = m.deviceStyle.Render(fmt.Sprintf(" (%d)", devices))
		}

		return m.runningStyle.Render(fmt.Sprintf("⚡ %s", url)) + deviceStr
	}

	return m.stoppedStyle.Render("⊘ Sync: Off")
}

// IsRunning returns whether sync is active
func (m *SyncStatusModel) IsRunning() bool {
	return m.syncServer != nil && m.syncServer.IsRunning()
}

// GetLocalURL returns the local sync URL
func (m *SyncStatusModel) GetLocalURL() string {
	if m.syncServer == nil {
		return ""
	}
	return m.syncServer.GetLocalURL()
}

// GetConnectedDevices returns the number of connected devices
func (m *SyncStatusModel) GetConnectedDevices() int {
	if m.syncServer == nil {
		return 0
	}
	return m.syncServer.GetConnectedDevices()
}
