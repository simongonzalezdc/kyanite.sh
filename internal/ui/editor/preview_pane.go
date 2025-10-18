package editor

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/lyricforge/internal/ui/styles"
)

// PreviewPaneModel handles the markdown preview pane
type PreviewPaneModel struct {
	content   string
	width     int
	height    int
	focused   bool
	scrollPos int
	zoomLevel float64

	// Glamour renderer
	renderer *glamour.TermRenderer

	// Styles
	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
	borderStyle  lipgloss.Style
	titleStyle   lipgloss.Style

	// Control states
	showRefresh bool
	lastError   string

	// Smooth scrolling animation
	scrollSpring  harmonica.Spring
	targetScroll  float64
	currentScroll float64
	isScrolling   bool

	// Keyboard shortcuts
	shortcutManager *ShortcutManager

	// Real-time preview features
	realtimeManager     *RealTimePreviewManager
	lastScrollPositions map[string]int // Track scroll positions by content hash
	scrollSyncEnabled   bool

	// Advanced preview features
	showWordCount   bool
	showReadingTime bool
	showTOC         bool
	previewStats    PreviewStats
	tocEntries      []TOCEntry

	// Performance tracking
	renderCache  map[string]string
	cacheMutex   sync.RWMutex
	maxCacheSize int
}

// NewPreviewPaneModel creates a new preview pane model
func NewPreviewPaneModel() *PreviewPaneModel {
	// Create custom dark theme for "Midnight Jazz" aesthetic
	styleConfig := createMidnightJazzStyle()

	// Initialize Glamour renderer with custom style
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(styleConfig),
	)

	if err != nil {
		// Fallback to basic dark style if custom style fails
		renderer, _ = glamour.NewTermRenderer(
			glamour.WithStylePath("dark"),
		)
	}

	// Initialize smooth scrolling spring
	scrollSpring := harmonica.NewSpring(
		harmonica.FPS(60),
		8.0, // Higher frequency for responsive scrolling
		0.6, // Balanced damping
	)

	model := &PreviewPaneModel{
		content:         "Preview will appear here...",
		focused:         false,
		scrollPos:       0,
		zoomLevel:       1.0,
		renderer:        renderer,
		showRefresh:     false,
		lastError:       "",
		scrollSpring:    scrollSpring,
		targetScroll:    0.0,
		currentScroll:   0.0,
		isScrolling:     false,
		shortcutManager: NewShortcutManager(),
		focusedStyle:    styles.BorderActive,
		blurredStyle:    styles.Border,
		borderStyle:     styles.Border,
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.TextPrimary).
			Background(styles.Dark3).
			Padding(0, 1).
			Width(20),
		// Initialize real-time preview features
		realtimeManager:     NewRealTimePreviewManager(DefaultPreviewUpdateConfig()),
		lastScrollPositions: make(map[string]int),
		scrollSyncEnabled:   true,
		showWordCount:       true,
		showReadingTime:     true,
		showTOC:             true,
		renderCache:         make(map[string]string),
		maxCacheSize:        50,
	}

	// Set up real-time preview callbacks
	model.realtimeManager.SetCallbacks(
		model.onContentUpdate,
		model.onUpdateStart,
		model.onUpdateComplete,
		model.onValidationError,
		model.onStatsUpdate,
	)

	return model
}

// createMidnightJazzStyle creates a custom dark theme for the "Midnight Jazz" aesthetic
func createMidnightJazzStyle() string {
	return fmt.Sprintf(`
{
	"document": {
		"background": "%s",
		"foreground": "%s",
		"margin": 2
	},
	"block_quote": {
		"background": "%s",
		"border-left": "4px solid %s",
		"foreground": "%s",
		"margin": "1em 0"
	},
	"code": {
		"background": "%s",
		"foreground": "%s",
		"padding": "0.2em 0.4em",
		"border-radius": "3px"
	},
	"code_block": {
		"background": "%s",
		"foreground": "%s",
		"margin": "1em 0",
		"padding": "1em",
		"border": "1px solid %s",
		"border-radius": "5px"
	},
	"em": {
		"foreground": "%s",
		"font-style": "italic"
	},
	"heading": {
		"foreground": "%s",
		"font-weight": "bold"
	},
	"h1": {
		"foreground": "%s",
		"font-size": "1.8em",
		"margin": "0.5em 0"
	},
	"h2": {
		"foreground": "%s",
		"font-size": "1.6em",
		"margin": "0.4em 0"
	},
	"h3": {
		"foreground": "%s",
		"font-size": "1.4em",
		"margin": "0.3em 0"
	},
	"h4": {
		"foreground": "%s",
		"font-size": "1.2em",
		"margin": "0.2em 0"
	},
	"h5": {
		"foreground": "%s",
		"font-size": "1.1em",
		"margin": "0.1em 0"
	},
	"h6": {
		"foreground": "%s",
		"font-size": "1em",
		"margin": "0.1em 0"
	},
	"hr": {
		"background": "%s",
		"margin": "1em 0"
	},
	"image": {
		"margin": "1em 0"
	},
	"link": {
		"foreground": "%s",
		"text-decoration": "underline"
	},
	"list": {
		"margin": "0.5em 0",
		"padding-left": "1.5em"
	},
	"table": {
		"background": "%s",
		"border": "1px solid %s",
		"border-collapse": "collapse",
		"margin": "1em 0"
	},
	"thead": {
		"background": "%s",
		"font-weight": "bold"
	},
	"th": {
		"border": "1px solid %s",
		"padding": "0.5em"
	},
	"td": {
		"border": "1px solid %s",
		"padding": "0.5em"
	},
	"strong": {
		"foreground": "%s",
		"font-weight": "bold"
	}
}
`,
		// Background colors
		styles.ToHex(styles.Dark1),
		styles.ToHex(styles.TextPrimary),
		styles.Dark2,
		styles.BorderColor,
		styles.TextSecondary,
		styles.Dark2,
		styles.Accent,
		styles.Dark1,
		styles.TextPrimary,
		styles.BorderColor,
		styles.TextAccent,
		styles.Primary,
		styles.Success,
		styles.Secondary,
		styles.TextAccent,
		styles.Accent,
		styles.Info,
		styles.TextSecondary,
		styles.BorderColor,
		styles.Dark2,
		styles.BorderColor,
		styles.Dark3,
		styles.BorderColor,
		styles.BorderColor,
		styles.BorderColor,
		styles.Primary,
	)
}

// Init initializes the preview pane
func (m *PreviewPaneModel) Init() tea.Cmd {
	return m.startScrollAnimation()
}

// startScrollAnimation starts the scroll animation loop
func (m *PreviewPaneModel) startScrollAnimation() tea.Cmd {
	return func() tea.Msg {
		// Animation tick for smooth scrolling
		for m.isScrolling {
			time.Sleep(time.Second / 60) // 60 FPS

			// Update spring animation
			m.currentScroll, _ = m.scrollSpring.Update(m.currentScroll, 0.0, m.targetScroll)

			// Check if we've reached the target
			if abs(m.targetScroll-m.currentScroll) < 0.1 {
				m.currentScroll = m.targetScroll
				m.scrollPos = int(m.currentScroll)
				m.isScrolling = false
				break
			}

			m.scrollPos = int(m.currentScroll)
		}
		return nil
	}
}

// Update handles messages for the preview pane
func (m *PreviewPaneModel) Update(msg tea.Msg) (*PreviewPaneModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.focused {
				m.scrollUp()
				cmds = append(cmds, m.startScrollAnimation())
			}
		case "down", "j":
			if m.focused {
				m.scrollDown()
				cmds = append(cmds, m.startScrollAnimation())
			}
		case "pgup":
			if m.focused {
				m.scrollPageUp()
				cmds = append(cmds, m.startScrollAnimation())
			}
		case "pgdown":
			if m.focused {
				m.scrollPageDown()
				cmds = append(cmds, m.startScrollAnimation())
			}
		case "ctrl+r":
			if m.focused {
				m.refreshPreview()
			}
		case "ctrl++", "ctrl+=":
			if m.focused {
				m.zoomIn()
			}
		case "ctrl+-":
			if m.focused {
				m.zoomOut()
			}
		case "ctrl+0":
			if m.focused {
				m.resetZoom()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the preview pane
func (m *PreviewPaneModel) View() string {
	var style lipgloss.Style

	if m.focused {
		style = m.focusedStyle
	} else {
		style = m.blurredStyle
	}

	// Render content with Glamour
	renderedContent, err := m.renderContentWithGlamour()
	if err != nil {
		renderedContent = m.renderError(err)
	}

	// Calculate visible lines
	visibleHeight := m.height - 8 // Account for padding, borders, title, and controls
	start := m.scrollPos
	end := start + visibleHeight

	// Split content into lines
	lines := strings.Split(renderedContent, "\n")

	// Show only visible portion
	var visibleLines []string
	if start < len(lines) {
		if end > len(lines) {
			end = len(lines)
		}
		visibleLines = lines[start:end]
	}

	displayContent := strings.Join(visibleLines, "\n")

	// Create title bar with controls
	title := m.createTitleBar()

	// Add real-time update indicator
	updateIndicator := ""
	if m.realtimeManager != nil {
		updateIndicator = m.realtimeManager.GetUpdateIndicatorView()
	}

	// Add scroll indicator if needed
	scrollIndicator := ""
	if len(lines) > visibleHeight {
		scrollPercent := float64(m.scrollPos) / float64(len(lines)-visibleHeight)
		if scrollPercent > 1.0 {
			scrollPercent = 1.0
		}
		progressBar := strings.Repeat("█", int(scrollPercent*20)) + strings.Repeat("░", 20-int(scrollPercent*20))
		scrollIndicator = lipgloss.NewStyle().
			Foreground(styles.TextMuted).
			Align(lipgloss.Center).
			Width(m.width - 4).
			Render("↓ " + progressBar + " ↑")
	}

	// Add preview statistics if enabled
	statsInfo := ""
	if (m.showWordCount || m.showReadingTime) && m.previewStats.WordCount > 0 {
		var statsParts []string
		if m.showWordCount {
			statsParts = append(statsParts, fmt.Sprintf("%d words", m.previewStats.WordCount))
		}
		if m.showReadingTime {
			statsParts = append(statsParts, fmt.Sprintf("%s read", m.formatDuration(m.previewStats.ReadingTime)))
		}
		if len(statsParts) > 0 {
			statsInfo = lipgloss.NewStyle().
				Foreground(styles.TextMuted).
				Align(lipgloss.Center).
				Width(m.width - 4).
				Render(strings.Join(statsParts, " • "))
		}
	}

	// Add table of contents if enabled and available
	tocInfo := ""
	if m.showTOC && len(m.tocEntries) > 0 {
		tocInfo = m.renderTOC()
	}

	// Add controls info if focused
	controlsInfo := ""
	if m.focused {
		controlsInfo = lipgloss.NewStyle().
			Foreground(styles.TextMuted).
			Align(lipgloss.Center).
			Width(m.width - 4).
			Render("Controls: ↑↓/jk: scroll | Ctrl+R: refresh | Ctrl+±: zoom | Ctrl+0: reset zoom")
	}

	// Combine all elements
	fullContent := lipgloss.JoinVertical(lipgloss.Left, title)
	if updateIndicator != "" {
		fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, updateIndicator)
	}
	if displayContent != "" {
		fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, displayContent)
	}
	if scrollIndicator != "" {
		fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, scrollIndicator)
	}
	if statsInfo != "" {
		fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, statsInfo)
	}
	if tocInfo != "" {
		fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, tocInfo)
	}
	if controlsInfo != "" {
		fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, controlsInfo)
	}

	return style.Width(m.width).Height(m.height).Render(fullContent)
}

// SetDimensions sets the pane dimensions
func (m *PreviewPaneModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// Focus focuses the preview pane
func (m *PreviewPaneModel) Focus() {
	m.focused = true
}

// Blur blurs the preview pane
func (m *PreviewPaneModel) Blur() {
	m.focused = false
}

// SetContent sets the markdown content to preview with real-time updates
func (m *PreviewPaneModel) SetContent(content string) {
	if m.realtimeManager != nil {
		m.realtimeManager.UpdateContent(content, ChangeSourceExternal)
	} else {
		m.content = content
	}
}

// GetContent returns the current content
func (m *PreviewPaneModel) GetContent() string {
	return m.content
}

// Scroll methods with smooth animation
func (m *PreviewPaneModel) scrollUp() {
	if m.scrollPos > 0 {
		m.targetScroll = float64(m.scrollPos - 1)
		m.isScrolling = true
	}
}

func (m *PreviewPaneModel) scrollDown() {
	visibleHeight := m.height - 6
	lines := strings.Split(m.renderContent(), "\n")
	maxScroll := len(lines) - visibleHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollPos < maxScroll {
		m.targetScroll = float64(m.scrollPos + 1)
		m.isScrolling = true
	}
}

func (m *PreviewPaneModel) scrollPageUp() {
	visibleHeight := m.height - 6
	newPos := m.scrollPos - visibleHeight
	if newPos < 0 {
		newPos = 0
	}
	m.targetScroll = float64(newPos)
	m.isScrolling = true
}

func (m *PreviewPaneModel) scrollPageDown() {
	visibleHeight := m.height - 6
	lines := strings.Split(m.renderContent(), "\n")
	maxScroll := len(lines) - visibleHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	newPos := m.scrollPos + visibleHeight
	if newPos > maxScroll {
		newPos = maxScroll
	}
	m.targetScroll = float64(newPos)
	m.isScrolling = true
}

// renderContentWithGlamour renders the content using Glamour
func (m *PreviewPaneModel) renderContentWithGlamour() (string, error) {
	if m.renderer == nil {
		return m.renderBasicContent(), nil
	}

	// Pre-process content for lyric-specific formatting
	processedContent := m.preprocessLyricContent()

	// Render with Glamour
	rendered, err := m.renderer.Render(processedContent)
	if err != nil {
		m.lastError = err.Error()
		return m.renderBasicContent(), err
	}

	m.lastError = ""
	return rendered, nil
}

// renderBasicContent renders content with basic formatting as fallback
func (m *PreviewPaneModel) renderBasicContent() string {
	content := m.content

	// Basic markdown-like formatting
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		// Headers
		if strings.HasPrefix(line, "# ") {
			lines[i] = strings.Replace(line, "# ", "📝 ", 1)
		} else if strings.HasPrefix(line, "## ") {
			lines[i] = strings.Replace(line, "## ", "📋 ", 1)
		} else if strings.HasPrefix(line, "### ") {
			lines[i] = strings.Replace(line, "### ", "📌 ", 1)
		}

		// Bold text
		lines[i] = strings.ReplaceAll(lines[i], "**", "")
		// Italic text
		lines[i] = strings.ReplaceAll(lines[i], "*", "•")

		// Code blocks
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			lines[i] = "┌─ Code Block ─────────────────────────────────"
		}
	}

	content = strings.Join(lines, "\n")
	content = strings.ReplaceAll(content, "\n\n", "\n")

	return content
}

// renderError renders an error message
func (m *PreviewPaneModel) renderError(err error) string {
	errorMsg := fmt.Sprintf("Error rendering markdown: %s", err.Error())
	return lipgloss.NewStyle().
		Foreground(styles.Error).
		Render("⚠️  " + errorMsg)
}

// createTitleBar creates the title bar with status information
func (m *PreviewPaneModel) createTitleBar() string {
	title := "Preview"
	if m.focused {
		title = "Preview (Focused)"
	}

	// Add zoom level indicator
	if m.zoomLevel != 1.0 {
		title += fmt.Sprintf(" | Zoom: %.0f%%", m.zoomLevel*100)
	}

	// Add error indicator if present
	if m.lastError != "" {
		title += " | ⚠️ Error"
	}

	// Add refresh indicator if needed
	if m.showRefresh {
		title += " | ⟳ Refreshing..."
	}

	titleBar := m.titleStyle.Render(title)
	return titleBar
}

// preprocessLyricContent processes content for lyric-specific formatting
func (m *PreviewPaneModel) preprocessLyricContent() string {
	content := m.content
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Handle verse markers [Verse 1], [Chorus], etc.
		if strings.Contains(line, "[Verse") {
			lines[i] = strings.ReplaceAll(line, "[Verse", "**[Verse**")
			lines[i] = strings.ReplaceAll(lines[i], "]", "]**")
		} else if strings.Contains(line, "[Chorus]") {
			lines[i] = "**[Chorus]**"
		} else if strings.Contains(line, "[Bridge]") {
			lines[i] = "**[Bridge]**"
		} else if strings.Contains(line, "[Outro]") {
			lines[i] = "**[Outro]**"
		} else if strings.Contains(line, "[Intro]") {
			lines[i] = "**[Intro]**"
		}

		// Handle section markers with numbers
		if strings.Contains(line, "[Chorus 1]") {
			lines[i] = "**[Chorus 1]**"
		} else if strings.Contains(line, "[Verse 2]") {
			lines[i] = "**[Verse 2]**"
		}
	}

	return strings.Join(lines, "\n")
}

// Control methods

// refreshPreview refreshes the preview content
func (m *PreviewPaneModel) refreshPreview() {
	m.showRefresh = true
	// Re-render the current content to refresh the preview
	_, err := m.renderContentWithGlamour()
	if err != nil {
		m.lastError = err.Error()
	}
	m.showRefresh = false
}

// zoomIn increases the zoom level
func (m *PreviewPaneModel) zoomIn() {
	if m.zoomLevel < 2.0 {
		m.zoomLevel += 0.1
		// Note: Status bar update will be handled by parent component
	}
}

// zoomOut decreases the zoom level
func (m *PreviewPaneModel) zoomOut() {
	if m.zoomLevel > 0.5 {
		m.zoomLevel -= 0.1
		// Note: Status bar update will be handled by parent component
	}
}

// resetZoom resets zoom to 100%
func (m *PreviewPaneModel) resetZoom() {
	m.zoomLevel = 1.0
	// Note: Status bar update will be handled by parent component
}

// renderContent renders the content with basic formatting (legacy method for compatibility)
func (m *PreviewPaneModel) renderContent() string {
	content, _ := m.renderContentWithGlamour()
	return content
}

// formatDuration formats a duration for display
func (m *PreviewPaneModel) formatDuration(d time.Duration) string {
	if d.Hours() >= 1 {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	return fmt.Sprintf("%.0fs", d.Seconds())
}

// renderTOC renders the table of contents
func (m *PreviewPaneModel) renderTOC() string {
	if len(m.tocEntries) == 0 {
		return ""
	}

	var tocLines []string
	tocLines = append(tocLines, "Table of Contents:")

	for _, entry := range m.tocEntries {
		// Create indentation based on header level
		indent := strings.Repeat("  ", entry.Level-1)

		// Create clickable-style entry
		tocEntry := fmt.Sprintf("%s• %s", indent, entry.Title)
		tocLines = append(tocLines, tocEntry)
	}

	tocContent := strings.Join(tocLines, "\n")

	return lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Align(lipgloss.Left).
		Width(m.width-4).
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		Padding(0, 1).
		Render(tocContent)
}

// SetShortcutContext sets the keyboard shortcut context
func (m *PreviewPaneModel) SetShortcutContext(context KeyContext) {
	if m.shortcutManager != nil {
		m.shortcutManager.SetContext(context)
	}
}

// GetShortcutManager returns the shortcut manager for external access
func (m *PreviewPaneModel) GetShortcutManager() *ShortcutManager {
	return m.shortcutManager
}

// GetZoomLevel returns the current zoom level as a percentage
func (m *PreviewPaneModel) GetZoomLevel() int {
	return int(m.zoomLevel * 100)
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Real-time preview callback methods

// onContentUpdate handles content updates from the real-time preview manager
func (m *PreviewPaneModel) onContentUpdate(content string) {
	m.content = content

	// Maintain scroll position if enabled
	if m.scrollSyncEnabled {
		m.maintainScrollPosition(content)
	}

	// Clear render cache for new content
	m.clearRenderCache()
}

// onUpdateStart handles preview update start events
func (m *PreviewPaneModel) onUpdateStart() {
	// Show visual indicator that update is starting
	m.showRefresh = true
}

// onUpdateComplete handles preview update completion events
func (m *PreviewPaneModel) onUpdateComplete(duration time.Duration) {
	// Hide refresh indicator
	m.showRefresh = false

	// Update performance metrics if needed
	// Could be used for adaptive throttling based on performance
}

// onValidationError handles markdown validation errors
func (m *PreviewPaneModel) onValidationError(errors []ValidationError) {
	if len(errors) > 0 {
		// Show first validation error
		m.lastError = errors[0].Message
	} else {
		m.lastError = ""
	}
}

// onStatsUpdate handles preview statistics updates
func (m *PreviewPaneModel) onStatsUpdate(stats PreviewStats) {
	m.previewStats = stats

	// Update TOC if enabled
	if m.showTOC && m.realtimeManager.tocGenerator != nil {
		m.tocEntries = m.realtimeManager.tocGenerator.GenerateTOC(m.content)
	}
}

// maintainScrollPosition maintains scroll position during content updates
func (m *PreviewPaneModel) maintainScrollPosition(newContent string) {
	if m.realtimeManager == nil {
		return
	}

	// Get current content hash for position tracking
	currentHash := m.realtimeManager.hashContent(m.content)

	// Store current scroll position for current content
	if currentHash != "" {
		m.lastScrollPositions[currentHash] = m.scrollPos
	}

	// Try to restore scroll position for new content
	newHash := m.realtimeManager.hashContent(newContent)
	if storedPos, exists := m.lastScrollPositions[newHash]; exists {
		// Check if the stored position is still valid for the new content
		lines := strings.Split(m.renderContent(), "\n")
		maxScroll := len(lines) - (m.height - 6)
		if maxScroll < 0 {
			maxScroll = 0
		}

		if storedPos <= maxScroll {
			m.scrollPos = storedPos
		}
	}
}

// clearRenderCache clears the render cache when content changes significantly
func (m *PreviewPaneModel) clearRenderCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	m.renderCache = make(map[string]string)
}

// getCachedRender attempts to get a cached render for the given content hash
// TODO: Uncomment when caching functionality is needed
// func (m *PreviewPaneModel) getCachedRender(contentHash string) (string, bool) {
// 	m.cacheMutex.RLock()
// 	defer m.cacheMutex.RUnlock()
//
// 	cached, exists := m.renderCache[contentHash]
// 	return cached, exists
// }

// setCachedRender stores a rendered result in the cache
// TODO: Uncomment when caching functionality is needed
// func (m *PreviewPaneModel) setCachedRender(contentHash, rendered string) {
// 	m.cacheMutex.Lock()
// 	defer m.cacheMutex.Unlock()
//
// 	// Implement LRU-style cache management
// 	if len(m.renderCache) >= m.maxCacheSize {
// 		// Remove oldest entries (simple implementation)
// 		for key := range m.renderCache {
// 			delete(m.renderCache, key)
// 			break
// 		}
// 	}
//
// 	m.renderCache[contentHash] = rendered
// }

// SetContentImmediate sets content immediately without debouncing
func (m *PreviewPaneModel) SetContentImmediate(content string) {
	m.content = content
	m.clearRenderCache()
}

// GetRealtimeManager returns the real-time preview manager
func (m *PreviewPaneModel) GetRealtimeManager() *RealTimePreviewManager {
	return m.realtimeManager
}

// ToggleScrollSync toggles scroll position synchronization
func (m *PreviewPaneModel) ToggleScrollSync() {
	m.scrollSyncEnabled = !m.scrollSyncEnabled
}

// ToggleWordCount toggles word count display
func (m *PreviewPaneModel) ToggleWordCount() {
	m.showWordCount = !m.showWordCount
}

// ToggleReadingTime toggles reading time display
func (m *PreviewPaneModel) ToggleReadingTime() {
	m.showReadingTime = !m.showReadingTime
}

// ToggleTOC toggles table of contents display
func (m *PreviewPaneModel) ToggleTOC() {
	m.showTOC = !m.showTOC
	if m.showTOC && m.realtimeManager != nil && m.realtimeManager.tocGenerator != nil {
		m.tocEntries = m.realtimeManager.tocGenerator.GenerateTOC(m.content)
	}
}

// GetPreviewStats returns current preview statistics
func (m *PreviewPaneModel) GetPreviewStats() PreviewStats {
	return m.previewStats
}

// GetTOC returns the table of contents entries
func (m *PreviewPaneModel) GetTOC() []TOCEntry {
	return m.tocEntries
}

// renderContentWithCache renders content with caching for performance
// TODO: Uncomment when caching functionality is needed
// func (m *PreviewPaneModel) renderContentWithCache() (string, error) {
// 	contentHash := m.realtimeManager.hashContent(m.content)
//
// 	// Try to get from cache first
// 	if cached, exists := m.getCachedRender(contentHash); exists {
// 		return cached, nil
// 	}
//
// 	// Render content
// 	rendered, err := m.renderContentWithGlamour()
//
// 	// Cache the result if successful
// 	if err == nil && rendered != "" {
// 		m.setCachedRender(contentHash, rendered)
// 	}
//
// 	return rendered, err
// }
