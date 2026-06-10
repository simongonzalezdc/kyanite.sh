# Export Formats Documentation

This document describes the export formats available in noise.sh and how to use them.

## Overview

noise.sh supports multiple export formats for your songs and lyrics, allowing you to share your work in different formats suitable for various use cases. The export system is designed to preserve your song structure, chords, and metadata while converting to the appropriate format.

## Available Export Formats

### 1. Markdown (.md)

Markdown export provides a clean, readable format that preserves the structure of your song with proper formatting.

**Features:**
- Song title as a level 1 header
- Section headers (Verse, Chorus, etc.) as level 2 headers
- Chord lines formatted with code formatting for visibility
- BPM metadata displayed
- Clean, readable formatting

**Use Cases:**
- Sharing songs on platforms that support Markdown (GitHub, GitLab, etc.)
- Documentation and README files
- Web content that can be easily converted to HTML

**Keyboard Shortcut:** `Ctrl+Shift+M`

**Example Output:**
```markdown
# My Song

**BPM:** 120

## Verse

`C        G        Am        F`
This is the first line of the verse
`G        C        G         C`
This is the second line of the verse

## Chorus

`F        C        G        Am`
This is the chorus line
`F        C        G        C`
This is the second chorus line
```

### 2. Plain Text (.txt)

Plain text export provides a clean, format-free version of your lyrics with minimal formatting.

**Features:**
- Removes all markdown formatting
- Skips chord lines for cleaner lyrics
- Preserves section structure
- Compatible with any text editor

**Use Cases:**
- Copy-pasting lyrics into other applications
- Email sharing
- Printing lyrics without chord clutter
- Compatibility with legacy systems

**Keyboard Shortcut:** `Ctrl+Shift+T` (Note: This conflicts with theme switching. Use the export menu `Ctrl+E` instead.)

**Example Output:**
```
My Song

This is the first line of the verse
This is the second line of the verse

This is the chorus line
This is the second chorus line
```

### 3. ChordPro (.cho)

ChordPro export creates a standardized format for lyrics with chords that's widely supported by music software.

**Features:**
- ChordPro metadata directives ({title}, {tempo}, {key})
- Section directives ({start_of_verse}, {end_of_verse})
- Preserves chord positioning
- Compatible with chord sheet software

**Use Cases:**
- Importing into chord sheet software (ChordPro, OnSong, etc.)
- Sharing with other musicians
- Creating printable chord sheets
- Band collaboration

**Keyboard Shortcut:** `Ctrl+Shift+P`

**Example Output:**
```chordpro
{title:My Song}
{tempo:120}
{key:C}

{start_of_verse}
C        G        Am        F
This is the first line of the verse
G        C        G         C
This is the second line of the verse
{end_of_verse}

{start_of_chorus}
F        C        G        Am
This is the chorus line
F        C        G        C
This is the second chorus line
{end_of_chorus}
```

## How to Export

### Using the Export Menu

1. From any screen, press `Ctrl+E` to open the export menu
2. Select your desired export format using arrow keys
3. Press `Enter` to export
4. The file will be saved with a timestamp in the filename

### Using Keyboard Shortcuts

You can bypass the export menu and directly export to specific formats:

- `Ctrl+Shift+M`: Export as Markdown
- `Ctrl+Shift+T`: Export as Plain Text
- `Ctrl+Shift+P`: Export as ChordPro

### Programmatic Export

If you're using the export service programmatically:

```go
import "github.com/kyanite/noise/internal/export"

// Create export service
service := export.NewExportService("/path/to/output/directory")

// Export as Markdown
markdownPath, err := service.ExportToMarkdown(content, "Song Title")

// Export as Plain Text
textPath, err := service.ExportToPlainText(content, "Song Title")

// Export as ChordPro
chordProPath, err := service.ExportToChordPro(content, "Song Title")
```

## Format Detection and Processing

The export system automatically detects and processes various song elements:

### Section Headers

The system recognizes section headers in multiple formats:
- `[Verse]`, `[Chorus]`, `[Bridge]`, `[Intro]`, `[Outro]`
- `Verse:`, `Chorus:`, `Bridge:`, etc.
- Plain text section names

### Chord Lines

Chord lines are identified by:
- Multiple chord symbols (C, G, Am, F, etc.)
- Short lines with primarily chord content
- Standard chord notation including sharps, flats, and quality indicators

### Metadata

The system extracts and preserves:
- **BPM/Tempo**: Detected from "BPM: 120" or "tempo: 120" patterns
- **Key**: Detected from the most common chord root in the song
- **Title**: Provided by the user or extracted from content

## File Naming

Exported files use the following naming convention:
```
{sanitized_title}_{timestamp}.{extension}
```

For example:
- `My_Song_20231018_143022.md`
- `My_Song_20231018_143022.txt`
- `My_Song_20231018_143022.cho`

## Integration with Existing Features

The new export formats integrate seamlessly with existing noise.sh features:

### AI Integration

Export formats work with AI-generated content, preserving:
- AI-suggested chord progressions
- AI-generated lyrics
- Quality analysis and suggestions

### Theme System

Exported content is independent of the current UI theme, ensuring consistent output regardless of the visual theme.

### Autosave

Export files are separate from autosave files and are only created when explicitly requested by the user.

## Troubleshooting

### Common Issues

1. **Chords not detected properly**
   - Ensure chords use standard notation (C, Am, G7, etc.)
   - Check that chord lines are separate from lyric lines

2. **Section headers not recognized**
   - Try using bracket notation: `[Verse]`, `[Chorus]`
   - Ensure section headers are on their own lines

3. **BPM not detected**
   - Use format "BPM: 120" or "tempo: 120"
   - Place BPM metadata at the beginning of the file

4. **Key detection is incorrect**
   - The system uses the most common chord root
   - You can manually specify the key in ChordPro format

### Export Quality Tips

- Use consistent formatting throughout your song
- Separate chord lines from lyric lines
- Place metadata at the beginning of the file
- Use standard chord notation

## Future Enhancements

Planned improvements to the export system:

1. **Additional Formats**
   - PDF export with professional formatting
   - HTML export with embedded styling
   - MusicXML export for notation software

2. **Advanced Formatting Options**
   - Customizable chord placement
   - Font and styling options
   - Header and footer customization

3. **Batch Export**
   - Export multiple songs at once
   - Export entire projects
   - Custom naming patterns

4. **Cloud Integration**
   - Direct export to cloud storage
   - Sharing links for exported files
   - Collaboration features

## Contributing

To contribute to the export system:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Update documentation
5. Submit a pull request

The export system is located in `internal/export/` and the UI components are in `internal/ui/export.go`.