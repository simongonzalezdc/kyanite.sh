package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// MenuModel handles the main menu screen
type MenuModel struct {
	list   list.Model
	width  int
	height int
}

// NewMenuModel creates a new menu model
func NewMenuModel() *MenuModel {
	items := []list.Item{
		item{title: "New Song", desc: "Create a new song", screen: screenEditor},
		item{title: "Open Song", desc: "Open an existing song", screen: screenEditor},
		item{title: "Theory Tools", desc: "Music theory and rhyme tools", screen: screenTheory},
		item{title: "Audio Tools", desc: "Metronome and chord playback", screen: screenAudio},
		item{title: "Project Manager", desc: "Manage songs and projects", screen: screenManager},
		item{title: "Settings", desc: "Application settings", screen: screenSettings},
		item{title: "Exit", desc: "Exit LyricForge", screen: screenSplash},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "🎵 LyricForge"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return &MenuModel{
		list: l,
	}
}

// Init initializes the menu model
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the menu
func (m *MenuModel) Update(msg tea.Msg) (*MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				// Here you would navigate to the selected screen
				// For now, just return the model
				_ = selectedItem // Prevent unused variable warning
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

// View renders the menu
func (m *MenuModel) View() string {
	return m.list.View()
}

// item represents a menu item
type item struct {
	title  string
	desc   string
	screen screen
}

// FilterValue returns the filter value for the item
func (i item) FilterValue() string {
	return i.title
}

// Title returns the title of the item
func (i item) Title() string {
	return i.title
}

// Description returns the description of the item
func (i item) Description() string {
	return i.desc
}
