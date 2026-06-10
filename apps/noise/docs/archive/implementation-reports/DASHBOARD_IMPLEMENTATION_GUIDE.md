# noise.sh Dashboard Implementation Guide

## Technical Implementation Details

This guide provides the technical specifications and code structure for implementing the impressive dashboard design.

## Core Architecture

### 1. Dashboard Model Structure

```go
// internal/ui/dashboard/dashboard.go

package dashboard

import (
    "github.com/kyanite/noise/internal/theme"
    "github.com/kyanite/noise/internal/ui"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// DashboardModel represents the main dashboard
type DashboardModel struct {
    // Layout and sizing
    width     int
    height    int
    
    // Theme management
    currentTheme theme.Theme
    themeManager *ThemePreviewModel
    
    // Component panels
    quickActions  *QuickActionsModel
    recentWork    *RecentWorkModel
    musicTools    *MusicToolsModel
    aiAssistant   *AIAssistantModel
    systemInfo    *SystemInfoModel
    
    // UI state
    focusedPanel string
    animations   *ui.AnimationManager
    
    // Header and footer
    header *HeaderModel
    footer *FooterModel
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel() *DashboardModel {
    return &DashboardModel{
        currentTheme:  theme.GetManager().Current(),
        themeManager:  NewThemePreviewModel(),
        quickActions:  NewQuickActionsModel(),
        recentWork:    NewRecentWorkModel(),
        musicTools:    NewMusicToolsModel(),
        aiAssistant:   NewAIAssistantModel(),
        systemInfo:    NewSystemInfoModel(),
        animations:    ui.NewAnimationManager(),
        focusedPanel:  "quickactions",
        header:        NewHeaderModel(),
        footer:        NewFooterModel(),
    }
}
```

### 2. Theme Preview Component

```go
// internal/ui/dashboard/theme_preview.go

package dashboard

import (
    "fmt"
    "github.com/kyanite/noise/internal/theme"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// ThemePreviewModel showcases the current theme
type ThemePreviewModel struct {
    width       int
    height      int
    currentIdx  int
    themes      []theme.Theme
    animation   *ui.AnimationManager
    showDetails bool
}

// NewThemePreviewModel creates a new theme preview
func NewThemePreviewModel() *ThemePreviewModel {
    themeIDs := theme.ListThemes()
    themes := make([]theme.Theme, len(themeIDs))
    for i, id := range themeIDs {
        themes[i] = theme.GetTheme(id)
    }
    
    return &ThemePreviewModel{
        currentIdx: 0,
        themes:     themes,
        animation:  ui.NewAnimationManager(),
        showDetails: false,
    }
}

// View renders the theme preview
func (m *ThemePreviewModel) View() string {
    if m.width == 0 {
        return "Theme Preview"
    }
    
    currentTheme := m.themes[m.currentIdx]
    
    // Create theme showcase
    title := lipgloss.NewStyle().
        Foreground(currentTheme.Primary).
        Bold(true).
        Align(lipgloss.Center).
        Render(currentTheme.Name)
    
    // Color palette display
    palette := m.renderColorPalette(currentTheme)
    
    // Sample UI elements
    samples := m.renderSampleElements(currentTheme)
    
    // Navigation
    nav := m.renderNavigation()
    
    return lipgloss.JoinVertical(lipgloss.Center,
        title,
        "",
        palette,
        "",
        samples,
        "",
        nav,
    )
}

func (m *ThemePreviewModel) renderColorPalette(t theme.Theme) string {
    colors := []struct {
        name  string
        color lipgloss.Color
    }{
        {"Primary", t.Primary},
        {"Secondary", t.Secondary},
        {"Accent", t.Accent},
        {"Background", t.Background},
        {"Text", t.Text},
        {"Success", t.Success},
        {"Warning", t.Warning},
        {"Error", t.Error},
    }
    
    var rows []string
    for i := 0; i < len(colors); i += 4 {
        var row []string
        for j := i; j < i+4 && j < len(colors); j++ {
            colorBox := lipgloss.NewStyle().
                Background(colors[j].color).
                Foreground(t.Text).
                Width(8).
                Height(3).
                Align(lipgloss.Center).
                Bold(true).
                Render(colors[j].name[:3])
            row = append(row, colorBox)
        }
        rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, row...))
    }
    
    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *ThemePreviewModel) renderSampleElements(t theme.Theme) string {
    button := lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Primary).
        Padding(0, 2).
        Bold(true).
        Render("Primary Button")
    
    secondary := lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Secondary).
        Padding(0, 2).
        Render("Secondary")
    
    accent := lipgloss.NewStyle().
        Foreground(t.Background).
        Background(t.Accent).
        Padding(0, 2).
        Bold(true).
        Render("Accent")
    
    return lipgloss.JoinHorizontal(lipgloss.Left, button, " ", secondary, " ", accent)
}
```

### 3. Quick Actions Grid

```go
// internal/ui/dashboard/quick_actions.go

package dashboard

import (
    "github.com/kyanite/noise/internal/ui"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Action represents a quick action
type Action struct {
    ID          string
    Title       string
    Description string
    Icon        string
    Shortcut    string
    Screen      ui.Screen
    Animated    bool
}

// QuickActionsModel manages the quick actions grid
type QuickActionsModel struct {
    width     int
    height    int
    actions   []Action
    selected  int
    animation *ui.AnimationManager
}

// NewQuickActionsModel creates a new quick actions model
func NewQuickActionsModel() *QuickActionsModel {
    actions := []Action{
        {ID: "new", Title: "New Song", Description: "Create a new song", Icon: "🎵", Shortcut: "1", Screen: ui.ScreenEditor, Animated: true},
        {ID: "open", Title: "Open Project", Description: "Browse existing projects", Icon: "📁", Shortcut: "2", Screen: ui.ScreenManager, Animated: true},
        {ID: "ai", Title: "AI Brainstorm", Description: "Get AI assistance", Icon: "🤖", Shortcut: "3", Screen: ui.ScreenEditor, Animated: true},
        {ID: "export", Title: "Export", Description: "Export your work", Icon: "📤", Shortcut: "4", Screen: ui.ScreenExport, Animated: true},
        {ID: "theory", Title: "Theory Tools", Description: "Music theory reference", Icon: "🎼", Shortcut: "5", Screen: ui.ScreenTheory, Animated: true},
        {ID: "audio", Title: "Audio Tools", Description: "Metronome & playback", Icon: "🎧", Shortcut: "6", Screen: ui.ScreenAudio, Animated: true},
    }
    
    return &QuickActionsModel{
        actions:   actions,
        selected:  0,
        animation: ui.NewAnimationManager(),
    }
}

// View renders the quick actions grid
func (m *QuickActionsModel) View() string {
    if m.width == 0 {
        return "Quick Actions"
    }
    
    t := theme.GetManager().Current()
    
    // Calculate grid layout
    cols := 2
    if m.width > 100 {
        cols = 3
    }
    
    var rows []string
    for i := 0; i < len(m.actions); i += cols {
        var row []string
        for j := i; j < i+cols && j < len(m.actions); j++ {
            action := m.actions[j]
            card := m.renderActionCard(action, j == m.selected)
            row = append(row, card)
        }
        rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, row...))
    }
    
    title := lipgloss.NewStyle().
        Foreground(t.Primary).
        Bold(true).
        Render("Quick Actions")
    
    return lipgloss.JoinVertical(lipgloss.Left,
        title,
        "",
        lipgloss.JoinVertical(lipgloss.Left, rows...),
    )
}

func (m *QuickActionsModel) renderActionCard(action Action, selected bool) string {
    t := theme.GetManager().Current()
    
    // Base card style
    baseStyle := lipgloss.NewStyle().
        Width(m.width/3-4).
        Height(6).
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Margin(1)
    
    if selected {
        baseStyle = baseStyle.
            BorderForeground(t.Primary).
            Background(lipgloss.Color(string(t.Primary))).
            Foreground(t.Background)
    } else {
        baseStyle = baseStyle.
            BorderForeground(t.Secondary).
            Background(t.Background).
            Foreground(t.Text)
    }
    
    // Apply animation if active
    if action.Animated {
        progress := m.animation.GetAnimationProgress("hover_" + action.ID)
        if progress > 0 && progress < 1 {
            // Pulsing effect
            intensity := 0.8 + 0.2*progress
            if selected {
                baseStyle = baseStyle.
                    Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", 
                        uint(intensity*uint8(t.Primary>>16)), 
                        uint(intensity*uint8(t.Primary>>8)), 
                        uint(intensity*uint8(t.Primary)))))
            }
        }
    }
    
    content := lipgloss.JoinVertical(lipgloss.Left,
        lipgloss.NewStyle().Bold(true).Render(action.Icon + " " + action.Title),
        "",
        lipgloss.NewStyle().Faint(true).Render(action.Description),
        "",
        lipgloss.NewStyle().Align(lipgloss.Right).Render("[" + action.Shortcut + "]"),
    )
    
    return baseStyle.Render(content)
}
```

### 4. Animation System Extensions

```go
// internal/ui/animations/dashboard_animations.go

package animations

import (
    "github.com/kyanite/noise/internal/ui"
    "github.com/charmbracelet/lipgloss"
)

// Dashboard-specific animations
type DashboardAnimationManager struct {
    *ui.AnimationManager
}

// NewDashboardAnimationManager creates a dashboard-specific animation manager
func NewDashboardAnimationManager() *DashboardAnimationManager {
    return &DashboardAnimationManager{
        AnimationManager: ui.NewAnimationManager(),
    }
}

// StaggeredEntrance creates a staggered entrance animation for dashboard panels
func (dam *DashboardAnimationManager) StaggeredEntrance(panelIDs []string) {
    for i, id := range panelIDs {
        go func(panelID string, delay int) {
            time.Sleep(time.Duration(delay) * 100 * time.Millisecond)
            dam.SlideTransition(panelID+"_entrance", 1.0)
        }(id, i)
    }
}

// ThemeTransition creates a smooth theme transition animation
func (dam *DashboardAnimationManager) ThemeTransition(fromTheme, toTheme theme.Theme) {
    // Animate color changes
    dam.FadeTransition("theme_primary", 1.0)
    dam.FadeTransition("theme_secondary", 1.0)
    dam.FadeTransition("theme_accent", 1.0)
    dam.FadeTransition("theme_background", 1.0)
}

// PulsePanel creates a pulsing effect for a panel
func (dam *DashboardAnimationManager) PulsePanel(panelID string) {
    dam.PulseAnimation(panelID+"_pulse", 1.0)
}

// ApplyThemeTransition applies theme transition effects to a style
func ApplyThemeTransition(base lipgloss.Style, progress float64, fromTheme, toTheme theme.Theme) lipgloss.Style {
    if progress <= 0 {
        return base
    }
    if progress >= 1 {
        return base
    }
    
    // Interpolate colors
    interpolateColor := func(from, to lipgloss.Color) lipgloss.Color {
        // Simple color interpolation (would need proper RGB extraction)
        return to // Simplified for example
    }
    
    return base.
        Foreground(interpolateColor(fromTheme.Text, toTheme.Text)).
        Background(interpolateColor(fromTheme.Background, toTheme.Background))
}
```

### 5. Responsive Layout System

```go
// internal/ui/dashboard/responsive_layout.go

package dashboard

import (
    "github.com/charmbracelet/lipgloss"
)

// LayoutConfig defines the layout configuration for different screen sizes
type LayoutConfig struct {
    PanelWidths  map[string]int
    PanelHeights map[string]int
    GridCols     int
    GridRows     int
    ShowPanels   map[string]bool
}

// GetLayoutConfig returns the appropriate layout configuration for the given screen size
func GetLayoutConfig(width, height int) LayoutConfig {
    switch {
    case width >= 120 && height >= 30:
        // Large terminals - full 3x3 grid
        return LayoutConfig{
            PanelWidths: map[string]int{
                "theme":     width / 3,
                "actions":   width / 3,
                "recent":    width / 3,
                "tools":     width / 3,
                "ai":        width / 3,
                "info":      width / 3,
            },
            PanelHeights: map[string]int{
                "theme":     height / 3,
                "actions":   height / 3,
                "recent":    height / 3,
                "tools":     height / 3,
                "ai":        height / 3,
                "info":      height / 3,
            },
            GridCols:   3,
            GridRows:   3,
            ShowPanels: map[string]bool{
                "theme": true, "actions": true, "recent": true,
                "tools": true, "ai": true, "info": true,
            },
        }
        
    case width >= 80 && height >= 24:
        // Medium terminals - 2x3 grid
        return LayoutConfig{
            PanelWidths: map[string]int{
                "theme":     width / 2,
                "actions":   width / 2,
                "recent":    width / 2,
                "tools":     width / 2,
                "ai":        width / 2,
                "info":      width / 2,
            },
            PanelHeights: map[string]int{
                "theme":     height / 3,
                "actions":   height / 3,
                "recent":    height / 3,
                "tools":     height / 3,
                "ai":        height / 3,
                "info":      height / 3,
            },
            GridCols:   2,
            GridRows:   3,
            ShowPanels: map[string]bool{
                "theme": true, "actions": true, "recent": true,
                "tools": true, "ai": true, "info": true,
            },
        }
        
    case width >= 60 && height >= 20:
        // Small terminals - single column with tabs
        return LayoutConfig{
            PanelWidths: map[string]int{
                "active": width - 4,
            },
            PanelHeights: map[string]int{
                "active": height - 8,
            },
            GridCols:   1,
            GridRows:   1,
            ShowPanels: map[string]bool{
                "active": true,
            },
        }
        
    default:
        // Minimal terminals - essential features only
        return LayoutConfig{
            PanelWidths: map[string]int{
                "menu": width - 2,
            },
            PanelHeights: map[string]int{
                "menu": height - 4,
            },
            GridCols:   1,
            GridRows:   1,
            ShowPanels: map[string]bool{
                "menu": true,
            },
        }
    }
}

// RenderLayout renders the dashboard layout based on configuration
func (dm *DashboardModel) RenderLayout() string {
    config := GetLayoutConfig(dm.width, dm.height)
    t := theme.GetManager().Current()
    
    // Header
    header := dm.header.View()
    
    // Main content based on layout
    var content string
    switch config.GridCols {
    case 3:
        // Full 3x3 grid
        topRow := lipgloss.JoinHorizontal(lipgloss.Top,
            dm.themeManager.View(),
            dm.quickActions.View(),
            dm.recentWork.View(),
        )
        middleRow := lipgloss.JoinHorizontal(lipgloss.Top,
            dm.musicTools.View(),
            dm.aiAssistant.View(),
            dm.systemInfo.View(),
        )
        content = lipgloss.JoinVertical(lipgloss.Left, topRow, middleRow)
        
    case 2:
        // 2x3 grid
        topRow := lipgloss.JoinHorizontal(lipgloss.Top,
            dm.themeManager.View(),
            dm.quickActions.View(),
        )
        middleRow := lipgloss.JoinHorizontal(lipgloss.Top,
            dm.recentWork.View(),
            dm.musicTools.View(),
        )
        bottomRow := lipgloss.JoinHorizontal(lipgloss.Top,
            dm.aiAssistant.View(),
            dm.systemInfo.View(),
        )
        content = lipgloss.JoinVertical(lipgloss.Left, topRow, middleRow, bottomRow)
        
    case 1:
        // Single column with active panel
        content = dm.renderActivePanel()
        
    default:
        // Minimal menu
        content = dm.renderMinimalMenu()
    }
    
    // Footer
    footer := dm.footer.View()
    
    return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}
```

## Integration with Existing System

### 1. Root Model Updates

```go
// internal/ui/root.go - Add dashboard integration

// Add screenDashboard to the screen enum
const (
    screenSplash screen = iota
    screenDashboard  // New dashboard screen
    screenMenu
    // ... other screens
)

// Update RootModel to include dashboard
type RootModel struct {
    // ... existing fields
    dashboard *DashboardModel
}

// Update initialization to start with dashboard
func (m *RootModel) initializeChildModels() {
    // ... existing initialization
    
    m.dashboard = NewDashboardModel()
    
    // ... rest of initialization
}

// Update screen routing
func (m *RootModel) updateCurrentScreen(msg tea.Msg) tea.Cmd {
    switch m.currentScreen {
    case screenDashboard:
        return m.updateDashboard(msg)
    // ... other screen cases
    }
    return nil
}

// Add dashboard update method
func (m *RootModel) updateDashboard(msg tea.Msg) tea.Cmd {
    if m.dashboard != nil {
        return m.dashboard.Update(msg)
    }
    return nil
}

// Update view to render dashboard
func (m *RootModel) View() string {
    switch m.currentScreen {
    case screenDashboard:
        return m.dashboard.View()
    // ... other screen cases
    }
    return "Unknown screen"
}

// Update initialization to start with dashboard instead of menu
case initSuccessMsg:
    // ... existing code
    if m.quickStartConfig != nil {
        m.currentScreen = screenEditor
    } else {
        m.currentScreen = screenDashboard  // Changed from screenMenu
    }
    // ... rest of initialization
```

### 2. Theme Manager Integration

```go
// internal/theme/manager.go - Add dashboard integration

// Add dashboard notification method
func (m *Manager) NotifyDashboard() {
    // This would be called to notify the dashboard of theme changes
    // Implementation depends on message passing system
}

// Add theme preview support
func (m *Manager) GetThemePreview(id string) ThemePreview {
    theme := GetTheme(id)
    return ThemePreview{
        Name:    theme.Name,
        Colors:  extractColors(theme),
        Samples: generateSamples(theme),
    }
}

type ThemePreview struct {
    Name    string
    Colors  []string
    Samples []SampleElement
}

type SampleElement struct {
    Type  string
    Style lipgloss.Style
    Text  string
}
```

## Performance Optimizations

### 1. Lazy Loading Strategy

```go
// internal/ui/dashboard/lazy_loading.go

package dashboard

type LazyPanel struct {
    PanelID    string
    Loaded     bool
    Loader     func() ui.Model
    Model      ui.Model
    Loading    bool
    Error      error
}

func (lp *LazyPanel) EnsureLoaded() {
    if !lp.Loaded && !lp.Loading {
        lp.Loading = true
        go func() {
            defer func() {
                lp.Loading = false
            }()
            lp.Model = lp.Loader()
            lp.Loaded = true
        }()
    }
}

func (lp *LazyPanel) View() string {
    if lp.Loading {
        return "Loading..."
    }
    if lp.Error != nil {
        return "Error: " + lp.Error.Error()
    }
    if lp.Loaded {
        return lp.Model.View()
    }
    return ""
}
```

### 2. Render Optimization

```go
// internal/ui/dashboard/render_cache.go

package dashboard

type RenderCache struct {
    cache map[string]CachedRender
    mutex sync.RWMutex
}

type CachedRender struct {
    Content string
    Time    time.Time
    Valid   bool
}

func (rc *RenderCache) Get(key string, maxAge time.Duration) (string, bool) {
    rc.mutex.RLock()
    defer rc.mutex.RUnlock()
    
    if cached, exists := rc.cache[key]; exists {
        if time.Since(cached.Time) < maxAge {
            return cached.Content, true
        }
    }
    return "", false
}

func (rc *RenderCache) Set(key string, content string) {
    rc.mutex.Lock()
    defer rc.mutex.Unlock()
    
    rc.cache[key] = CachedRender{
        Content: content,
        Time:    time.Now(),
        Valid:   true,
    }
}
```

## Testing Strategy

### 1. Component Tests

```go
// internal/ui/dashboard/dashboard_test.go

package dashboard

import (
    "testing"
    "github.com/kyanite/noise/internal/theme"
)

func TestDashboardModel(t *testing.T) {
    model := NewDashboardModel()
    
    // Test initialization
    if model == nil {
        t.Fatal("NewDashboardModel() returned nil")
    }
    
    // Test theme integration
    if model.currentTheme.Name == "" {
        t.Error("Current theme not initialized")
    }
    
    // Test component initialization
    if model.themeManager == nil {
        t.Error("Theme manager not initialized")
    }
    if model.quickActions == nil {
        t.Error("Quick actions not initialized")
    }
}

func TestThemePreview(t *testing.T) {
    model := NewThemePreviewModel()
    
    // Test theme loading
    if len(model.themes) != 10 {
        t.Errorf("Expected 10 themes, got %d", len(model.themes))
    }
    
    // Test theme switching
    oldIdx := model.currentIdx
    model.NextTheme()
    if model.currentIdx == oldIdx {
        t.Error("Theme index did not change")
    }
}

func TestResponsiveLayout(t *testing.T) {
    tests := []struct {
        width  int
        height int
        cols   int
        rows   int
    }{
        {120, 30, 3, 3},
        {80, 24, 2, 3},
        {60, 20, 1, 1},
        {40, 15, 1, 1},
    }
    
    for _, tt := range tests {
        config := GetLayoutConfig(tt.width, tt.height)
        if config.GridCols != tt.cols {
            t.Errorf("Width %d: expected %d cols, got %d", tt.width, tt.cols, config.GridCols)
        }
        if config.GridRows != tt.rows {
            t.Errorf("Height %d: expected %d rows, got %d", tt.height, tt.rows, config.GridRows)
        }
    }
}
```

### 2. Integration Tests

```go
// integration/dashboard_integration_test.go

package integration

import (
    "testing"
    "github.com/kyanite/noise/internal/ui"
    "github.com/kyanite/noise/internal/ui/dashboard"
)

func TestDashboardIntegration(t *testing.T) {
    // Test dashboard integration with root model
    rootModel := ui.NewRootModel(nil)
    
    // Initialize dashboard
    rootModel.initializeChildModels()
    
    // Test dashboard is accessible
    if rootModel.dashboard == nil {
        t.Error("Dashboard not initialized in root model")
    }
    
    // Test theme switching
    originalTheme := theme.GetManager().Current()
    theme.GetManager().Next()
    newTheme := theme.GetManager().Current()
    
    if originalTheme.Name == newTheme.Name {
        t.Error("Theme did not change")
    }
    
    // Test dashboard responds to theme changes
    if rootModel.dashboard.currentTheme.Name != newTheme.Name {
        t.Error("Dashboard did not update to new theme")
    }
}
```

## Deployment Strategy

### Phase 1: Core Infrastructure (Week 1-2)
1. Implement basic dashboard model structure
2. Create responsive layout system
3. Integrate with existing root model
4. Basic theme integration

### Phase 2: Visual Components (Week 3-4)
1. Implement theme preview component
2. Create quick actions grid
3. Add animation system integration
4. Basic visual polish

### Phase 3: Feature Integration (Week 5-6)
1. Connect all existing features
2. Implement recent work panel
3. Add AI assistant integration
4. Create system info display

### Phase 4: Polish & Optimization (Week 7-8)
1. Advanced animations
2. Performance optimizations
3. Accessibility features
4. Testing and bug fixes

This implementation guide provides the technical foundation for creating the impressive dashboard that showcases the Kyanite theme system while providing immediate access to all application features.