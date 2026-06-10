# RALPLAN: Cross-App Memory, AI Features, Config File

**Run ID**: 2026-06-10-2143-de97  
**Status**: pending-approval  
**Scope**: All 4 kyanite.sh apps + pkg/ai + new pkg/config

---

## Principles

1. **Feature must be felt** — Every deliverable must produce a user-visible change. No infrastructure-only phases.
2. **Offline-first** — Every AI feature degrades gracefully when NUCBox is unreachable. Fallback to local-only behavior.
3. **Clean as you go** — Each phase removes dead code, fixes layer boundaries, and extracts shared patterns. No phase ships without leaving the codebase architecturally cleaner.
4. **One brain, many surfaces** — The Brain (pkg/ai) is the single source of inference. Apps provide surfaces (views, keybindings, commands) that call into it.
5. **Incremental delivery** — Each phase is independently shippable. No big-bang releases.

## Decision Drivers

1. **User impact per line changed** — Maximize felt improvements, minimize code churn.
2. **Cross-app leverage** — Features that benefit all 4 apps (config, memory, STT) ship before app-specific features.
3. **Architectural debt reduction** — Every phase must retire at least one known divergence (from ARCHITECTURE.md "Known Divergence" table).

---

## Recommended Plan

### Phase 1: Config File + Dead Code Purge (shared foundation)

**Why first**: Config is needed by everything else. Dead code removal reduces merge conflicts in later phases.

**Config file** (`pkg/config/`):
- Use koanf with YAML backend (`~/.config/kyanite/config.yaml`)
- Merge order: defaults → config file → env vars → CLI flags
- XDG-compliant paths (`~/.config/kyanite/`)
- Single struct covering all apps + brain config
- `kyanite config init` subcommand to scaffold
- `kyanite config show` to display resolved config
- All KYANITE_* env vars continue working (backward compat)
- Config struct replaces the scattered env-var reads in pkg/ai/config.go

**Config structure**:
```yaml
brain:
  ollama_url: http://nucbox:11434
  model: gemma4:12b
  timeout: 60s
  whisper_bin: whisper-stream
  whisper_model: ~/.local/share/achiote-voice/models/ggml-large-v3-turbo.bin
  db_host: nucbox
  db_port: 5432
  db_name: kyanite
  db_user: kyanite

focus:
  theme: amber-night
  journal_dir: ~/.focus

noise:
  theme: amber-night
  samples_dir: ~/noise/samples

syntax:
  theme: amber-night
  stories_dir: ~/syntax/stories

prism:
  theme: amber-night
  palettes_dir: ~/.local/share/prism
```

**Dead code purge** (in this phase):
- Remove focus legacy Ollama/OpenRouter code (~1200 lines in manager.go)
- Remove deprecated files (ollama.go, ollama_provider.go, openrouter_provider.go, fallback_provider.go)
- Make BrainProvider the sole provider path in focus
- Remove pkg/ai/prompts.go dead focus prompts (incompatible JSON shapes — focus uses its own PromptBuilder)
- Consolidate noise internal/infra/ollama/ (deprecated, only used by tests)
- Remove noise/internal/infra/voice/model_manager.go dead whisper model download

**Acceptance criteria**:
- [ ] `~/.config/kyanite/config.yaml` is read by all 4 apps + kyanite launcher
- [ ] Env vars still override config file values
- [ ] `kyanite config init` creates a valid config file
- [ ] Focus no longer imports or references legacy Ollama/OpenRouter code
- [ ] All 55+ test packages pass with -race
- [ ] 5 binaries build cleanly
- [ ] Net line count decreases from current baseline

---

### Phase 2: Cross-App Memory (resume where you left off)

**Why second**: Memory needs config (Phase 1) for DB connection. Doesn't need new views — surfaces through existing UIs.

**Memory service** (`pkg/ai/memory.go` additions):
- `SaveSession(app, sessionID, state)` — saves app state snapshot on exit
- `LoadSession(app, sessionID)` — restores app state on launch
- `GetRecentSessions(app, limit)` — lists recent sessions for "resume" menu
- `GetCrossAppContext(app)` — loads relevant context from other apps (e.g., "user was editing 'Dark Forest' chapter in syntax" when opening focus)
- `SearchMemory(app, query)` — full-text search across stored context
- Auto-save on app exit (via Brain.Close)
- Auto-load on app start (via brain init)

**New table**: `kyanite_sessions`
```sql
CREATE TABLE kyanite_sessions (
  id SERIAL PRIMARY KEY,
  app TEXT NOT NULL,
  session_id TEXT NOT NULL,
  title TEXT,
  state JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(app, session_id)
);
```

**New table**: `kyanite_cross_app_context`
```sql
CREATE TABLE kyanite_cross_app_context (
  id SERIAL PRIMARY KEY,
  source_app TEXT NOT NULL,
  target_app TEXT NOT NULL,
  context_type TEXT NOT NULL,
  summary TEXT NOT NULL,
  relevance_score REAL DEFAULT 0.5,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Per-app memory surfaces**:
- **Focus**: On launch, show "Welcome back — last session: 3 tasks completed, 2 in progress" (status bar)
- **Noise**: On launch, restore last project state (song, chords, lyrics)
- **Syntax**: On launch, "Resume 'Dark Forest' chapter 3?" prompt
- **Prism**: On launch, "Recent palettes: Sunburst, Ocean Deep, Midnight" in sidebar
- **Launcher**: Show recent activity per app in the launcher menu

**Architecture cleanup in this phase**:
- Extract shared session/state serialization to `pkg/session/`
- Standardize app lifecycle hooks (OnStartup, OnShutdown) across all 4 apps
- Remove per-app ad-hoc state loading (focus has viper-based state, others use file I/O)

**Acceptance criteria**:
- [ ] Closing and reopening an app restores meaningful state
- [ ] Launcher shows recent activity per app
- [ ] Cross-app context surfaces (e.g., focus knows about syntax stories)
- [ ] Memory degrades gracefully when PostgreSQL is unreachable
- [ ] No app hangs on startup when memory is slow/unavailable

---

### Phase 3: AI-Powered Views (new screens/panels)

**Why third**: Needs config + memory foundation. This is where the user feels the biggest leap.

#### Focus — AI Dashboard View
- New `Tab:AI` in the TUI (alongside tasks, journal, calendar)
- **Daily Briefing**: "You have 3 high-priority tasks due today. 'Finish proposal' is overdue."
- **Smart Suggestions**: "Based on your journal, you've been stressed about deadlines. Consider breaking down 'Ship v2' into subtasks."
- **Natural Language Task Entry**: Press `a`, speak/type "remind me to call Sarah about the marketing budget tomorrow at 3pm" → parsed into structured task
- **Weekly Summary**: AI-generated summary of what was accomplished, what's stalled, trends

#### Noise — AI Copilot Panel
- New split-pane view (toggle with `Ctrl+A`)
- **Lyric Continuation**: Type 2 lines, press `Tab`, get 4 suggested continuation lines (streaming)
- **Chord Suggestions**: "What chords go with Am → F → C in the key of C?" → suggestions
- **Mood Board**: Type a mood ("melancholic rainy day"), get: chord progression + tempo suggestion + lyric fragments + key recommendation
- **Voice-to-lyrics**: Press `v`, speak a lyric idea → transcribed and inserted at cursor

#### Syntax — AI Writing Partner
- New side panel (toggle with `Ctrl+P`)
- **Story Bible**: AI reads all chapters, maintains consistency (character names, locations, timeline)
- **Continue Writing**: Place cursor, press `Ctrl+Space`, get 3 continuation options (streaming)
- **Consistency Check**: "Does chapter 5 contradict what was established in chapter 2?"
- **Character Voice**: "Rewrite this dialogue in Marcus's voice based on previous chapters"
- **Voice dictation**: Press `v`, speak → text inserted at cursor

#### Prism — AI Design Partner
- Expand existing 'p' keybinding into full AI mode (`Ctrl+A`)
- **Palette from Image Description**: "Describe the artwork for your album cover" → palette
- **Mood Palettes**: "90s grunge" → 5 curated palettes with rationale
- **Accessibility Advisor**: "This palette fails WCAG AA for small text. Here are 3 fixes."
- **Palette Stories**: "Generate a short visual story using this palette" (creative exploration)

**Architecture cleanup in this phase**:
- Extract shared AI panel/view component to `pkg/tui/aipanel/` (reusable Bubble Tea component for all 4 apps)
- Consolidate prompt templates — remove per-app prompt code, use centralized pkg/ai/prompts.go
- Add streaming response renderer component (shows tokens as they arrive)
- Remove deprecated AI client code from all apps (old ollama providers, old syntax client.go)

**Acceptance criteria**:
- [ ] All 4 apps have at least one new AI-powered view/panel
- [ ] Streaming responses render token-by-token in the TUI
- [ ] Press-to-talk STT works in focus, noise, syntax (press `v`, speak, release)
- [ ] All AI features degrade gracefully when NUCBox is offline
- [ ] No deprecated AI client code remains in any app

---

### Phase 4: Cross-App Intelligence (the unified experience)

**Why last**: Needs all prior phases. This is the capstone.

**Kyanite Hub** (`cmd/kyanite` launcher upgrade):
- Show activity feed from all apps (recent tasks, song progress, writing stats, palettes created)
- "Today's Focus" — AI summary of what needs attention across all apps
- Cross-app suggestions: "You wrote about 'the forest' in syntax — here's a forest-green palette in prism"
- Quick actions: "Create task in focus about finishing the chapter you're writing in syntax"

**Shared STT component** (`pkg/tui/stt/`):
- Reusable Bubble Tea component for press-to-talk
- Visual indicator (recording dot, waveform)
- Works in focus, noise, syntax (prism doesn't need it)
- Configurable activation key per app

**Cross-app context awareness**:
- Focus surfaces "you were writing 'Dark Forest' in syntax" as context for task suggestions
- Noise surfaces "you created a 'melancholy' palette in prism" as mood reference
- Syntax surfaces "you have 3 overdue tasks in focus" as writing context
- Prism surfaces story themes from syntax as palette inspiration

**Architecture cleanup in this phase**:
- Extract shared Bubble Tea components to `pkg/tui/` (panel, status bar, AI panel, STT widget)
- Consolidate error handling — all apps use pkg/ai typed errors consistently
- Remove remaining app-specific theme code (anything not using pkg/design directly)

**Acceptance criteria**:
- [ ] Launcher shows cross-app activity feed
- [ ] Cross-app suggestions appear in at least 2 apps
- [ ] Shared STT component used by 3+ apps
- [ ] No app-specific theme code remains (all use pkg/design)
- [ ] Shared TUI components in pkg/tui/ used by all 4 apps

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| PostgreSQL schema migration breaks existing data | Low | High | Versioned migrations with rollback |
| koanf dependency conflicts | Low | Medium | Pin version, single shared module |
| NUCBox unreachable during demo | Medium | Low | All features degrade offline |
| Brain streaming latency on tailnet | Medium | Medium | Timeout + retry + local fallback |
| Prompt engineering doesn't produce good results | Medium | Medium | Iterative prompt tuning per app |
| Focus legacy code removal breaks existing behavior | Medium | High | Comprehensive test coverage before removal |

---

## RALPLAN-DR Summary

| Dimension | Position |
|---|---|
| **Principles** | 1. Feature must be felt 2. Offline-first 3. Clean as you go 4. One brain, many surfaces 5. Incremental delivery |
| **Drivers** | 1. User impact per line 2. Cross-app leverage 3. Debt reduction |
| **Options** | **A (Recommended)**: 4-phase incremental — config → memory → AI views → cross-app intelligence. Each phase ships independently. |
| | **B**: 2-phase big-bang — infrastructure (config+memory) then all features. Simpler but larger blast radius. |
| **Why A over B** | Option A delivers user value in Phase 1 (config + cleanup) and Phase 2 (resume). Option B delivers nothing until both phases complete. Option A also reduces merge conflict risk by keeping changes bounded per phase. |

---

## Phase Summary

| Phase | Ships | User Feels | Lines Changed (est.) | Debt Reduced |
|---|---|---|---|---|
| 1 | Config file + dead code purge | No more env vars, faster builds | +800 / -2000 | Focus legacy Ollama, dead prompts, deprecated files |
| 2 | Cross-app memory | "Resume where you left off" | +600 / -200 | Per-app state loading, session hooks |
| 3 | AI-powered views | New panels, streaming, press-to-talk | +3000 / -500 | Deprecated clients, duplicated prompts |
| 4 | Cross-app intelligence | Unified launcher, cross-app context | +1200 / -300 | App-specific theme code, duplicated TUI components |
