# Enhancement #3+5: Mobile Companion App (PWA)
**noise.sh - Capture Tool for On-the-Go Ideas**

**Version:** 3.0 (Companion App)  
**Status:** READY FOR IMPLEMENTATION  
**Date:** October 17, 2025  
**Owner:** Simon (Puente Labs)  
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
✅ **Capture device** — Record voice memos, scan handwritten notes, tap tempo  
✅ **Visual companion** — See prosody, syllables, waveforms (terminal can't show)  
✅ **Mobile tool** — Works offline on phone, optimized for quick capture  
✅ **Sync bridge** — Saves to files that TUI can open  

### What This IS NOT
❌ **Full writing environment** — No AI panel, no knowledge base, no multi-song management  
❌ **TUI replacement** — Complement, not compete  
❌ **Complex editor** — Minimal editing only  
❌ **Desktop-first** — Mobile is the focus, desktop is bonus  

---

## Division of Labor: PWA vs TUI

### PWA Does (Capture & Feedback)

**Capture Tools:**
1. 🎤 **Voice memo recording** — Hum melodies, capture ideas
2. 🎤 **Voice-to-text dictation** — Speak lyrics rapidly
3. 📷 **Camera OCR scanning** — Digitize handwritten notebooks
4. 🎯 **Tap tempo input** — Set BPM by tapping rhythm

**Visual Feedback:**
5. 📊 **Visual prosody overlay** — See stress marks (˘ ˊ) above syllables
6. 📈 **Live syllable counter** — Real-time count as you type
7. 🎵 **Inline waveforms** — View voice memos visually
8. 🎨 **Rhyme highlighting** — Color-coded rhyme schemes

**Mobile Sharing:**
9. 📤 **Rich exports** — Styled PNG/PDF for sharing
10. 📱 **Share sheet integration** — Send via Messages, Mail, AirDrop
11. 📋 **Formatted clipboard** — Paste with styling

**Basic Editing:**
12. ✏️ **Minimal text editor** — Refine captured text
13. 💾 **Invisible auto-sync** — Saves & syncs automatically (like Apple Notes)

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
📱 Subway (no internet)
↓
Open PWA from home screen
↓
Tap mic button → hum melody (10 sec)
↓
Type quick lyric idea:
  "The scent of cedar still haunts her bed"
  └─ 8 syllables ✓ (live counter shows)
↓
Tap tempo to rhythm (84 BPM detected)
↓
Auto-saves to local storage (IndexedDB)
↓
Status shows: ⚠ Offline (will sync later)
```

**Step 2: Automatic Sync (Invisible)**
```
📱 Exit subway, phone gets signal
↓
PWA detects connection
↓
Background sync starts automatically
↓
Status briefly shows: ↻ Syncing...
↓
Status changes to: ✓ Synced
↓
File appears in Dropbox/Drive instantly

You never tapped a sync button.
You never thought about it.
It just worked.
```

**Step 3: Open on Desktop TUI (Instant)**
```
💻 Open TUI (minutes or hours later)
↓
See new draft "Lost My Dog" immediately
↓
Play voice memo (melody idea)
↓
Use AI to expand verse
↓
Add knowledge base insights
↓
Polish with full editor
↓
Export final version
```

**Total User Actions for Sync:** Zero  
**Sync happens:** Automatically in background  
**Like:** Apple Notes, Google Docs, Notion

---

## Mobile PWA Interface (Minimal)

### Quick Capture Screen

```
┌─────────────────────────────────────┐
│ 🎵 Quick Capture    [SKETCH]  ✓    │ ← Status: ✓ Synced
├─────────────────────────────────────┤
│                                     │
│ The scent of cedar still...        │ ← Basic text editor
│ haunts her bed                      │   Auto-saves on keystroke
│ ˘    ˊ    ˘  ˊ  ˘   ˊ     ˘   ˊ    │ ← Prosody overlay
│ 8 syllables ✓                       │ ← Live counter
│                                     │
│ I find her collar in the drawer...  │
│█                                    │ ← Blinking cursor
│                                     │
│ ┌─────────────────────────────────┐ │
│ │ Voice Memo 1       ▁▂▃▅▇█▇▅▃▂▁ │ │ ← Inline waveform
│ │ [▶️ Play] 0:08                  │ │
│ └─────────────────────────────────┘ │
│                                     │
├─────────────────────────────────────┤
│ [🎤] [📷] [♪] [📤]                  │ ← Capture tools only
│  Mic  Scan Tempo Share              │   No manual sync button
└─────────────────────────────────────┘

Sync Status Indicators (top-right):
✓    = Synced to TUI (green)
↻    = Syncing... (blue, spinning)
⚠    = Offline, will sync when online (yellow)
✗    = Sync failed, retrying... (red, rare)

Interactions:
- Type → Auto-saves locally (500ms debounce)
- Save → Auto-syncs to Dropbox/Drive
- No buttons needed, just works
```

### What's NOT Included
- ❌ No AI panel (use TUI for that)
- ❌ No song library (single-draft focus)
- ❌ No knowledge base search
- ❌ No complex menus/settings
- ❌ No multi-window layout
- ❌ No manual sync button (automatic)

**Design Philosophy:** "Just enough to capture and refine, syncs invisibly to TUI."

---

## Core Features (Detailed)

### 1. 🎤 Voice Memo Recording

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
┌─────────────────────────────────────┐
│ Recording... ⬤ 0:08                 │
│ ▁▂▃▅▇█▇▅▃▂▁▂▃▅▇█▇▅▃▂▁              │ ← Live waveform
│                                     │
│ [■ Stop]  [✓ Keep]  [✗ Cancel]     │
└─────────────────────────────────────┘
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

### 2. 🎤 Voice-to-Text Dictation

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

### 3. 📷 Camera OCR Scanning

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
┌─────────────────────────────────────┐
│ [Photo of notebook page]            │
│                                     │
│ Detected Text (tap to edit):        │
│ ┌─────────────────────────────────┐ │
│ │ The scent of cedar still        │ │
│ │ haunts her bed                   │ │
│ │ I find her collar in the drawer  │ │
│ └─────────────────────────────────┘ │
│                                     │
│ [✓ Import]  [📷 Retake]  [✗ Cancel] │
└─────────────────────────────────────┘
```

**Technical:**
- MediaDevices API (camera access)
- Tesseract.js (client-side OCR)
- Target: 85%+ accuracy for printed text
- 70%+ accuracy for handwriting

---

### 4. 🎯 Tap Tempo Input

**Use Case:** Set song BPM naturally

**Flow:**
```
1. Tap tempo button (♪)
2. Tap tempo mode activates
3. Tap screen to rhythm (4+ taps)
4. BPM calculated live
5. Confidence meter shows accuracy
6. Tap "Set Tempo" to save
```

**UI:**
```
┌─────────────────────────────────────┐
│ Tap to the Rhythm                   │
│                                     │
│         124 BPM                     │ ← Live updating
│       ████████░░                    │ ← Confidence
│                                     │
│ [Tap anywhere on screen]            │
│                                     │
│ [✓ Set Tempo]  [⟳ Reset]           │
└─────────────────────────────────────┘
```

**Technical:**
- DeviceMotion API or touch events
- Calculate average interval
- Min 4 taps, max 16 taps
- Saves to YAML: `tempo: 124`

---

### 5. 📊 Visual Prosody Overlay

**Use Case:** Check syllable stress while writing

**Display:**
```
The scent of cedar still haunts her bed
˘    ˊ    ˘  ˊ  ˘   ˊ     ˘    ˊ
8 syllables ✓

I find her collar in the drawer sometimes
˘ ˊ   ˘   ˊ   ˘   ˘  ˊ     ˘   ˊ
10 syllables ⚠️ (target: 8-9)
```

**Features:**
- Live stress marks (˘ = unstressed, ˊ = stressed)
- Color-coded feedback (green = good, yellow = warning)
- Syllable count per line
- Toggle on/off (Settings)

**Technical:**
- Pronouncing library (or port to JS)
- CSS overlays for stress marks
- Real-time analysis (debounced)

---

### 6. 📈 Live Syllable Counter

**Use Case:** Stay in meter while typing

**Display:**
```
[VERSE 1]                    Target: 8-9 syllables
The scent of cedar still haunts her bed
└─ 8 syllables ✓

I find her collar in the drawer sometimes yeah
└─ 12 syllables ⚠️ Too long
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

### 7. 🎵 Inline Waveform Display

**Use Case:** See voice memos visually

**Display:**
```
[Voice Memo 1]
┌─────────────────────────────────────┐
│ ▁▂▃▅▇█▇▅▃▂▁▂▃▅▇█▇▅▃▂▁  0:12 / 0:18 │
│          ▲                          │ ← Scrubber
│ [▶️ Play] [🗑️] [↓ Download]         │
└─────────────────────────────────────┘
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

### 8. 🎨 Rhyme Highlighting

**Use Case:** See rhyme patterns visually

**Display:**
```
[CHORUS]
I still see her shadow by the door    A ┐
Her collar jingles in my dreams        B │ Color-coded
The scent of cedar never felt so       C │ (subtle tint)
The house is quiet, but it's not       A ┘
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

### 9. 📤 Rich Export Formats

**Use Case:** Share draft with bandmate

**Formats:**
1. **PNG Image** — Styled lyrics (terminal aesthetic)
2. **PDF** — Print-friendly (US Letter)
3. **Plain Text** — Universal compatibility
4. **Markdown** — For GitHub, notes apps

**Example PNG Export:**
```
┌─────────────────────────────────────┐
│ noise.sh                    Draft │ ← Terminal header
│─────────────────────────────────────│
│ Lost My Dog                         │
│ Key: C | Tempo: 84 BPM              │
│─────────────────────────────────────│
│                                     │
│ [VERSE 1]                           │
│ The scent of cedar still haunts her │
│ bed                                 │
│ I find her collar in the drawer...  │
│                                     │
│─────────────────────────────────────│
│ Written Oct 17, 2025                │
│ noise.sh.app                      │
└─────────────────────────────────────┘
```

**Technical:**
- html-to-image library
- CSS print media queries (PDF)
- Native file download API

---

### 10. 📱 Share Sheet Integration

**Use Case:** Send to Messages, Mail, AirDrop

**Flow:**
```
1. Tap share button (📤)
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

### 11. 💾 Automatic Sync (Invisible)

**Use Case:** Zero-effort file sync with desktop

**How It Works (Completely Automatic):**
```
1. You type or capture → Auto-saves locally (500ms debounce)
2. When online → Syncs to Dropbox/Drive (background)
3. TUI sees changes instantly
4. No buttons, no thinking, just works
```

**User Never Thinks About Sync:**
```
Write on phone → Appears on desktop
Record voice memo → Instantly in TUI
Scan with camera → Available everywhere

Like Apple Notes or Google Docs:
- No save button
- No sync button
- Just works
```

**Sync Status (Subtle):**
```
Top-right corner (tiny, unobtrusive):
✓ Synced          (green dot)
↻ Syncing...      (spinning)
⚠ Offline (saved) (yellow dot)
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
↓
500ms debounce (wait for pause)
↓
Auto-save to IndexedDB (local, instant)
↓
Queue sync to Dropbox/Drive (background)
↓
When online: Push to cloud (Background Sync API)
↓
TUI sees file change (instant)
↓
Status indicator: ✓ Synced (briefly)
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
- **✓ Synced** — Green dot, everything current
- **↻ Syncing...** — Blue spinner, upload in progress
- **⚠ Offline** — Yellow dot, queued for later
- **✗ Error** — Red dot, retry automatically

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
- **Web Audio API** — Voice recording, waveforms
- **MediaRecorder API** — Audio encoding
- **Web Speech API** — Voice-to-text
- **MediaDevices API** — Camera access
- **DeviceMotion API** — Tap tempo
- **Vibration API** — Haptic feedback (optional)

**PWA APIs:**
- **Service Worker** — Offline caching
- **IndexedDB** — Local storage
- **Web App Manifest** — Install to home screen
- **Background Sync API** — Auto-sync when online
- **Web Share API** — Native sharing
- **File System Access API** — Dropbox/Drive sync

**Libraries:**
- **Tesseract.js** — Client-side OCR
- **html-to-image** — Export to PNG
- **syllable-counter** — Syllable counting
- **Canvas API** — Waveform visualization

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
✅ Lightweight capture tool  
✅ Visual feedback companion  
✅ Mobile-first experience  
✅ Sync to TUI  

### What We're NOT Building
❌ Full writing environment  
❌ AI features (use TUI)  
❌ Knowledge base (use TUI)  
❌ Multi-song management  
❌ Complex editor  
❌ Desktop-first design  

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
1. **Explicit opt-in** — Never assume permissions
2. **Local-first** — All processing on device
3. **No tracking** — Zero analytics
4. **User control** — Easy to disable features

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

- **PRD.md** — Product requirements (TUI)
- **TDD.md** — Technical design (TUI)
- **Enhancement #2** — Rapid Prototyping (TUI)
- **DEVELOPMENT_ROADMAP.md** — Timeline

---

**Document Status:** ✅ APPROVED

**Core Vision:** Lightweight mobile companion that captures ideas (voice, camera, tempo) and provides visual feedback (prosody, syllables, waveforms). **Syncs invisibly and automatically** to TUI for serious writing work.

**Sync Philosophy:** Like Apple Notes or Google Docs — you never think about it, never tap a button, it just works.

**Timeline:** 11 weeks to production-ready PWA

**Budget:** $0/month ongoing costs (no backend)

---

_Capture on mobile. Syncs automatically. Create on desktop. The perfect workflow._ 🎵📱💻
