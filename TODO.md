# kyanite.sh Architecture Audit — Session Handoff

## Session: 2026-06-11 (session 4)

**HEAD**: `bb91e9e` (docs update pending)
**Push**: origin ❌ (needs `git push origin main`), foge ❌ (NUCBox tailnet still down)
**Net session**: +43/-28 lines (+15 LoC)
**All tests green** with `-race` across all 4 apps

## Completed This Session

| Commit | Items | Summary |
|--------|-------|---------|
| `bb91e9e` | T6-01..04, T6-06 | Cross-app keybinding standardization: Ctrl+Q quit, Ctrl+A AI panel, Ctrl+Shift+T theme, ?+Ctrl+H help; brain init warnings in all 4 apps |

## Remaining Backlog (leverage-ranked)

### Quick Wins (< 1 hr each)

1. **T6-05** (1 hr) — Cross-app context save targets (pkg/appnames)
2. **T7-04** (1 hr) — pkg/config error-path tests

### High Leverage (2-4 hr each)

3. **T5-05** (3 hr) — Shared app bootstrapper
4. **T5-08** (2 hr) — Shared toast/notification
5. **T5-04 + T5-03** (5 hr) — Error paradigm + message type unification

### Medium Priority

6. **T8-04, T8-06..T8-12** (3 hr) — Remaining CI hardening
7. **T9-04..T9-07** (2 hr) — Remaining security items
8. **T10-01..T10-09** (4 hr) — Data flow & state polish

### Test Coverage (can be done incrementally)

9. **T7-01** — Focus TUI tests (2317 lines, zero coverage)
10. **T7-02** — Focus CLI command handler tests
11. **T7-03** — Focus Engine edge-case tests

## Key Decisions

- **FocusConfig inlined, not embedded**: koanf/v2 silently drops embedded struct tags
- **Two-phase JSON unmarshal**: LLM returns `"deadline":""` which can't parse as `time.Time`; strip empty before unmarshal
- **MockBrain wraps all responses**: Must use Ollama `{"message":{"content":"..."},"done":true}` format; raw JSON gives empty `Message.Content`
- **resolvePath kept simple**: Always joins with BaseDir + Clean instead of returning error
- **music_tools ctrl+shift+t → ctrl+shift+m**: Theme cycling takes priority; music tools templates moved to non-conflicting key

## Pending

- `git push foge main` — 18 commits ahead, retry when NUCBox tailnet is up
