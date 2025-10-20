# noise.sh Launch Guide - Kyanite Theme System Testing

This guide provides comprehensive instructions for building, launching, and testing the noise.sh application with the new Kyanite theme system.

## Quick Start

### Prerequisites
- Go 1.21+ installed
- Terminal/Command Prompt
- (Optional) Ollama running locally for AI features

### Building the Application

#### Option 1: Using Make (Recommended)
```bash
# Build the application
make build

# Run the application
make run
```

#### Option 2: Direct Go Commands
```bash
# Build the application
go build -o noise.exe ./cmd/noise

# Run the application
./noise.exe
```

#### Option 3: Using the Helper Script (Windows)
```bash
# Run the build and launch script
.\scripts\build_and_launch.bat
```

## Launch Options

### Basic Launch
```bash
# Start with default settings
./noise.exe

# Start with debug logging
./noise.exe --debug

# Start in scratch mode with instant brainstorming
./noise.exe quick

# Start with a specific theme for brainstorming
./noise.exe quick "lost love"
```

### Opening Specific Files
```bash
# Open a specific song file
./noise.exe song.md

# Open with debug mode
./noise.exe --debug song.md
```

## Kyanite Theme System Testing

### Available Themes (10 Total)
1. **Monochrome** - Classic black and white
2. **Amber Night** - Warm amber tones (default)
3. **Twilight Mist** - Soft purple gradients
4. **Indigo Depths** - Deep blue ocean
5. **Forest Path** - Natural green tones
6. **Clay Earth** - Warm earth colors
7. **Iron Forge** - Industrial reds
8. **Sunlight** - Bright golden yellows
9. **Cyan Wave** - Cool cyan blues
10. **Electric Rose** - Vibrant pink and cyan

### Theme Switching Shortcuts

| Action | Shortcut | Description |
|--------|----------|-------------|
| Cycle Themes | `Ctrl+Shift+T` | Switch to next theme in sequence |
| Next Theme | `Ctrl+Shift+N` | Jump to next theme |
| Previous Theme | `Ctrl+Shift+P` | Jump to previous theme |

### Testing Theme System

#### 1. Basic Theme Cycling Test
1. Launch the application: `./noise.exe`
2. Press `Ctrl+Shift+T` to cycle through all 10 themes
3. Verify each theme displays correctly with proper colors
4. Check that UI elements remain readable in each theme

#### 2. Theme Persistence Test
1. Switch to a specific theme (e.g., `Ctrl+Shift+T` until "Electric Rose")
2. Exit the application with `Ctrl+Q`
3. Relaunch the application: `./noise.exe`
4. Verify the selected theme is preserved

#### 3. Theme-Specific UI Elements Test
For each theme, verify:
- **Primary elements** (buttons, active elements) are visible
- **Secondary elements** (inactive elements) have good contrast
- **Text** is readable against the background
- **Success/Warning/Error messages** are appropriately colored
- **Accent colors** highlight important elements

#### 4. Theme Integration with Features Test
1. Open a song file or start in quick mode
2. Test theme switching while using different features:
   - Editor pane with text
   - Chord picker (`Ctrl+F`)
   - BPM tapper (`Ctrl+Shift+B`)
   - AI assistance (`Alt+G`, `Alt+R`, `Alt+V`, `Alt+C`)
   - Export functionality (`Ctrl+E`)

## Comprehensive Testing Workflow

### 1. Launch and Basic Functionality
```bash
# Build and run with debug
make build
./noise.exe --debug
```

**Verification Checklist:**
- [ ] Application launches successfully
- [ ] Main menu appears
- [ ] Default theme (Amber Night) is applied
- [ ] UI elements are properly styled

### 2. Theme System Verification
1. **Cycle through all themes** using `Ctrl+Shift+T`
2. **Test theme persistence** by restarting the app
3. **Verify theme-specific styling** in different contexts

**Theme Testing Checklist:**
- [ ] All 10 themes load without errors
- [ ] Theme switching is smooth (< 100ms)
- [ ] Color contrast is adequate for readability
- [ ] Theme preference is saved and restored
- [ ] UI elements adapt correctly to theme changes

### 3. Feature Integration Testing
1. **Editor Testing:**
   - Create a new song with `Ctrl+N`
   - Type some lyrics
   - Switch themes while editing
   - Verify text remains readable

2. **AI Integration Testing:**
   - Use `Alt+G` for AI unstick
   - Use `Alt+R` for AI spark
   - Use `Alt+V` for AI tweak
   - Use `Alt+C` for AI check
   - Verify AI responses display correctly in all themes

3. **Export Testing:**
   - Create content and export with `Ctrl+E`
   - Test different export formats:
     - `Ctrl+Shift+M` for Markdown
     - `Ctrl+Shift+P` for ChordPro

### 4. Performance Testing
1. **Theme Switching Performance:**
   - Time theme switching between different themes
   - Should be under 100ms per switch

2. **Memory Usage:**
   - Monitor memory usage during theme switching
   - No significant memory leaks should occur

## Helper Scripts

### Windows Batch Script
Run `.\scripts\build_and_launch.bat` for automated building and launching.

### Test Script
Run `.\scripts\test_themes.go` for automated theme testing.

## Troubleshooting

### Common Issues

#### Build Errors
```bash
# Ensure Go is properly installed
go version

# Clean and rebuild
make clean
make build
```

#### Theme Not Applying
1. Check if theme preference file exists:
   - Windows: `%USERPROFILE%\.config\noise\theme.json`
   - Linux/Mac: `~/.config/noise/theme.json`

2. Reset to default theme:
   ```bash
   # Delete theme preference file
   rm ~/.config/noise/theme.json  # Linux/Mac
   del %USERPROFILE%\.config\noise\theme.json  # Windows
   ```

#### UI Elements Not Visible
1. Check terminal color support
2. Try a different terminal emulator
3. Ensure terminal supports 256 colors or true color

### Debug Mode
Enable debug mode for detailed logging:
```bash
./noise.exe --debug
```

### Getting Help
Press `F1` or `?` in the application to see all available shortcuts.

## Advanced Testing

### Automated Theme Testing
For comprehensive automated testing, run:
```bash
go run test_enhancement_065_comprehensive.go
```

### Performance Benchmarking
Run the performance tests to verify theme switching meets requirements:
```bash
go test -run TestPerformanceRequirements ./...
```

## Feedback and Issues

If you encounter issues with the theme system:
1. Note the specific theme where the issue occurs
2. Describe the UI elements affected
3. Include your terminal type and OS
4. Report issues on the GitHub repository

## Summary

The Kyanite theme system provides 10 carefully crafted themes for noise.sh. Use this guide to:
- Build and launch the application
- Test all theme functionality
- Verify theme persistence and performance
- Ensure proper integration with all features

All themes should provide excellent readability and a consistent user experience across the entire application.