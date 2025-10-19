package ai

import (
	"fmt"
	"regexp"
	"strings"
)

// ContentType represents the type of content being analyzed
type ContentType string

const (
	ContentTypeLyrics   ContentType = "lyrics"
	ContentTypePatterns ContentType = "patterns"
	ContentTypeMixed    ContentType = "mixed"
	ContentTypeUnknown  ContentType = "unknown"
)

// ContextDetector analyzes content to determine if it's lyrics, patterns, or mixed
type ContextDetector struct {
	// Patterns for identifying different content types
	lyricPatterns           []*regexp.Regexp
	patternPatterns         []*regexp.Regexp
	lyricStructuralPatterns map[*regexp.Regexp]bool
	lyricStrongPatterns     map[*regexp.Regexp]bool
	lyricLexicalPatterns    map[*regexp.Regexp]bool
	patternStrongPatterns   map[*regexp.Regexp]bool
	mixedThreshold          float64
	analysisWindow          int // Number of lines to analyze for context
}

// NewContextDetector creates a new context detector with default patterns
func NewContextDetector() *ContextDetector {
	detector := &ContextDetector{
		mixedThreshold: 0.45, // Require stronger presence from both types to be considered mixed
		analysisWindow: 10,  // Analyze last 10 lines for context
	}

	// Initialize lyric detection patterns
	detector.lyricPatterns = []*regexp.Regexp{
		// Common lyrical structures
		regexp.MustCompile(`(?i)^\s*(verse|chorus|bridge|intro|outro|pre-chorus)\s*[:\-\d]*`),
		regexp.MustCompile(`(?i)^\s*\[.*\]\s*$`), // Section headers in brackets
		regexp.MustCompile(`(?m)^\s*[A-Z][^.!?]*[.!?]`), // Sentences starting with capital
		regexp.MustCompile(`(?i)\b(I'm|you're|we're|they're|can't|won't|don't)\b`), // Contractions
		regexp.MustCompile(`(?i)\b(love|heart|night|day|sky|rain|sun|moon|stars)\b`), // Common lyrical words
		regexp.MustCompile(`(?m)^\s*[A-Z][a-z]+\s+[a-z]+\s+[a-z]+\s+`), // Rhyming pattern detection
		regexp.MustCompile(`(?i)\b(baby|honey|darling|sweetheart)\b`), // Terms of endearment
		regexp.MustCompile(`(?m)^\s*\w+\s+\w+\s+\w+\s+\w+\s*$`), // Short lines typical of lyrics
	}

	// Initialize pattern detection patterns
	detector.patternPatterns = []*regexp.Regexp{
		// Chord progressions
		regexp.MustCompile(`(?i)^\s*([A-G][#b]?(m|maj|min|dim|aug)?\s*[-|/]\s*)+[A-G][#b]?(m|maj|min|dim|aug)?\s*$`),
		regexp.MustCompile(`(?i)^\s*(?:[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?)(?:\s+(?:[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?))*\s*$`), // Space-separated chords
		regexp.MustCompile(`(?i)^\s*(?:verse|chorus|bridge|intro|outro|pre-chorus|prechorus)\s*:\s*(?:[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?)(?:\s*[-|/ ]\s*[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?)*\s*$`), // Section-labelled chords
		regexp.MustCompile(`(?i)^\s*([IVX]+[\/]?\s*)+[IVX]+\s*$`), // Roman numerals
		regexp.MustCompile(`(?i)^\s*(?:[IVX]+[ivx]?)(?:\s*[-|/ ]\s*(?:[IVX]+[ivx]?))*\s*$`),                          // Roman numeral sequences with separators
		regexp.MustCompile(`(?i)^\s*(?:verse|chorus|bridge|intro|outro|pre-chorus|prechorus)\s*:\s*(?:[IVX]+[ivx]?)(?:\s*[-|/ ]\s*(?:[IVX]+[ivx]?))*\s*$`), // Section-labelled roman numerals
		regexp.MustCompile(`(?i)^\s*\|(?:\s*[IVX]+[ivx]?\s*\|)+\s*$`),                                                 // Roman numeral grids with bars
		// Musical notation
		regexp.MustCompile(`(?i)^\s*(verse|chorus)\s*:\s*[A-G]`),
		regexp.MustCompile(`(?i)^\s*tempo\s*:\s*\d+\s*bpm`),
		regexp.MustCompile(`(?i)^\s*key\s*:\s*[A-G][#b]?\s*(major|minor)?`),
		regexp.MustCompile(`(?i)^\s*time\s*signature\s*:\s*\d+\/\d+`),
		regexp.MustCompile(`(?i)^\s*(bpm|tempo)\s*[:=]\s*\d+`),
		// Drum patterns
		regexp.MustCompile(`(?i)^\s*[xXoO\-\|]+\s*$`), // Drum notation
		regexp.MustCompile(`(?i)^\s*(kick|snare|hihat|hi-hat|ride|tom|floor tom|cymbal)\s*[:=]\s*[xXoO\-\|]+`),
		// Pattern structures
		regexp.MustCompile(`(?i)^\s*pattern\s*\d*\s*[:\-=]`),
		regexp.MustCompile(`(?i)^\s*loop\s*\d*\s*[:\-=]`),
		regexp.MustCompile(`(?i)^\s*beat\s*\d*\s*[:\-=]`),
		regexp.MustCompile(`(?i)^\s*rhythm\s*\d*\s*[:\-=]`),
	}

	detector.lyricStructuralPatterns = map[*regexp.Regexp]bool{
		detector.lyricPatterns[0]: true,
		detector.lyricPatterns[1]: true,
	}

	detector.lyricStrongPatterns = map[*regexp.Regexp]bool{
		detector.lyricPatterns[2]: true,
		detector.lyricPatterns[3]: true,
		detector.lyricPatterns[4]: true,
		detector.lyricPatterns[5]: true,
		detector.lyricPatterns[6]: true,
	}

	detector.lyricLexicalPatterns = map[*regexp.Regexp]bool{
		detector.lyricPatterns[3]: true,
		detector.lyricPatterns[4]: true,
		detector.lyricPatterns[6]: true,
	}

	detector.patternStrongPatterns = map[*regexp.Regexp]bool{
		detector.patternPatterns[0]:  true,
		detector.patternPatterns[1]:  true,
		detector.patternPatterns[2]:  true,
		detector.patternPatterns[3]:  true,
		detector.patternPatterns[4]:  true,
		detector.patternPatterns[5]:  true,
		detector.patternPatterns[6]:  true,
		detector.patternPatterns[12]: true,
		detector.patternPatterns[13]: true,
	}

	return detector
}

// AnalyzeContent determines the content type of the provided text
func (cd *ContextDetector) AnalyzeContent(content string) ContentType {
	if content == "" {
		return ContentTypeUnknown
	}

	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ContentTypeUnknown
	}

	// Get the last N lines for analysis (most relevant context)
	startLine := max(0, len(lines)-cd.analysisWindow)
	recentLines := lines[startLine:]

	if len(recentLines) == 0 {
		return ContentTypeUnknown
	}

	// Count matches for each content type
	lyricScore := 0.0
	patternScore := 0.0
	totalLines := len(recentLines)
	strongLyricLines := 0
	strongPatternLines := 0
	lyricLexicalHits := 0

	for _, line := range recentLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lineLyricMatched := false
		lineLyricStrong := false
		lineLyricStructural := false
		lineLyricLexical := false
		for _, pattern := range cd.lyricPatterns {
			if pattern.MatchString(line) {
				lineLyricMatched = true
				if cd.lyricStructuralPatterns != nil && cd.lyricStructuralPatterns[pattern] {
					lineLyricStructural = true
				}
				if cd.lyricStrongPatterns != nil && cd.lyricStrongPatterns[pattern] {
					lineLyricStrong = true
				}
				if cd.lyricLexicalPatterns != nil && cd.lyricLexicalPatterns[pattern] {
					lineLyricLexical = true
				}
			}
		}

		linePatternMatched := false
		linePatternStrong := false
		for _, pattern := range cd.patternPatterns {
			if pattern.MatchString(line) {
				linePatternMatched = true
				if cd.patternStrongPatterns != nil && cd.patternStrongPatterns[pattern] {
					linePatternStrong = true
				}
			}
		}

		if linePatternMatched {
			if linePatternStrong {
				patternScore += 1.5
				strongPatternLines++
			} else {
				patternScore += 1.0
			}
		}

		if lineLyricMatched {
			if linePatternMatched && lineLyricStructural && !lineLyricStrong {
				continue
			}

			if lineLyricStrong {
				if linePatternMatched {
					lyricScore += 1.2
				} else {
					lyricScore += 1.5
				}
				strongLyricLines++
			} else if !linePatternMatched {
				lyricScore += 0.6
			} else {
				lyricScore += 0.4
			}
		}

		if lineLyricLexical {
			lyricLexicalHits++
		}
	}

	// Normalize scores
	rawLyricRatio := 0.0
	rawPatternRatio := 0.0
	if totalLines > 0 {
		rawLyricRatio = lyricScore / float64(totalLines)
		rawPatternRatio = patternScore / float64(totalLines)
	}

	lyricRatio := min(rawLyricRatio, 1.0)
	patternRatio := min(rawPatternRatio, 1.0)

	// Determine content type based on ratios
	if lyricRatio == 0 && patternRatio == 0 {
		return ContentTypeUnknown
	}

	if totalLines <= 2 {
		if patternRatio == 0 && rawLyricRatio < 0.85 {
			return ContentTypeUnknown
		}
		if lyricRatio == 0 && rawPatternRatio < 0.85 {
			return ContentTypeUnknown
		}
		if patternRatio == 0 && lyricLexicalHits == 0 {
			return ContentTypeUnknown
		}
	}

	if lyricScore > 0 && lyricLexicalHits == 0 && strongLyricLines == 0 && patternScore == 0 {
		return ContentTypeUnknown
	}

	if strongLyricLines == 0 && patternScore == 0 && rawLyricRatio < 0.7 {
		return ContentTypeUnknown
	}

	if strongPatternLines == 0 && lyricScore == 0 && rawPatternRatio < 0.7 {
		return ContentTypeUnknown
	}

	strongPresenceThreshold := 0.75
	secondaryPresenceThreshold := 0.2
	dominanceMargin := 0.15
	nearMixedThreshold := cd.mixedThreshold - 0.1
	if nearMixedThreshold < 0.3 {
		nearMixedThreshold = 0.3
	}

	// High-confidence single type detection
	if lyricRatio >= strongPresenceThreshold && patternRatio <= secondaryPresenceThreshold {
		return ContentTypeLyrics
	}

	if patternRatio >= strongPresenceThreshold && lyricRatio <= secondaryPresenceThreshold {
		return ContentTypePatterns
	}

	// Mixed detection requires meaningful presence from both, with a softer fallback threshold
	if lyricRatio >= cd.mixedThreshold && patternRatio >= cd.mixedThreshold {
		return ContentTypeMixed
	}

	if lyricRatio >= nearMixedThreshold && patternRatio >= nearMixedThreshold {
		return ContentTypeMixed
	}

	// Apply dominance margin to reduce minor lyric priority bias
	if lyricRatio-patternRatio >= dominanceMargin {
		return ContentTypeLyrics
	}

	if patternRatio-lyricRatio >= dominanceMargin {
		return ContentTypePatterns
	}

	// Use absolute scores to resolve close calls
	if patternScore > lyricScore {
		return ContentTypePatterns
	}

	if lyricScore > patternScore {
		return ContentTypeLyrics
	}

	// Fallbacks when ratios are very close
	if lyricRatio > 0 && patternRatio > 0 {
		return ContentTypeMixed
	}

	if patternRatio > 0 {
		return ContentTypePatterns
	}

	return ContentTypeLyrics
}

// GetContextAnalysis provides detailed analysis of content
func (cd *ContextDetector) GetContextAnalysis(content string) *ContextAnalysis {
	if content == "" {
		return &ContextAnalysis{
			ContentType: ContentTypeUnknown,
			Confidence:  0.0,
			Details:     "Empty content",
		}
	}

	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return &ContextAnalysis{
			ContentType: ContentTypeUnknown,
			Confidence:  0.0,
			Details:     "No lines to analyze",
		}
	}

	// Get the last N lines for analysis
	startLine := max(0, len(lines)-cd.analysisWindow)
	recentLines := lines[startLine:]

	if len(recentLines) == 0 {
		return &ContextAnalysis{
			ContentType: ContentTypeUnknown,
			Confidence:  0.0,
			Details:     "No recent lines to analyze",
		}
	}

	// Detailed analysis
	lyricMatches := 0
	patternMatches := 0
	emptyLines := 0
	totalLines := len(recentLines)
	var matchedLyricPatterns []string
	var matchedPatternPatterns []string

	for _, line := range recentLines {
		line = strings.TrimSpace(line)
		if line == "" {
			emptyLines++
			continue
		}

		// Check for lyric patterns
		for _, pattern := range cd.lyricPatterns {
			if pattern.MatchString(line) {
				lyricMatches++
				matchedLyricPatterns = append(matchedLyricPatterns, pattern.String())
				break
			}
		}

		// Check for pattern patterns
		for _, pattern := range cd.patternPatterns {
			if pattern.MatchString(line) {
				patternMatches++
				matchedPatternPatterns = append(matchedPatternPatterns, pattern.String())
				break
			}
		}
	}

	// Calculate ratios and confidence
	lyricRatio := float64(lyricMatches) / float64(totalLines)
	patternRatio := float64(patternMatches) / float64(totalLines)
	
	contentType := cd.AnalyzeContent(content)
	
	// Calculate confidence based on how dominant the detected type is
	var confidence float64
	switch contentType {
	case ContentTypeLyrics:
		confidence = lyricRatio
	case ContentTypePatterns:
		confidence = patternRatio
	case ContentTypeMixed:
		confidence = min(lyricRatio, patternRatio) * 2 // Confidence in mixed is lower
	default:
		confidence = 0.0
	}

	// Cap confidence at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return &ContextAnalysis{
		ContentType:           contentType,
		Confidence:            confidence,
		TotalLines:            totalLines,
		LyricMatches:          lyricMatches,
		PatternMatches:        patternMatches,
		EmptyLines:            emptyLines,
		MatchedLyricPatterns:  matchedLyricPatterns,
		MatchedPatternPatterns: matchedPatternPatterns,
		Details:               fmt.Sprintf("Analyzed %d lines, found %d lyric matches and %d pattern matches", totalLines, lyricMatches, patternMatches),
	}
}

// ContextAnalysis provides detailed information about content analysis
type ContextAnalysis struct {
	ContentType           ContentType
	Confidence            float64
	TotalLines            int
	LyricMatches          int
	PatternMatches        int
	EmptyLines            int
	MatchedLyricPatterns  []string
	MatchedPatternPatterns []string
	Details               string
}

// IsLyricContent returns true if the content is primarily lyrics
func (ca *ContextAnalysis) IsLyricContent() bool {
	return ca.ContentType == ContentTypeLyrics || 
		(ca.ContentType == ContentTypeMixed && ca.LyricMatches >= ca.PatternMatches)
}

// IsPatternContent returns true if the content is primarily patterns
func (ca *ContextAnalysis) IsPatternContent() bool {
	return ca.ContentType == ContentTypePatterns || 
		(ca.ContentType == ContentTypeMixed && ca.PatternMatches > ca.LyricMatches)
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}