# Phase 0 Audit: noise.sh Codebase Analysis

## Overview
This audit identifies the features that actually exist in the noise.sh codebase so that Enhancement&nbsp;#6.5 can build on verified functionality instead of assumptions.

## Audit Status
- **Started**: 2025-10-18
- **Completed**: 2025-10-18
- **Scope**: Editor experience, AI helpers, theory tools, theming, knowledge base expectations

---

## Executive Summary
- The editor surface in `internal/ui/editor/` is notably complete: split panes, multiple layout modes, auxiliary tools (BPM tapper, chord picker, status bar) and rich shortcut coverage are already present and wired together.
- Rhyme and syllable tooling exists but is static: a hard-coded rhyme dictionary and heuristic syllable counter power the “Rhymes” tab, with no external knowledge base integration.
- All four AI agents are implemented only as placeholders that return canned responses; no LLM or RAG wiring is in place.
- Theme infrastructure is production-ready with 12 registered themes and runtime switching support.
- No general knowledge base implementation was located beyond the chord progression dataset; RAG/KB work remains outstanding.

---

## Detailed Findings

### 1. Editor Components and Functionality — **Status: Exists (Functional)**
**Key files**
- [`editor_pane.go`](internal/ui/editor/editor_pane.go)
- [`split_pane.go`](internal/ui/editor/split_pane.go)
- [`preview_pane.go`](internal/ui/editor/preview_pane.go)
- [`status_bar.go`](internal/ui/editor/status_bar.go)
- [`layout_sketch.go`](internal/ui/editor/layout_sketch.go), [`layout_draft.go`](internal/ui/editor/layout_draft.go), [`layout_polish.go`](internal/ui/editor/layout_polish.go)
- [`bpm_tapper.go`](internal/ui/editor/bpm_tapper.go)
- [`chord_picker.go`](internal/ui/editor/chord_picker.go)
- [`help_pane.go`](internal/ui/editor/help_pane.go)

**Highlights**
- Split-pane shell (`split_pane.go`) orchestrates editor, preview, theory, AI overlays, and quick tools based on the active mode (Sketch/Draft/Polish).
- `editor_pane.go` layers syntax highlighting scaffolding, auto-save integration, AI assist overlays (continue/variation/brainstorm), and wiring for quick tools.
- Layout variants dynamically resize panes per prototyping mode, already matching Enhancement&nbsp;#6 goals.
- Supporting components (status bar, BPM tapper, in-editor chord picker, help pane) are implemented with styling and message plumbing.
- Tests exist for many subcomponents (`test/…_test.go` files) indicating active coverage.

### 2. Keyboard Shortcuts and Handling — **Status: Exists (Functional)**
**Key files**
- [`shortcuts.go`](internal/ui/editor/shortcuts.go)
- [`help_pane.go`](internal/ui/editor/help_pane.go)
- [`test/shortcuts_test.go`](test/shortcuts_test.go)

**Highlights**
- `ShortcutManager` centralizes bindings for navigation, editing, tools (chord picker, BPM tapper), theme switching, and application-level actions.
- Context-aware help overlay (`help_pane.go`) renders shortcuts with search/filtering and responsive layouts.
- Tests confirm shortcut routing and help presentation.

### 3. Rhyme & Syllable Tooling — **Status: Exists (Partial / Static)**
**Key files**
- [`theory.go`](internal/app/theory.go)
- [`theory.go`](internal/ui/theory.go)

**Highlights**
- `TheoryService` implements `FindRhymes`, `CountSyllables`, and `AnalyzeProsody` using a static dictionary and simple heuristics (no external pronunciation data).
- UI Theory screen exposes a “Rhymes” tab plus inputs for scales/chords/progressions.
- Functionality is usable but limited to canned datasets; expansion would require a real lexicon or API.

### 4. AI Agents — **Status: Exists (Placeholder)**  
**Key files**
- [`ai.go`](internal/app/ai.go)
- [`continuation_agent.go`](internal/app/ai/continuation_agent.go)
- [`variation_agent.go`](internal/app/ai/variation_agent.go)
- [`rapid_brainstorm_agent.go`](internal/app/ai/rapid_brainstorm_agent.go)
- [`quality_agent.go`](internal/app/ai/quality_agent.go)

**Highlights**
- AI service wires four agents (continuation, variation, rapid brainstorm, quality), but each agent returns canned strings.
- No calls to Ollama/OpenAI or RAG retrieval exist; enhancement effort must implement real model integration and prompt handling.
- Quality checks simulate results without true analysis, though tiered entry points (red flag/basic/full) are defined.

### 5. Frontmatter Handling — **Status: Exists (Partial)**
**Key files**
- [`editor.go`](internal/app/editor.go)

**Highlights**
- Editor service serializes `SongMetadata` into YAML frontmatter (`serializeSongToMarkdown`) and parses versions back via `parseMarkdownToSong`.
- Versioning, milestones, and auto-save leverage repository interfaces with retry/error handling.
- Frontmatter parsing is minimal (single-section fallback, no validation of custom fields); Enh.6.5 should consider strengthening structure and UI exposure.

### 6. Chord Functionality — **Status: Exists (Functional)**
**Key files**
- [`chord_picker.go`](internal/ui/editor/chord_picker.go)
- [`chord_picker.go`](internal/ui/chord_picker.go)
- [`chord_progressions.go`](internal/data/chord_progressions.go)
- [`data/chord_progressions.json`](data/chord_progressions.json)

**Highlights**
- UI chord picker supports filtering by mood, random selection, and callback insertion into the editor.
- Data loader reads JSON dataset with caching, filtering by mood, ID lookup, and global helpers.
- Integration is complete from data→UI→editor insertion; dataset can be extended without code changes.

### 7. Theme System — **Status: Exists (Functional)**
**Key files**
- [`registry.go`](internal/theme/registry.go)
- [`manager.go`](internal/theme/manager.go)
- [`loader.go`](internal/theme/loader.go)
- [`theme.go`](internal/theme/theme.go)
- [`theme.go`](internal/ui/styles/theme.go)

**Highlights**
- Registry defines 12 named themes (plus aliases) meeting the Enhancement&nbsp;#6 requirement.
- Manager supports thread-safe switching, default fallback, and config-driven initialization.
- UI styles layer (`internal/ui/styles/…`) consumes theme values for consistent rendering.
- Theme cycling shortcuts already exist (`ctrl+shift+n/p`, etc.).

### 8. Knowledge Base / RAG Implementation — **Status: Missing**
**Findings**
- No knowledge base, embeddings store, or retrieval API found under `internal/app`, `internal/data`, or documentation directories.
- The only domain “dataset” is chord progressions; no lyric knowledge, tutorials, or AI reference corpora exist.
- Enhancement&nbsp;#6.5 must plan and implement KB/RAG infrastructure from scratch.

---

## Summary Matrix

| Component | Exists | Completeness | Notes |
|-----------|--------|--------------|-------|
| Editor components | ✅ | High | Multi-pane editor, quick tools, and layouts are production-ready. |
| Keyboard shortcuts | ✅ | High | Comprehensive bindings with contextual help and tests. |
| Rhyme / syllable tools | ✅ | Medium | Functional but static dictionary and heuristic syllable counts. |
| AI agents | ✅ | Low | Stub responses only; no real AI integration. |
| Frontmatter handling | ✅ | Medium | YAML serialization/parsing exists; structure is minimal. |
| Chord functionality | ✅ | High | Data loader, picker UI, and editor insertion fully wired. |
| Theme system (12 themes) | ✅ | High | Runtime switching, config integration, and aliases implemented. |
| Knowledge base / RAG | ❌ | N/A | No KB infrastructure present. |

---

## Recommendations & Next Steps
1. **Prioritize AI integration**: replace placeholder agents with real model calls (including prompt templates, streaming, and error handling) before building new flows in Enhancement&nbsp;#6.5.
2. **Design the knowledge base**: define data sources, storage, and retrieval APIs to support future AI features; none exist today.
3. **Upgrade rhyme tooling**: integrate a richer dataset or API to move beyond the static dictionary if enhancement requirements need it.
4. **Harden frontmatter**: expand parsing/validation and surface metadata editing in the UI to leverage existing serialization logic.
5. **Document existing editor capabilities**: ensure enhancement planning accounts for current quick tools, layouts, and shortcuts so rework is avoided.

---

## Audit Log
- **2025-10-18**: Repository assessment and component verification completed.