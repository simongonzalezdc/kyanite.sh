# focus.sh Codebase Analysis Report

## 📊 Current State vs Target State Analysis

### Executive Summary
The current codebase is **partially aligned** with the enhancement document but has **critical gaps** in major features. Module naming and basic structure are correct, but the **journal feature is completely missing** and the **theme system requires major expansion**.

---

## ✅ **ALIGNED COMPONENTS**

### 1. Module Path & Build Configuration ✅
- **Current**: `github.com/kyanite/focus` ✅ (matches target)
- **Build scripts**: Output `focus.exe`/`focus` ✅
- **Directory structure**: `cmd/focus/main.go` ✅
- **Storage path**: `~/.focus/` ✅

### 2. Core Infrastructure ✅  
- **Task management**: add/list/complete/delete ✅
- **AI integration**: Ollama + OpenRouter fallback ✅
- **TUI system**: Bubble Tea implementation ✅
- **Configuration**: Viper + YAML ✅

---

## ❌ **CRITICAL GAPS**

### 1. Theme System - **MAJOR DEFICIENCY** 

**Current State:**
- Only **3 themes** (synthwave, light, plain)
- Simple color variables in `pkg/styles/themes.go`
- No theme registry structure

**Required State:**
- **13 Kyanite themes** with specific hex codes
- **Default: Amber Night** theme
- Complete theme registry with all color definitions

**Impact:** 🔴 **HIGH** - Core branding requirement

### 2. Journal Feature - **COMPLETELY MISSING**

**Required Components:**
- ❌ `models/journal.go` - Entry data structure
- ❌ `internal/journal/export.go` - syntax.sh export pipeline  
- ❌ CLI commands: `focus journal new/list/view/search/export`
- ❌ TUI journal screen with navigation
- ❌ Storage: `~/.focus/journal.json`
- ❌ Export to `~/syntax/imports/` directory

**Impact:** 🔴 **CRITICAL** - Major feature gap

### 3. Kyanite Suite Integration - **MISSING**

**Current Issues:**
- Help text mentions "cyberpunk" instead of Kyanite Suite
- No cross-tool references (noise.sh, syntax.sh, prism.sh, wave.sh)
- Missing suite positioning in documentation

**Impact:** 🟡 **MEDIUM** - Brand alignment issue

---

## ⚠️ **MINOR ISSUES**

### 1. Branding Language Inconsistencies
**File:** `internal/cli/root.go`
- Lines 36-42: Uses "cyberpunk", "synthwave" language
- Should reference Kyanite Suite instead

### 2. Old References (5 found)
- `main.go:35`: `todoMain` function
- `test/integration_test.go:14`: "todo_test" temp dir  
- `pkg/utils/storage.go:28`: "~/.neon/" migration reference
- `cmd/focus-mcp/main.go:101`: "crush-golangci-lint"
- `internal/cli/calendar.go.bak`: Old import paths

### 3. Test Naming
- Test files still use "todo" naming convention
- Should be updated to "focus"

---

## 📋 **IMPLEMENTATION PRIORITY MATRIX**

### 🔴 **CRITICAL (Blockers)**
1. **Journal Feature Implementation** (2-3 weeks)
2. **Theme System Expansion** (1-2 weeks)

### 🟡 **HIGH (Important)**  
3. **Help Text Updates** (2-3 days)
4. **Kyanite Suite Integration** (3-5 days)
5. **Old Reference Cleanup** (1-2 days)

### 🟢 **MEDIUM (Polish)**
6. **Test Naming Updates** (1 day)
7. **Documentation Updates** (2-3 days)

---

## 🛠 **DETAILED NEXT STEPS**

### **Phase 1: Theme System Expansion (Week 1-2)**

**1.1 Create Theme Registry**
```bash
# Create new file
touch internal/theme/registry.go
```

**Required Content (from document):**
- All 13 theme definitions with exact hex codes
- Amber Night as default theme
- Theme switching with Ctrl+T functionality

**1.2 Update Theme Integration**
- Replace current 3-theme system with 13-theme registry
- Update `pkg/styles/themes.go` to use registry
- Ensure theme persistence across restarts

### **Phase 2: Journal Feature Implementation (Week 2-3)**

**2.1 Data Model**
```bash
# Create journal structure
mkdir -p internal/journal
touch models/journal.go
touch internal/journal/export.go
```

**2.2 CLI Commands**
- Add `internal/cli/journal.go` with all subcommands
- Implement new/list/view/search/export functionality
- Storage in `~/.focus/journal.json`

**2.3 Export Pipeline**
- Implement `ExportToSyntax()` function
- Support 4 export types: character, dialogue, scene, research
- Output to `~/syntax/imports/` directory

**2.4 TUI Integration**
- Create `internal/tui/journal_view.go`
- Add journal navigation to main TUI
- Implement shortcuts: J, ←/→, N, S, E, M

### **Phase 3: Branding Alignment (Week 3)**

**3.1 Help Text Updates**
- Update `internal/cli/root.go` lines 36-47
- Replace "cyberpunk" with Kyanite Suite references
- Add cross-tool awareness

**3.2 Clean Old References**
- Fix 5 remaining old brand references
- Update test naming from "todo" to "focus"
- Remove backup files

**3.3 Documentation**
- Update README.md with focus.sh branding
- Add journal feature documentation
- Document all 13 themes

---

## 📊 **EFFORT ESTIMATION**

| Feature | Estimated Time | Priority |
|---------|----------------|----------|
| Theme System (13 themes) | 5-7 days | 🔴 Critical |
| Journal Data Model | 2-3 days | 🔴 Critical |
| Journal CLI Commands | 3-4 days | 🔴 Critical |
| Journal Export Pipeline | 2-3 days | 🔴 Critical |
| Journal TUI Screen | 3-4 days | 🔴 Critical |
| Help Text Updates | 1-2 days | 🟡 High |
| Old Reference Cleanup | 1 day | 🟡 High |
| Documentation Updates | 2-3 days | 🟢 Medium |

**Total Estimated Effort: 19-27 days (3-4 weeks)**

---

## 🎯 **SUCCESS CRITERIA CHECKLIST**

### Before Release v2.0:

**Theme System:**
- [ ] All 13 Kyanite themes implemented
- [ ] Amber Night is default theme
- [ ] Theme switching with Ctrl+T works
- [ ] Theme persistence across restarts

**Journal Feature:**
- [ ] `focus journal new` creates entry
- [ ] `focus journal list` shows entries  
- [ ] `focus journal search` finds entries
- [ ] `focus journal export --to-syntax` works
- [ ] Journal TUI screen functional
- [ ] Export creates files in `~/syntax/imports/`

**Branding:**
- [ ] All help text references Kyanite Suite
- [ ] Zero old brand references in user-facing text
- [ ] Cross-tool integration mentioned

**Quality:**
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] All themes load without error
- [ ] Migration from v1 data works

---

## 🚨 **IMMEDIATE ACTION REQUIRED**

**Blockers for Release:**
1. **Journal feature is 100% missing** - Start with data model
2. **Theme system needs 10 additional themes** - Use hex codes from document

**Recommended Next Steps:**
1. **Week 1**: Implement complete theme system with all 13 themes
2. **Week 2-3**: Build entire journal feature
3. **Week 4**: Branding alignment and polish

**Estimated Release:** 4 weeks from start of implementation

---

## 📝 **NOTES**

- **Migration logic exists** ✅ in `pkg/utils/storage.go`
- **Core features work** ✅ (tasks, AI, TUI)
- **Module naming correct** ✅
- **Build system ready** ✅

The foundation is solid, but **two major feature gaps** prevent release readiness.