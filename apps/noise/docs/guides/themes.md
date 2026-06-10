# Themes Guide

noise.sh includes 10 built-in themes and supports custom theme creation.

## Built-in Themes

| # | Theme | Description |
|---|-------|-------------|
| 1 | Monochrome | Classic black and white |
| 2 | Amber Night | Warm amber tones (default) |
| 3 | Twilight Mist | Soft purple gradients |
| 4 | Indigo Depths | Deep blue ocean |
| 5 | Forest Path | Natural green tones |
| 6 | Clay Earth | Warm earth colors |
| 7 | Iron Forge | Industrial reds |
| 8 | Sunlight | Bright golden yellows |
| 9 | Cyan Wave | Cool cyan blues |
| 10 | Electric Rose | Vibrant pink and cyan |

## Switching Themes

| Action | Shortcut |
|--------|----------|
| Cycle through themes | `Ctrl+Shift+T` |
| Next theme | `Ctrl+Shift+N` |
| Previous theme | `Ctrl+Shift+P` |

## Theme Persistence

Your selected theme is saved automatically and restored when you relaunch the app.

Theme preference is stored at:
- **Linux/macOS:** `~/.config/noise/theme.json`
- **Windows:** `%USERPROFILE%\.config\noise\theme.json`

## Custom Themes

You can create custom themes by adding TOML files to the themes directory.

### Theme Directory
```
~/.config/noise/themes/
```

### Theme File Format

Create a `.toml` file with your theme definition:

```toml
# mytheme.toml
name = "My Theme"
description = "A custom theme"

[colors]
primary = "#FF6B6B"
secondary = "#4ECDC4"
background = "#1A1A2E"
surface = "#16213E"
text = "#EAEAEA"
text_muted = "#888888"
accent = "#FFE66D"
success = "#95E1A3"
warning = "#FFD93D"
error = "#FF6B6B"
```

### Color Properties

| Property | Usage |
|----------|-------|
| primary | Main UI elements, active states |
| secondary | Secondary UI elements |
| background | Main background color |
| surface | Card/panel backgrounds |
| text | Primary text color |
| text_muted | Secondary/dim text |
| accent | Highlights and accents |
| success | Success messages |
| warning | Warning messages |
| error | Error messages |

## Accessibility

All built-in themes are designed to meet WCAG AA contrast requirements:
- Text contrast ratio ≥ 4.5:1
- UI element contrast ratio ≥ 3:1

When creating custom themes, use a contrast checker to ensure readability.

## Terminal Compatibility

Themes work best with terminals that support:
- 256 colors or true color (24-bit)
- Unicode characters

If colors look wrong, check your terminal's color settings.

### Recommended Terminals
- **macOS:** iTerm2, Terminal.app (with true color enabled)
- **Linux:** GNOME Terminal, Kitty, Alacritty
- **Windows:** Windows Terminal, ConEmu

## Troubleshooting

### Theme Not Applying
1. Delete the theme preference file and restart:
   ```bash
   rm ~/.config/noise/theme.json
   ```
2. Verify your terminal supports the required color depth

### Colors Look Wrong
- Check terminal color settings
- Try a different terminal emulator
- Ensure terminal isn't overriding colors

### Custom Theme Not Loading
- Verify TOML syntax is correct
- Check file is in the correct directory
- Ensure file extension is `.toml`
- Check app logs for loading errors

## Theme API (For Developers)

Themes are managed by `internal/theme/`:
- `theme.go` - Theme definitions
- `manager.go` - Theme switching and persistence
- `validator.go` - Contrast validation

See [reference/design-system.md](../reference/design-system.md) for implementation details.
