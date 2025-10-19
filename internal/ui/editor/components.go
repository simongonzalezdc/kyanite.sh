package editor

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/noise/internal/app"
	"github.com/puente-labs/noise/internal/app/ai"
	"github.com/puente-labs/noise/internal/domain"
	"github.com/puente-labs/noise/internal/export"
	"github.com/puente-labs/noise/internal/infra/files"
)

// EditorMode represents the different editing modes for rapid prototyping
type EditorMode int

const (
	ModeSketch EditorMode = iota
	ModeDraft
	ModePolish
)

// EditorState manages content, files, and editor state
type EditorState struct {
	textarea         *textarea.Model
	width            int
	height           int
	focused          bool
	showLineNumbers  bool
	lineNumbersWidth int
	wordWrap         bool
	autoIndent       bool
	bracketMatching  bool
	searchMode       bool
	searchQuery      string
	replaceQuery     string
	searchMatches    []int
	currentMatch     int
	cursorLine       int
	cursorColumn     int

	// Content and file management
	currentSong     *domain.Song
	fileService     *files.Service
	currentFilePath string
	autoSaveService *app.AutoSaveService
	lastContent     string
	lastSaveStatus  app.AutoSaveStatus
	exportService   *export.ExportService
	themeManager    interface{}

	// Editor modes
	scratchMode bool
	editorMode  EditorMode
}

// NewEditorState creates a new editor state component
func NewEditorState(textarea *textarea.Model) *EditorState {
	return &EditorState{
		textarea:         textarea,
		focused:          true,
		showLineNumbers:  true,
		lineNumbersWidth: 4,
		wordWrap:         true,
		autoIndent:       true,
		bracketMatching:  true,
		cursorLine:       0,
		cursorColumn:     0,
		lastContent:      "",
		lastSaveStatus:   app.AutoSaveIdle,
		currentSong:      nil,
		fileService:      nil,
		currentFilePath:  "",
		autoSaveService:  nil,
		exportService:    nil,
		themeManager:     nil,
		scratchMode:      false,
		editorMode:       ModeSketch,
	}
}

// EditorShortcuts handles keyboard shortcuts and actions
type EditorShortcuts struct {
	shortcutManager *ShortcutManager
}

// NewEditorShortcuts creates a new editor shortcuts component
func NewEditorShortcuts() *EditorShortcuts {
	return &EditorShortcuts{
		shortcutManager: NewShortcutManager(),
	}
}

// EditorAI handles AI integration features
type EditorAI struct {
	contextDetector     *ai.ContextDetector
	lastContentType     string
	aiAgent             *ai.QuickIdeaAgent
	rapidBrainstorm     bool
	brainstormTheme     string
	brainstormAngles    []string
	continueMode        bool
	continueSuggestions []string
	variationMode       bool
	variationOriginal   string
	variationOptions    []string
}

// NewEditorAI creates a new editor AI component
func NewEditorAI() *EditorAI {
	return &EditorAI{
		contextDetector: ai.NewContextDetector(),
		lastContentType: "Unknown",
		aiAgent:         ai.NewQuickIdeaAgent(),
		rapidBrainstorm: false,
		continueMode:    false,
		variationMode:   false,
	}
}

// EditorMetrics handles status tracking and metrics
type EditorMetrics struct {
	statusBar *StatusBarModel
	width     int
	height    int
}

// NewEditorMetrics creates a new editor metrics component
func NewEditorMetrics() *EditorMetrics {
	return &EditorMetrics{
		statusBar: NewStatusBarModel(),
		width:     0,
		height:    0,
	}
}

// Component interfaces for loose coupling

// StateManagerInterface defines the interface for state management
type StateManagerInterface interface {
	GetText() string
	SetText(text string)
	GetCurrentFilePath() string
	SetCurrentFilePath(path string)
	GetSong() *domain.Song
	SetSong(song *domain.Song)
	IsScratchMode() bool
	SetScratchMode(scratchMode bool)
	GetEditorMode() EditorMode
	SetEditorMode(mode EditorMode)
	UpdateCursorPosition()
	HandleAutoSave()
	ForceSave() error
	GetAutoSaveStatus() app.AutoSaveStatus
	GetLastSaveTime() time.Time
	SaveSong(isMilestone bool, name string) error
	CreateMilestone(name string) error
	RecoverFromLastSave() error
	ToggleLineNumbers()
	ToggleWordWrap()
	ToggleAutoIndent()
	ToggleBracketMatching()
	SetSearchMode(enabled bool)
	IsSearchMode() bool
	SetSearchQuery(query string)
	GetSearchQuery() string
	SetReplaceQuery(query string)
	GetReplaceQuery() string
	NextSearchMatch()
	PreviousSearchMatch()
	NewFile()
	OpenFile(filename string) error
	SaveAs(filename string) error
	CloseFile()
	GetCursorLine() int
	GetCursorColumn() int
	ShowLineNumbers() bool
	WordWrapEnabled() bool
	AutoIndentEnabled() bool
	BracketMatchingEnabled() bool
}

// ShortcutsInterface defines the interface for shortcut handling
type ShortcutsInterface interface {
	GetShortcutManager() *ShortcutManager
	HandleKey(msg tea.KeyMsg) (ShortcutAction, bool)
	HandleShortcutAction(action ShortcutAction, state StateManagerInterface) tea.Cmd
	SetShortcutContext(context KeyContext)
	GetShortcutHints() string
}

// AIInterface defines the interface for AI features
type AIInterface interface {
	StartRapidBrainstorm(theme string)
	GetBrainstormAngles() []string
	SelectBrainstormAngle(index int, state StateManagerInterface)
	StartContinueMode()
	GetContinueSuggestions() []string
	SelectContinueSuggestion(index int, state StateManagerInterface)
	CancelContinueMode()
	StartVariationMode(selectedText string)
	GetVariationOptions() []string
	SelectVariation(index int, state StateManagerInterface)
	CancelVariationMode()
	PerformQualityCheck(state StateManagerInterface)
	UpdateKnowledgeBaseStatus(metrics *EditorMetrics)
	SetAIAgent(agent *ai.QuickIdeaAgent)
	GetAIAgent() *ai.QuickIdeaAgent
	AnalyzeContentType(content string) string
	RenderOverlays(width int) string
}

// MetricsInterface defines the interface for metrics and status
type MetricsInterface interface {
	UpdateStatusBar(state StateManagerInterface, ai *EditorAI)
	RenderStatusBar() string
	SetDimensions(width, height int)
	GetDimensions() (int, int)
	UpdateKnowledgeBaseStatus(available bool, statusText string)
}

// Service injection interfaces

// ServiceInjector defines services that can be injected into components
type ServiceInjector interface {
	SetAutoSaveService(service *app.AutoSaveService)
	SetFileService(service *files.Service)
	SetExportService(service *export.ExportService)
	SetThemeManager(manager interface{}) // Using interface{} to avoid import cycle
}

// Event handling interfaces

// EventHandler defines how components handle Bubble Tea events
type EventHandler interface {
	Update(msg tea.Msg) tea.Cmd
}

// Renderer defines how components render themselves
type Renderer interface {
	View() string
}
