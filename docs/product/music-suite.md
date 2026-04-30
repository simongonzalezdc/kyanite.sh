# noise.sh Music Suite

`noise.sh` is the canonical app and suite hub for private, local-first songwriting and music creation.

Imported modules preserved with git history:

- `modules/voxforge/` from `Pastorsimon1798/VoxForge`
- `modules/lyrics-engine/` from `Pastorsimon1798/lyrics-engine`

## Suite map

1. **noise.sh** — local-first rapid-capture songwriting studio, TUI-first, with local AI and voice dictation.
2. **VoxForge** — browser voice-to-song arranger: pitch/BPM/key detection, drums, bass, chords, stems, MIDI, AI lyrics.
3. **lyrics-engine** — Python lyric/prosody library and CLI: rhymes, syllables, stress/flow analysis, structure templates, AI generation.

## Boundary rule

Modules remain isolated under `modules/` until integrations are specified and tested. Do not mix `VoxForge` or `lyrics-engine` internals into the core `noise.sh` app without a written design for overlap with existing prosody/AI/song models.

## Likely integration path

1. Compare `lyrics-engine` against existing `noise.sh` prosody/knowledge-base code.
2. Expose `lyrics-engine` as a local service or embedded CLI used by noise.sh.
3. Extract VoxForge voice-to-arrangement as a companion PWA/music engine.
4. Keep privacy-first/local-first constraints as the suite-wide product rule.
