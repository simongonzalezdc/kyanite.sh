package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/logging"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FileDialogType represents the type of file dialog
type FileDialogType int

const (
	DialogOpen FileDialogType = iota
	DialogSave
	DialogSaveAs
)

// FileDialogModel represents a file dialog model
type FileDialogModel struct {
	// Dialog configuration
	dialogType  FileDialogType
	title       string
	defaultPath string
	allowedExts []string
	showHidden  bool
	width       int
	height      int

	// UI components
	list      list.Model
	textinput textinput.Model
	visible   bool
	focused   bool

	// State
	currentDir   string
	selectedFile string
	items        []list.Item
	err          error

	// Callbacks
	onConfirm func(string) error
	onCancel  func()

	// Styles
	titleStyle    lipgloss.Style
	borderStyle   lipgloss.Style
	errorStyle    lipgloss.Style
	infoStyle     lipgloss.Style
	selectedStyle lipgloss.Style
}

// FileItem represents a file or directory in the dialog
type FileItem struct {
	name    string
	path    string
	isDir   bool
	size    int64
	modTime time.Time
}

// NewFileDialogModel creates a new file dialog model
func NewFileDialogModel(dialogType FileDialogType, title, defaultPath string, allowedExts []string) *FileDialogModel {
	// Initialize list
	listModel := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false)
	listModel.SetShowHelp(false)

	// Initialize text input for filename
	textInput := textinput.New()
	textInput.Placeholder = "Enter filename..."
	textInput.Focus()

	model := &FileDialogModel{
		dialogType:   dialogType,
		title:        title,
		defaultPath:  defaultPath,
		allowedExts:  allowedExts,
		showHidden:   false,
		list:         listModel,
		textinput:    textInput,
		visible:      false,
		focused:      false,
		currentDir:   "",
		selectedFile: "",
		items:        make([]list.Item, 0),
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1),
		borderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()),
		errorStyle: lipgloss.NewStyle().
			Bold(true),
		infoStyle: lipgloss.NewStyle(),
		selectedStyle: lipgloss.NewStyle().
			Bold(true),
	}

	// Set initial directory
	if defaultPath != "" {
		if filepath.Ext(defaultPath) != "" {
			// It's a file path, get the directory
			model.currentDir = filepath.Dir(defaultPath)
			model.textinput.SetValue(filepath.Base(defaultPath))
		} else {
			// It's a directory path
			model.currentDir = defaultPath
		}
	} else {
		// Use current working directory
		if cwd, err := os.Getwd(); err == nil {
			model.currentDir = cwd
		} else {
			model.currentDir = "."
		}
	}

	return model
}

// Show makes the file dialog visible
func (m *FileDialogModel) Show() {
	m.visible = true
	m.focused = true
	m.loadDirectory()
}

// Hide hides the file dialog
func (m *FileDialogModel) Hide() {
	m.visible = false
	m.focused = false
}

// IsVisible returns whether the dialog is visible
func (m *FileDialogModel) IsVisible() bool {
	return m.visible
}

// SetDimensions sets the dialog dimensions
func (m *FileDialogModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height

	// Adjust list dimensions
	listHeight := height - 10 // Account for title, input, and borders
	if listHeight < 5 {
		listHeight = 5
	}
	listWidth := width - 4 // Account for borders
	if listWidth < 20 {
		listWidth = 20
	}

	m.list.SetWidth(listWidth)
	m.list.SetHeight(listHeight)
	m.textinput.Width = listWidth - 2
}

// SetConfirmCallback sets the callback for when a file is confirmed
func (m *FileDialogModel) SetConfirmCallback(callback func(string) error) {
	m.onConfirm = callback
}

// SetCancelCallback sets the callback for when the dialog is cancelled
func (m *FileDialogModel) SetCancelCallback(callback func()) {
	m.onCancel = callback
}

// Init initializes the file dialog model
func (m *FileDialogModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the file dialog
func (m *FileDialogModel) Update(msg tea.Msg) (*FileDialogModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.onCancel != nil {
				m.onCancel()
			}
			m.Hide()
			return m, nil

		case tea.KeyEnter:
			// Handle based on focus
			if m.list.FilterState() == list.Filtering {
				// If filtering, select the first filtered item
				if m.list.SelectedItem() != nil {
					m.selectItem(m.list.SelectedItem())
				}
			} else if m.textinput.Focused() {
				// Text input is focused, confirm with current input
				m.confirmSelection()
			} else {
				// List is focused, select current item
				if m.list.SelectedItem() != nil {
					m.selectItem(m.list.SelectedItem())
				}
			}
			return m, nil

		case tea.KeyTab:
			// Toggle focus between list and text input
			if m.textinput.Focused() {
				m.textinput.Blur()
				// List doesn't have Focus/Blur methods, so we just track focus state
			} else {
				m.textinput.Focus()
				// List doesn't have Focus/Blur methods, so we just track focus state
			}
			return m, nil

		case tea.KeyCtrlH:
			// Toggle hidden files
			m.showHidden = !m.showHidden
			m.loadDirectory()
			return m, nil

		case tea.KeyBackspace:
			// If in text input and at beginning, go up directory
			if m.textinput.Focused() && m.textinput.Value() == "" {
				m.goUpDirectory()
				return m, nil
			}
		}
	}

	// Update child components
	if m.textinput.Focused() {
		m.textinput, cmd = m.textinput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

// View renders the file dialog
func (m *FileDialogModel) View() string {
	if !m.visible {
		return ""
	}

	// Calculate dialog dimensions
	dialogWidth := m.width - 20
	dialogHeight := m.height - 10
	if dialogWidth < 40 {
		dialogWidth = 40
	}
	if dialogHeight < 15 {
		dialogHeight = 15
	}

	// Get current theme
	t := theme.GetManager().Current()

	// Update styles with theme colors
	titleStyle := m.titleStyle.Foreground(t.Primary).Background(t.Background)
	borderStyle := m.borderStyle.BorderForeground(t.Secondary)
	errorStyle := m.errorStyle.Foreground(t.Error)
	infoStyle := m.infoStyle.Foreground(t.Text)

	// Build dialog content
	var content strings.Builder

	// Title
	title := titleStyle.Render(m.title)
	content.WriteString(title)
	content.WriteString("\n\n")

	// Current directory info
	dirInfo := infoStyle.Render("Directory: " + m.currentDir)
	content.WriteString(dirInfo)
	content.WriteString("\n\n")

	// File list
	if len(m.items) > 0 {
		content.WriteString(m.list.View())
	} else {
		content.WriteString("No files found")
	}
	content.WriteString("\n")

	// File input
	inputLabel := "File: "
	if m.dialogType == DialogSave || m.dialogType == DialogSaveAs {
		inputLabel = "Save as: "
	}
	content.WriteString(inputLabel)
	content.WriteString(m.textinput.View())
	content.WriteString("\n")

	// Help text
	helpText := "Enter: Confirm | Esc: Cancel | Tab: Switch focus | Ctrl+H: Toggle hidden"
	if m.dialogType == DialogOpen {
		helpText += " | Backspace: Parent directory"
	}
	content.WriteString("\n")
	content.WriteString(infoStyle.Render(helpText))

	// Error message
	if m.err != nil {
		content.WriteString("\n")
		content.WriteString(errorStyle.Render("Error: " + m.err.Error()))
	}

	// Wrap in border
	dialogContent := borderStyle.
		Width(dialogWidth).
		Height(dialogHeight).
		Render(content.String())

	return dialogContent
}

// loadDirectory loads the current directory contents
func (m *FileDialogModel) loadDirectory() {
	m.err = nil
	m.items = make([]list.Item, 0)

	// Add parent directory entry
	if m.currentDir != "/" && m.currentDir != "." {
		parentItem := FileItem{
			name:  "..",
			path:  filepath.Dir(m.currentDir),
			isDir: true,
		}
		m.items = append(m.items, parentItem)
	}

	// Read directory
	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		m.err = fmt.Errorf("failed to read directory: %w", err)
		logging.Errorf("Failed to read directory %s: %v", m.currentDir, err)
		return
	}

	// Process entries
	for _, entry := range entries {
		// Skip hidden files unless showHidden is true
		if !m.showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			logging.Warnf("Failed to get file info for %s: %v", entry.Name(), err)
			continue
		}

		fullPath := filepath.Join(m.currentDir, entry.Name())

		// Check file extension for open dialogs
		if !entry.IsDir() && len(m.allowedExts) > 0 {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			allowed := false
			for _, allowedExt := range m.allowedExts {
				if ext == allowedExt {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}

		item := FileItem{
			name:    entry.Name(),
			path:    fullPath,
			isDir:   entry.IsDir(),
			size:    info.Size(),
			modTime: info.ModTime(),
		}

		m.items = append(m.items, item)
	}

	// Update list
	m.list.SetItems(m.items)
	logging.Debugf("Loaded directory %s with %d items", m.currentDir, len(m.items))
}

// selectItem handles selection of an item from the list
func (m *FileDialogModel) selectItem(item list.Item) {
	fileItem, ok := item.(FileItem)
	if !ok {
		return
	}

	if fileItem.isDir {
		// Navigate to directory
		m.currentDir = fileItem.path
		m.textinput.SetValue("")
		m.loadDirectory()
	} else {
		// Select file
		m.selectedFile = fileItem.path
		m.textinput.SetValue(fileItem.name)
		m.confirmSelection()
	}
}

// confirmSelection confirms the current selection
func (m *FileDialogModel) confirmSelection() {
	filename := m.textinput.Value()
	if filename == "" {
		m.err = fmt.Errorf("please enter a filename")
		return
	}

	// Validate filename
	if err := m.validateFilename(filename); err != nil {
		m.err = err
		return
	}

	// Construct full path
	fullPath := filepath.Join(m.currentDir, filename)

	// For save dialogs, check if file exists
	if m.dialogType == DialogSave || m.dialogType == DialogSaveAs {
		if _, err := os.Stat(fullPath); err == nil {
			// File exists, this is okay for save operations
			logging.Debugf("Overwriting existing file: %s", fullPath)
		}
	}

	// For open dialogs, check if file exists
	if m.dialogType == DialogOpen {
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			m.err = fmt.Errorf("file does not exist: %s", filename)
			return
		}
	}

	// Call confirm callback
	if m.onConfirm != nil {
		if err := m.onConfirm(fullPath); err != nil {
			m.err = fmt.Errorf("failed to confirm selection: %w", err)
			logging.Errorf("Failed to confirm file selection: %v", err)
			return
		}
	}

	m.Hide()
}

// goUpDirectory navigates to the parent directory
func (m *FileDialogModel) goUpDirectory() {
	parent := filepath.Dir(m.currentDir)
	if parent != m.currentDir {
		m.currentDir = parent
		m.textinput.SetValue("")
		m.loadDirectory()
	}
}

// validateFilename validates the filename for security and correctness
func (m *FileDialogModel) validateFilename(filename string) error {
	// Check for empty filename
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Check for invalid characters
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	// Check for path traversal attempts
	if strings.Contains(filename, "..") {
		return fmt.Errorf("filename cannot contain parent directory references")
	}

	// Check for absolute paths
	if filepath.IsAbs(filename) {
		return fmt.Errorf("filename must be relative to current directory")
	}

	// Check file extension for save dialogs
	if (m.dialogType == DialogSave || m.dialogType == DialogSaveAs) && len(m.allowedExts) > 0 {
		ext := strings.ToLower(filepath.Ext(filename))
		allowed := false
		for _, allowedExt := range m.allowedExts {
			if ext == allowedExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file extension %s is not allowed", ext)
		}
	}

	return nil
}

// FileItem interface methods for list.Item

func (f FileItem) Title() string {
	suffix := ""
	if f.isDir {
		suffix = "/"
	} else {
		// Add file size for files
		if f.size < 1024 {
			suffix = fmt.Sprintf(" (%dB)", f.size)
		} else if f.size < 1024*1024 {
			suffix = fmt.Sprintf(" (%.1fKB)", float64(f.size)/1024)
		} else {
			suffix = fmt.Sprintf(" (%.1fMB)", float64(f.size)/(1024*1024))
		}
	}
	return f.name + suffix
}

func (f FileItem) Description() string {
	if f.isDir {
		return "Directory"
	}
	return f.modTime.Format("2006-01-02 15:04:05")
}

func (f FileItem) FilterValue() string {
	return f.name
}
