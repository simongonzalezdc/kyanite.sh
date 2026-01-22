package ai

import (
	"strings"
	"testing"
	"time"
)

// TestContextDetector_ComprehensiveEdgeCases tests edge cases for context detection
func TestContextDetector_ComprehensiveEdgeCases(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	tests := []struct {
		name     string
		content  string
		expected ContentType
	}{
		{
			name:     "Empty content",
			content:  "",
			expected: ContentTypeUnknown,
		},
		{
			name:     "Only whitespace",
			content:  "   \n  \n   \t  ",
			expected: ContentTypeUnknown,
		},
		{
			name:     "Single line with no clear indicators",
			content:  "Just some random text here",
			expected: ContentTypeUnknown,
		},
		{
			name:     "Single lyric line",
			content:  "I love you more than words can say",
			expected: ContentTypeLyrics,
		},
		{
			name:     "Single chord line",
			content:  "C - G - Am - F",
			expected: ContentTypePatterns,
		},
		{
			name:     "Minimal mixed content",
			content:  "C G\nI love you",
			expected: ContentTypeMixed,
		},
		{
			name:     "Only section headers",
			content:  "[Verse]\n[Chorus]\n[Bridge]",
			expected: ContentTypeUnknown, // Section headers alone aren't enough
		},
		{
			name:     "Only musical notation",
			content:  "120 BPM\n4/4\nC Major",
			expected: ContentTypePatterns,
		},
		{
			name:     "Drum notation only",
			content:  "x-x-x-x-x-x-x-x-\n----x-------x---",
			expected: ContentTypePatterns,
		},
		{
			name:     "Roman numerals only",
			content:  "I - V - vi - IV",
			expected: ContentTypePatterns,
		},
		{
			name:     "Very short lyric line",
			content:  "Love",
			expected: ContentTypeLyrics, // Short emotional word should be detected as lyric
		},
		{
			name:     "Very short pattern",
			content:  "C G",
			expected: ContentTypePatterns, // Short chord pattern should be detected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.AnalyzeContent(tt.content)
			helper.AssertEqual(tt.expected, result)
		})
	}
}

// TestContextDetector_AmbiguousContent tests ambiguous content that could be either type
func TestContextDetector_AmbiguousContent(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		content     string
		expected    ContentType
		description string
	}{
		{
			name:        "Short emotional text",
			content:     "Love tonight\nStars shine bright",
			expected:    ContentTypeLyrics,
			description: "Short but clear lyrical indicators",
		},
		{
			name:        "Minimal chord pattern",
			content:     "C G Am F\nF C G Am",
			expected:    ContentTypeMixed, // Short lines might be detected as mixed
			description: "Short chord pattern detected as mixed",
		},
		{
			name:        "Mixed short content",
			content:     "C G\nLove tonight",
			expected:    ContentTypeMixed,
			description: "Clear mix of both types",
		},
		{
			name:        "Generic text with one lyric word",
			content:     "The thing about love",
			expected:    ContentTypeLyrics, // Emotional words are detected as lyrics
			description: "Single emotional word detected as lyric",
		},
		{
			name:        "Generic text with one chord",
			content:     "The thing about C",
			expected:    ContentTypeUnknown,
			description: "Single chord not enough",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.AnalyzeContent(tt.content)
			helper.AssertEqual(tt.expected, result)
			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestContextDetector_LargeContent tests performance with large text inputs
func TestContextDetector_LargeContent(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	// Generate large lyric content
	largeLyricContent := strings.Builder{}
	for i := 0; i < 100; i++ {
		largeLyricContent.WriteString("[Verse ")
		largeLyricContent.WriteString(string(rune('A' + (i % 26))))
		largeLyricContent.WriteString("]\n")
		largeLyricContent.WriteString("This is line ")
		largeLyricContent.WriteString(string(rune(i + 1)))
		largeLyricContent.WriteString(" of the song\n")
		largeLyricContent.WriteString("With some lyrical content and emotion\n")
		largeLyricContent.WriteString("Love and heartbreak and feeling strong\n\n")
	}

	// Generate large pattern content
	largePatternContent := strings.Builder{}
	for i := 0; i < 100; i++ {
		largePatternContent.WriteString("Section ")
		largePatternContent.WriteString(string(rune('A' + (i % 26))))
		largePatternContent.WriteString(": C - G - Am - F\n")
		largePatternContent.WriteString("Tempo: 120 BPM\n")
		largePatternContent.WriteString("Key: C Major\n\n")
	}

	// Generate large mixed content
	largeMixedContent := strings.Builder{}
	for i := 0; i < 50; i++ {
		largeMixedContent.WriteString("[Verse ")
		largeMixedContent.WriteString(string(rune('A' + (i % 26))))
		largeMixedContent.WriteString("]\n")
		largeMixedContent.WriteString("C G Am F\n")
		largeMixedContent.WriteString("This is line ")
		largeMixedContent.WriteString(string(rune(i + 1)))
		largeMixedContent.WriteString(" of the song\n\n")
	}

	tests := []struct {
		name        string
		content     string
		expected    ContentType
		description string
	}{
		{
			name:        "Large lyric content",
			content:     largeLyricContent.String(),
			expected:    ContentTypeLyrics,
			description: "Should detect lyrics in large content",
		},
		{
			name:        "Large pattern content",
			content:     largePatternContent.String(),
			expected:    ContentTypeMixed, // Mixed content due to section headers
			description: "Should detect mixed content with section headers",
		},
		{
			name:        "Large mixed content",
			content:     largeMixedContent.String(),
			expected:    ContentTypeLyrics, // Content leans toward lyrics
			description: "Lyrics detected in large mixed content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			result := detector.AnalyzeContent(tt.content)
			elapsed := time.Since(start)

			helper.AssertEqual(tt.expected, result)
			helper.AssertTrue(elapsed < 100*time.Millisecond, "Analysis should complete quickly even for large content")
			t.Logf("Analysis completed in %v for %d characters", elapsed, len(tt.content))
			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestContextDetector_IntegrationWithKnowledgeBase tests integration with knowledge base
func TestContextDetector_IntegrationWithKnowledgeBase(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	// Test content that would benefit from knowledge base integration
	tests := []struct {
		name        string
		content     string
		expected    ContentType
		description string
	}{
		{
			name: "Lyrical content with musical terms",
			content: `[Verse]
C G Am F
I'm walking down the street tonight
The city lights are shining bright
With a minor key feeling`,
			expected:    ContentTypeLyrics, // Leans toward lyrics with emotional content
			description: "Should detect lyrics with musical terms",
		},
		{
			name: "Musical content with emotional descriptions",
			content: `Verse: C - G - Am - F
Mood: melancholic and reflective
Tempo: slow ballad
Emotional journey through the progression`,
			expected:    ContentTypeMixed,
			description: "Should detect mixed content with emotional descriptions",
		},
		{
			name: "Complex song structure",
			content: `[Intro]
Piano: Cmaj7 - Am7
Mood: hopeful

[Verse 1]
Walking through the morning light
C G Am F progression
Feeling the rhythm of life

[Chorus]
Emotional peak with soaring melody
F C G Am - powerful resolution`,
			expected:    ContentTypeLyrics,
			description: "Should detect lyrics in complex song structure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.AnalyzeContent(tt.content)
			helper.AssertEqual(tt.expected, result)
			t.Logf("Description: %s", tt.description)

			// Also test detailed analysis
			analysis := detector.GetContextAnalysis(tt.content)
			helper.AssertEqual(tt.expected, analysis.ContentType)
			helper.AssertTrue(analysis.Confidence > 0, "Confidence should be greater than 0")
			helper.AssertNotEmpty(analysis.Details)
		})
	}
}

// TestContextDetector_ErrorHandling tests error handling for malformed inputs
func TestContextDetector_ErrorHandling(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		content     string
		expected    ContentType
		description string
	}{
		{
			name:        "Nil content equivalent",
			content:     "",
			expected:    ContentTypeUnknown,
			description: "Should handle empty content gracefully",
		},
		{
			name:        "Only newlines",
			content:     "\n\n\n\n\n\n",
			expected:    ContentTypeUnknown,
			description: "Should handle content with only newlines",
		},
		{
			name:        "Unicode content",
			content:     "♫ ♪ ♫ ♪\n音乐之歌\n🎵 🎶 🎵",
			expected:    ContentTypeUnknown,
			description: "Should handle unicode content gracefully",
		},
		{
			name:        "Special characters only",
			content:     "!@#$%^&*()\n~`-_=+[]{}|;:'\",.<>/?",
			expected:    ContentTypeUnknown,
			description: "Should handle special characters gracefully",
		},
		{
			name:        "Malformed chord notation",
			content:     "C#b---G///Am??F!!!",
			expected:    ContentTypeUnknown,
			description: "Should handle malformed chord notation",
		},
		{
			name:        "Extremely long single line",
			content:     strings.Repeat("This is a very long line that should still be processed correctly ", 1000),
			expected:    ContentTypeUnknown,
			description: "Should handle extremely long lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.AnalyzeContent(tt.content)
			helper.AssertEqual(tt.expected, result)
			t.Logf("Description: %s", tt.description)

			// Ensure no panic occurs with detailed analysis
			analysis := detector.GetContextAnalysis(tt.content)
			helper.AssertEqual(tt.expected, analysis.ContentType)
		})
	}
}

// TestContextDetector_TableDrivenTests uses table-driven tests for multiple scenarios
func TestContextDetector_TableDrivenTests(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	tests := []struct {
		name            string
		content         string
		expectedType    ContentType
		expectedMinConf float64
		description     string
	}{
		// Strong lyric indicators
		{
			name: "Strong verse/chorus structure",
			content: `[Verse 1]
I'm walking down the street tonight
The city lights are shining bright
You left me here all alone
Now I'm trying to find my way home

[Chorus]
Love is like a burning fire
Desire takes me higher
Baby can't you see the light
You're my one and only delight`,
			expectedType:    ContentTypeLyrics,
			expectedMinConf: 0.8,
			description:     "Strong verse/chorus structure with emotional content",
		},
		{
			name: "Lyrical with contractions and emotional words",
			content: `I can't believe you're gone tonight
Won't you come back to me darling
My heart is breaking in two
Love is all I need from you honey`,
			expectedType:    ContentTypeLyrics,
			expectedMinConf: 0.7,
			description:     "Strong lyrical indicators with contractions and emotional words",
		},
		{
			name: "Rhyming pattern detection",
			content: `Rain falls down
On the ground
All around
Love is found`,
			expectedType:    ContentTypeLyrics,
			expectedMinConf: 0.5, // Lower confidence for short rhyming content
			description:     "Clear rhyming pattern in short lines",
		},

		// Strong pattern indicators
		{
			name: "Comprehensive musical notation",
			content: `Tempo: 120 BPM
Key: C Major
Time Signature: 4/4
Verse: C - G - Am - F
Chorus: F - C - G - Am
Bridge: Dm - Am - E - Am
Outro: C - G - Am - F`,
			expectedType:    ContentTypePatterns,
			expectedMinConf: 0.9,
			description:     "Comprehensive musical notation with all elements",
		},
		{
			name: "Roman numeral progressions",
			content: `Verse: I - vi - IV - V
Pre-Chorus: ii - V - I
Chorus: IV - I - V - vi
Bridge: vi - IV - I - V`,
			expectedType:    ContentTypePatterns,
			expectedMinConf: 0.8,
			description:     "Roman numeral progressions with clear structure",
		},
		{
			name: "Drum patterns and rhythm",
			content: `Kick: x---x---x---x---
Snare: ----x-------x---
Hi-hat: x-x-x-x-x-x-x-x-
Ride: ----x---x---x---`,
			expectedType:    ContentTypePatterns,
			expectedMinConf: 0.7,
			description:     "Drum patterns with clear rhythm notation",
		},

		// Mixed content
		{
			name: "Lyrics with chord symbols",
			content: `[Verse]
C G Am F
I'm walking down the street tonight
The city lights are shining bright
Am F C G
You left me here all alone

[Chorus]
F C G Am
Now I'm trying to find my way home
Love is like a burning fire`,
			expectedType:    ContentTypeMixed,
			expectedMinConf: 0.5,
			description:     "Clear mix of lyrics and chord symbols",
		},
		{
			name: "Musical description with lyrics",
			content: `Intro: Cmaj7 arpeggio
Mood: melancholic
"I'm standing in the rain"
Verse progression: Am - G - C - F
With emotional weight in each chord`,
			expectedType:    ContentTypeLyrics,
			expectedMinConf: 0.4,
			description:     "Lyrical content with musical descriptions",
		},

		// Edge cases
		{
			name: "Minimal but clear lyric",
			content: `Love is in the air
Tonight I care`,
			expectedType:    ContentTypeLyrics,
			expectedMinConf: 0.3,
			description:     "Minimal but clear lyrical pattern",
		},
		{
			name: "Minimal but clear pattern",
			content: `Verse: C G Am F
Chorus: F C G Am
Tempo: 120 BPM`,
			expectedType:    ContentTypePatterns,
			expectedMinConf: 0.3,
			description:     "Minimal but clear pattern structure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic analysis
			result := detector.AnalyzeContent(tt.content)
			helper.AssertEqual(tt.expectedType, result)

			// Test detailed analysis
			analysis := detector.GetContextAnalysis(tt.content)
			helper.AssertEqual(tt.expectedType, analysis.ContentType)
			helper.AssertTrue(analysis.Confidence >= tt.expectedMinConf,
				"Confidence %f should be >= %f", analysis.Confidence, tt.expectedMinConf)

			t.Logf("Description: %s", tt.description)
			t.Logf("Confidence: %f", analysis.Confidence)
			t.Logf("Details: %s", analysis.Details)
		})
	}
}

// TestContextDetector_PerformanceBenchmarks benchmarks the context detector performance
func TestContextDetector_PerformanceBenchmarks(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	// Create test content of different sizes
	sizes := []struct {
		name  string
		lines int
	}{
		{"Small", 10},
		{"Medium", 100},
		{"Large", 1000},
		{"Extra Large", 5000},
	}

	contentTypes := []struct {
		name      string
		generator func(int) string
		expected  ContentType
	}{
		{
			name: "Lyrics",
			generator: func(n int) string {
				var content strings.Builder
				for i := 0; i < n; i++ {
					content.WriteString("[Verse]\n")
					content.WriteString("This is line ")
					content.WriteString(string(rune(i + 1)))
					content.WriteString(" of the song with love and emotion\n")
				}
				return content.String()
			},
			expected: ContentTypeLyrics,
		},
		{
			name: "Patterns",
			generator: func(n int) string {
				var content strings.Builder
				for i := 0; i < n; i++ {
					content.WriteString("Section ")
					content.WriteString(string(rune('A' + (i % 26))))
					content.WriteString(": C - G - Am - F\n")
					content.WriteString("Tempo: 120 BPM\n")
					content.WriteString("Key: C Major\n")
				}
				return content.String()
			},
			expected: ContentTypePatterns,
		},
		{
			name: "Mixed",
			generator: func(n int) string {
				var content strings.Builder
				for i := 0; i < n; i++ {
					content.WriteString("[Verse]\n")
					content.WriteString("C G Am F\n")
					content.WriteString("This is line ")
					content.WriteString(string(rune(i + 1)))
					content.WriteString(" of the song\n\n")
				}
				return content.String()
			},
			expected: ContentTypeMixed,
		},
	}

	for _, size := range sizes {
		for _, contentType := range contentTypes {
			testName := size.name + " " + contentType.name
			t.Run(testName, func(t *testing.T) {
				content := contentType.generator(size.lines)

				// Benchmark the analysis
				start := time.Now()
				result := detector.AnalyzeContent(content)
				elapsed := time.Since(start)

				// Verify correct detection
				helper.AssertEqual(contentType.expected, result)

				// Performance assertions
				maxTime := time.Duration(size.lines) * time.Millisecond * 2 // 2ms per line
				helper.AssertTrue(elapsed < maxTime,
					"Analysis should complete in less than %v for %d lines, took %v",
					maxTime, size.lines, elapsed)

				t.Logf("Analyzed %d %s lines in %v (%.2fμs per line)",
					size.lines, contentType.name, elapsed,
					float64(elapsed.Nanoseconds())/1000/float64(size.lines))
			})
		}
	}
}

// TestContextDetector_ContextAnalysisMethods tests the ContextAnalysis methods
func TestContextDetector_ContextAnalysisMethods(t *testing.T) {
	detector := NewContextDetector()
	helper := NewTestHelper(t)

	tests := []struct {
		name                string
		content             string
		expectedContentType ContentType
		isLyricContent      bool
		isPatternContent    bool
	}{
		{
			name:                "Clear lyric content",
			content:             "[Verse]\nI love you more than words can say\nYou are my everything today",
			expectedContentType: ContentTypeLyrics,
			isLyricContent:      true,
			isPatternContent:    false,
		},
		{
			name:                "Clear pattern content",
			content:             "Verse: C - G - Am - F\nChorus: F - C - G - Am",
			expectedContentType: ContentTypePatterns,
			isLyricContent:      false,
			isPatternContent:    true,
		},
		{
			name:                "Mixed content with more lyrics",
			content:             "[Verse]\nC G Am F\nI love you\nThe city lights",
			expectedContentType: ContentTypeLyrics,
			isLyricContent:      true,
			isPatternContent:    false,
		},
		{
			name:                "Mixed content with more patterns",
			content:             "[Verse]\nC G Am F\nC - G - Am - F\nI - V - vi - IV",
			expectedContentType: ContentTypePatterns,
			isLyricContent:      false,
			isPatternContent:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := detector.GetContextAnalysis(tt.content)

			helper.AssertEqual(tt.expectedContentType, analysis.ContentType)
			helper.AssertEqual(tt.isLyricContent, analysis.IsLyricContent())
			helper.AssertEqual(tt.isPatternContent, analysis.IsPatternContent())

			// Test analysis details
			helper.AssertNotEmpty(analysis.Details)
			helper.AssertTrue(analysis.TotalLines > 0, "TotalLines should be greater than 0")
		})
	}
}
