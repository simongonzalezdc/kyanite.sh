# noise.sh Launch Guide - Kyanite Theme System Testing

This guide provides comprehensive instructions for building, launching, and testing the noise.sh application with the new Kyanite theme system.

## Quick Start

### Prerequisites
- Go 1.25+ installed
- Terminal/Command Prompt
- (Optional) Ollama running locally for AI features
- Voice models download automatically on first use (no manual setup needed)

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
```cmd
# Run the build and launch script
.\scripts\build_and_launch.bat
```

#### Option 4: Quick Launch Scripts (Windows)
```cmd
# Simple batch script launcher
.\launch_noise.bat

# PowerShell launcher (more verbose)
powershell -ExecutionPolicy Bypass -File .\launch_noise.ps1

# PowerShell launcher with parameters
powershell -ExecutionPolicy Bypass -File .\launch_noise.ps1 --debug
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

## Voice-to-Text (Dictation)

noise.sh includes built-in voice dictation powered by whisper.cpp. Everything sets up automatically on first use.

### How It Works
1. **First Use**: When you press `Ctrl+D` for the first time, the app downloads the whisper model (~142MB)
2. **Recording**: Hold `Ctrl+D` to start recording, release or press again to stop
3. **Transcription**: Your speech is transcribed locally and inserted at the cursor

### Voice Shortcuts
| Action | Shortcut |
|--------|----------|
| Start/Stop Dictation | `Ctrl+D` |
| Cancel Recording | `Esc` |

### Voice Settings
Access voice settings through the settings menu to:
- Change the whisper model (tiny, base, small)
- Test microphone levels
- Download additional models

### Troubleshooting Voice
- **No microphone detected**: Check system permissions for microphone access
- **Model download fails**: Check internet connection; models can be pre-downloaded with `make download-model`
- **Poor transcription**: Try the "small" model for better accuracy (larger download)

## PWA Sync (Mobile Companion)

noise.sh includes an embedded sync server for capturing ideas from your phone.

### Enabling Sync
Add to your config (`~/.config/noise/noise.yaml`):
```yaml
sync:
  enabled: true
  port: 8765
  auto_start: true
```

### Pairing Your Phone
1. Enable sync in noise.sh settings
2. Generate a pairing code
3. Open the PWA on your phone and enter the code
4. Start capturing voice memos, photos, and tempo ideas

### Sync Features
- **Voice memos**: Quick audio recordings synced to your inbox
- **Photos**: Snap lyric ideas, setlists, or inspiration
- **Tap tempo**: Capture BPM ideas on the go
- **Text ideas**: Quick text notes

## Kyanite Theme System Testing

### Available Themes (10 Total)
**Note:** The README mentions "twelve curated themes" but the current implementation includes 10 themes. This documentation reflects the actual implementation.
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
Run `.\scripts\build_and_launch.bat` for automated building and launching with interactive options.

### Alternative Windows Launchers
- **`launch_noise.bat`** - Simple batch script that automatically builds if needed
- **`launch_noise.ps1`** - PowerShell script with colored output and better error handling

### Test Script
Run `.\scripts\test_themes.go` for automated theme testing.

## Troubleshooting

### Common Issues

#### Script Not Found Error
If you get "command not recognized" when running `.\scripts\build_and_launch.bat`:

1. **Check file exists**: Verify the script is in the scripts directory
2. **Use full path**: Try the full path: `c:\path\to\noise.sh\scripts\build_and_launch.bat`
3. **Use alternative launcher**: Try `.\launch_noise.bat` instead
4. **PowerShell execution policy**: For PowerShell scripts, run:
   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

#### Build Errors
```bash
# Ensure Go is properly installed
go version

# Clean and rebuild
make clean
make build

# Or use the batch script
.\scripts\build_and_launch.bat
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

## Windows Launch Methods Summary

### Recommended Methods (in order of preference):

1. **Interactive Launcher** (Recommended for testing):
   ```cmd
   .\scripts\build_and_launch.bat
   ```
   - Provides interactive menu with theme testing options
   - Automatically builds if needed
   - Best for first-time users and theme testing

2. **Simple Batch Launcher** (Quick launch):
   ```cmd
   .launch_noise.bat
   ```
   - Minimal, fast launcher
   - Automatically builds if needed
   - Good for daily use

3. **PowerShell Launcher** (Verbose output):
   ```powershell
   powershell -ExecutionPolicy Bypass -File .\launch_noise.ps1
   ```
   - Colored output and detailed error messages
   - Automatically builds if needed
   - Good for troubleshooting

4. **Direct Go Commands** (For developers):
   ```cmd
   go build -o noise.exe ./cmd/noise
   .\noise.exe
   ```
   - Full control over build process
   - Requires manual building
   - Best for development

5. **Make Commands** (If Make is installed):
   ```cmd
   make build
   make run
   ```
   - Standard build system
   - Requires Make on Windows
   - Good for cross-platform development

## Summary

The Kyanite theme system provides 10 carefully crafted themes for noise.sh. Use this guide to:
- Build and launch the application using multiple methods
- Test all theme functionality
- Verify theme persistence and performance
- Ensure proper integration with all features
- Troubleshoot common Windows launch issues

All themes should provide excellent readability and a consistent user experience across the entire application.