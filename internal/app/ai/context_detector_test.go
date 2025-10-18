package ai

import (
	"testing"
)

func TestContextDetector_AnalyzeContent(t *testing.T) {
	detector := NewContextDetector()

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
			name: "Pure lyrics with verse structure",
			content: `[Verse 1]
I'm walking down the street tonight
The city lights are shining bright
You left me here all alone
Now I'm trying to find my way home`,
			expected: ContentTypeLyrics,
		},
		{
			name: "Pure chord progressions",
			content: `Verse: C - G - Am - F
Chorus: F - C - G - Am
Bridge: Dm - Am - E - Am`,
			expected: ContentTypePatterns,
		},
		{
			name: "Musical notation with tempo",
			content: `Tempo: 120 BPM
Key: C Major
Time Signature: 4/4
Verse: I - V - vi - IV`,
			expected: ContentTypePatterns,
		},
		{
			name: "Mixed content with lyrics and chords",
			content: `[Verse]
C G Am F
I'm walking down the street tonight
The city lights are shining bright`,
			expected: ContentTypeMixed,
		},
		{
			name: "Drum patterns",
			content: `Kick: x---x---x---x---
Snare: ----x-------x---
Hi-hat: x-x-x-x-x-x-x-x-`,
			expected: ContentTypePatterns,
		},
		{
			name: "Lyrical content with emotional words",
			content: `Love is in the air tonight
My heart is beating for you
Baby can't you see the light
Honey I'm so lost without you`,
			expected: ContentTypeLyrics,
		},
		{
			name: "Roman numeral progressions",
			content: `Verse: I - vi - IV - V
Pre-chorus: ii - V - I
Chorus: IV - I - V - vi`,
			expected: ContentTypePatterns,
		},
		{
			name: "Short lyrical lines",
			content: `Rain falls down
Tears hit the ground
Love fades away
Night turns to day`,
			expected: ContentTypeLyrics,
		},
		{
			name: "Pattern with loop structure",
			content: `Loop 1: C - Am - F - G
Loop 2: Am - F - C - G
Pattern: driving rhythm`,
			expected: ContentTypePatterns,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.AnalyzeContent(tt.content)
			if result != tt.expected {
				t.Errorf("AnalyzeContent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestContextDetector_GetContextAnalysis(t *testing.T) {
	detector := NewContextDetector()

	tests := []struct {
		name            string
		content         string
		expectedType    ContentType
		expectedMinConf float64
	}{
		{
			name: "Strong lyrical content",
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
			expectedMinConf: 0.5,
		},
		{
			name: "Strong pattern content",
			content: `Tempo: 120 BPM
Key: C Major
Verse: C - G - Am - F
Chorus: F - C - G - Am
Bridge: Dm - Am - E - Am
Time Signature: 4/4`,
			expectedType:    ContentTypePatterns,
			expectedMinConf: 0.5,
		},
		{
			name: "Mixed content",
			content: `[Verse]
C G Am F
I'm walking down the street tonight
The city lights are shining bright
Am F C G
You left me here all alone`,
			expectedType:    ContentTypeMixed,
			expectedMinConf: 0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := detector.GetContextAnalysis(tt.content)
			
			if analysis.ContentType != tt.expectedType {
				t.Errorf("GetContextAnalysis().ContentType = %v, expected %v", analysis.ContentType, tt.expectedType)
			}
			
			if analysis.Confidence < tt.expectedMinConf {
				t.Errorf("GetContextAnalysis().Confidence = %v, expected >= %v", analysis.Confidence, tt.expectedMinConf)
			}
			
			if analysis.TotalLines == 0 {
				t.Error("GetContextAnalysis().TotalLines should be > 0 for non-empty content")
			}
			
			if analysis.Details == "" {
				t.Error("GetContextAnalysis().Details should not be empty")
			}
		})
	}
}

func TestContextAnalysis_HelperMethods(t *testing.T) {
	tests := []struct {
		name      string
		analysis  *ContextAnalysis
		isLyric   bool
		isPattern bool
	}{
		{
			name: "Lyric content analysis",
			analysis: &ContextAnalysis{
				ContentType:    ContentTypeLyrics,
				LyricMatches:   5,
				PatternMatches: 1,
			},
			isLyric:   true,
			isPattern: false,
		},
		{
			name: "Pattern content analysis",
			analysis: &ContextAnalysis{
				ContentType:    ContentTypePatterns,
				LyricMatches:   1,
				PatternMatches: 5,
			},
			isLyric:   false,
			isPattern: true,
		},
		{
			name: "Mixed content with more lyrics",
			analysis: &ContextAnalysis{
				ContentType:    ContentTypeMixed,
				LyricMatches:   4,
				PatternMatches: 2,
			},
			isLyric:   true,
			isPattern: false,
		},
		{
			name: "Mixed content with more patterns",
			analysis: &ContextAnalysis{
				ContentType:    ContentTypeMixed,
				LyricMatches:   2,
				PatternMatches: 4,
			},
			isLyric:   false,
			isPattern: true,
		},
		{
			name: "Mixed content with equal matches",
			analysis: &ContextAnalysis{
				ContentType:    ContentTypeMixed,
				LyricMatches:   3,
				PatternMatches: 3,
			},
			isLyric:   true,
			isPattern: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.analysis.IsLyricContent() != tt.isLyric {
				t.Errorf("IsLyricContent() = %v, expected %v", tt.analysis.IsLyricContent(), tt.isLyric)
			}
			
			if tt.analysis.IsPatternContent() != tt.isPattern {
				t.Errorf("IsPatternContent() = %v, expected %v", tt.analysis.IsPatternContent(), tt.isPattern)
			}
		})
	}
}

func TestContextDetector_EdgeCases(t *testing.T) {
	detector := NewContextDetector()

	tests := []struct {
		name     string
		content  string
		expected ContentType
	}{
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
			name: "Single lyric line",
			content: `I love you more than words can say`,
			expected: ContentTypeLyrics,
		},
		{
			name: "Single chord line",
			content: `C - G - Am - F`,
			expected: ContentTypePatterns,
		},
		{
			name: "Minimal mixed content",
			content: `C G
I love you`,
			expected: ContentTypeMixed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.AnalyzeContent(tt.content)
			if result != tt.expected {
				t.Errorf("AnalyzeContent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestContextDetector_Performance(t *testing.T) {
	detector := NewContextDetector()
	
	// Large content test
	largeContent := ""
	for i := 0; i < 100; i++ {
		largeContent += "[Verse]\n"
		largeContent += "This is line " + string(rune(i)) + " of the song\n"
		largeContent += "With some lyrical content and emotion\n"
		largeContent += "Love and heartbreak and feeling strong\n\n"
	}
	
	// This should complete quickly even with large content
	t.Run("Large content analysis", func(t *testing.T) {
		result := detector.AnalyzeContent(largeContent)
		if result != ContentTypeLyrics {
			t.Errorf("Expected ContentTypeLyrics for large lyrical content, got %v", result)
		}
	})
}

func TestContextDetector_Consistency(t *testing.T) {
	detector := NewContextDetector()
	
	content := `[Verse]
C G Am F
I'm walking down the street tonight
The city lights are shining bright`
	
	// Multiple calls should return the same result
	t.Run("Consistent results", func(t *testing.T) {
		result1 := detector.AnalyzeContent(content)
		result2 := detector.AnalyzeContent(content)
		result3 := detector.AnalyzeContent(content)
		
		if result1 != result2 || result2 != result3 {
			t.Error("AnalyzeContent() should return consistent results for the same input")
		}
	})
}