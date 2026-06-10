# RALPLAN Init

## Brief
Cross-app memory, new AI-powered views/screens for all 4 kyanite.sh apps, config file system. No always-on voice. Leverage the unified pkg/ai Brain (LLM on NUCBox gemma4:12b, STT via local whisper.cpp, PostgreSQL memory on NUCBox) to build features that were impossible when each app had its own tiny local model.

## Apps
- focus: task manager (ADHD-friendly)
- noise: music creation (songwriting studio)
- syntax: fiction editor (interactive fiction)
- prism: color palette tool (WCAG contrast)

## Stack
Go, Bubble Tea, Lipgloss, pkg/design/, pkg/ai/

## Constraints
- No always-on voice
- All inference via pkg/ai Brain
- PostgreSQL on NUCBox for memory
- whisper.cpp for STT (on-demand, press-to-talk)
- Must work offline (graceful degradation)
