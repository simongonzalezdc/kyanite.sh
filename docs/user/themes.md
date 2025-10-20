# Themes Guide

Focus.sh includes a beautiful theming system that lets you customize the visual appearance of the terminal interface. From cyberpunk synthwave to clean minimal themes, you can personalize your task management experience.

## Built-in Themes

### Synthwave (Default)
The signature Focus.sh theme with cyberpunk aesthetics.

- **Colors**: Neon pink, cyan, purple with dark backgrounds
- **Vibe**: Retro-futuristic, 80s synthwave
- **Best for**: Dark terminals, evening use
- **Usage**: `focus theme synthwave`

### Light
Clean, minimal theme for daytime use.

- **Colors**: Light grays, subtle blues, high contrast
- **Vibe**: Professional, clean, readable
- **Best for**: Daytime, bright environments
- **Usage**: `focus theme light`

### Plain
Monochrome theme for maximum compatibility.

- **Colors**: Black and white only
- **Vibe**: Minimal, distraction-free
- **Best for**: Basic terminals, scripting, accessibility
- **Usage**: `focus theme plain`

## Theme Commands

### Switching Themes

```bash
# Set theme by name
focus theme synthwave
focus theme light
focus theme plain

# List available themes
focus theme list

# Preview a theme before switching
focus theme preview synthwave

# Use config command
focus config set theme light
```

### Theme Information

```bash
# Show current theme
focus theme current

# Get theme details
focus config get theme
```

## Creating Custom Themes

You can create your own themes by adding JSON files to `~/.config/focus/themes/`.

### Theme Structure

```json
{
  "name": "my-custom-theme",
  "description": "A beautiful custom theme",
  "colors": {
    "primary": "#00ff00",
    "secondary": "#ff00ff",
    "background": "#000000",
    "text": "#ffffff",
    "accent": "#00ffff",
    "success": "#00ff00",
    "warning": "#ffff00",
    "error": "#ff0000",
    "muted": "#666666",
    "border": "#444444"
  },
  "styles": {
    "border": "single",
    "padding": 2,
    "margin": 1,
    "title_style": "bold",
    "border_color": "border"
  },
  "ui": {
    "dashboard_header": true,
    "show_dividers": true,
    "compact_mode": false,
    "animations": true
  }
}
```

### Color Definitions

- **primary**: Main accent color for buttons and highlights
- **secondary**: Secondary accent color
- **background**: Terminal background color
- **text**: Default text color
- **accent**: Emphasis and highlighting
- **success**: Success messages and completed tasks
- **warning**: Warnings and pending items
- **error**: Error messages and failed operations
- **muted**: Subtle text and secondary information
- **border**: Borders and dividers

### Style Options

#### Border Types
- `single` - Single line borders (`│ ─ ┌ ┐ └ ┘`)
- `double` - Double line borders (`║ ╔ ╗ ╚ ╝`)
- `rounded` - Rounded corners (`│ ─ ╭ ╮ ╯ ╰`)
- `hidden` - No borders
- `thick` - Thick block borders (`┃ ━ ┏ ┓ ┗ ┛`)

#### Text Styles
- `normal` - Regular text
- `bold` - Bold text
- `italic` - Italic text (terminal dependent)
- `underline` - Underlined text
- `blink` - Blinking text (use sparingly)
- `dim` - Dimmed text

### Example Custom Themes

#### Ocean Theme
```json
{
  "name": "ocean",
  "description": "Deep blue ocean theme",
  "colors": {
    "primary": "#00bfff",
    "secondary": "#4682b4",
    "background": "#001f3f",
    "text": "#e6f3ff",
    "accent": "#87ceeb",
    "success": "#00ff7f",
    "warning": "#ffd700",
    "error": "#ff6b6b",
    "muted": "#4a5568",
    "border": "#2c5282"
  },
  "styles": {
    "border": "rounded",
    "padding": 2,
    "margin": 1
  }
}
```

#### Forest Theme
```json
{
  "name": "forest",
  "description": "Natural forest colors",
  "colors": {
    "primary": "#228b22",
    "secondary": "#8fbc8f",
    "background": "#0d2818",
    "text": "#f0fff0",
    "accent": "#90ee90",
    "success": "#32cd32",
    "warning": "#daa520",
    "error": "#cd5c5c",
    "muted": "#556b2f",
    "border": "#2e7d32"
  },
  "styles": {
    "border": "single",
    "padding": 1,
    "margin": 1
  }
}
```

#### Sunset Theme
```json
{
  "name": "sunset",
  "description": "Warm sunset colors",
  "colors": {
    "primary": "#ff6b35",
    "secondary": "#f77f00",
    "background": "#2b1f17",
    "text": "#ffecb3",
    "accent": "#fcab64",
    "success": "#81c784",
    "warning": "#ffb74d",
    "error": "#e57373",
    "muted": "#8d6e63",
    "border": "#d84315"
  },
  "styles": {
    "border": "double",
    "padding": 2,
    "margin": 2
  }
}
```

## Installing Custom Themes

### Method 1: Create Theme File

```bash
# Create themes directory
mkdir -p ~/.config/focus/themes

# Create your theme file
cat > ~/.config/focus/themes/my-theme.json << 'EOF'
{
  "name": "my-theme",
  "description": "My custom theme",
  "colors": {
    "primary": "#00ff00",
    "secondary": "#ff00ff",
    "background": "#000000",
    "text": "#ffffff"
  }
}
EOF

# Apply your theme
focus theme my-theme
```

### Method 2: Save Current Settings

```bash
# Apply various settings
focus config set ui.dashboard_header false
focus config set ui.compact_mode true

# Save as theme
focus theme save my-minimal-theme

# Apply later
focus theme my-minimal-theme
```

## Theme Development

### Testing Your Theme

```bash
# Preview theme without applying
focus theme preview my-theme

# Apply temporarily
focus theme my-theme

# Test different views
focus list
focus unified
focus calendar today

# Get feedback
focus theme current
```

### Debugging Themes

Enable debug mode to see how colors are applied:

```bash
# Debug theme loading
focus --debug theme synthwave

# Check theme file
cat ~/.config/focus/themes/my-theme.json | jq .

# Validate JSON syntax
jsonlint ~/.config/focus/themes/my-theme.json
```

## Terminal Compatibility

### Color Support

Focus.sh automatically detects terminal color capabilities:

- **True Color (24-bit)**: Full 16.7 million colors
- **256 Colors**: Extended color palette
- **16 Colors**: Basic ANSI colors
- **Monochrome**: Fallback for limited terminals

### Compatibility Tips

```bash
# Check terminal color support
echo $TERM
echo $COLORTERM

# Force 256-color mode
export TERM=xterm-256color

# Disable colors for compatibility
export NO_COLOR=1
focus list

# Use plain theme for basic terminals
focus theme plain
```

### Terminal-Specific Notes

#### Windows Terminal
- Full color support recommended
- Use `focus theme synthwave` for best experience

#### macOS Terminal
- Good color support
- May need `export TERM=xterm-256color`

#### Linux Terminal
- Varies by terminal emulator
- Most modern terminals support full colors

#### Remote SSH
- Color support depends on connection
- Use `focus theme plain` if colors are broken

## Theme Inspiration

### Color Palettes

Get color ideas from these sources:
- **Coolors.co**: Online color palette generator
- **Adobe Color**: Professional color tools
- **Material Design**: Google's color guidelines
- **Tailwind CSS**: Modern web color palette

### Popular Theme Concepts

#### Dark Themes
- `midnight` - Deep blues and purples
- `matrix` - Green on black
- `vampire` - Red and black
- `cyberpunk` - Neon and dark

#### Light Themes
- `paper` - White and light grays
- `cream` - Warm off-whites
- `sky` - Light blues
- `mint` - Light greens

#### Seasonal Themes
- `autumn` - Oranges, reds, browns
- `winter` - Cool blues and whites
- `spring` - Fresh greens and pinks
- `summer` - Bright yellows and blues

## Sharing Themes

### Export Theme

```bash
# Export theme to share
cp ~/.config/focus/themes/my-theme.json ./my-theme.json

# Include documentation
cat > my-theme-README.md << 'EOF'
# My Awesome Theme

A beautiful custom theme for Focus.sh with cyberpunk aesthetics.

## Installation
```bash
cp my-theme.json ~/.config/focus/themes/
focus theme my-awesome-theme
```

## Colors
- Primary: Neon green (#00ff00)
- Background: Dark gray (#1a1a1a)
EOF
```

### Contributing Themes

Share your themes with the community:

1. Create theme file
2. Test thoroughly
3. Add documentation
4. Submit pull request to `themes/` directory
5. Include screenshot if possible

## Troubleshooting Themes

### Common Issues

#### "Theme not found"
```bash
# Check theme directory
ls -la ~/.config/focus/themes/

# Validate JSON syntax
cat ~/.config/focus/themes/my-theme.json | python -m json.tool

# Use fallback theme
focus theme synthwave
```

#### "Colors look wrong"
```bash
# Check terminal capabilities
echo $TERM $COLORTERM

# Force 256-color mode
export TERM=xterm-256color

# Try different theme
focus theme light
focus theme plain
```

#### "Borders display incorrectly"
```bash
# Check font and terminal
# Use simpler border style
echo '{"styles":{"border":"single"}}' > ~/.config/focus/themes/simple.json
focus theme simple
```

### Reset Theme

```bash
# Reset to default theme
focus theme synthwave

# Reset theme configuration
focus config reset theme

# Clear custom themes (backup first!)
mv ~/.config/focus/themes ~/.config/focus/themes.backup
mkdir ~/.config/focus/themes
```

## Advanced Customization

### Dynamic Themes

Create themes that change based on time of day:

```bash
# In your shell config (.bashrc, .zshrc)
theme_by_time() {
    hour=$(date +%H)
    if [ "$hour" -ge 6 ] && [ "$hour" -lt 18 ]; then
        focus theme light
    else
        focus theme synthwave
    fi
}

# Run on shell startup
theme_by_time
```

### Environment-Based Themes

```bash
# Set theme based on environment
case "$ENVIRONMENT" in
    "work")
        focus theme light
        ;;
    "home")
        focus theme synthwave
        ;;
    "production")
        focus theme plain
        ;;
esac
```

---

Experiment with themes to create your perfect task management environment! 🎨