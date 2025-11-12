# DEPRECATED

## 1. Deprecation Policy
When features, documents, or configuration flags no longer align with the rapid-capture direction summarized in [`docs/roadmap.md`](./roadmap.md), they are retired from the primary docs set, moved into [`docs/archive/`](./archive/), and linked from this page with current replacements. This mirrors the rename guidance in `docs/archive/noise_sh_rename_doc.txt`.

## 2. Recently Deprecated Features
- **LyricForge branding assets** → Use the new “noise.sh — rapid-capture songwriting studio” messaging in [`README.md`](../README.md).
- **Legacy four-agent AI pipeline** → Replaced by QuickIdeaAgent (see [`docs/context_aware_ai_guide.md`](./context_aware_ai_guide.md)).
- **Legacy CLI flows (`lyricforge …`)** → Use `noise` and `noise quick [theme]` commands documented in README. Export is accessed via UI (`Ctrl+E`).
- **Old export path assumptions (pre JSON + theme-ready exports)** → Use the Enhancement 6 export pipeline in [`docs/export_formats.md`](./export_formats.md).

## 3. Legacy AI Agent System
Archived spec: [`docs/archive/enhancement_02_rapid.txt`](./archive/enhancement_02_rapid.txt) describes the four-agent cadence (Continuation, Variation, Quality, Brainstorm). That architecture is superseded by the Week 3 right-sizing work; rely on QuickIdeaAgent for unstick, spark, tweak, and check flows.

## 4. Theme System (Pre-Week 3)
Earlier styling guidance lived in outdated docs. The authoritative reference is [`docs/design_system_reference.md`](./design_system_reference.md) plus the theme registry in `internal/theme/`. Archived historical notes: see `docs/archive/enhancement_02_rapid.txt` and other pre-week-3 materials.

## 5. Export Workflow (Pre-Week 2)
Legacy export tooling lacked pattern JSON and theme-ready Markdown. Use the current export commands (`noise export pattern`, `noise export draft`) and JSON schema described in Enhancement 6 Week 2 plus [`docs/export_formats.md`](./export_formats.md). Archived reference: `docs/archive/enhancement_02_rapid.txt`.

## 6. Dictionary & Knowledge Stubs
Early dictionary experiments (pre Week 2 CMU enhancement) are documented in `docs/archive` and [`docs/cmu_dictionary_enhancements.md`](./cmu_dictionary_enhancements.md). Production builds should depend on the enhanced syllable/rhyme services captured in Enhancement 6 (Week 2) and not the original stubbed lists.

## 7. File Layout Changes
The Lyric → noise.sh rename plan, as well as repository moves, appear in [`docs/archive/noise_sh_rename_doc.txt`](./archive/noise_sh_rename_doc.txt). Current repo location is `C:\Users\Simon\noise.sh`; adjust personal forks and scripts accordingly.

## 8. Deprecated Config/Env Flags
Remove or avoid:
- `LYRICFORGE_*` environment variables (use `NOISE_*` equivalents).
- Legacy CLI aliases (`lyricforge`, `lyricforge quick`, etc.).
- Old theme keys or hard-coded palette references that preceded the Amber Night default.

## 9. Historical CLI Reference
| Deprecated Command | Replacement | Notes |
| --- | --- | --- |
| `lyricforge quick [theme]` | `noise quick [theme]` | Launches scratch mode with auto-brainstorm. |
| `lyricforge live` | Not implemented | Pattern coding mode is planned but not yet available. |
| `lyricforge export …` | `Ctrl+E` in UI | Export menu provides Markdown, Plain Text, ChordPro, JSON, HTML, and PDF formats. |
| `lyricforge brainstorm` | `noise` + QuickIdeaAgent shortcuts | Use `Alt+G/R/V/C` triggers in editor. |

## 10. Migration Checklist
- [ ] Update shell aliases and scripts to call `noise` binaries.
- [ ] Move or copy legacy specs into [`docs/archive/`](./archive/); consult `docs/archive/README.md` for rationale.
- [ ] Replace outdated README snippets with rapid-capture messaging.
- [ ] Switch configuration to the theme manager and QuickIdeaAgent options described in current docs.
- [ ] Verify exports use the Week 2 JSON + theme-ready Markdown format.
- [ ] Remove deprecated env flags (`LYRICFORGE_*`) and adopt `NOISE_*` variants.
