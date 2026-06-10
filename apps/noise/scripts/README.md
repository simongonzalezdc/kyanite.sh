# noise.sh Scripts Directory

This directory contains helper scripts for building, launching, and testing the noise.sh application with the Kyanite theme system.

## Available Scripts

### Build and Launch Scripts

#### Windows
- **`build_and_launch.bat`** - Windows batch script for building and launching noise.sh with theme testing options

#### Linux/macOS
- **`build_and_launch.sh`** - Shell script for building and launching noise.sh with theme testing options

### Testing Scripts

#### Theme Testing
- **`../test_themes.go`** - Automated theme system testing script (run with `go run scripts/test_themes.go`)

#### Tools
- **`../tools/theme_test.go`** - Comprehensive theme testing tool with interactive interface

### Documentation

- **`theme_shortcuts_reference.md`** - Quick reference card for theme shortcuts and testing

## Usage

### Quick Start

#### Windows
```cmd
# Run the build and launch script
.\scripts\build_and_launch.bat
```

#### Linux/macOS
```bash
# Make the script executable and run it
chmod +x scripts/build_and_launch.sh
./scripts/build_and_launch.sh
```

### Using Make

```bash
# Build and run
make run

# Build and run with debug mode
make launch-debug

# Build and run in quick mode
make launch-quick

# Run automated theme tests
make test-themes

# Run comprehensive theme testing
make comprehensive-test
```

### Theme Testing

#### Automated Testing
```bash
# Run automated theme system tests
go run scripts/test_themes.go

# Or using make
make test-themes
```

#### Comprehensive Testing
```bash
# Run interactive comprehensive testing
go run tools/theme_test.go

# Or using make
make comprehensive-test
```

## Script Descriptions

### build_and_launch.bat / build_and_launch.sh
These scripts provide an easy way to build and launch noise.sh with various options:

1. **Normal launch** - Standard application launch
2. **Debug mode** - Launch with debug logging enabled
3. **Quick mode** - Launch in scratch mode for instant brainstorming
4. **Theme testing** - Launch with instructions for theme testing

### test_themes.go
Automated testing script that verifies:
- All 10 themes load correctly
- Theme switching performance meets requirements (< 100ms)
- Theme persistence works correctly
- Theme validation and migration functions properly

### theme_test.go (in tools/)
Interactive testing tool that provides:
- Automated testing integration
- Manual testing instructions
- Theme information display
- Interactive launch options
- Theme switching demonstration

## Theme System Overview

The Kyanite theme system includes 10 carefully crafted themes:

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

## Theme Switching Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+Shift+T` | Cycle through themes |
| `Ctrl+Shift+N` | Next theme |
| `Ctrl+Shift+P` | Previous theme |

## Testing Workflow

### 1. Automated Testing
```bash
# Run automated tests first
make test-themes
```

### 2. Manual Testing
```bash
# Launch for manual testing
make launch-debug
```

### 3. Comprehensive Testing
```bash
# Run full testing suite
make comprehensive-test
```

## Troubleshooting

### Build Issues
- Ensure Go 1.25+ is installed
- Check that all dependencies are available
- Run `go mod tidy` to update dependencies

### Theme Issues
- Verify terminal supports 256 colors or true color
- Try a different terminal emulator
- Check terminal color settings

### Script Issues
- Windows: Ensure batch files have execution permissions
- Linux/macOS: Run `chmod +x` on shell scripts
- Check that scripts are in the correct directory

## Contributing

When adding new scripts:
1. Follow the existing naming conventions
2. Include appropriate error handling
3. Add documentation to this README
4. Test on multiple platforms if possible
5. Update the Makefile with new targets if needed

## Support

For issues with the scripts or theme system:
1. Check the troubleshooting section above
2. Review the main [LAUNCH_GUIDE.md](../LAUNCH_GUIDE.md)
3. Report issues on the GitHub repository