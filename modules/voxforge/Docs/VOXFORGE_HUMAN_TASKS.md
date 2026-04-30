# VOXFORGE - TASKS FOR HUMANS (NON-CODING)
**Version:** 2.0  
**Last Updated:** November 9, 2025  
**Purpose:** All the things YOU need to do (not the AI coding agent)

---

## OVERVIEW

While your AI coding agent handles the development, **these tasks require a human.** They include:
- Testing and feedback
- Content creation
- Account setup
- Design decisions
- Marketing activities

**Estimated Total Time:** 10-15 hours across 12 days

---

## SETUP TASKS (Before Day 1)

### Task 1: Get API Keys (Optional, for Lyrics Feature)

**Why:** Needed for lyrics generation  
**When:** Before Day 11 (or skip if not implementing)  
**Estimated Time:** 15 minutes  
**Cost:** OpenAI API ~$5-20/month depending on usage

**Steps:**

**Option A: OpenAI (Recommended)**
1. Go to https://platform.openai.com/
2. Sign up / Log in
3. Click "API keys" in sidebar
4. Click "Create new secret key"
5. Name it "VoxForge"
6. Copy the key (starts with `sk-`)
7. Save in `.env.local` file:
   ```
   OPENAI_API_KEY=sk-...your-key-here
   ```

**Option B: Anthropic Claude**
1. Go to https://console.anthropic.com/
2. Sign up / Log in
3. Click "API Keys"
4. Click "Create Key"
5. Copy the key (starts with `sk-ant-`)
6. Save in `.env.local` file:
   ```
   ANTHROPIC_API_KEY=sk-ant-...your-key-here
   ```

**Option C: Skip It**
- Don't implement lyrics feature
- Save $5-20/month
- Lose a cool feature

**Decision:** ⬜ OpenAI | ⬜ Anthropic | ⬜ Skip

---

### Task 2: Set Up GitHub Repository

**Why:** Version control and deployment  
**When:** Before Day 1  
**Estimated Time:** 10 minutes  
**Cost:** Free

**Steps:**
1. Go to https://github.com/new
2. Repository name: `voxforge`
3. Description: "Voice-to-song tool that transforms vocals into music"
4. Select: ⬜ Public | ⬜ Private
5. Check "Add README"
6. Create repository
7. Copy the repository URL

**In terminal:**
```bash
cd voxforge
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin YOUR_REPO_URL
git push -u origin main
```

**Repository URL:** _______________________________

---

### Task 3: Choose Deployment Platform

**Why:** Need a place to host the live app  
**When:** Before Day 12  
**Estimated Time:** 5 minutes  
**Cost:** Free tier available for all options

**Options:**

**Railway (Recommended)**
- ✅ Easy to use
- ✅ Generous free tier ($5 credit)
- ✅ Auto-deploys from GitHub
- ✅ Good for fullstack apps
- Website: https://railway.app/

**Vercel**
- ✅ Made by Next.js creators
- ✅ Very fast CDN
- ✅ Great for frontend
- ⚠️ Request timeouts (10s) may affect lyrics API
- Website: https://vercel.com/

**Netlify**
- ✅ Easy setup
- ✅ Continuous deployment
- ✅ Good free tier
- Website: https://www.netlify.com/

**Decision:** ⬜ Railway | ⬜ Vercel | ⬜ Netlify

**Account created:** ⬜ Yes | ⬜ No

---

## TESTING TASKS (Throughout Development)

### Task 4: Record Test Vocals

**Why:** Need consistent test cases for development  
**When:** Day 1  
**Estimated Time:** 30 minutes  
**Cost:** Free (use your voice!)

**What to Record:**

1. **Simple Scale (8 seconds)**
   - Sing: "Do Re Mi Fa Sol La Ti Do"
   - Or: "C D E F G A B C"
   - Save as: `test-scale.webm`

2. **Happy Birthday (15 seconds)**
   - First verse only
   - Clear pronunciation
   - Save as: `test-happybirthday.webm`

3. **Your Original Melody (30 seconds)**
   - Make up something
   - Can be nonsense syllables
   - Save as: `test-original.webm`

4. **Humming (20 seconds)**
   - No words, just hum a tune
   - Save as: `test-humming.webm`

5. **Beatboxing (15 seconds)**
   - Basic beat pattern
   - Save as: `test-beatbox.webm`

**Where to Save:** `public/test-recordings/` folder

**Tips:**
- Use headphones to hear yourself
- Record in quiet environment
- Don't worry about quality - these are test files

**Completed:** ⬜ Yes | ⬜ No

---

### Task 5: Daily Manual Testing

**Why:** Catch bugs early  
**When:** Every day after coding session  
**Estimated Time:** 15-30 minutes per day  
**Cost:** Free

**Test Checklist:**

**Day 1: Recording**
- [ ] Click "Record" button
- [ ] Grant microphone permission
- [ ] Speak into mic, see waveform move
- [ ] Click "Stop"
- [ ] See static waveform
- [ ] Click "Play Back"
- [ ] Hear your recording

**Day 2: Pitch Detection**
- [ ] Record simple scale
- [ ] See analysis complete
- [ ] Verify notes detected (C, D, E, F, G, A, B, C)
- [ ] Check average frequency makes sense
- [ ] No console errors

**Day 3: BPM & Key**
- [ ] Clap steady rhythm
- [ ] BPM detected (±5 of actual)
- [ ] Sing in C Major
- [ ] Key detected correctly
- [ ] Try different keys

**Day 4: Drums**
- [ ] Generate drums
- [ ] Click play
- [ ] Drums play at detected BPM
- [ ] Try simple/moderate/busy patterns
- [ ] Metronome syncs

**Day 5: Bass & Chords**
- [ ] Full arrangement plays
- [ ] All instruments in sync
- [ ] Sounds musical
- [ ] No clipping/distortion

**Day 6: Instrument Selection**
- [ ] Toggle instruments on/off
- [ ] Changes apply immediately
- [ ] Can't unselect all

**Day 7: Sections**
- [ ] Record multiple sections
- [ ] Each displays correctly
- [ ] Delete works
- [ ] Rename works

**Day 8: Arrangement**
- [ ] Sequential mode: sections play in order
- [ ] Layered mode: sections play together
- [ ] Both sound correct

**Day 9: Export**
- [ ] Export each stem type
- [ ] Files download
- [ ] Files open in music app
- [ ] Stems are aligned

**Day 10: MIDI**
- [ ] MIDI exports
- [ ] Opens in DAW
- [ ] Notes match melody

**Day 11: Lyrics (Optional)**
- [ ] Lyrics generate
- [ ] Match melody rhythm
- [ ] Copy to clipboard works

**Day 12: Full Test**
- [ ] Complete workflow start to finish
- [ ] No errors
- [ ] Smooth UX

---

### Task 6: Browser Compatibility Testing

**Why:** Users have different browsers  
**When:** Day 12  
**Estimated Time:** 1 hour  
**Cost:** Free

**Browsers to Test:**

**Desktop:**
- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (Mac only)
- [ ] Edge (latest)

**Mobile (if accessible):**
- [ ] Mobile Chrome (Android)
- [ ] Mobile Safari (iOS)

**For Each Browser:**
1. Open the app
2. Try to record
3. Note if it works or fails
4. Document any issues

**Results:**

| Browser | Works? | Issues |
|---------|--------|--------|
| Chrome  |        |        |
| Firefox |        |        |
| Safari  |        |        |
| Edge    |        |        |
| Mobile Chrome |  |        |
| Mobile Safari |  |        |

---

### Task 7: Get Feedback from Friends

**Why:** Fresh eyes catch issues  
**When:** Day 10-11  
**Estimated Time:** 2 hours  
**Cost:** Free (maybe buy them coffee?)

**Who to Ask:**
- 3-5 friends
- Mix of musicians and non-musicians
- Mix of technical and non-technical

**What to Do:**
1. Send them the link
2. Ask them to try making a song
3. **DO NOT HELP** (watch them struggle)
4. Take notes on where they get stuck
5. Ask for feedback after

**Questions to Ask:**
- What was confusing?
- What did you like?
- What did you expect to happen that didn't?
- Would you use this again?
- Any features missing?

**Feedback Notes:**

**Friend 1:**
- Name: _______________
- Musician? ⬜ Yes | ⬜ No
- Technical? ⬜ Yes | ⬜ No
- Notes:

**Friend 2:**
- Name: _______________
- Musician? ⬜ Yes | ⬜ No
- Technical? ⬜ Yes | ⬜ No
- Notes:

**Friend 3:**
- Name: _______________
- Musician? ⬜ Yes | ⬜ No
- Technical? ⬜ Yes | ⬜ No
- Notes:

---

## CONTENT CREATION TASKS

### Task 8: Write README

**Why:** First thing people see on GitHub  
**When:** Day 12  
**Estimated Time:** 1 hour  
**Cost:** Free

**Sections to Include:**

```markdown
# VoxForge

Transform your voice into music

## Features
- Voice recording with real-time visualization
- Automatic pitch, BPM, and key detection
- Generate drums, bass, and chords
- Export stems and MIDI
- (Optional) AI-powered lyrics assistance

## Demo
[Link to live site]

## Screenshots
[Add 3-5 screenshots]

## How to Use
1. Click "Record" and sing/hum a melody
2. Wait for analysis
3. Click "Generate Music"
4. Export your stems

## Tech Stack
- Next.js 15
- Tone.js (music generation)
- pitchfinder (pitch detection)
- realtime-bpm-analyzer (tempo)
- Tonal.js (music theory)

## Installation
```bash
npm install
npm run dev
```

## License
MIT
```

**Completed:** ⬜ Yes | ⬜ No

---

### Task 9: Take Screenshots

**Why:** Show off your work!  
**When:** Day 12  
**Estimated Time:** 30 minutes  
**Cost:** Free

**Screenshots Needed:**

1. **Hero/Landing**
   - Full interface, clean state
   - Show app title and description

2. **Recording in Progress**
   - Live waveform
   - Record button active

3. **Analysis Results**
   - Pitch detection display
   - BPM and key showing

4. **Music Generation**
   - Playback controls
   - Instrument picker
   - Metronome

5. **Export Panel**
   - All export options visible
   - Professional looking

**Tips:**
- Use full window (not just browser)
- Clean up desktop background
- Hide personal info
- Good lighting (if showing your face)
- Use macOS screenshot (Cmd+Shift+4) or Windows (Win+Shift+S)

**Screenshots saved to:** `docs/screenshots/`

**Completed:** ⬜ Yes | ⬜ No

---

### Task 10: Record Demo Video

**Why:** Show how it works  
**When:** Day 12  
**Estimated Time:** 1-2 hours  
**Cost:** Free

**Video Structure (2-3 minutes):**

**0:00 - 0:15 - Intro**
- "Hi, I'm [name]"
- "I built VoxForge, a voice-to-song tool"
- Show the website

**0:15 - 0:45 - Demo**
- Click record
- Sing a simple melody
- Show analysis happening
- Point out detected BPM, key

**0:45 - 1:15 - Music Generation**
- Click generate music
- Play the result
- Toggle instruments
- Show it sounds good

**1:15 - 1:45 - Multiple Sections**
- Record a second section
- Show arrangement modes
- Play sequential
- Play layered

**1:45 - 2:15 - Export**
- Export stems
- Show files download
- (Optional) Import into DAW

**2:15 - 2:30 - Outro**
- "Try it yourself at [URL]"
- "Built with [tech stack]"
- "Thanks for watching!"

**Tools:**
- **Screen Recording:** 
  - Mac: QuickTime (Cmd+Shift+5)
  - Windows: Xbox Game Bar (Win+G)
  - Cross-platform: OBS Studio (free)

- **Video Editing:**
  - iMovie (Mac, free)
  - DaVinci Resolve (cross-platform, free)
  - Clipchamp (web-based, free)

**Tips:**
- Write a script
- Practice a few times
- Use good microphone
- Edit out mistakes
- Add captions (optional)
- Background music (optional, copyright-free)

**Video uploaded to:** ⬜ YouTube | ⬜ Vimeo | ⬜ Other

**Video URL:** _______________________________

**Completed:** ⬜ Yes | ⬜ No

---

## DEPLOYMENT TASKS

### Task 11: Deploy to Production

**Why:** Share it with the world!  
**When:** Day 12  
**Estimated Time:** 30 minutes  
**Cost:** Free (on free tiers)

**Railway Deployment:**

1. Go to https://railway.app/
2. Sign up with GitHub
3. Click "New Project"
4. Select "Deploy from GitHub repo"
5. Choose `voxforge` repository
6. Wait for build (5-10 minutes)
7. Go to Settings → Domains
8. Click "Generate Domain"
9. Your app is live!

**Environment Variables (if using lyrics):**
1. Go to project
2. Click "Variables"
3. Add `OPENAI_API_KEY` or `ANTHROPIC_API_KEY`
4. Paste your key
5. Redeploy

**Vercel Deployment:**

```bash
npx vercel
```

Follow prompts, done!

**Netlify Deployment:**

```bash
npx netlify-cli deploy --prod
```

**Live URL:** _______________________________

**Completed:** ⬜ Yes | ⬜ No

---

### Task 12: Set Up Custom Domain (Optional)

**Why:** Looks more professional  
**When:** After deployment  
**Estimated Time:** 30 minutes  
**Cost:** $10-15/year for domain

**Steps:**

1. **Buy Domain:**
   - Namecheap: https://www.namecheap.com/
   - GoDaddy: https://www.godaddy.com/
   - Google Domains: https://domains.google/
   - Suggested: `voxforge.app`, `voxforge.io`, `voxforge.com`

2. **Point Domain to App:**
   - In Railway/Vercel/Netlify: Add custom domain
   - In domain registrar: Add DNS records
   - Wait 24 hours for DNS propagation

**Custom Domain:** _______________________________

**Completed:** ⬜ Yes | ⬜ No | ⬜ Skipped

---

## MARKETING TASKS (Optional)

### Task 13: Share on Social Media

**Why:** Get users and feedback  
**When:** After deployment  
**Estimated Time:** 1 hour  
**Cost:** Free

**Where to Share:**

**Twitter/X:**
```
🎵 Just launched VoxForge!

Transform your voice into music with:
✨ Auto pitch detection
🥁 Generated drums, bass, chords
📊 BPM & key detection
💾 Export stems & MIDI

Try it: [URL]

Built with Next.js + Tone.js + AI
#buildinpublic #indiedev #musictech
```

**Reddit:**
- r/webdev
- r/javascript
- r/musicproduction
- r/WeAreTheMusicMakers

**Post Template:**
```
I built VoxForge - a voice-to-song tool

[Demo GIF/Video]

Features:
- Record voice
- Auto-detect pitch, BPM, key
- Generate music arrangement
- Export stems + MIDI

Tech stack: Next.js, Tone.js, pitchfinder, Tonal.js

Try it: [URL]
GitHub: [URL]

Feedback welcome!
```

**Product Hunt (Advanced):**
- Launch when you have 20+ upvotes guaranteed
- Prepare graphics and copy
- Launch on Tuesday-Thursday morning
- Guide: https://www.producthunt.com/ship

**LinkedIn:**
- Share your journey
- Tag relevant people/companies
- Use professional tone

**Completed:**
- [ ] Twitter/X
- [ ] Reddit
- [ ] LinkedIn
- [ ] Product Hunt
- [ ] Hacker News

---

### Task 14: Write Blog Post

**Why:** SEO, portfolio, deeper explanation  
**When:** After launch  
**Estimated Time:** 2-4 hours  
**Cost:** Free

**Blog Post Outline:**

**Title Ideas:**
- "Building VoxForge: A Voice-to-Song Tool"
- "How I Built a Music Generator with Web Audio APIs"
- "From Voice to Music: My 12-Day Build Challenge"

**Sections:**
1. **Introduction**
   - What is VoxForge?
   - Why I built it
   - Demo GIF

2. **The Challenge**
   - Problem: Turning ideas into music is hard
   - Solution: Browser-based voice-to-song tool

3. **Technical Deep Dive**
   - Audio recording with Web Audio API
   - Pitch detection with pitchfinder
   - BPM detection
   - Music generation with Tone.js
   - Challenges faced

4. **Key Learnings**
   - What worked well
   - What was hard
   - What I'd do differently

5. **Try It Yourself**
   - Link to app
   - Link to GitHub
   - Call to action

**Where to Publish:**
- Dev.to (technical audience)
- Medium (broader audience)
- Your personal blog
- Hashnode (dev community)

**Blog Post URL:** _______________________________

**Completed:** ⬜ Yes | ⬜ No | ⬜ Skipped

---

## LEGAL TASKS (Optional but Recommended)

### Task 15: Add License

**Why:** Clarify usage rights  
**When:** Before making repo public  
**Estimated Time:** 5 minutes  
**Cost:** Free

**Recommended License: MIT**

Create `LICENSE` file:

```
MIT License

Copyright (c) 2025 [Your Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

**Other Options:**
- Apache 2.0 (more corporate-friendly)
- GPL (copyleft)
- No license (all rights reserved)

**Decision:** ⬜ MIT | ⬜ Apache | ⬜ GPL | ⬜ None

**Completed:** ⬜ Yes | ⬜ No

---

### Task 16: Add Privacy Policy (If Collecting Data)

**Why:** Legal requirement in many jurisdictions  
**When:** Before public launch  
**Estimated Time:** 30 minutes  
**Cost:** Free (use generator) or $50+ (lawyer)

**Note:** VoxForge doesn't collect data in MVP, so this may not be needed.

**If you add user accounts or analytics later:**
- Use privacy policy generator: https://www.privacypolicygenerator.info/
- Or hire lawyer: https://www.rocketlawyer.com/

**Completed:** ⬜ Yes | ⬜ No | ⬜ Not Needed

---

## ONGOING TASKS (After Launch)

### Task 17: Monitor and Respond

**Why:** Build community, fix bugs  
**When:** Ongoing  
**Estimated Time:** 30 min/day  
**Cost:** Free

**What to Monitor:**
- [ ] GitHub Issues
- [ ] Twitter mentions
- [ ] Reddit comments
- [ ] Email (if you added contact form)

**How to Respond:**
- Thank people for feedback
- Ask clarifying questions
- Fix reported bugs
- Consider feature requests

---

### Task 18: Iterate Based on Feedback

**Why:** Make it better  
**When:** Ongoing  
**Estimated Time:** Varies  
**Cost:** Time

**Common Requests to Expect:**
- More instrument options
- Better sound quality
- User accounts / save projects
- Mobile support
- Collaboration features
- Different music styles

**Prioritize:**
1. Bug fixes
2. Most requested features
3. Personal interest
4. Monetization potential

---

## CHECKLIST SUMMARY

### Pre-Launch
- [ ] Get API keys (optional)
- [ ] Create GitHub repo
- [ ] Choose deployment platform
- [ ] Record test vocals
- [ ] Test daily during development

### Launch Day
- [ ] Complete browser testing
- [ ] Get friend feedback
- [ ] Write README
- [ ] Take screenshots
- [ ] Record demo video
- [ ] Deploy to production
- [ ] (Optional) Set up custom domain

### Post-Launch
- [ ] Share on social media
- [ ] Write blog post
- [ ] Add license
- [ ] (Optional) Add privacy policy
- [ ] Monitor feedback
- [ ] Iterate and improve

---

## FINAL NOTES

**Time Investment Summary:**
- Setup: 1 hour
- Testing: 3-5 hours (spread across 12 days)
- Content Creation: 4-6 hours
- Deployment: 1 hour
- Marketing: 2-4 hours (optional)
- **Total: 11-17 hours**

**This is manageable!** Most tasks are quick, and you can do them while the AI codes.

**Remember:**
- These tasks make the difference between "I built something" and "I launched a product"
- Don't skip testing - it's the most important
- Content creation pays off long-term
- Marketing is optional but recommended

**Good luck! 🚀**

---

**END OF HUMAN TASKS**
