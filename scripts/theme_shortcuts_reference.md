# noise.sh Theme System - Quick Reference Card

## Theme Switching Shortcuts

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+Shift+T` | Cycle Themes | Switch to next theme in sequence |
| `Ctrl+Shift+N` | Next Theme | Jump to next theme |
| `Ctrl+Shift+P` | Previous Theme | Jump to previous theme |

## Available Themes (Kyanite Suite)

| ID | Name | Description |
|----|------|-------------|
| `monochrome` | Monochrome | Classic black and white theme |
| `amber-night` | Amber Night | Warm amber tones (default) |
| `twilight-mist` | Twilight Mist | Soft purple gradients |
| `indigo-depths` | Indigo Depths | Deep blue ocean colors |
| `forest-path` | Forest Path | Natural green tones |
| `clay-earth` | Clay Earth | Warm earth colors |
| `iron-forge` | Iron Forge | Industrial reds and grays |
| `sunlight` | Sunlight | Bright golden yellows |
| `cyan-wave` | Cyan Wave | Cool cyan blues |
| `electric-rose` | Electric Rose | Vibrant pink and cyan |

## Essential Application Shortcuts

| Shortcut | Action | Context |
|----------|--------|---------|
| `F1` or `?` | Help | Global |
| `Ctrl+Q` | Quit | Global |
| `Esc` | Back/Menu | Global |
| `Ctrl+N` | New File | Global |
| `Ctrl+O` | Open File | Global |
| `Ctrl+S` | Save | Global |
| `Ctrl+E` | Export | Global |

## Editor Shortcuts

| Shortcut | Action | Context |
|----------|--------|---------|
| `Alt+G` | AI Unstick | Editor |
| `Alt+R` | AI Spark | Editor |
| `Alt+V` | AI Tweak | Editor |
| `Alt+C` | AI Check | Editor |
| `Ctrl+F` | Chord Picker | Editor |
| `Ctrl+Shift+B` | BPM Tapper | Editor |
| `Tab` | Next Pane | Global |
| `Shift+Tab` | Previous Pane | Global |

## Testing Checklist

### Visual Verification
- [ ] All 10 themes load without errors
- [ ] Text is readable in all themes
- [ ] UI elements are properly styled
- [ ] Color contrast is adequate
- [ ] No visual artifacts or glitches

### Functional Testing
- [ ] Theme switching works smoothly
- [ ] Theme preference is saved/restored
- [ ] All features work in all themes
- [ ] Performance is acceptable (< 100ms switching)

### Integration Testing
- [ ] Editor works correctly in all themes
- [ ] AI assistance displays properly
- [ ] Chord picker is usable in all themes
- [ ] Export functionality works
- [ ] Help system displays correctly

## Troubleshooting

### Theme Not Applying
1. Check terminal color support
2. Try a different terminal emulator
3. Ensure 256-color or true color support
4. Reset theme preference:
   ```bash
   # Windows
   del %USERPROFILE%\.config\noise\theme.json
   # Linux/Mac
   rm ~/.config/noise/theme.json
   ```

### Performance Issues
1. Check system resources
2. Close other applications
3. Restart the application
4. Run with debug mode: `./noise.exe --debug`

### Visual Issues
1. Verify terminal supports ANSI colors
2. Check terminal font settings
3. Try different terminal
4. Report issue with theme name and terminal type

## Testing Commands

```bash
# Build and test themes
make build
go run scripts/test_themes.go

# Launch with debug mode
./noise.exe --debug

# Quick theme testing
./noise.exe quick "test theme"
```

Remember: Use `Ctrl+Shift+T` to quickly cycle through all themes during testing!