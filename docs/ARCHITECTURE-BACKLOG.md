# kyanite.sh Architecture Audit Backlog

Generated: 2026-06-11
Methodology: 10 architect subagents with distinct personas, 180 findings, consolidated and ranked by (Risk × Leverage) / Effort.
Updated: 2026-06-11 (session 4)
HEAD: `bb91e9e`

## Completed

### Tiers 1-3 (`22c4d06`)

- T1-01: CI Go version 1.25→1.26
- T1-02: Dual DI containers → single `di.GetContainer()`
- T1-05: Config file permissions 0644→0600
- T1-07: Store fallback warning
- T2-01: Deleted noise/collaboration/ (~1,700 lines)
- T2-03: Deleted noise/plugins/ (~3,500 lines)
- T2-05: Deleted noise/knowledge/ (~1,200 lines, stubbed in ai/knowledge_stub.go)
- T2-06: Deleted noise/agent/ (~2,700 lines, museAgent commented out in split_pane.go)
- T2-07: Deleted focus/recurrence/ (~530 lines)
- T2-08: Deleted focus/command/ (~600 lines, removed from root.go)
- T2-09: Deleted pkg/tui/stt/ (128 lines)
- T2-10: Deleted noise generic_list.go, generic_cache.go (77 lines)
- T2-11: Deleted noise errors/generic_result.go (170 lines)
- T3-01: Deleted pkg/design/icons/width.go (155 lines)
- T3-02: Deleted pkg/design/border.go (33 lines)
- T3-03: Removed RegisterIcons/SetStyle/CurrentIcon from icons.go
- T3-04: Removed 9 unused validators from focus/validation
- T3-05: Removed unused exports from pkg/ai (IsError, WithStopWords, WithTimeout, NewSTTClient)
- T3-06: Removed LoadAllContext/DeleteContext/ContextEntry from pkg/ai
- T3-07: Removed MustGet from pkg/design
- T3-08: Removed unused exports from focus/errors
- T3-09: Removed gum.Spin

**Net: -31,443 lines deleted.**

### Session 2 (`073afa6` → `8953dfb`)

- T4-01 (`033fddb`): 5 fire-and-forget goroutines bounded with 5s timeout + `defer recover()`
- T4-02 (`2d029e5`): 6 AI call sites wrapped with 30s timeout via `cli/withAITimeout()`
- T4-05 (`e3aa878`): `AudioPlayer.enabled` → `sync/atomic.Bool`
- T4-06 (`e3aa878`): `defer recover()` on 4 naked goroutines
- T5-02 (`0167a64` + `361cb28`): Focus dual config → shared `pkg/config`, deleted `focus/pkg/config` (-696 net)
- T7-05 (`8953dfb`): MockBrain wired to focus AI tests, `NewWithBrain()` DI, two-phase JSON unmarshal
- T8-02 + T8-03 + T8-05 (`4c7e1c1`): CI: `-race` on all test steps, `pkg-and-cmd` job, `vulncheck` job
- T9-01 + T9-02 + T9-03 (`35612b3`): SSLMode require, file perms 0600/0700, resolvePath containment

**Session 2 net: ~-500 LoC**

### Session 3 (`8953dfb` → `b19ace8`)

- T4-03 (`f50c5c7`): `writeMutex` scoped to DB calls in autosave executeSave/executeSaveWithVersioning/CleanupOldVersions; lastSaveTime data race fixed
- T4-04 (`f50c5c7`): `flushCache` → `flushCacheLocked` in engine.go + journal_storage.go; all mutations now hold mutex through flush (TOCTOU eliminated)
- T5-06 (`d0b59ce`): Top-level Bubble Tea models (RootModel, MainModel, synthwave Model) → value receivers; noise/cmd type assertion updated
- T5-01 (`b19ace8`): Shared `design.Manager` in `pkg/design/manager.go`; all 4 apps delegate via thin wrappers (-133 LoC)

**Session 3 net: -136 LoC**

### Session 4 (`022592b` → `bb91e9e`)

- T6-01 (`bb91e9e`): Quit keybinding standardized to Ctrl+Q across all 4 apps; removed `q` from focus main/synthwave, noise settings, focus unified dashboard
- T6-02 (`bb91e9e`): AI panel toggle standardized to Ctrl+A; syntax ctrl+p→ctrl+a, text editor AI suggestion moved to ctrl+p
- T6-03 (`bb91e9e`): Theme cycling standardized to Ctrl+Shift+T; focus ctrl+t→ctrl+shift+t, noise global handler added, music_tools conflict resolved
- T6-04 (`bb91e9e`): Help keybinding standardized to ? + Ctrl+H; added ctrl+h to noise, ? to prism, removed h from syntax
- T6-06 (`bb91e9e`): Brain init error handling — all 4 apps now log warning on failure (focus, syntax, prism added; noise already had it)

**Session 4 net: +15 LoC (16 files, keybinding standardization)**
---

## Remaining Backlog (Tiers 4-10)

### TIER 4 — Concurrency Bugs (all resolved ✅)

| ID | Finding | Files | Effort |
|----|---------|-------|--------|
| ~~T4-01~~ | ~~goroutine leaks~~ | ~~noise/autosave.go, theme.go, files.go, voice, errors~~ | ~~2 hr~~ ✅ |
| ~~T4-02~~ | ~~missing context timeouts on AI ops~~ | ~~focus/cli, tui~~ | ~~30 min~~ ✅ |
| ~~T4-03~~ | ~~Mutex held during I/O~~ | ~~noise/autosave.go, focus/journal_storage.go~~ | ~~1 hr~~ ✅ |
| ~~T4-04~~ | ~~Engine lock-modify-unlock-flush TOCTOU race~~ | ~~focus/engine.go, journal_storage.go~~ | ~~30 min~~ ✅ |
| ~~T4-05~~ | ~~AudioPlayer.enabled data race~~ | ~~focus/audio/audio.go~~ | ~~10 min~~ ✅ |
| ~~T4-06~~ | ~~naked goroutines without recover()~~ | ~~focus/tui, audio, noise/autosave~~ | ~~20 min~~ ✅ |

---

### TIER 5 — API Surface Unification (~18 hr remaining, MEDIUM risk, HIGH leverage)

| ID | Finding | Effort | Notes |
|----|---------|--------|-------|
| ~~T5-01~~ | ~~Unify 4× theme managers~~ | ~~4 hr~~ ✅ | Shared `design.Manager` in `pkg/design/manager.go`; apps delegate via thin wrappers (-133 LoC) |
| ~~T5-02~~ | ~~Unify focus dual config~~ | ~~3 hr~~ | ✅ Deleted `focus/pkg/config`, using shared `pkg/config` |
| T5-03 | Unify ChatMessage/Message types (3 types → 1 in pkg/ai) | 2 hr | ai.Message, app.ChatMessage, agent.ChatMessage — different field counts |
| T5-04 | Unify error paradigm (sentinels vs structured, not both) | 3 hr | Requires noise/errors/ simplification (already partially done) |
| T5-05 | Create shared app bootstrapper (`pkg/tui/bootstrap`) | 3 hr | Standardize: flags → config → brain → session → model → run → shutdown |
| ~~T5-06~~ | ~~Standardize Bubble Tea model receivers~~ | ~~2 hr~~ ✅ | Top-level models (RootModel, MainModel, synthwave Model) → value receivers; noise sub-models deferred |
| T5-07 | Extract shared export infrastructure (`pkg/export`) | 2 hr | Noise: 6-format ExportModel, Prism: palette exporters, Syntax: story exporters |
| T5-08 | Extract shared toast/notification (`pkg/tui/toast`) | 2 hr | Noise has ToastModel + NotificationManager (overlapping), other apps have nothing |

**T5-08 is the highest leverage remaining item** — shared toast/notification across all apps.

---

### TIER 6 — Cross-App Consistency (~3 hr, LOW risk)

| ID | Finding | Detail |
|----|---------|--------|
| ~~T6-01~~ | ~~Standardize quit keybinding~~ | ~~ctrl+q everywhere~~ ✅ |
| ~~T6-02~~ | ~~Standardize AI panel toggle~~ | ~~ctrl+a everywhere~~ ✅ |
| ~~T6-03~~ | ~~Standardize theme cycling key~~ | ~~ctrl+shift+t everywhere~~ ✅ |
| ~~T6-04~~ | ~~Standardize help keybinding~~ | ~~? + ctrl+h everywhere~~ ✅ |
| T6-05 | Cross-app context save targets | Use `pkg/appnames.All` constant instead of hardcoded lists per app |
| ~~T6-06~~ | ~~Brain init error handling~~ | ~~All 4 apps log warning on brain init failure~~ ✅ |
---

### TIER 7 — Test Coverage Gaps (~14 hr remaining, MEDIUM risk)

| ID | Finding | Detail |
|----|---------|--------|
| T7-01 | Focus TUI — 2317 lines, zero tests | State machine transitions, chat, filter, calendar, timers all untested |
| T7-02 | Focus CLI — 20+ command handlers untested | Only cli_test.go tests imported validation functions |
| T7-03 | Focus Engine — AddSubtask/RestoreTask/UpdateTask untested | Real parent-child linking, duplicate-ID, timestamp logic |
| T7-04 | pkg/config — no error-path tests | Missing: malformed YAML, permission errors, directory-as-path |
| ~~T7-05~~ | ~~pkg/testutil/mock_brain exists but unused~~ | ~~Focus AI tests still hit network (nucbox:11434)~~ ✅ |
| T7-06 | Focus store_test.go asserts on unexported filePath field | Tests implementation, not behavior |
| T7-07 | Focus ai/manager_test.go tests 5 unexported functions | Tests implementation details |

---

### TIER 8 — Build/CI Hardening (~2 hr remaining, MEDIUM risk)

| ID | Finding | Detail |
|----|---------|--------|
| ~~T8-02~~ | ~~Add CI test jobs for pkg modules + cmd/kyanite~~ | ~~CI only tests 4 apps, misses 7 pkg modules~~ ✅ |
| ~~T8-03~~ | ~~Add -race flag to CI test step~~ | ~~Currently `go test ./...` without -race~~ ✅ |
| T8-04 | Add test coverage enforcement | No -coverprofile, no threshold gate |
| ~~T8-05~~ | ~~Add govulncheck step to CI~~ | ~~No vulnerability scanning~~ ✅ |
| T8-06 | Align dependency versions (go work sync) | x/sys at 3 versions, yaml at 2, pkg/tui has drift |
| T8-07 | Fix pkg/tui dep drift | lipgloss 1.1.0 vs 1.1.1-pre, runewidth 0.0.16 vs 0.0.20 |
| T8-08 | Add go.sum for pkg/cache | Missing after extraction |
| T8-09 | Add dependabot for 5 missing modules | pkg/design, pkg/tui, pkg/cache, pkg/session, pkg/testutil |
| T8-10 | Fix Makefile test target | Uses `cd $$dir && go test` (breaks on failure), missing -race |
| T8-11 | Add Makefile targets: tidy, vulncheck, verify, fmt | Missing standard targets |
| T8-12 | Add golangci linters: gosec, exhaustive, revive | Only 7 linters enabled |

---

### TIER 9 — Security Hardening (~2 hr remaining, MEDIUM risk)

| ID | Finding | Detail |
|----|---------|--------|
| ~~T9-01~~ | ~~PostgreSQL sslmode default to 'require'~~ | ~~Currently 'disable' — cleartext credentials over network~~ ✅ |
| ~~T9-02~~ | ~~File permissions: 0600/0700 for user data~~ | ~~Song files, media, DB dirs still world-readable~~ ✅ |
| ~~T9-03~~ | ~~Noise resolvePath — add path containment check~~ | ~~Absolute paths bypass BaseDir, enabling traversal~~ ✅ |
| T9-04 | Noise upload — validate magic bytes, not just Content-Type | Accepts any file with spoofed MIME type |
| T9-05 | Standardize on modernc.org/sqlite, remove CGO mattn/go-sqlite3 | Dual drivers increase attack surface |
| T9-06 | Mask pairing codes in logs (Debug, not Info) | Currently logs at Info |
| T9-07 | Don't return full filesystem paths in upload responses | Leaks internal directory structure |

---

### TIER 10 — Data Flow & State Polish (~4 hr, LOW risk)

| ID | Finding | Detail |
|----|---------|--------|
| T10-01 | Focus AI cache has no invalidation | Stale suggestions persist up to 24h TTL after task mutation |
| T10-02 | Noise settings theme change not propagated | Requires restart to see theme changes |
| T10-03 | Focus UnifiedDashboard styles stale after theme cycle | updateTheme() only touches main.go's package-level vars |
| T10-04 | LRU cache: ExpiresAt field dead, Set() doesn't update order | pkg/cache/lru.go |
| T10-05 | Noise RootModel.config pointer stale after save | Settings saves to disk but in-memory pointer unchanged |
| T10-06 | Focus CLI silently swallows AI fallback | User can't distinguish AI-chosen vs heuristic-chosen priorities |
| T10-07 | Focus edit.go deadline parse failure silently swallowed | User's deadline edit ignored without feedback |
| T10-08 | Focus bulk.go Scanln error discarded | Misleading "Deletion cancelled" on read error |
| T10-09 | 8 packages missing doc.go | pkg/cache, pkg/session, noise/{glm,sync,voice,export,pattern} |

---

## Recommended Next Session Order

1. **T5-05** (3 hr) — Shared app bootstrapper.
2. **T5-08** (2 hr) — Shared toast/notification.
3. **T5-04 + T5-03** (5 hr) — Error paradigm + message type unification.
4. **T6-05** (1 hr) — Cross-app context save targets (pkg/appnames).
5. Remaining T8 CI, T9 security, T10 data-flow items.

---

## Key Context for Continuation

- **Local repo**: ~/workspaces/kyanite-labs/kyanite.sh, branch `main`, HEAD `bb91e9e`
- **GitHub**: pushed to origin/main (8 commits since `073afa6`)
- **Forgejo**: 18 commits ahead, push pending (NUCBox tailnet down)
- **9 go.mod files**: apps/{focus,noise,syntax,prism}, pkg/{design,ai,config,tui,cache,session,testutil}, cmd/kyanite
- **go.work**: includes all modules + pkg/cache + pkg/session + pkg/testutil
- **Knowledge stub**: noise's deleted knowledge/ package is stubbed at `internal/app/ai/knowledge_stub.go`
- **Agent stub**: noise's deleted agent/ package is commented out in `internal/ui/editor/split_pane.go` (museAgent = interface{}(nil))
- **DI singleton**: focus CLI uses `di.GetContainer()` from helpers.go, provider.go has `GetContainer()` with `sync.Once`
- **MockBrain**: wraps all responses in Ollama format; `NewWithBrain()` uses memory-only cache
- **Two-phase unmarshal**: manager.go strips empty `"deadline":""` before parsing into `time.Time`
