# kyanite.sh Architecture Audit — Session Handoff

## Session: 2026-06-11 (session 3)

**HEAD**: `f50c5c7` (1 commit since session 2, not yet pushed to origin)
**Push**: origin ❌ (needs `git push origin main`), foge ❌ (NUCBox tailnet still down)
**Net session**: +63/-76 lines (-13 LoC)
**All tests green** with `-race` across all 4 apps

## Completed This Session (Session 2 + 3)

| Commit | Items | Summary |
|--------|-------|---------|
| `2d029e5` | T4-02 | 30s timeout on 6 AI call sites via `cli/withAITimeout()` |
| `e3aa878` | T4-05, T4-06 | `sync/atomic.Bool` for audio.enabled + `defer recover()` on 4 naked goroutines |
| `4c7e1c1` | T8-02, T8-03, T8-05 | CI: `-race`, `pkg-and-cmd` job, `vulncheck` job |
| `0167a64` + `361cb28` | T5-02 | Focus dual config → shared `pkg/config`, deleted `focus/pkg/config` (-696 net) |
| `033fddb` | T4-01 | 5 fire-and-forget goroutines: 5s timeout + `defer recover()` |
| `35612b3` | T9-01, T9-02, T9-03 | SSLMode require, file perms 0600/0700, resolvePath containment |
| `8953dfb` | T7-05 | MockBrain wired to focus AI tests, two-phase JSON unmarshal for empty deadlines |
| `f50c5c7` | T4-03, T4-04 | TOCTOU race fix (flushCacheLocked) + writeMutex scoped to DB calls; lastSaveTime race fix |
| `d0b59ce` | T5-06 | Top-level Bubble Tea models → value receivers (RootModel, MainModel, synthwave) |

## Remaining Backlog (leverage-ranked)

### Quick Wins (< 1 hr each)

1. **T7-04** (1 hr) — pkg/config error-path tests

### High Leverage (2-4 hr each)

5. **T5-01** (4 hr) — Unify 4× theme managers into `pkg/design/manager.go` (highest leverage remaining)
6. **T5-05** (3 hr) — Shared app bootstrapper
7. **T5-04 + T5-03** (5 hr) — Error paradigm + message type unification
8. **T5-08** (2 hr) — Shared toast/notification

### Medium Priority

9. **T6-01..06** (3 hr) — Cross-app UX consistency (quit keys, theme cycling, help)
10. **T8-04, T8-06..T8-12** (3 hr) — Remaining CI hardening
11. **T9-04..T9-07** (2 hr) — Remaining security items
12. **T10-01..T10-09** (4 hr) — Data flow & state polish

### Test Coverage (can be done incrementally)

13. **T7-01** — Focus TUI tests (2317 lines, zero coverage)
14. **T7-02** — Focus CLI command handler tests
15. **T7-03** — Focus Engine edge-case tests

## Key Decisions

- **FocusConfig inlined, not embedded**: koanf/v2 silently drops embedded struct tags
- **Two-phase JSON unmarshal**: LLM returns `"deadline":""` which can't parse as `time.Time`; strip empty before unmarshal
- **MockBrain wraps all responses**: Must use Ollama `{"message":{"content":"..."},"done":true}` format; raw JSON gives empty `Message.Content`
- **resolvePath kept simple**: Always joins with BaseDir + Clean instead of returning error

## Pending

- `git push foge main` — 18 commits ahead, retry when NUCBox tailnet is up
- CyberWitches `.gjc/` kyanite state leak — ✅ cleaned this session
