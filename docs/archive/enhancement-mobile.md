# Enhancement #3+5: Mobile Companion App (PWA)
**noise.sh - Capture Tool for On-the-Go Ideas**

**Version:** 3.0 (Companion App)  
**Status:** READY FOR IMPLEMENTATION  
**Date:** October 17, 2025  
**Owner:** Simon (Kyanite)  
**Priority:** High - Start after Enhancement #2 (Rapid Prototyping) ships

---

## Executive Summary

Create a **lightweight mobile companion PWA** that does ONLY what the TUI can't:
- **Capture ideas via sensors** (voice, camera, tap tempo)
- **Visual feedback** (prosody, syllable counts, waveforms)
- **Mobile sharing** (export, share sheet)
- **Minimal editing** (just enough to refine captured text)
- **Sync to TUI** (auto-sync to Dropbox/Drive)

**Key Principle:** This is NOT a full writing environment. It's a capture tool that feeds the TUI. The TUI remains the primary workspace for serious writing.

---

## Product Vision

### What This IS
âœ… **Capture device** â€” Record voice memos, scan handwritten notes, tap tempo  
âœ… **Visual companion** â€” See prosody, syllables, waveforms (terminal can't show)  
âœ… **Mobile tool** â€” Works offline on phone, optimized for quick capture  
âœ… **Sync bridge** â€” Saves to files that TUI can open  

### What This IS NOT
âŒ **Full writing environment** â€” No AI panel, no knowledge base, no multi-song management  
âŒ **TUI replacement** â€” Complement, not compete  
âŒ **Complex editor** â€” Minimal editing only  
âŒ **Desktop-first** â€” Mobile is the focus, desktop is bonus  

---

## Division of Labor: PWA vs TUI

### PWA Does (Capture & Feedback)

**Capture Tools:**
1. ðŸŽ¤ **Voice memo recording** â€” Hum melodies, capture ideas
2. ðŸŽ¤ **Voice-to-text dictation** â€” Speak lyrics rapidly
3. ðŸ“· **Camera OCR scanning** â€” Digitize handwritten notebooks
4. ðŸŽ¯ **Tap tempo input** â€” Set BPM by tapping rhythm

**Visual Feedback:**
5. ðŸ“Š **Visual prosody overlay** â€” See stress marks (Ë˜ ËŠ) above syllables
6. ðŸ“ˆ **Live syllable counter** â€” Real-time count as you type
7. ðŸŽµ **Inline waveforms** â€” View voice memos visually
8. ðŸŽ¨ **Rhyme highlighting** â€” Color-coded rhyme schemes

**Mobile Sharing:**
9. ðŸ“¤ **Rich exports** â€” Styled PNG/PDF for sharing
10. ðŸ“± **Share sheet integration** â€” Send via Messages, Mail, AirDrop
11. ðŸ“‹ **Formatted clipboard** â€” Paste with styling

**Basic Editing:**
12. âœï¸ **Minimal text editor** â€” Refine captured text
13. ðŸ’¾ **Invisible auto-sync** â€” Saves & syncs automatically (like Apple Notes)

---

### TUI Does (Create & Polish)

**AI Features:**
- AI brainstorming panel
- AI continuation suggestions
- AI variations generator
- Knowledge base integration

**Power Editing:**
- Full keyboard shortcuts
- Multi-song management
- Complex text manipulation
- Professional workflow

**Analysis:**
- Prosody analysis (beyond visual)
- Rhyme scheme detection
- Meter analysis
- Quality scoring

**The Rule:** If it requires serious focus and keyboard work, it belongs in the TUI.

---

## Workflow Example

### Scenario: Idea Strikes on Subway

**Step 1: Capture on Mobile PWA (Offline)**
```
ðŸ“± Subway (no internet)
â†“
Open PWA from home screen
â†“
Tap mic button â†’ hum melody (10 sec)
â†“
Type quick lyric idea:
  "The scent of cedar still haunts her bed"
  â””â”€ 8 syllables âœ“ (live counter shows)
â†“
Tap tempo to rhythm (84 BPM detected)
â†“
Auto-saves to local storage (IndexedDB)
â†“
Status shows: âš  Offline (will sync later)
```

**Step 2: Automatic Sync (Invisible)**
```
ðŸ“± Exit subway, phone gets signal
â†“
PWA detects connection
â†“
Background sync starts automatically
â†“
Status briefly shows: â†» Syncing...
â†“
Status changes to: âœ“ Synced
â†“
File appears in Dropbox/Drive instantly

You never tapped a sync button.
You never thought about it.
It just worked.
```

**Step 3: Open on Desktop TUI (Instant)**
```
ðŸ’» Open TUI (minutes or hours later)
â†“
See new draft "Lost My Dog" immediately
â†“
Play voice memo (melody idea)
â†“
Use AI to expand verse
â†“
Add knowledge base insights
â†“
Polish with full editor
â†“
Export final version
```

**Total User Actions for Sync:** Zero  
**Sync happens:** Automatically in background  
**Like:** Apple Notes, Google Docs, Notion

---

## Mobile PWA Interface (Minimal)

### Quick Capture Screen

```
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚ ðŸŽµ Quick Capture    [SKETCH]  âœ“    â”‚ â† Status: âœ“ Synced
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚                                     â”‚
â”‚ The scent of cedar still...        â”‚ â† Basic text editor
â”‚ haunts her bed                      â”‚   Auto-saves on keystroke
â”‚ Ë˜    ËŠ    Ë˜  ËŠ  Ë˜   ËŠ     Ë˜   ËŠ    â”‚ â† Prosody overlay
â”‚ 8 syllables âœ“                       â”‚ â† Live counter
â”‚                                     â”‚
â”‚ I find her collar in the drawer...  â”‚
â”‚â–ˆ                                    â”‚ â† Blinking cursor
â”‚                                     â”‚
â”‚ â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â” â”‚
â”‚ â”‚ Voice Memo 1       â–â–‚â–ƒâ–…â–‡â–ˆâ–‡â–…â–ƒâ–‚â– â”‚ â”‚ â† Inline waveform
â”‚ â”‚ [â–¶ï¸ Play] 0:08                  â”‚ â”‚
â”‚ â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜ â”‚
â”‚                                     â”‚
â”œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¤
â”‚ [ðŸŽ¤] [ðŸ“·] [â™ª] [ðŸ“¤]                  â”‚ â† Capture tools only
â”‚  Mic  Scan Tempo Share              â”‚   No manual sync button
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜

Sync Status Indicators (top-right):
âœ“    = Synced to TUI (green)
â†»    = Syncing... (blue, spinning)
âš     = Offline, will sync when online (yellow)
âœ—    = Sync failed, retrying... (red, rare)

Interactions:
- Type â†’ Auto-saves locally (500ms debounce)
- Save â†’ Auto-syncs to Dropbox/Drive
- No buttons needed, just works
```

### What's NOT Included
- âŒ No AI panel (use TUI for that)
- âŒ No song library (single-draft focus)
- âŒ No knowledge base search
- âŒ No complex menus/settings
- âŒ No multi-window layout
- âŒ No manual sync button (automatic)

**Design Philosophy:** "Just enough to capture and refine, syncs invisibly to TUI."

---

## Core Features (Detailed)

### 1. ðŸŽ¤ Voice Memo Recording

**Use Case:** Melody idea strikes while walking

**Flow:**
```
1. Tap mic button (bottom bar)
2. Recording starts (red dot indicator)
3. Hum melody (5-30 seconds)
4. Tap stop
5. Waveform appears inline
6. Tap play to review
7. Auto-saves to song metadata
```

**UI:**
```
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚ Recording... â¬¤ 0:08                 â”‚
â”‚ â–â–‚â–ƒâ–…â–‡â–ˆâ–‡â–…â–ƒâ–‚â–â–‚â–ƒâ–…â–‡â–ˆâ–‡â–…â–ƒâ–‚â–              â”‚ â† Live waveform
â”‚                                     â”‚
â”‚ [â–  Stop]  [âœ“ Keep]  [âœ— Cancel]     â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
```

**Technical:**
- Web Audio API (MediaRecorder)
- Store as .webm or .m4a
- Attach to YAML frontmatter
- Max 60 seconds per memo

**Syncs to TUI:**
```yaml
---
title: "Lost My Dog"
attachments:
  voice_memos:
    - melody_idea_1.webm
---
```

---

### 2. ðŸŽ¤ Voice-to-Text Dictation

**Use Case:** Ideas flow faster than typing

**Flow:**
```
1. Tap-hold mic button (1 second)
2. Speak lyrics
3. Text appears live
4. Pause = new line
5. Say "new verse" = section break
6. Say "scratch that" = undo
7. Release to stop
```

**Smart Features:**
- Auto-capitalization
- Pause detection = line break
- Voice commands ("new verse", "chorus", "scratch that")
- Live text display

**Technical:**
- Web Speech API
- Continuous recognition
- Interim results for live display

---

### 3. ðŸ“· Camera OCR Scanning

**Use Case:** Handwritten lyrics in notebook

**Flow:**
```
1. Tap camera button
2. Photo napkin/notebook
3. OCR processing (2-3 seconds)
4. Review detected text
5. Quick edit to fix errors
6. Tap "Import" to add to editor
```

**UI (Photo Review):**
```
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚ [Photo of notebook page]            â”‚
â”‚                                     â”‚
â”‚ Detected Text (tap to edit):        â”‚
â”‚ â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â” â”‚
â”‚ â”‚ The scent of cedar still        â”‚ â”‚
â”‚ â”‚ haunts her bed                   â”‚ â”‚
â”‚ â”‚ I find her collar in the drawer  â”‚ â”‚
â”‚ â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜ â”‚
â”‚                                     â”‚
â”‚ [âœ“ Import]  [ðŸ“· Retake]  [âœ— Cancel] â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
```

**Technical:**
- MediaDevices API (camera access)
- Tesseract.js (client-side OCR)
- Target: 85%+ accuracy for printed text
- 70%+ accuracy for handwriting

---

### 4. ðŸŽ¯ Tap Tempo Input

**Use Case:** Set song BPM naturally

**Flow:**
```
1. Tap tempo button (â™ª)
2. Tap tempo mode activates
3. Tap screen to rhythm (4+ taps)
4. BPM calculated live
5. Confidence meter shows accuracy
6. Tap "Set Tempo" to save
```

**UI:**
```
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚ Tap to the Rhythm                   â”‚
â”‚                                     â”‚
â”‚         124 BPM                     â”‚ â† Live updating
â”‚       â–ˆâ–ˆâ–ˆâ–ˆâ–ˆâ–ˆâ–ˆâ–ˆâ–‘â–‘                    â”‚ â† Confidence
â”‚                                     â”‚
â”‚ [Tap anywhere on screen]            â”‚
â”‚                                     â”‚
â”‚ [âœ“ Set Tempo]  [âŸ³ Reset]           â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
```

**Technical:**
- DeviceMotion API or touch events
- Calculate average interval
- Min 4 taps, max 16 taps
- Saves to YAML: `tempo: 124`

---

### 5. ðŸ“Š Visual Prosody Overlay

**Use Case:** Check syllable stress while writing

**Display:**
```
The scent of cedar still haunts her bed
Ë˜    ËŠ    Ë˜  ËŠ  Ë˜   ËŠ     Ë˜    ËŠ
8 syllables âœ“

I find her collar in the drawer sometimes
Ë˜ ËŠ   Ë˜   ËŠ   Ë˜   Ë˜  ËŠ     Ë˜   ËŠ
10 syllables âš ï¸ (target: 8-9)
```

**Features:**
- Live stress marks (Ë˜ = unstressed, ËŠ = stressed)
- Color-coded feedback (green = good, yellow = warning)
- Syllable count per line
- Toggle on/off (Settings)

**Technical:**
- Pronouncing library (or port to JS)
- CSS overlays for stress marks
- Real-time analysis (debounced)

---

### 6. ðŸ“ˆ Live Syllable Counter

**Use Case:** Stay in meter while typing

**Display:**
```
[VERSE 1]                    Target: 8-9 syllables
The scent of cedar still haunts her bed
â””â”€ 8 syllables âœ“

I find her collar in the drawer sometimes yeah
â””â”€ 12 syllables âš ï¸ Too long
```

**Features:**
- Live count as you type
- Color indicators (green/yellow/red)
- Configurable target range
- Per-line analysis

**Technical:**
- Syllable-counting library
- Keystroke debouncing (300ms)
- Instant visual feedback

---

### 7. ðŸŽµ Inline Waveform Display

**Use Case:** See voice memos visually

**Display:**
```
[Voice Memo 1]
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚ â–â–‚â–ƒâ–…â–‡â–ˆâ–‡â–…â–ƒâ–‚â–â–‚â–ƒâ–…â–‡â–ˆâ–‡â–…â–ƒâ–‚â–  0:12 / 0:18 â”‚
â”‚          â–²                          â”‚ â† Scrubber
â”‚ [â–¶ï¸ Play] [ðŸ—‘ï¸] [â†“ Download]         â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
```

**Features:**
- Visual waveform (Canvas/SVG)
- Scrubbing (tap to jump)
- Play/pause control
- Download to camera roll

**Technical:**
- Web Audio API (analyzer)
- Canvas for waveform rendering
- Touch events for scrubbing

---

### 8. ðŸŽ¨ Rhyme Highlighting

**Use Case:** See rhyme patterns visually

**Display:**
```
[CHORUS]
I still see her shadow by the door    A â”
Her collar jingles in my dreams        B â”‚ Color-coded
The scent of cedar never felt so       C â”‚ (subtle tint)
The house is quiet, but it's not       A â”˜
```

**Features:**
- Detect end rhymes
- Assign colors per rhyme (A=pink, B=blue, etc.)
- Very subtle background tint (5% opacity)
- Uses TUI color palette
- Toggle on/off

**Technical:**
- Pronouncing library (rhyme detection)
- CSS background colors
- Terminal color scheme

---

### 9. ðŸ“¤ Rich Export Formats

**Use Case:** Share draft with bandmate

**Formats:**
1. **PNG Image** â€” Styled lyrics (terminal aesthetic)
2. **PDF** â€” Print-friendly (US Letter)
3. **Plain Text** â€” Universal compatibility
4. **Markdown** â€” For GitHub, notes apps

**Example PNG Export:**
```
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚ noise.sh                    Draft â”‚ â† Terminal header
â”‚â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”‚
â”‚ Lost My Dog                         â”‚
â”‚ Key: C | Tempo: 84 BPM              â”‚
â”‚â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”‚
â”‚                                     â”‚
â”‚ [VERSE 1]                           â”‚
â”‚ The scent of cedar still haunts her â”‚
â”‚ bed                                 â”‚
â”‚ I find her collar in the drawer...  â”‚
â”‚                                     â”‚
â”‚â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”‚
â”‚ Written Oct 17, 2025                â”‚
â”‚ noise.sh.app                      â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
```

**Technical:**
- html-to-image library
- CSS print media queries (PDF)
- Native file download API

---

### 10. ðŸ“± Share Sheet Integration

**Use Case:** Send to Messages, Mail, AirDrop

**Flow:**
```
1. Tap share button (ðŸ“¤)
2. Native share sheet appears
3. Choose: Messages, Mail, AirDrop, Notes, etc.
4. Shares as text or image (user choice)
```

**Technical:**
- Web Share API (`navigator.share()`)
- Share text, URL, or files
- Fallback: copy to clipboard

**iOS Share Options:**
- Messages
- Mail
- AirDrop
- Notes
- Copy
- Save to Files

---

### 11. ðŸ’¾ Automatic Sync (Invisible)

**Use Case:** Zero-effort file sync with desktop

**How It Works (Completely Automatic):**
```
1. You type or capture â†’ Auto-saves locally (500ms debounce)
2. When online â†’ Syncs to Dropbox/Drive (background)
3. TUI sees changes instantly
4. No buttons, no thinking, just works
```

**User Never Thinks About Sync:**
```
Write on phone â†’ Appears on desktop
Record voice memo â†’ Instantly in TUI
Scan with camera â†’ Available everywhere

Like Apple Notes or Google Docs:
- No save button
- No sync button
- Just works
```

**Sync Status (Subtle):**
```
Top-right corner (tiny, unobtrusive):
âœ“ Synced          (green dot)
â†» Syncing...      (spinning)
âš  Offline (saved) (yellow dot)
```

**File Format (YAML + Markdown):**
```yaml
---
title: "Lost My Dog"
key: "C"
tempo: 84
created_at: 2025-10-17T14:32:00Z
attachments:
  voice_memos:
    - melody_idea_1.webm
---

[VERSE 1]
The scent of cedar still haunts her bed
I find her collar in the drawer sometimes
...
```

**Technical:**
- Auto-save on keystroke (500ms debounce)
- Background Sync API (syncs when online)
- Dropbox/Drive API (one-time OAuth setup)
- Conflict resolution (last-write-wins)
- Queue pending changes when offline
- Retry failed syncs automatically

---

## Automatic Sync Architecture

### How Sync Works (Invisible to User)

**Like Apple Notes or Google Docs:**
```
User types "The scent of cedar..."
â†“
500ms debounce (wait for pause)
â†“
Auto-save to IndexedDB (local, instant)
â†“
Queue sync to Dropbox/Drive (background)
â†“
When online: Push to cloud (Background Sync API)
â†“
TUI sees file change (instant)
â†“
Status indicator: âœ“ Synced (briefly)
```

**User Experience:**
- No save button
- No sync button  
- No manual actions
- Just type and it's everywhere

**Technical Flow:**
```javascript
// On keystroke (debounced 500ms)
const autoSave = debounce(async (content) => {
  // 1. Save locally (instant)
  await saveToIndexedDB(content);
  
  // 2. Queue sync (background)
  if ('sync' in registration) {
    await registration.sync.register('sync-lyrics');
  }
  
  // 3. Sync happens automatically when online
}, 500);

// Background sync (Service Worker)
self.addEventListener('sync', async (event) => {
  if (event.tag === 'sync-lyrics') {
    event.waitUntil(syncToDropbox());
  }
});
```

**Sync States:**
- **âœ“ Synced** â€” Green dot, everything current
- **â†» Syncing...** â€” Blue spinner, upload in progress
- **âš  Offline** â€” Yellow dot, queued for later
- **âœ— Error** â€” Red dot, retry automatically

**Setup (One-Time Only):**
```
First launch:
1. Tap "Connect Dropbox" (or Google Drive)
2. OAuth login (once)
3. Choose folder: /noise.sh/
4. Done. Never think about it again.
```

**No Sync Button Anywhere:**
- PWA interface: No sync button
- Settings: Just shows sync status
- User never manually syncs
- It just works

---

## Design System (TUI Aesthetic)

### Color Palette (Exact Match)

```css
/* Terminal Colors from Bubble Tea TUI */
--bg-primary:    #1a1b26;  /* Main background */
--bg-black:      #000000;  /* OLED optimization */
--bg-elevated:   #24283b;  /* Panels */

/* Mode Colors (EXACT) */
--mode-sketch:   #f7768e;  /* Pink */
--mode-draft:    #7aa2f7;  /* Blue */
--mode-polish:   #9ece6a;  /* Green */

/* Text Colors */
--text-primary:   #c0caf5;
--text-secondary: #a9b1d6;
--text-tertiary:  #565f89;

/* Borders */
--border-default: #3b4261;
```

### Typography

```css
--font-mono: 'JetBrains Mono', 'SF Mono', monospace;
--font-sans: -apple-system, BlinkMacSystemFont, sans-serif;

--text-base: 16px;  /* Never <16px on mobile */
--text-lg:   18px;
```

### Layout (Terminal-Inspired)

- Box-drawing characters for borders
- Monospace for lyrics
- Grid-aligned (8px base)
- Minimal chrome
- Touch targets: 44px minimum

---

## Technical Stack

### Frontend (PWA)

**Framework:**
- **Next.js 14** (App Router, SSG for offline)
- **React 18** (minimal components)
- **TypeScript** (type safety)

**Styling:**
- **Tailwind CSS** (TUI color tokens)
- **CSS Modules** (terminal components)

**Sensor APIs:**
- **Web Audio API** â€” Voice recording, waveforms
- **MediaRecorder API** â€” Audio encoding
- **Web Speech API** â€” Voice-to-text
- **MediaDevices API** â€” Camera access
- **DeviceMotion API** â€” Tap tempo
- **Vibration API** â€” Haptic feedback (optional)

**PWA APIs:**
- **Service Worker** â€” Offline caching
- **IndexedDB** â€” Local storage
- **Web App Manifest** â€” Install to home screen
- **Background Sync API** â€” Auto-sync when online
- **Web Share API** â€” Native sharing
- **File System Access API** â€” Dropbox/Drive sync

**Libraries:**
- **Tesseract.js** â€” Client-side OCR
- **html-to-image** â€” Export to PNG
- **syllable-counter** â€” Syllable counting
- **Canvas API** â€” Waveform visualization

---

### Backend (Minimal/None)

**Option 1: No Backend (Recommended)**
- Everything client-side
- OCR via Tesseract.js (in browser)
- Sync via Dropbox/Drive API (OAuth)
- No server costs

**Option 2: Optional Server (If OCR too slow)**
- Python FastAPI endpoint
- Tesseract server-side processing
- Only for OCR, everything else client-side

---

### Performance Budget

```
Bundle Size:     <75KB gzipped
First Load:      <1.5 seconds (3G)
Time Interactive: <2 seconds
Offline:         100% functional
Memory:          <40MB
Battery:         <5% per 30min session
```

---

## Implementation Roadmap

### Phase 0: Design & Prototype (2 weeks)

**Activities:**
- [ ] Extract exact TUI colors to Tailwind config
- [ ] Design terminal-styled components (Figma)
- [ ] Build clickable prototype (Vercel)
- [ ] Test sensor APIs (mic, camera, motion)
- [ ] User testing (5 people)

**Deliverables:**
- Tailwind config with TUI colors
- Component library (minimal)
- Working prototype
- Sensor accuracy tests

---

### Phase 1: Core Capture (3 weeks)

**Features:**
- [ ] Minimal text editor (TUI colors)
- [ ] Voice memo recording
- [ ] Live syllable counter
- [ ] File save/load (IndexedDB)
- [ ] PWA setup (offline mode)

**Quality Bar:**
- Works 100% offline
- TUI aesthetic perfect
- <75KB bundle
- Smooth 60fps

---

### Phase 2: Visual Enhancements (2 weeks)

**Features:**
- [ ] Visual prosody overlay
- [ ] Inline waveforms
- [ ] Rhyme highlighting
- [ ] Tap tempo mode

**Quality Bar:**
- Visuals feel native, not gimmicky
- Toggle on/off
- No performance impact

---

### Phase 3: Capture Tools (2 weeks)

**Features:**
- [ ] Voice-to-text dictation
- [ ] Camera OCR scanning
- [ ] Haptic feedback (optional)

**Quality Bar:**
- Clear permissions
- OCR >85% accuracy
- Voice recognition >90%

---

### Phase 4: Invisible Sync (1 week)

**Features:**
- [ ] Auto-save on keystroke (500ms debounce)
- [ ] Dropbox/Drive OAuth (one-time setup)
- [ ] Background sync (automatic, no user action)
- [ ] Sync status indicator (subtle, top-right)
- [ ] Share sheet integration
- [ ] Rich exports (PNG, PDF)

**Quality Bar:**
- Zero-effort sync (user never thinks about it)
- Works exactly like Apple Notes
- No sync button anywhere
- Status indicator subtle, unobtrusive
- Conflict resolution automatic

---

### Phase 5: Polish & Launch (1 week)

**Activities:**
- [ ] Beta testing (10+ users)
- [ ] Performance optimization
- [ ] Accessibility audit
- [ ] Documentation
- [ ] Deploy to production

**Quality Bar:**
- Zero critical bugs
- Lighthouse >95
- Works on all devices

---

**Total Timeline:** 11 weeks

---

## Success Metrics

### Adoption
- [ ] >50% of users install PWA to home screen
- [ ] Average 2+ captures per week per user
- [ ] >98% of captures sync automatically (no user action)
- [ ] <1% of users experience sync failures

### Feature Usage
- [ ] >60% use voice memo feature
- [ ] >40% use OCR feature
- [ ] >80% use visual prosody overlay
- [ ] >30% use tap tempo

### Quality
- [ ] OCR accuracy >85%
- [ ] Voice recognition >90%
- [ ] Auto-sync success rate >98%
- [ ] Sync happens <3 seconds when online
- [ ] Battery drain <5% per 30min
- [ ] Lighthouse score >95

### User Sentiment
- [ ] "Perfect companion to TUI" feedback
- [ ] "I never think about syncing, it just works"
- [ ] "Capture on subway, polish at desk"
- [ ] "Like Apple Notes but for songwriting"

---

## Scope Control

### What We're Building
âœ… Lightweight capture tool  
âœ… Visual feedback companion  
âœ… Mobile-first experience  
âœ… Sync to TUI  

### What We're NOT Building
âŒ Full writing environment  
âŒ AI features (use TUI)  
âŒ Knowledge base (use TUI)  
âŒ Multi-song management  
âŒ Complex editor  
âŒ Desktop-first design  

**The Rule:** If it takes more than 2 minutes, hand it off to the TUI.

---

## File Format (Same as TUI)

```yaml
---
title: "Lost My Dog"
key: "C"
tempo: 84             # From tap tempo
created_at: 2025-10-17T14:32:00Z

attachments:
  voice_memos:
    - melody_idea_1.webm    # From mic
  scanned_lyrics:
    - notebook_page_1.jpg    # From camera
---

[VERSE 1]
The scent of cedar still haunts her bed
I find her collar in the drawer sometimes
...
```

**Storage:** Same Dropbox/Drive folder as TUI

---

## Privacy

### Permissions Required

**iOS:**
```xml
<key>NSMicrophoneUsageDescription</key>
<string>Record melody ideas and voice notes</string>

<key>NSCameraUsageDescription</key>
<string>Scan handwritten lyrics</string>

<key>NSMotionUsageDescription</key>
<string>Detect tap tempo</string>
```

**Android:**
```xml
<uses-permission android:name="android.permission.RECORD_AUDIO" />
<uses-permission android:name="android.permission.CAMERA" />
```

### Privacy Principles
1. **Explicit opt-in** â€” Never assume permissions
2. **Local-first** â€” All processing on device
3. **No tracking** â€” Zero analytics
4. **User control** â€” Easy to disable features

---

## Next Steps

### This Week
- [ ] Review spec with developer
- [ ] Extract TUI colors to CSS
- [ ] Choose Next.js final (confirmed?)
- [ ] Set up Tailwind config

### Phase 0 Start (Week 1-2)
- [ ] Figma design system
- [ ] Terminal-styled components
- [ ] Working prototype
- [ ] Sensor API tests

### Development Start (Week 3)
- [ ] Next.js scaffolding
- [ ] Minimal editor
- [ ] Voice memo recording
- [ ] PWA setup

---

## Related Documents

- **PRD.md** â€” Product requirements (TUI)
- **TDD.md** â€” Technical design (TUI)
- **Enhancement #2** â€” Rapid Prototyping (TUI)
- **DEVELOPMENT_ROADMAP.md** â€” Timeline

---

**Document Status:** âœ… APPROVED

**Core Vision:** Lightweight mobile companion that captures ideas (voice, camera, tempo) and provides visual feedback (prosody, syllables, waveforms). **Syncs invisibly and automatically** to TUI for serious writing work.

**Sync Philosophy:** Like Apple Notes or Google Docs â€” you never think about it, never tap a button, it just works.

**Timeline:** 11 weeks to production-ready PWA

**Budget:** $0/month ongoing costs (no backend)

---

_Capture on mobile. Syncs automatically. Create on desktop. The perfect workflow._ ðŸŽµðŸ“±ðŸ’»
