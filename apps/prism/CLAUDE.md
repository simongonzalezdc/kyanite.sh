# Prism.sh

## Project
A terminal-based color palette design and theory tool with interactive TUI, comprehensive color theory education, and WCAG accessibility validation

## Rules
This project follows KyaniteLabs engineering standards.

### Git and PR Hygiene
- Commits tell *why*, not *what*
- PRs under 400 lines
- Rebase, don't merge on feature branches
- Never skip CI

### Git Workspace Hygiene
Agents must leave the repository in the same clean state they found it. No exceptions.

- **Delete feature branches** after merge — whether you merged via PR or locally. A merged branch that still exists is a failure.
- **Remove worktrees** when done. If you created a git worktree for isolated work, remove it when the work is complete or abandoned. No orphaned worktrees.
- **Clean up stale references** — prune remote-tracking branches that no longer exist upstream (`git remote prune origin`).
- **No abandoned work left behind** — if you started a branch, pushed it, and then the task was cancelled or superseded, delete the branch. Dead branches accumulate and confuse everyone.
- **Local branches stay current** — regularly rebase local tracking branches against their upstream. If a local branch has diverged and is no longer needed, delete it.
- **Worktree state matches intent** — if you're done with a task, the working directory should be clean (`git status` shows nothing). Staged changes, untracked files, and modified files from abandoned work are not acceptable.

The rule is simple: finish your work, merge or discard it, clean up every artifact you created. The main branch is the only permanent artifact. Everything else is temporary and must be treated as such.

### Security Non-Negotiables
- No `shell=True` with user input
- No hardcoded secrets/tokens/API keys
- Validate file paths (no path traversal)
- Use HTTPS for external calls

### Testing
- Test behavior, not implementation
- Test the unhappy path
- One assertion per concept

### Code Quality
- No orphaned code — every function called, every file imported
- No commented-out code blocks
- Pin exact dependency versions

## Epoch Data Tracking

Every project must use Epoch (KyaniteLabs/Epoch) for time estimation and actively feed it data. Epoch only becomes useful when it has real estimate-vs-actual data from real projects.

### What This Means
- **Before starting a task**, get a time estimate from Epoch (via MCP, REST API at `localhost:3099`, or CLI)
- **After completing a task**, record the actual time spent using `record_actual` or `POST /v1/feedback/record-actual`
- **Include context** — task type, complexity, tools used, anything that helps the model learn
- **Batch submissions are fine** — use `batch_record_actuals` for multiple estimates at once

### Integration
- MCP: add `@puenteworks/epoch` to your project's `.mcp.json`
- REST API: `epoch serve --port 3099`
- CLI: `npx @puenteworks/epoch pert-estimate ...`

### Why This Is Non-Negotiable
Epoch's accuracy improves with data. Without estimate-vs-actual feedback from real projects, it's just a calculator with uncalibrated assumptions. Every project that uses Epoch and reports back makes every other project's estimates better. This is a collective intelligence system — it only works if everyone contributes.

<!-- EMPOWER_ORCHESTRATOR:START -->
## Empower the Orchestrator

This repository is governed by the Empower Orchestrator law. Every top-level/orchestrator Claude Code session is an audition to improve the system, not only finish the current task.

When you notice a repeatable task done 3+ times or a recurring agent failure mode, consider shipping the smallest durable artifact that prevents the repetition: a tool, skill, slash command, hook, guardrail, memory entry, test, verifier, or doctrine doc.

This applies to top-level/orchestrator sessions. Background workers execute their assigned slice and do not independently widen scope.

Before dispatching automation or creating a durable system change, state the four-question blast-radius check in chat:

1. Scale: one file/workspace/all sessions?
2. Severity: minor friction/broken workflow/data loss or leaked content?
3. Reversibility: single revert/manual cleanup/surgery?
4. Predictability: bounded failure mode/guessing/unknown?

All green permits auto mode. Any yellow requires inline human approval. Any red means do not dispatch; do the work inline or escalate.

Worker discipline: isolated worktree/sandbox, one artifact equals one commit/change unit, verify before commit, register through the target tool's native discovery surface, and never write outside the assigned scope.

Success line: “I noticed X, found a better way. The system just got an upgrade.”

Full recipe: `docs/agent-law/empower-orchestrator.md`.
<!-- EMPOWER_ORCHESTRATOR:END -->
