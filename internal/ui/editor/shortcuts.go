package editor

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyContext represents the current context for keyboard shortcuts
type KeyContext int

const (
	ContextGlobal KeyContext = iota
	ContextEditor
	ContextPreview
	ContextSearch
	ContextMenu
	ContextHelp
)

// KeyBinding represents a keyboard shortcut binding
type KeyBinding struct {
	Key         key.Binding
	Description string
	Context     KeyContext
	Category    string
}

// ShortcutManager manages all keyboard shortcuts for the editor
type ShortcutManager struct {
	bindings map[string]*KeyBinding
	context  KeyContext
	helpMode bool
}

// NewShortcutManager creates a new shortcut manager
func NewShortcutManager() *ShortcutManager {
	sm := &ShortcutManager{
		bindings: make(map[string]*KeyBinding),
		context:  ContextGlobal,
		helpMode: false,
	}

	sm.initializeDefaultBindings()
	return sm
}

// initializeDefaultBindings sets up all default keyboard shortcuts
func (sm *ShortcutManager) initializeDefaultBindings() {
	// Navigation shortcuts
	sm.registerBinding("tab", key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next pane")), "Switch to next pane", ContextGlobal, "Navigation")
	sm.registerBinding("shift+tab", key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev pane")), "Switch to previous pane", ContextGlobal, "Navigation")
	sm.registerBinding("ctrl+j", key.NewBinding(key.WithKeys("ctrl+j"), key.WithHelp("ctrl+j", "next pane")), "Next pane", ContextGlobal, "Navigation")
	sm.registerBinding("ctrl+k", key.NewBinding(key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "prev pane")), "Previous pane", ContextGlobal, "Navigation")

	// Editor navigation
	sm.registerBinding("home", key.NewBinding(key.WithKeys("home"), key.WithHelp("home", "start of line")), "Start of line", ContextEditor, "Navigation")
	sm.registerBinding("end", key.NewBinding(key.WithKeys("end"), key.WithHelp("end", "end of line")), "End of line", ContextEditor, "Navigation")
	sm.registerBinding("ctrl+home", key.NewBinding(key.WithKeys("ctrl+home"), key.WithHelp("ctrl+home", "start of file")), "Start of file", ContextEditor, "Navigation")
	sm.registerBinding("ctrl+end", key.NewBinding(key.WithKeys("ctrl+end"), key.WithHelp("ctrl+end", "end of file")), "End of file", ContextEditor, "Navigation")
	sm.registerBinding("page up", key.NewBinding(key.WithKeys("pgup"), key.WithHelp("pgup", "page up")), "Page up", ContextEditor, "Navigation")
	sm.registerBinding("page down", key.NewBinding(key.WithKeys("pgdn"), key.WithHelp("pgdn", "page down")), "Page down", ContextEditor, "Navigation")

	// Word navigation
	sm.registerBinding("ctrl+left", key.NewBinding(key.WithKeys("ctrl+left"), key.WithHelp("ctrl+←", "prev word")), "Previous word", ContextEditor, "Navigation")
	sm.registerBinding("ctrl+right", key.NewBinding(key.WithKeys("ctrl+right"), key.WithHelp("ctrl+→", "next word")), "Next word", ContextEditor, "Navigation")
	sm.registerBinding("alt+left", key.NewBinding(key.WithKeys("alt+left"), key.WithHelp("alt+←", "prev word")), "Previous word", ContextEditor, "Navigation")
	sm.registerBinding("alt+right", key.NewBinding(key.WithKeys("alt+right"), key.WithHelp("alt+→", "next word")), "Next word", ContextEditor, "Navigation")

	// Text editing shortcuts
	sm.registerBinding("ctrl+c", key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "copy")), "Copy", ContextEditor, "Edit")
	sm.registerBinding("ctrl+v", key.NewBinding(key.WithKeys("ctrl+v"), key.WithHelp("ctrl+v", "paste")), "Paste", ContextEditor, "Edit")
	sm.registerBinding("ctrl+x", key.NewBinding(key.WithKeys("ctrl+x"), key.WithHelp("ctrl+x", "cut")), "Cut", ContextEditor, "Edit")
	sm.registerBinding("ctrl+z", key.NewBinding(key.WithKeys("ctrl+z"), key.WithHelp("ctrl+z", "undo")), "Undo", ContextEditor, "Edit")
	sm.registerBinding("ctrl+y", key.NewBinding(key.WithKeys("ctrl+y"), key.WithHelp("ctrl+y", "redo")), "Redo", ContextEditor, "Edit")

	// Text selection
	sm.registerBinding("ctrl+a", key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "select all")), "Select all", ContextEditor, "Edit")
	sm.registerBinding("shift+home", key.NewBinding(key.WithKeys("shift+home"), key.WithHelp("shift+home", "sel to start")), "Select to start of line", ContextEditor, "Edit")
	sm.registerBinding("shift+end", key.NewBinding(key.WithKeys("shift+end"), key.WithHelp("shift+end", "sel to end")), "Select to end of line", ContextEditor, "Edit")
	sm.registerBinding("ctrl+shift+home", key.NewBinding(key.WithKeys("ctrl+shift+home"), key.WithHelp("ctrl+shift+home", "sel to start")), "Select to start of file", ContextEditor, "Edit")
	sm.registerBinding("ctrl+shift+end", key.NewBinding(key.WithKeys("ctrl+shift+end"), key.WithHelp("ctrl+shift+end", "sel to end")), "Select to end of file", ContextEditor, "Edit")
	sm.registerBinding("shift+left", key.NewBinding(key.WithKeys("shift+left"), key.WithHelp("shift+←", "sel left")), "Select left", ContextEditor, "Edit")
	sm.registerBinding("shift+right", key.NewBinding(key.WithKeys("shift+right"), key.WithHelp("shift+→", "sel right")), "Select right", ContextEditor, "Edit")
	sm.registerBinding("shift+up", key.NewBinding(key.WithKeys("shift+up"), key.WithHelp("shift+↑", "sel up")), "Select up", ContextEditor, "Edit")
	sm.registerBinding("shift+down", key.NewBinding(key.WithKeys("shift+down"), key.WithHelp("shift+↓", "sel down")), "Select down", ContextEditor, "Edit")

	// Search and replace
	sm.registerBinding("ctrl+f", key.NewBinding(key.WithKeys("ctrl+f"), key.WithHelp("ctrl+f", "find")), "Find", ContextEditor, "Search")
	sm.registerBinding("ctrl+h", key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "replace")), "Replace", ContextEditor, "Search")
	sm.registerBinding("f3", key.NewBinding(key.WithKeys("f3"), key.WithHelp("f3", "find next")), "Find next", ContextEditor, "Search")
	sm.registerBinding("shift+f3", key.NewBinding(key.WithKeys("shift+f3"), key.WithHelp("shift+f3", "find prev")), "Find previous", ContextEditor, "Search")
	sm.registerBinding("ctrl+g", key.NewBinding(key.WithKeys("ctrl+g"), key.WithHelp("ctrl+g", "goto line")), "Go to line", ContextEditor, "Search")

	// File operations
	sm.registerBinding("ctrl+n", key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "new file")), "New file", ContextGlobal, "File")
	sm.registerBinding("ctrl+o", key.NewBinding(key.WithKeys("ctrl+o"), key.WithHelp("ctrl+o", "open file")), "Open file", ContextGlobal, "File")
	sm.registerBinding("ctrl+s", key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")), "Save", ContextGlobal, "File")
	sm.registerBinding("ctrl+shift+s", key.NewBinding(key.WithKeys("ctrl+shift+s"), key.WithHelp("ctrl+shift+s", "save as")), "Save as", ContextGlobal, "File")
	sm.registerBinding("ctrl+w", key.NewBinding(key.WithKeys("ctrl+w"), key.WithHelp("ctrl+w", "close")), "Close file", ContextGlobal, "File")
	sm.registerBinding("ctrl+e", key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "export")), "Export file", ContextGlobal, "File")

	// Editor features
	sm.registerBinding("ctrl+l", key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "toggle line nums")), "Toggle line numbers", ContextEditor, "View")
	sm.registerBinding("ctrl+w", key.NewBinding(key.WithKeys("ctrl+w"), key.WithHelp("ctrl+w", "toggle wrap")), "Toggle word wrap", ContextEditor, "View")
	sm.registerBinding("ctrl+shift+w", key.NewBinding(key.WithKeys("ctrl+shift+w"), key.WithHelp("ctrl+shift+w", "toggle wrap")), "Toggle word wrap", ContextEditor, "View")
	sm.registerBinding("ctrl+i", key.NewBinding(key.WithKeys("ctrl+i"), key.WithHelp("ctrl+i", "toggle indent")), "Toggle auto-indent", ContextEditor, "View")
	sm.registerBinding("ctrl+b", key.NewBinding(key.WithKeys("ctrl+b"), key.WithHelp("ctrl+b", "toggle bracket")), "Toggle bracket matching", ContextEditor, "View")

	// Application navigation
	sm.registerBinding("esc", key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back/menu")), "Back to menu", ContextGlobal, "Application")
	sm.registerBinding("f1", key.NewBinding(key.WithKeys("f1"), key.WithHelp("f1", "help")), "Show help", ContextGlobal, "Application")
	sm.registerBinding("?", key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")), "Show help", ContextGlobal, "Application")
	sm.registerBinding("ctrl+q", key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("ctrl+q", "quit")), "Quit application", ContextGlobal, "Application")
	sm.registerBinding("ctrl+,", key.NewBinding(key.WithKeys("ctrl+,"), key.WithHelp("ctrl+,", "settings")), "Settings", ContextGlobal, "Application")

	// Theory tools (when available)
	sm.registerBinding("ctrl+t", key.NewBinding(key.WithKeys("ctrl+t"), key.WithHelp("ctrl+t", "theory tools")), "Theory tools", ContextGlobal, "Tools")
	sm.registerBinding("ctrl+m", key.NewBinding(key.WithKeys("ctrl+m"), key.WithHelp("ctrl+m", "audio tools")), "Audio tools", ContextGlobal, "Tools")

	// Preview navigation
	sm.registerBinding("up", key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "up")), "Up", ContextPreview, "Navigation")
	sm.registerBinding("down", key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "down")), "Down", ContextPreview, "Navigation")
	sm.registerBinding("left", key.NewBinding(key.WithKeys("left"), key.WithHelp("←", "left")), "Left", ContextPreview, "Navigation")
	sm.registerBinding("right", key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "right")), "Right", ContextPreview, "Navigation")
	sm.registerBinding("space", key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "next page")), "Next page", ContextPreview, "Navigation")
	sm.registerBinding("shift+space", key.NewBinding(key.WithKeys("shift+space"), key.WithHelp("shift+space", "prev page")), "Previous page", ContextPreview, "Navigation")
}

// registerBinding registers a new key binding
func (sm *ShortcutManager) registerBinding(keyStr string, keyBinding key.Binding, description string, context KeyContext, category string) {
	sm.bindings[keyStr] = &KeyBinding{
		Key:         keyBinding,
		Description: description,
		Context:     context,
		Category:    category,
	}
}

// SetContext sets the current keyboard context
func (sm *ShortcutManager) SetContext(context KeyContext) {
	sm.context = context
}

// GetContext returns the current keyboard context
func (sm *ShortcutManager) GetContext() KeyContext {
	return sm.context
}

// SetHelpMode enables or disables help mode
func (sm *ShortcutManager) SetHelpMode(enabled bool) {
	sm.helpMode = enabled
}

// IsHelpMode returns whether help mode is active
func (sm *ShortcutManager) IsHelpMode() bool {
	return sm.helpMode
}

// HandleKey handles a key press and returns the appropriate action
func (sm *ShortcutManager) HandleKey(msg tea.KeyMsg) (ShortcutAction, bool) {
	keyStr := msg.String()

	// Check for exact key matches first
	if binding, exists := sm.bindings[keyStr]; exists {
		if sm.isBindingActive(binding) {
			return sm.createActionFromBinding(binding, keyStr), true
		}
	}

	// Check for help mode toggle
	if keyStr == "f1" || keyStr == "?" {
		sm.helpMode = !sm.helpMode
		return ShortcutAction{Type: ActionToggleHelp}, true
	}

	return ShortcutAction{}, false
}

// isBindingActive checks if a binding is active in the current context
func (sm *ShortcutManager) isBindingActive(binding *KeyBinding) bool {
	// Global context bindings are always active
	if binding.Context == ContextGlobal {
		return true
	}

	// Context-specific bindings are only active when in that context
	return binding.Context == sm.context
}

// createActionFromBinding creates a shortcut action from a key binding
func (sm *ShortcutManager) createActionFromBinding(binding *KeyBinding, keyStr string) ShortcutAction {
	action := ShortcutAction{
		Type:        ActionUnknown,
		Description: binding.Description,
		Category:    binding.Category,
	}

	switch {
	// Navigation actions
	case keyStr == "tab" || keyStr == "ctrl+j":
		action.Type = ActionNextPane
	case keyStr == "shift+tab" || keyStr == "ctrl+k":
		action.Type = ActionPrevPane
	case keyStr == "esc":
		action.Type = ActionBackToMenu
	case keyStr == "home":
		action.Type = ActionStartOfLine
	case keyStr == "end":
		action.Type = ActionEndOfLine
	case keyStr == "ctrl+home":
		action.Type = ActionStartOfFile
	case keyStr == "ctrl+end":
		action.Type = ActionEndOfFile
	case keyStr == "page up" || keyStr == "pgup":
		action.Type = ActionPageUp
	case keyStr == "page down" || keyStr == "pgdn":
		action.Type = ActionPageDown
	case keyStr == "ctrl+left" || keyStr == "alt+left":
		action.Type = ActionPrevWord
	case keyStr == "ctrl+right" || keyStr == "alt+right":
		action.Type = ActionNextWord

	// Text editing actions
	case keyStr == "ctrl+c":
		action.Type = ActionCopy
	case keyStr == "ctrl+v":
		action.Type = ActionPaste
	case keyStr == "ctrl+x":
		action.Type = ActionCut
	case keyStr == "ctrl+z":
		action.Type = ActionUndo
	case keyStr == "ctrl+y":
		action.Type = ActionRedo

	// Text selection actions
	case keyStr == "ctrl+a":
		action.Type = ActionSelectAll
	case keyStr == "shift+home":
		action.Type = ActionSelectToStartOfLine
	case keyStr == "shift+end":
		action.Type = ActionSelectToEndOfLine
	case keyStr == "ctrl+shift+home":
		action.Type = ActionSelectToStartOfFile
	case keyStr == "ctrl+shift+end":
		action.Type = ActionSelectToEndOfFile
	case keyStr == "shift+left":
		action.Type = ActionSelectLeft
	case keyStr == "shift+right":
		action.Type = ActionSelectRight
	case keyStr == "shift+up":
		action.Type = ActionSelectUp
	case keyStr == "shift+down":
		action.Type = ActionSelectDown

	// Search actions
	case keyStr == "ctrl+f":
		action.Type = ActionFind
	case keyStr == "ctrl+h":
		action.Type = ActionReplace
	case keyStr == "f3":
		action.Type = ActionFindNext
	case keyStr == "shift+f3":
		action.Type = ActionFindPrev
	case keyStr == "ctrl+g":
		action.Type = ActionGoToLine

	// File operations
	case keyStr == "ctrl+n":
		action.Type = ActionNewFile
	case keyStr == "ctrl+o":
		action.Type = ActionOpenFile
	case keyStr == "ctrl+s":
		action.Type = ActionSave
	case keyStr == "ctrl+shift+s":
		action.Type = ActionSaveAs
	case keyStr == "ctrl+e":
		action.Type = ActionExport
	case keyStr == "ctrl+w":
		action.Type = ActionCloseFile

	// Editor features
	case keyStr == "ctrl+l":
		action.Type = ActionToggleLineNumbers
	case keyStr == "ctrl+w" && sm.context == ContextEditor:
		action.Type = ActionToggleWordWrap
	case keyStr == "ctrl+i":
		action.Type = ActionToggleAutoIndent
	case keyStr == "ctrl+b":
		action.Type = ActionToggleBracketMatching

	// Application actions
	case keyStr == "ctrl+q":
		action.Type = ActionQuit
	case keyStr == "ctrl+,":
		action.Type = ActionSettings
	case keyStr == "ctrl+t":
		action.Type = ActionTheoryTools
	case keyStr == "ctrl+m":
		action.Type = ActionAudioTools

	// Preview actions
	case keyStr == "up":
		action.Type = ActionPreviewUp
	case keyStr == "down":
		action.Type = ActionPreviewDown
	case keyStr == "left":
		action.Type = ActionPreviewLeft
	case keyStr == "right":
		action.Type = ActionPreviewRight
	case keyStr == "space":
		action.Type = ActionPreviewNextPage
	case keyStr == "shift+space":
		action.Type = ActionPreviewPrevPage
	}

	return action
}

// GetBindingsForContext returns all bindings for a specific context
func (sm *ShortcutManager) GetBindingsForContext(context KeyContext) []*KeyBinding {
	var bindings []*KeyBinding

	for _, binding := range sm.bindings {
		if binding.Context == context || binding.Context == ContextGlobal {
			bindings = append(bindings, binding)
		}
	}

	return bindings
}

// GetBindingsByCategory returns all bindings in a specific category
func (sm *ShortcutManager) GetBindingsByCategory(category string) []*KeyBinding {
	var bindings []*KeyBinding

	for _, binding := range sm.bindings {
		if binding.Category == category && (binding.Context == sm.context || binding.Context == ContextGlobal) {
			bindings = append(bindings, binding)
		}
	}

	return bindings
}

// GetAllBindings returns all registered bindings
func (sm *ShortcutManager) GetAllBindings() []*KeyBinding {
	var bindings []*KeyBinding
	for _, binding := range sm.bindings {
		bindings = append(bindings, binding)
	}
	return bindings
}

// GetHelpText returns formatted help text for shortcuts
func (sm *ShortcutManager) GetHelpText() string {
	var sections []string

	categories := []string{"Navigation", "Edit", "Search", "File", "View", "Application", "Tools"}

	for _, category := range categories {
		bindings := sm.GetBindingsByCategory(category)
		if len(bindings) > 0 {
			var categoryBindings []string
			categoryBindings = append(categoryBindings, fmt.Sprintf("\n%s:", category))

			for _, binding := range bindings {
				keyHelp := binding.Key.Help().Key
				desc := binding.Description
				categoryBindings = append(categoryBindings, fmt.Sprintf("  %-15s %s", keyHelp, desc))
			}

			sections = append(sections, strings.Join(categoryBindings, "\n"))
		}
	}

	return strings.Join(sections, "\n")
}

// GetStatusBarHints returns shortcut hints for the status bar
func (sm *ShortcutManager) GetStatusBarHints() string {
	var hints []string

	// Context-specific hints
	switch sm.context {
	case ContextEditor:
		editorHints := []string{"Ctrl+F:Find", "Ctrl+S:Save", "Tab:Next Pane"}
		hints = append(hints, editorHints...)
	case ContextPreview:
		previewHints := []string{"↑↓:Navigate", "Space:Next Page", "Tab:Next Pane"}
		hints = append(hints, previewHints...)
	case ContextGlobal:
		globalHints := []string{"F1:Help", "Esc:Menu", "Ctrl+Q:Quit"}
		hints = append(hints, globalHints...)
	}

	// Always show help hint
	if !sm.helpMode {
		hints = append(hints, "F1:Help")
	}

	return strings.Join(hints, " | ")
}

// ShortcutAction represents an action triggered by a keyboard shortcut
type ShortcutAction struct {
	Type        ShortcutActionType
	Description string
	Category    string
	Data        interface{}
}

// ShortcutActionType represents different types of shortcut actions
type ShortcutActionType int

const (
	ActionUnknown ShortcutActionType = iota
	ActionNextPane
	ActionPrevPane
	ActionBackToMenu
	ActionStartOfLine
	ActionEndOfLine
	ActionStartOfFile
	ActionEndOfFile
	ActionPageUp
	ActionPageDown
	ActionPrevWord
	ActionNextWord
	ActionCopy
	ActionPaste
	ActionCut
	ActionUndo
	ActionRedo
	ActionSelectAll
	ActionSelectToStartOfLine
	ActionSelectToEndOfLine
	ActionSelectToStartOfFile
	ActionSelectToEndOfFile
	ActionSelectLeft
	ActionSelectRight
	ActionSelectUp
	ActionSelectDown
	ActionFind
	ActionReplace
	ActionFindNext
	ActionFindPrev
	ActionGoToLine
	ActionNewFile
	ActionOpenFile
	ActionSave
	ActionSaveAs
	ActionExport
	ActionCloseFile
	ActionToggleLineNumbers
	ActionToggleWordWrap
	ActionToggleAutoIndent
	ActionToggleBracketMatching
	ActionQuit
	ActionSettings
	ActionTheoryTools
	ActionAudioTools
	ActionToggleHelp
	ActionPreviewUp
	ActionPreviewDown
	ActionPreviewLeft
	ActionPreviewRight
	ActionPreviewNextPage
	ActionPreviewPrevPage
)

// String returns a string representation of the action type
func (sat ShortcutActionType) String() string {
	switch sat {
	case ActionNextPane:
		return "NextPane"
	case ActionPrevPane:
		return "PrevPane"
	case ActionBackToMenu:
		return "BackToMenu"
	case ActionStartOfLine:
		return "StartOfLine"
	case ActionEndOfLine:
		return "EndOfLine"
	case ActionStartOfFile:
		return "StartOfFile"
	case ActionEndOfFile:
		return "EndOfFile"
	case ActionPageUp:
		return "PageUp"
	case ActionPageDown:
		return "PageDown"
	case ActionPrevWord:
		return "PrevWord"
	case ActionNextWord:
		return "NextWord"
	case ActionCopy:
		return "Copy"
	case ActionPaste:
		return "Paste"
	case ActionCut:
		return "Cut"
	case ActionUndo:
		return "Undo"
	case ActionRedo:
		return "Redo"
	case ActionSelectAll:
		return "SelectAll"
	case ActionSelectToStartOfLine:
		return "SelectToStartOfLine"
	case ActionSelectToEndOfLine:
		return "SelectToEndOfLine"
	case ActionSelectToStartOfFile:
		return "SelectToStartOfFile"
	case ActionSelectToEndOfFile:
		return "SelectToEndOfFile"
	case ActionSelectLeft:
		return "SelectLeft"
	case ActionSelectRight:
		return "SelectRight"
	case ActionSelectUp:
		return "SelectUp"
	case ActionSelectDown:
		return "SelectDown"
	case ActionFind:
		return "Find"
	case ActionReplace:
		return "Replace"
	case ActionFindNext:
		return "FindNext"
	case ActionFindPrev:
		return "FindPrev"
	case ActionGoToLine:
		return "GoToLine"
	case ActionNewFile:
		return "NewFile"
	case ActionOpenFile:
		return "OpenFile"
	case ActionSave:
		return "Save"
	case ActionSaveAs:
		return "SaveAs"
	case ActionExport:
		return "Export"
	case ActionCloseFile:
		return "CloseFile"
	case ActionToggleLineNumbers:
		return "ToggleLineNumbers"
	case ActionToggleWordWrap:
		return "ToggleWordWrap"
	case ActionToggleAutoIndent:
		return "ToggleAutoIndent"
	case ActionToggleBracketMatching:
		return "ToggleBracketMatching"
	case ActionQuit:
		return "Quit"
	case ActionSettings:
		return "Settings"
	case ActionTheoryTools:
		return "TheoryTools"
	case ActionAudioTools:
		return "AudioTools"
	case ActionToggleHelp:
		return "ToggleHelp"
	case ActionPreviewUp:
		return "PreviewUp"
	case ActionPreviewDown:
		return "PreviewDown"
	case ActionPreviewLeft:
		return "PreviewLeft"
	case ActionPreviewRight:
		return "PreviewRight"
	case ActionPreviewNextPage:
		return "PreviewNextPage"
	case ActionPreviewPrevPage:
		return "PreviewPrevPage"
	default:
		return "Unknown"
	}
}
