package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/noise/internal/domain"
	"github.com/puente-labs/noise/internal/infra/db"
	"github.com/puente-labs/noise/internal/ui/styles"
)

// ManagerState represents different states of the project manager
type ManagerState int

const (
	StateProjectList ManagerState = iota
	StateProjectDetail
	StateCreateProject
	StateEditProject
	StateSongList
	StateProjectTemplates
	StateLoading
)

// ProjectManagerModel handles the comprehensive project management system
type ProjectManagerModel struct {
	// Core dependencies
	database *db.DB
	width    int
	height   int

	// State management
	state ManagerState

	// Project data
	projects       []*domain.Project
	currentProject *domain.Project
	songs          []*domain.Song

	// UI components
	projectList list.Model
	songList    list.Model
	spinner     spinner.Model

	// Input components
	projectNameInput textinput.Model
	projectDescInput textinput.Model
	searchInput      textinput.Model

	// Navigation and selection
	focusedComponent     string

	// Status and feedback
	statusMessage string
	errorMessage  string
	loading       bool

	// Animation system
	animation *AnimationManager

	// Project templates
	templates []ProjectTemplate

	// Recent projects for quick access
	recentProjects []*domain.Project
}

// ProjectTemplate represents a project template
type ProjectTemplate struct {
	Name        string
	Description string
	Icon        string
	Category    string
	SongCount   int
	Tags        []string
}

// NewProjectManagerModel creates a new comprehensive project manager model
func NewProjectManagerModel(database *db.DB) *ProjectManagerModel {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.Primary)

	// Initialize project list
	projectList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	projectList.Title = "Projects"
	projectList.SetShowStatusBar(false)
	projectList.SetFilteringEnabled(false)
	projectList.SetShowHelp(false)

	// Initialize song list
	songList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	songList.Title = "Songs"
	songList.SetShowStatusBar(false)
	songList.SetFilteringEnabled(false)
	songList.SetShowHelp(false)

	// Initialize input components
	projectNameInput := textinput.New()
	projectNameInput.Placeholder = "Project name"
	projectNameInput.Focus()
	projectNameInput.CharLimit = 100
	projectNameInput.Width = 50

	projectDescInput := textinput.New()
	projectDescInput.Placeholder = "Project description"
	projectDescInput.Width = 50

	searchInput := textinput.New()
	searchInput.Placeholder = "Search projects..."
	searchInput.Width = 30

	// Define project templates
	templates := []ProjectTemplate{
		{
			Name:        "Songwriter's Collection",
			Description: "A comprehensive collection for songwriters",
			Icon:        "🎵",
			Category:    "Music",
			SongCount:   0,
			Tags:        []string{"songwriting", "lyrics", "music"},
		},
		{
			Name:        "Album Project",
			Description: "Organize songs for an album release",
			Icon:        "💿",
			Category:    "Music",
			SongCount:   0,
			Tags:        []string{"album", "release", "collection"},
		},
		{
			Name:        "Poetry Collection",
			Description: "Collection of poetic works and lyrics",
			Icon:        "📝",
			Category:    "Literature",
			SongCount:   0,
			Tags:        []string{"poetry", "literature", "creative-writing"},
		},
		{
			Name:        "Collaborative Project",
			Description: "Project for collaborative songwriting",
			Icon:        "👥",
			Category:    "Collaboration",
			SongCount:   0,
			Tags:        []string{"collaboration", "team", "shared"},
		},
		{
			Name:        "Genre-Specific",
			Description: "Organize songs by specific genre",
			Icon:        "🎸",
			Category:    "Music",
			SongCount:   0,
			Tags:        []string{"genre", "style", "categorization"},
		},
	}

	return &ProjectManagerModel{
		database:         database,
		state:            StateLoading,
		spinner:          s,
		projectList:      projectList,
		songList:         songList,
		projectNameInput: projectNameInput,
		projectDescInput: projectDescInput,
		searchInput:      searchInput,
		focusedComponent: "projects",
		animation:        NewAnimationManager(),
		templates:        templates,
		loading:          true,
	}
}

// Init initializes the project manager model
func (m *ProjectManagerModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadProjects(),
		m.spinner.Tick,
	)
}

// loadProjects loads all projects from the database
func (m *ProjectManagerModel) loadProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := m.database.ListProjects()
		if err != nil {
			return projectLoadErrorMsg{err: err}
		}

		// Load recent projects (last 5)
		var recentProjects []*domain.Project
		if len(projects) > 5 {
			recentProjects = projects[:5]
		} else {
			recentProjects = projects
		}

		return projectLoadSuccessMsg{
			projects:       projects,
			recentProjects: recentProjects,
		}
	}
}

// loadProjectSongs loads songs for a specific project
func (m *ProjectManagerModel) loadProjectSongs(projectID int) tea.Cmd {
	return func() tea.Msg {
		project, err := m.database.GetProject(projectID)
		if err != nil {
			return songLoadErrorMsg{err: err}
		}

		var songs []*domain.Song
		for _, songID := range project.SongIDs {
			song, err := m.database.GetSong(songID)
			if err == nil {
				songs = append(songs, song)
			}
		}

		return songLoadSuccessMsg{
			project: project,
			songs:   songs,
		}
	}
}

// createProject creates a new project
func (m *ProjectManagerModel) createProject(name, description string) tea.Cmd {
	return func() tea.Msg {
		project := &domain.Project{
			Name:        name,
			Description: description,
			SongIDs:     []int{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		createdProject, err := m.database.CreateProject(project)
		if err != nil {
			return projectCreateErrorMsg{err: err}
		}

		return projectCreateSuccessMsg{project: createdProject}
	}
}

// Message types for project manager
type projectLoadSuccessMsg struct {
	projects       []*domain.Project
	recentProjects []*domain.Project
}

type projectLoadErrorMsg struct {
	err error
}

type songLoadSuccessMsg struct {
	project *domain.Project
	songs   []*domain.Song
}

type songLoadErrorMsg struct {
	err error
}

type projectCreateSuccessMsg struct {
	project *domain.Project
}

type projectCreateErrorMsg struct {
	err error
}

// OpenSongMsg is sent when a song should be opened in the main editor.
type OpenSongMsg struct {
	Song *domain.Song
}

// Update handles messages for the project manager
func (m *ProjectManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateListDimensions()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			switch m.state {
			case StateProjectDetail, StateSongList:
				m.state = StateProjectList
				m.focusedComponent = "projects"
				return m, m.loadProjects()
			case StateCreateProject, StateEditProject:
				m.state = StateProjectList
				m.focusedComponent = "projects"
			case StateProjectTemplates:
				m.state = StateProjectList
			}

		case "enter":
			switch m.state {
			case StateProjectList:
				if len(m.projects) > 0 {
					selected := m.projectList.SelectedItem()
					if selected != nil {
						m.state = StateProjectDetail
						m.currentProject = m.projects[m.projectList.Index()]
						return m, m.loadProjectSongs(m.currentProject.ID)
					}
				}

			case StateProjectDetail:
				if m.focusedComponent == "songs" && len(m.songs) > 0 {
					selected := m.songList.SelectedItem()
					if selected != nil {
						// Open selected song in editor by sending OpenSongMsg
						idx := m.songList.Index()
						if idx >= 0 && idx < len(m.songs) {
							song := m.songs[idx]
							return m, func() tea.Msg {
								return OpenSongMsg{Song: song}
							}
						}
					}
				}

			case StateCreateProject:
				if m.projectNameInput.Value() != "" {
					return m, m.createProject(m.projectNameInput.Value(), m.projectDescInput.Value())
				}

			case StateProjectTemplates:
				selected := m.projectList.SelectedItem()
				if selected != nil {
					template := m.getTemplateByName(selected.(projectTemplateItem).name)
					if template != nil {
						m.state = StateCreateProject
						m.projectNameInput.SetValue(template.Name)
						m.projectDescInput.SetValue(template.Description)
						m.focusedComponent = "name_input"
					}
				}
			}

		case "n", "ctrl+n":
			if m.state == StateProjectList {
				m.state = StateCreateProject
				m.projectNameInput.Focus()
				m.focusedComponent = "name_input"
			}

		case "t", "ctrl+t":
			if m.state == StateProjectList {
				m.state = StateProjectTemplates
				m.setupTemplateList()
			}

		case "tab":
			m.cycleFocus()

		case "shift+tab":
			m.cycleFocusReverse()

		case "up", "k":
			switch m.state {
			case StateProjectList:
				m.projectList, _ = m.projectList.Update(msg)
			case StateSongList:
				m.songList, _ = m.songList.Update(msg)
			case StateProjectTemplates:
				m.projectList, _ = m.projectList.Update(msg)
			}

		case "down", "j":
			switch m.state {
			case StateProjectList:
				m.projectList, _ = m.projectList.Update(msg)
			case StateSongList:
				m.songList, _ = m.songList.Update(msg)
			case StateProjectTemplates:
				m.projectList, _ = m.projectList.Update(msg)
			}

		case "backspace":
			if m.state == StateCreateProject || m.state == StateEditProject {
				switch m.focusedComponent {
				case "name_input":
					m.projectNameInput.SetValue("")
				case "desc_input":
					m.projectDescInput.SetValue("")
				}
			}
		}

	case projectLoadSuccessMsg:
		m.projects = msg.projects
		m.recentProjects = msg.recentProjects
		m.loading = false
		m.state = StateProjectList
		m.setupProjectList()
		m.statusMessage = fmt.Sprintf("Loaded %d projects", len(m.projects))

	case projectLoadErrorMsg:
		m.errorMessage = fmt.Sprintf("Failed to load projects: %v", msg.err)
		m.loading = false

	case songLoadSuccessMsg:
		m.currentProject = msg.project
		m.songs = msg.songs
		m.state = StateSongList
		m.setupSongList()
		m.statusMessage = fmt.Sprintf("Loaded %d songs", len(m.songs))

	case songLoadErrorMsg:
		m.errorMessage = fmt.Sprintf("Failed to load songs: %v", msg.err)

	case projectCreateSuccessMsg:
		m.state = StateProjectList
		m.focusedComponent = "projects"
		m.statusMessage = fmt.Sprintf("Created project: %s", msg.project.Name)
		return m, m.loadProjects()

	case projectCreateErrorMsg:
		m.errorMessage = fmt.Sprintf("Failed to create project: %v", msg.err)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update input components
	var nameCmd tea.Cmd
	var descCmd tea.Cmd
	var searchCmd tea.Cmd

	m.projectNameInput, nameCmd = m.projectNameInput.Update(msg)
	m.projectDescInput, descCmd = m.projectDescInput.Update(msg)
	m.searchInput, searchCmd = m.searchInput.Update(msg)

	if nameCmd != nil {
		cmds = append(cmds, nameCmd)
	}
	if descCmd != nil {
		cmds = append(cmds, descCmd)
	}
	if searchCmd != nil {
		cmds = append(cmds, searchCmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the project manager screen
func (m *ProjectManagerModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	switch m.state {
	case StateProjectList:
		return m.renderProjectList()
	case StateProjectDetail:
		return m.renderProjectDetail()
	case StateCreateProject:
		return m.renderCreateProject()
	case StateEditProject:
		return m.renderEditProject()
	case StateSongList:
		return m.renderSongList()
	case StateProjectTemplates:
		return m.renderProjectTemplates()
	default:
		return "Unknown state"
	}
}

// Helper methods

func (m *ProjectManagerModel) updateListDimensions() {
	availableHeight := m.height - 10 // Reserve space for header and footer
	listHeight := availableHeight

	if m.width > 100 {
		listHeight = availableHeight - 2
	}

	m.projectList.SetSize(m.width/2-2, listHeight)
	m.songList.SetSize(m.width/2-2, listHeight)
}

func (m *ProjectManagerModel) setupProjectList() {
	items := make([]list.Item, len(m.projects))
	for i, project := range m.projects {
		items[i] = projectItem{
			id:          project.ID,
			name:        project.Name,
			description: project.Description,
			songCount:   len(project.SongIDs),
			createdAt:   project.CreatedAt,
		}
	}
	m.projectList.SetItems(items)
}

func (m *ProjectManagerModel) setupSongList() {
	items := make([]list.Item, len(m.songs))
	for i, song := range m.songs {
		items[i] = songItem{
			id:     song.ID,
			title:  song.Metadata.Title,
			artist: song.Metadata.Artist,
			key:    song.Metadata.Key,
			tempo:  song.Metadata.Tempo,
			tags:   song.Metadata.Tags,
		}
	}
	m.songList.SetItems(items)
}

func (m *ProjectManagerModel) setupTemplateList() {
	items := make([]list.Item, len(m.templates))
	for i, template := range m.templates {
		items[i] = projectTemplateItem{
			name:        template.Name,
			description: template.Description,
			icon:        template.Icon,
			category:    template.Category,
		}
	}
	m.projectList.SetItems(items)
}

func (m *ProjectManagerModel) cycleFocus() {
	switch m.focusedComponent {
	case "projects":
		if m.state == StateProjectDetail {
			m.focusedComponent = "songs"
		}
	case "songs":
		m.focusedComponent = "projects"
	case "name_input":
		m.focusedComponent = "desc_input"
	case "desc_input":
		m.focusedComponent = "name_input"
	}
}

func (m *ProjectManagerModel) cycleFocusReverse() {
	switch m.focusedComponent {
	case "projects":
		if m.state == StateProjectDetail {
			m.focusedComponent = "songs"
		}
	case "songs":
		m.focusedComponent = "projects"
	case "name_input":
		m.focusedComponent = "desc_input"
	case "desc_input":
		m.focusedComponent = "name_input"
	}
}

func (m *ProjectManagerModel) getTemplateByName(name string) *ProjectTemplate {
	for _, template := range m.templates {
		if template.Name == name {
			return &template
		}
	}
	return nil
}

func (m *ProjectManagerModel) renderLoading() string {
	loadingStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height)

	content := "🎵 Project Manager 🎵\n\n"
	content += "Loading projects...\n\n"
	content += m.spinner.View()

	return loadingStyle.Render(content)
}

// List item types for the project manager

type projectItem struct {
	id          int
	name        string
	description string
	songCount   int
	createdAt   time.Time
}

func (p projectItem) FilterValue() string {
	return p.name
}

type songItem struct {
	id     int
	title  string
	artist string
	key    string
	tempo  int
	tags   []string
}

func (s songItem) FilterValue() string {
	return s.title
}

type projectTemplateItem struct {
	name        string
	description string
	icon        string
	category    string
}

func (p projectTemplateItem) FilterValue() string {
	return p.name
}

// Render methods for different states

func (m *ProjectManagerModel) renderProjectList() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles.TextPrimary)

	// Header
	header := headerStyle.Render("🎵 Project Manager")
	header += "\n\n"

	// Status bar
	statusBar := ""
	if m.statusMessage != "" {
		statusBar = lipgloss.NewStyle().
			Foreground(styles.Success).
			Render("✓ " + m.statusMessage)
	} else if m.errorMessage != "" {
		statusBar = lipgloss.NewStyle().
			Foreground(styles.Error).
			Render("✗ " + m.errorMessage)
	}

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Render("\n[n] New Project  [t] Templates  [Enter] Open  [Esc] Back")

	// Project list
	projectListView := m.projectList.View()

	// Layout
	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		header,
		projectListView,
		instructions,
	)

	// Recent projects sidebar (if space available)
	rightColumn := ""
	if m.width > 100 && len(m.recentProjects) > 0 {
		recentStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(styles.BorderColor).
			Padding(0, 1).
			Width(30)

		recentHeader := lipgloss.NewStyle().
			Foreground(styles.Primary).
			Bold(true).
			Render("Recent Projects")

		recentList := ""
		for _, project := range m.recentProjects {
			recentList += fmt.Sprintf("• %s\n", project.Name)
		}

		rightColumn = recentStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				recentHeader,
				recentList,
			),
		)
	}

	// Combine columns
	if rightColumn != "" {
		content := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
		return contentStyle.Render(content + "\n" + statusBar)
	}

	return contentStyle.Render(leftColumn + "\n" + statusBar)
}

func (m *ProjectManagerModel) renderProjectDetail() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles.TextPrimary)

	// Header
	header := headerStyle.Render(fmt.Sprintf("📁 %s", m.currentProject.Name))
	header += "\n\n"

	// Project info
	info := fmt.Sprintf("Description: %s\n", m.currentProject.Description)
	info += fmt.Sprintf("Songs: %d\n", len(m.currentProject.SongIDs))
	info += fmt.Sprintf("Created: %s\n", m.currentProject.CreatedAt.Format("2006-01-02 15:04"))
	info += fmt.Sprintf("Updated: %s\n", m.currentProject.UpdatedAt.Format("2006-01-02 15:04"))

	infoStyle := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Padding(0, 2).
		Render(info)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Render("\n[Tab] Switch to Songs  [Enter] Open Song  [Esc] Back")

	// Songs list
	songsHeader := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render("Songs in Project:")
	songsHeader += "\n"

	songListView := m.songList.View()

	// Layout
	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		header,
		infoStyle,
		songsHeader,
		songListView,
		instructions,
	)

	return contentStyle.Render(leftColumn)
}

func (m *ProjectManagerModel) renderCreateProject() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles.TextPrimary)

	// Header
	header := headerStyle.Render("📝 Create New Project")
	header += "\n\n"

	// Form
	form := "Project Name:\n"
	form += m.projectNameInput.View() + "\n\n"
	form += "Description:\n"
	form += m.projectDescInput.View() + "\n\n"

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Render("[Tab] Navigate  [Enter] Create  [Esc] Cancel")

	// Layout
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		form,
		instructions,
	)

	return contentStyle.Render(content)
}

func (m *ProjectManagerModel) renderEditProject() string {
	// Similar to create but with existing data
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles.TextPrimary)

	header := headerStyle.Render("✏️  Edit Project")
	header += "\n\n"

	form := "Project Name:\n"
	form += m.projectNameInput.View() + "\n\n"
	form += "Description:\n"
	form += m.projectDescInput.View() + "\n\n"

	instructions := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Render("[Tab] Navigate  [Enter] Save  [Esc] Cancel")

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		form,
		instructions,
	)

	return contentStyle.Render(content)
}

func (m *ProjectManagerModel) renderSongList() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles.TextPrimary)

	// Header
	header := headerStyle.Render(fmt.Sprintf("🎵 Songs in: %s", m.currentProject.Name))
	header += "\n\n"

	// Song list
	songListView := m.songList.View()

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Render("\n[Enter] Open Song  [Esc] Back")

	// Layout
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		songListView,
		instructions,
	)

	return contentStyle.Render(content)
}

func (m *ProjectManagerModel) renderProjectTemplates() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles.TextPrimary)

	// Header
	header := headerStyle.Render("📋 Project Templates")
	header += "\n\n"

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Render("\n[Enter] Use Template  [Esc] Back")

	// Template list
	templateListView := m.projectList.View()

	// Layout
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		templateListView,
		instructions,
	)

	return contentStyle.Render(content)
}

// Compatibility wrapper so existing code (RootModel) that expects ManagerModel/NewManagerModel continues to work.
// It forwards to the new ProjectManagerModel implementation.
type ManagerModel struct {
	inner *ProjectManagerModel
}

func NewManagerModel(database *db.DB) *ManagerModel {
	return &ManagerModel{
		inner: NewProjectManagerModel(database),
	}
}

func (m *ManagerModel) Init() tea.Cmd {
	if m == nil || m.inner == nil {
		return nil
	}
	return m.inner.Init()
}

// Update forwards the message to the inner ProjectManagerModel.
// Keeps the original signature used by RootModel.
func (m *ManagerModel) Update(msg tea.Msg) (*ManagerModel, tea.Cmd) {
	if m == nil || m.inner == nil {
		return m, nil
	}
	_, cmd := m.inner.Update(msg)
	return m, cmd
}

func (m *ManagerModel) View() string {
	if m == nil || m.inner == nil {
		return ""
	}
	return m.inner.View()
}
