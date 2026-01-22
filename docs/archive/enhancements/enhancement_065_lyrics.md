# Enhancement #6.5: Lyric Features Preservation & KB Integration
**Ensure existing lyric features still work after Enhancement #6 changes**

**Status:** Ready | **Priority:** Critical | **Effort:** 2 weeks | **Start:** After #6 Week 2

---

## Phase 0: Codebase Audit (Day 1)

**CRITICAL: Before starting work, audit actual codebase.**

### Audit Checklist

Run these checks:

```bash
# 1. Check if editor exists
find . -name "editor.go" -o -name "edit_mode.go"

# 2. Check for KB system
find . -path "*/knowledge/*" -name "*.go"
find . -name "*chroma*" -o -name "*embed*"

# 3. Check for rhyme/syllable
find . -path "*/rhyme/*" -name "*.go"
find . -path "*/syllable/*" -name "*.go"

# 4. Check for AI agent
find . -path "*/ai/*" -name "*.go"
grep -r "QuickIdeaAgent" .

# 5. Check data files
ls data/knowledge/ 2>/dev/null
ls data/cmudict* 2>/dev/null
```

### Document Findings

Create `AUDIT.md`:

```markdown
# Enhancement #6.5 Codebase Audit

## What Actually Exists

✅/❌ Split-pane editor: [path/to/file]
✅/❌ Knowledge base: [path/to/kb]
✅/❌ Rhyme dictionary: [path/to/rhyme]
✅/❌ Syllable counter: [path/to/syllable]
✅/❌ AI agent: [path/to/agent, current model: X]
✅/❌ YAML frontmatter: [path/to/parser]
✅/❌ Chord insertion: [path/to/chord]

## What's Missing

- List anything assumed to exist but doesn't
- Note what needs building vs integration

## Adjustments Needed

- Adjust Week 1/2 tasks based on findings
- Flag if scope is wrong
```

### Report & Adjust

If reality differs from spec:
1. Update enhancement scope
2. Reprioritize tasks
3. Report to Simon before proceeding

---

## Purpose

**NOT building new features.** Ensuring existing lyric editing system still works after Enhancement #6's changes, plus integrating with new systems (themes, 3B AI model).

---

## What EXISTS (Don't Rebuild)

✅ Split-pane markdown editor  
✅ Lyric editing mode  
✅ Knowledge base (RAG with ChromaDB)  
✅ Rhyme dictionary (CMU)  
✅ Syllable counter  
✅ YAML frontmatter parsing  
✅ Chord symbol insertion  
✅ AI agent (needs update to 3B model)

---

## What NEEDS WORK

### Week 1: Integration Testing

**Test existing features with #6 changes:**
- [ ] Editor still works with new theme system
- [ ] Knowledge base loads correctly
- [ ] Rhyme dictionary functional
- [ ] Syllable counter accurate
- [ ] YAML frontmatter parses
- [ ] Chord insertion works

**Update to new systems:**
- [ ] Editor uses new theme colors (12 themes)
- [ ] KB integrates with QuickIdeaAgent
- [ ] AI context detection (lyric vs pattern)

---

### Week 2: AI Integration & Export

**Update AI for lyrics:**
```go
// AI now detects context
func (a *QuickIdeaAgent) Generate(req QuickRequest) (*QuickResponse, error) {
    // Detect if lyric or pattern
    context := DetectContext(req.Context)
    
    // Query KB if lyrics
    var kbCards []Card
    if context == "lyric" {
        kbCards = a.kb.Search(req.Context, 3)
    }
    
    // Build prompt with KB context
    prompt := a.buildPrompt(req, kbCards, context)
    
    // Use qwen2.5:3b (same as patterns)
    return a.generate(prompt)
}
```

**Lyric-specific prompts:**
```go
const LyricUnstickPrompt = `Songwriting assistant. Next lyric line.

Knowledge Base Context:
%s

Current Context:
%s

Section: %s
Style: Concrete imagery, conversational

3 options (8-12 syllables):
1.
2.
3.`
```

**Export formats:**
- Markdown (.md)
- Plain text (.txt)
- ChordPro (.cho)

---

## Implementation

### Task 1: Theme Integration
```go
// Update editor to use new theme system
editor := NewEditor(theme.GetCurrent())
```

### Task 2: KB + AI Integration
```go
// internal/app/ai/context_detector.go
func DetectContext(content string, cursor int) string {
    if isPatternSyntax(content) {
        return "pattern"
    }
    return "lyric"
}

// Update QuickIdeaAgent to query KB for lyrics
```

### Task 3: Model Update
```go
// Update to qwen2.5:3b (from whatever it was)
model: "qwen2.5:3b"
```

### Task 4: Export Commands
```bash
noise export lyrics --format md -o song.md
noise export lyrics --format txt -o song.txt
noise export lyrics --format chordpro -o song.cho
```

---

## Testing Checklist

**Existing Features:**
- [ ] `noise edit [song]` launches
- [ ] Split-pane renders
- [ ] Markdown preview works
- [ ] KB loads on startup
- [ ] Rhyme lookup (Ctrl+R) works
- [ ] Syllable counts show
- [ ] Chord insertion (Ctrl+K) works
- [ ] YAML frontmatter parses

**New Integration:**
- [ ] All 12 themes work in editor
- [ ] AI unstick uses KB
- [ ] AI responds <2 sec
- [ ] Context detection works
- [ ] Export formats work

---

## Files Modified

```
internal/ui/editor/editor.go          (theme integration)
internal/app/ai/quick_idea_agent.go   (KB + context)
internal/app/ai/context_detector.go   (new)
internal/app/ai/prompts.go            (lyric prompts)
internal/export/lyrics.go             (export formats)
cmd/export.go                         (export command)
```

---

## Critical: Do NOT

❌ Rebuild editor from scratch  
❌ Rebuild KB system  
❌ Rebuild rhyme dictionary  
❌ Rebuild syllable counter  
❌ Change existing data formats

**Only integrate with new systems from #6.**

---

**Timeline:** 2 weeks after #6 Week 2 completes  
**Next:** Return to Enhancement #6 Week 3 (AI agent)