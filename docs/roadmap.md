# noise.sh Strategic Roadmap

## 1. Executive Summary
noise.sh is positioned as the rapid-capture songwriting studio that brings instant ideation, local-first privacy, and expert-grade craft into a single terminal workspace. The product differentiates itself through zero-cloud processing, embedded songwriting pedagogy, a showcase-quality TUI, and anti-generic AI guardrails captured in [`docs/prd_markdown.md`](./prd_markdown.md). Enhancement 6 realigns every workflow toward “idea to captured draft in under 60 seconds,” ensuring that contributorship, user trust, and long-term extensibility stay anchored in the same core vision.

## 2. Current Status
- **Live capture foundation (Weeks 1-2):** Pattern parser, live coding mode, and sample-backed audio playback are operating within the modular monolith, delivering immediate sketching capability.
- **Creative accelerants (Week 2 deliverables):** Quick chord picker, BPM tapper, JSON export, and the full 12-theme registry shipped alongside custom-theme loading, satisfying the capture-to-export loop.
- **AI right-sizing (Week 3):** The QuickIdeaAgent replaced the four-agent pipeline, keeping response times <2s while covering unstick, spark, tweak, and check flows.
- **Documentation sprint (Week 4 in progress):** Unified public docs, README refresh, and archival tasks are underway to remove legacy friction and clarify contributor pathways, as logged in [`docs/Enhancements/enhancement_06_complete (2).md`](./Enhancements/enhancement_06_complete%20(2).md).

## 3. Development Philosophy
Derived from [`docs/dev_roadmap.md`](./dev_roadmap.md):
1. **Speed over completeness** – ship the fastest viable capture workflow first.
2. **Capture before polish** – the sketchbook experience always precedes refinement tooling.
3. **Export over perfection** – ensure ideas leave the tool effortlessly, even if unfinished.
4. **Simple beats sophisticated** – prune complexity that slows the solo developer + AI loop.
5. **Local-first trust** – privacy, offline capability, and deterministic behavior remain non-negotiable.

## 4. Near-Term Roadmap (Week 4 Focus + Immediate Next Steps)
| Timeline | Objective | Highlights |
| --- | --- | --- |
| Week 4 | Documentation polish & contributor clarity | Publish unified public documents (this roadmap, ARCHITECTURE, DEPRECATED), align README messaging with rapid capture, and archive superseded assets. |
| Week 4 | Validation & stability | Run `go test ./...`, spot-check Markdown, and ensure export JSON plus theme registry regressions are covered. |
| Week 4 → Week 5 | Integration readiness | Close any outstanding TODOs from the integration roadmap in [`docs/enhancement_integration_roadmap.md`](./enhancement_integration_roadmap.md), then queue automation for export and theme persistence smoke tests. |

## 5. Mid-Term Outlook (Planned Weeks 5-8)
- **Theme system integration hardening:** Complete runtime theme switching QA, contrast validation, and persistence checks described in the integration roadmap.
- **AI polish and observability:** Instrument QuickIdeaAgent timing budgets, add lightweight logging, and rehearse fallback paths so AI remains deterministic under load.
- **Workflow stitching:** Ensure live mode, chord picker, export, and theme management behave coherently when chained, reducing context switches before the mobile companion workstream.

## 6. Long-Term Vision
Guided by the risk assessment in [`docs/enhancement_risk_assessment.md`](./enhancement_risk_assessment.md) and long-term plans:
- **Capture-first ecosystem:** noise.sh remains the personal sketchbook, exporting JSON patterns and Markdown drafts directly to wave.sh (production) and future tools like the mobile PWA.
- **Companion surfaces:** Ship the mobile capture PWA (voice, camera, tempo) and Spanish-language expansion once rapid capture adoption stabilizes, per Enhancement #4 research and mobile blueprint.
- **Plugin-ready core:** Keep the modular monolith clean so theme packs, exporters, and AI extensions live behind well-documented plugin interfaces.

## 7. Risk & Mitigation Watchlist
| Risk | Signal to Watch | Mitigation Strategy |
| --- | --- | --- |
| AI integration complexity | Latency >2s or degraded suggestion quality | Profile QuickIdeaAgent, cache prompts, provide offline fallbacks, and log timing budgets. |
| Theme performance | Theme switch lag or accessibility regressions | Preload palettes, enforce WCAG checks, and reuse theme manager patterns from Week 2. |
| CLI migration debt | Command friction during Cobra adoption | Maintain backward-compatible aliases, document new flow in README, and add smoke tests. |
| Documentation drift | Conflicting guidance between new and archived docs | Centralize authoritative docs (README + docs/), keep archive README updated, and sunset references during PR reviews. |

## 8. Success Metrics
Pulled from [`PROJECT_SNAPSHOT.md`](../PROJECT_SNAPSHOT.md):
- **Time-to-first-lyric:** <60 seconds from launch or `noise.exe quick`.
- **AI responsiveness:** <2-second suggestions with deterministic prompts.
- **Theme adoption:** >80% of users experiment with multiple themes; none report accessibility regressions.
- **Export reliability:** Pattern JSON and Markdown exports succeed without manual fixes.
- **User sentiment:** Feedback validates “rapid capture, instant ideation, local-first privacy” as lived experience, not aspirational copy.
