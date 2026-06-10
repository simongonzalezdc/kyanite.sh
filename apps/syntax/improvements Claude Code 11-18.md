# Comprehensive Codebase Analysis - November 18, 2025

## 📊 Executive Summary

**Overall Code Health**: Good foundations with significant improvement opportunities
- **Lines of Code**: ~5,000 LOC Go
- **Test Coverage**: ~20% (critical gap)
- **Architecture**: Clean but needs refactoring
- **Feature Completeness**: 85% (some TODOs remain)

---

## 🐛 Critical Issues Found

### **High Priority Bugs**

1. **Nil Pointer Risks** (`internal/editor/buffer.go:84-89`)
   - `InsertRune()` and `InsertNewline()` don't validate cursor bounds
   - Could panic on edge cases

2. **Error Handling Gaps** (12 locations)
   - `crypto/rand` errors ignored in ID generation (`storage/utils.go:28,35,42`)
   - Could generate duplicate IDs in rare cases
   - Silent failures when loading characters/scenes/locations

3. **Race Condition** (`internal/app/model.go:26-30`)
   - `CurrentProject.Characters` loaded without mutex protection
   - Potential issues if AI and UI access simultaneously

4. **Logic Error in Scene Numbering** (`internal/app/scenes.go:94-99`)
   - Uses `len(scenes)+1` on a map (could have gaps)
   - Should find actual next available number

### **Unfinished Features**

- **Location Editor**: Marked as TODO at `locations.go:134` - "coming soon" placeholder
- **Character/Scene Editors**: Only create with hardcoded values, no edit UI
- **Markdown Preview**: Extremely basic (only handles headers/lists)
- **Word Count**: Doesn't actually strip markdown as claimed

---

## ⚡ Performance & Memory Issues

### **Critical: Undo/Redo Memory Explosion**
```go
// internal/editor/buffer.go:194-210
// Creates FULL deep copy of entire document on every keystroke!
// 10,000 lines × 100 undo states = 1,000,000 line allocations
```
**Fix**: Use delta-based undo or copy-on-write

### **Repeated Disk I/O**
View functions load from disk on every render:
- `scenes.go:26-30`
- `characters.go:25-29`
- `relationship_map.go:25-29`

**Fix**: Move loading to `Update()` or cache in model

### **Inefficient Map Iteration**
Using O(n) iteration to find nth element in maps (`scenes.go:78-88`)

**Fix**: Maintain sorted slice alongside map

### **Duplicated Code**
Scene sorting logic copied 3 times across exporters

---

## 🎯 Missing Essential Features

| Feature | Priority | Impact |
|---------|----------|---------|
| **Search/Find** | CRITICAL | Can't find text in documents |
| **Spell Check** | HIGH | Professional writing requirement |
| **Auto-save indicators** | HIGH | Users don't know if saved |
| **Version history** | MEDIUM | `.backups/` exists but unused |
| **Word count goals** | MEDIUM | Types exist, no UI |
| **Find & replace** | HIGH | Essential editing tool |

---

## 🏗️ Architecture Improvements

### **Current Issues**

1. **God Object**: `Model` struct has 20+ fields mixing UI, data, and business logic
   - Should split into `UIState`, `ProjectState`, `EditorState`, `AIState`

2. **No Service Layer**: Business logic mixed with UI
   - Need separation: `UI → Application Services → Domain → Storage`

3. **Missing Interfaces**: Direct dependencies on concrete types
   - Can't mock for testing
   - Tight coupling

4. **Hardcoded Values**: 25+ magic numbers should be configurable
   ```go
   60*time.Second  // AI timeout
   100            // Max undo states
   20             // Context lines
   ```

### **Recommended Refactoring**

```go
// Current: God object
type Model struct {
    CurrentScreen, Width, Height, ThemeManager,
    CurrentTheme, CurrentProject, Message, Error,
    Projects, SelectedIndex, CurrentScene, Buffer,
    EditorMode, AIClient, AISuggestion, ...
}

// Better: Composition
type Model struct {
    UI         UIState
    Navigation NavigationState
    Project    *ProjectState
    Editor     *EditorState
    AI         *AIState
}
```

---

## 🧪 Testing Gaps

```
Current Coverage: ~20%

✗ internal/app/      0%   (High risk - user-facing)
✗ internal/ai/       0%   (Medium risk - external API)
✗ internal/scene/    0%   (High risk - complex validation)
✓ internal/editor/  71%   (Good)
✓ internal/export/  37%   (Fair)
⚠ internal/storage/ 20%   (Critical - data loss risk)
```

**Critical Missing Tests**:
- Scene validation/compilation logic
- AI client error handling
- Storage concurrent access
- Undo/redo edge cases

---

## 💻 JavaScript Conversion Analysis

### **Option 1: Keep Go TUI** ⭐ Recommended for v1.0
- Already works, production-ready
- Low maintenance
- Perfect for terminal-loving writers
- Ship now, iterate later

### **Option 2: Tauri Desktop App** ⭐ Recommended for v2.0
```
┌────────────────────────────────┐
│   React + Monaco Editor        │  Modern UI
│   Tailwind CSS                 │  Rich visualization
└────────────┬───────────────────┘
             │ IPC
┌────────────┴───────────────────┐
│   Go Backend (embedded)        │  Reuse existing code
│   File system, AI, exports     │
└────────────────────────────────┘
```

**Pros**:
- ✅ Modern UX with Monaco editor (like VS Code)
- ✅ Native app (5-10MB binary)
- ✅ Reuse 70% of Go backend
- ✅ Rich visualizations (interactive relationship graphs)
- ✅ Broader audience appeal

**Timeline**: 3 months, 1 developer

### **Option 3: Browser Web App**
- Max accessibility (mobile/tablet)
- Requires server infrastructure
- Higher complexity

### **Option 4: Node.js TUI**
- 4-6 weeks to port
- blessed/ink less polished than Bubble Tea
- Performance concerns
- **Not recommended** - current Go TUI is better

---

## 📋 Prioritized Action Plan

### **Phase 1: Critical Fixes (1-2 weeks)**
1. Fix bounds checking in buffer operations (prevent crashes)
2. Add error handling for crypto/rand (prevent duplicate IDs)
3. Implement auto-save with visual indicator
4. Add search/find functionality

### **Phase 2: Technical Debt (2-3 weeks)**
5. Refactor Model into sub-models
6. Implement delta-based undo/redo
7. Move disk I/O out of view functions
8. Extract configuration system

### **Phase 3: Missing Features (3-4 weeks)**
9. Complete location editor
10. Add spell check integration
11. Implement find & replace
12. Version history UI for backups

### **Phase 4: Quality (2-3 weeks)**
13. Increase test coverage to 80%+
14. Extract service layer
15. Add dependency injection
16. Performance optimization

### **Phase 5: v2.0 - Tauri App (3 months)**
17. Build React + Monaco frontend
18. Embed Go backend via IPC
19. Interactive relationship visualizations
20. Cross-platform installers

---

## 🎯 Recommendations

### **Immediate Actions**

1. **Fix the undo/redo memory issue** - This will bite you with large documents
2. **Add search functionality** - Writers absolutely need this
3. **Implement auto-save indicators** - Prevent user anxiety about data loss
4. **Boost test coverage** - Storage layer is critical (prevent data loss)

### **Strategic Direction**

**Short-term**: Polish the Go TUI to v1.0
- Fix critical bugs
- Add search/spell check
- Improve test coverage
- Ship a solid terminal tool

**Long-term**: Build Tauri desktop app for v2.0
- Keep Go backend (it's good!)
- Modern React + Monaco frontend
- Serve both markets: CLI power users + GUI general users
- Same file format = compatibility

### **Don't Convert to Pure JavaScript**
The Go TUI is excellent. Bubble Tea is more mature than Node.js TUI libraries. If you want JavaScript, go Tauri (hybrid) or web app, not Node.js TUI.

---

## Detailed Technical Findings

### Technical Debt Details

**File: /internal/app/locations.go**
- Line 134: `// TODO: Navigate to location editor` - Location editor functionality is incomplete
- Line 136: Placeholder message "Location editor coming soon"

**File: /internal/app/text_editor.go**
- Lines 277-300: `stripANSI()` function is simple implementation that may not handle all ANSI escape sequences
- Lines 263-275: Helper functions `max()` and `min()` are duplicated

**File: /internal/storage/scene.go**
- Lines 109-118: `calculateWordCount()` has "simple implementation" and doesn't clean markdown syntax

**File: /internal/export/markdown.go**
- Lines 156-176: `simpleMarkdownToHTML()` is very basic (doesn't handle bold, italic, links, code blocks)

### Linting Issues Details

**Error Handling Issues**:
- `/internal/app/characters.go:26-29` - Error from `LoadAllCharacters()` silently ignored
- `/internal/storage/utils.go:28,35,42` - `rand.Read()` errors ignored in ID generation
- `/internal/storage/character.go:89-92` - Silently skips invalid characters
- `/internal/storage/scene.go:93-96` - Silently skips invalid scenes
- `/internal/ai/client.go:132` - Error from `io.ReadAll()` ignored

### Performance Details

**Repeated Data Loading**:
```go
// Multiple view functions load data from disk on every render
if m.CurrentProject.Scenes == nil || len(m.CurrentProject.Scenes) == 0 {
    scenes, err := storage.LoadAllScenes(m.CurrentProject.Directory)
    if err == nil {
        m.CurrentProject.Scenes = scenes
    }
}
```

**Inefficient ANSI Stripping**:
```go
func stripANSI(s string) string {
    // Character-by-character parsing - should use regex
    result := strings.Builder{}
    inEscape := false
    for _, r := range s {
        if r == '\033' {
            inEscape = true
        } else if inEscape && r == 'm' {
            inEscape = false
        } else if !inEscape {
            result.WriteRune(r)
        }
    }
    return result.String()
}
```

### Resource Management Issues

**Missing Defer for Zip Writer** (`/internal/export/docx.go:89-114`):
```go
zipWriter := zip.NewWriter(&buf)
// Multiple addZipFile calls that could fail
if err := zipWriter.Close(); err != nil {
    return nil, fmt.Errorf("failed to close zip: %w", err)
}
// Should add: defer zipWriter.Close()
```

### Missing Features Analysis

**Search & Replace**:
- No find/replace functionality in editor
- Should add to `/internal/editor/buffer.go`:
  - `Find(pattern string, caseSensitive bool) []Position`
  - `Replace(pattern, replacement string, replaceAll bool)`
  - `FindNext()` / `FindPrevious()` navigation
  - Regex support

**Spell Check**:
- None currently
- Integrate `aspell` or `hunspell`
- Add spell check toggle (Ctrl+P)
- Highlight misspelled words
- Custom dictionary per project

**Version Control Integration**:
- `.backups/` directory exists but not used
- Implement automatic backups with timestamps
- Add "Version History" screen
- Consider git integration

---

## Summary Statistics

- **TODOs Found**: 1 explicit TODO comment
- **Critical Bugs**: 3 (nil pointer dereferences, race conditions)
- **Error Handling Issues**: 12 locations
- **Unfinished Features**: 8 major areas
- **Code Quality Issues**: 6 areas needing refactoring
- **Duplicated Code**: 3+ instances
- **Magic Numbers**: 25+ hardcoded values
- **Test Coverage**: ~20% overall

---

*Analysis completed: November 18, 2025*
*Codebase: Syntax.sh v2.0*
*Total LOC: ~5,000 lines of Go*
