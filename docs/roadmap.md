# noise.sh Strategic Roadmap

## 1. Executive Summary
noise.sh is positioned as the rapid-capture songwriting studio that brings instant ideation, local-first privacy, and expert-grade craft into a single terminal workspace. The product differentiates itself through zero-cloud processing, embedded songwriting pedagogy, a showcase-quality TUI, and anti-generic AI guardrails captured in [`docs/prd_markdown.md`](./prd_markdown.md). Enhancement 6 realigns every workflow toward “idea to captured draft in under 60 seconds,” ensuring that contributorship, user trust, and long-term extensibility stay anchored in the same core vision.

## 2. Current Status
- **Live capture foundation:** Pattern parser, live coding mode, and sample-backed audio playback operating within the modular monolith.
- **Creative accelerants:** Quick chord picker, BPM tapper, JSON export, and 10-theme registry with custom-theme loading.
- **AI right-sizing:** QuickIdeaAgent replaced the four-agent pipeline, keeping response times <2s.
- **Voice-to-text (NEW):** Push-to-talk dictation with whisper.cpp. Models auto-download on first use. No manual setup required.
- **PWA sync infrastructure (NEW):** Embedded HTTP/WebSocket server for companion app. Device pairing, media storage, and idea inbox ready.
- **Documentation:** Comprehensive docs updated to reflect current feature set.

## 3. Development Philosophy
Derived from [`docs/dev_roadmap.md`](./dev_roadmap.md):
1. **Speed over completeness** – ship the fastest viable capture workflow first.
2. **Capture before polish** – the sketchbook experience always precedes refinement tooling.
3. **Export over perfection** – ensure ideas leave the tool effortlessly, even if unfinished.
4. **Simple beats sophisticated** – prune complexity that slows the solo developer + AI loop.
5. **Local-first trust** – privacy, offline capability, and deterministic behavior remain non-negotiable.

## 4. Near-Term Roadmap
| Timeline | Objective | Highlights |
| --- | --- | --- |
| Current | PWA Companion Development | Build the Progressive Web App for mobile idea capture (voice memos, photos, tap tempo). |
| Current | Voice Polish | Model selection UI, microphone testing, transcription quality tuning. |
| Next | Sync Workflow | Idea inbox actions (assign to song, delete, preview media), conflict resolution. |
| Next | Cross-device Testing | Validate sync across different devices and network conditions. |

## 5. Mid-Term Outlook
- **PWA release:** Ship the companion Progressive Web App with full capture capabilities.
- **Voice model options:** Support additional languages and model sizes based on user feedback.
- **Sync reliability:** Add offline queue, retry logic, and conflict resolution for unreliable networks.
- **AI polish and observability:** Instrument QuickIdeaAgent timing budgets, add lightweight logging, and rehearse fallback paths.
- **Workflow stitching:** Ensure voice dictation, sync inbox, and editor work seamlessly together.

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
