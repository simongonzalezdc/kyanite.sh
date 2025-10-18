# Context-Aware AI System Guide

## Overview

The Context-Aware AI System is a Week 2 implementation for Enhancement #6.5 that automatically detects whether users are working with lyrics or musical patterns and provides appropriate AI suggestions based on the detected content type.

## Features

### 1. Context Detection System

The system automatically analyzes content to determine whether it's:
- **Lyrics**: Song lyrics with emotional content, imagery, and narrative
- **Patterns**: Musical patterns, chord progressions, and rhythmic structures
- **Mixed**: Content combining both lyrics and musical elements
- **Unknown**: Content that doesn't clearly fit any category

#### Detection Patterns

**Lyric Detection:**
- Section headers (Verse, Chorus, Bridge, etc.)
- Emotional and sensory language
- Common lyrical structures and themes
- Contractions and personal pronouns
- Rhyming patterns and meter

**Pattern Detection:**
- Chord progressions (C - G - Am - F, I - V - vi - IV, etc.)
- Musical notation (tempo, key, time signature)
- Drum patterns and rhythmic notation
- Musical structure keywords (loop, pattern, beat, rhythm)

### 2. Context-Aware AI Prompts

The AI system provides different suggestions based on detected content:

#### For Lyrics:
- **Continue Writing**: Suggestions that maintain tone, imagery, and emotional content
- **Generate Ideas**: Opening lines with vivid imagery and emotional resonance
- **Refine Content**: Variations that enhance imagery, emotion, and authenticity
- **Quality Check**: Feedback on imagery, emotion, rhythm, and originality

#### For Patterns:
- **Continue Writing**: Suggestions that build harmonic and rhythmic foundations
- **Generate Ideas**: Musical patterns appropriate for the mood
- **Refine Content**: Variations that improve harmonic movement and rhythmic interest
- **Quality Check**: Feedback on harmonic logic, rhythm, and playability

#### For Mixed Content:
- **Continue Writing**: Suggestions that enhance both lyrical and musical elements
- **Generate Ideas**: Creative starting points with both lyrical and musical potential
- **Refine Content**: Variations that improve flow between words and music
- **Quality Check**: Feedback on balance between lyrics and music

### 3. Visual Indicators

The status bar displays the detected content type with color coding:
- 🟢 **Lyrics**: Green indicator for lyrical content
- 🔵 **Patterns**: Blue indicator for musical patterns
- 🟡 **Mixed**: Yellow indicator for mixed content
- ⚪ **Unknown**: Gray indicator for unclear content

### 4. Integration Points

#### Editor Integration
- Real-time content analysis as you type
- Automatic status bar updates
- Context-aware AI shortcuts

#### AI Service Integration
- All AI methods (GenerateContinuations, GenerateVariations, etc.) use context detection
- Fallback suggestions when AI is unavailable
- Consistent context-aware behavior across all AI features

## Usage

### Automatic Detection
The system works automatically in the background:
1. Start typing lyrics or musical patterns
2. The system analyzes your content in real-time
3. The status bar shows the detected content type
4. AI suggestions are tailored to the detected type

### AI Shortcuts
Use the following shortcuts with context-aware AI:
- `Ctrl+G`: Continue writing (context-aware suggestions)
- `Ctrl+1`: Generate ideas (context-aware starting points)
- `Ctrl+2`: Refine content (context-aware variations)
- `Ctrl+3`: Quality check (context-aware feedback)

### Manual Override
If automatic detection doesn't match your intent:
1. The system adapts to content changes
2. Mixed content provides balanced suggestions
3. Unknown content falls back to general assistance

## Technical Implementation

### Core Components

#### ContextDetector
- Analyzes content using regex patterns
- Calculates confidence scores
- Provides detailed analysis information

#### ContextAwarePrompts
- Manages different prompt templates
- Renders prompts based on content type
- Provides fallback prompts for unknown content

#### QuickIdeaAgent Integration
- Automatically detects content type before generating suggestions
- Uses appropriate prompts based on content type
- Provides context-aware fallback suggestions

### Performance Considerations
- Content analysis is optimized for real-time use
- Analysis window focuses on recent content (last 10 lines)
- Throttled updates prevent performance issues

### Error Handling
- Graceful fallback when content type is unclear
- Robust error handling for AI service failures
- Consistent behavior even with limited AI availability

## Testing

### Test Coverage
- Context detection accuracy tests
- Context-aware prompt rendering tests
- Integration tests for AI services
- UI component tests for visual indicators
- Performance tests for large content

### Running Tests
```bash
go test ./internal/app/ai/...
```

## Configuration

### Customizing Detection Patterns
Modify the patterns in `ContextDetector` to adjust detection sensitivity:
- Add new regex patterns for specific use cases
- Adjust confidence thresholds
- Modify analysis window size

### Customizing Prompts
Edit prompt templates in `ContextAwarePrompts` to:
- Adjust AI behavior for specific content types
- Add new prompt variations
- Modify fallback suggestions

## Troubleshooting

### Common Issues

**Incorrect Content Type Detection**
- Ensure content follows expected patterns
- Check for mixed content indicators
- Verify status bar shows current detection

**AI Suggestions Not Context-Appropriate**
- Check content type indicator in status bar
- Try adding more content for better detection
- Use manual prompts if needed

**Performance Issues**
- Reduce analysis window size in ContextDetector
- Check for large content files
- Monitor update throttling settings

### Debug Information
Enable detailed analysis by checking:
- Context analysis results
- Confidence scores
- Matched patterns
- Prompt rendering output

## Future Enhancements

### Planned Improvements
- Machine learning-based content detection
- User preference learning
- More granular content type categories
- Advanced pattern recognition

### Extension Points
- Custom content type definitions
- Plugin system for new detection patterns
- User-trained detection models
- Integration with external AI services

## Examples

### Lyric Content Example
```
[Verse 1]
I'm walking down the street tonight
The city lights are shining bright
You left me here all alone
Now I'm trying to find my way home
```
- Detected as: **Lyrics**
- AI suggestions focus on: imagery, emotion, narrative

### Pattern Content Example
```
Tempo: 120 BPM
Key: C Major
Verse: C - G - Am - F
Chorus: F - C - G - Am
```
- Detected as: **Patterns**
- AI suggestions focus on: harmony, rhythm, structure

### Mixed Content Example
```
[Verse]
C G Am F
I'm walking down the street tonight
Am F C G
The city lights are shining bright
```
- Detected as: **Mixed**
- AI suggestions focus on: both lyrics and music

## Conclusion

The Context-Aware AI System provides intelligent, content-appropriate assistance that adapts to your creative workflow. By automatically detecting content type and tailoring AI suggestions accordingly, it helps songwriters and musicians stay in their creative flow while receiving relevant, helpful suggestions.