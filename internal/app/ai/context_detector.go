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
	lyricPatterns    []*regexp.Regexp
	patternPatterns  []*regexp.Regexp
	mixedThreshold   float64
	analysisWindow   int // Number of lines to analyze for context
}

// NewContextDetector creates a new context detector with default patterns
func NewContextDetector() *ContextDetector {
	detector := &ContextDetector{
		mixedThreshold: 0.3, // 30% of one type qualifies as mixed
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
		regexp.MustCompile(`(?i)^\s*([IVX]+[\/]?\s*)+[IVX]+\s*$`), // Roman numerals
		// Musical notation
		regexp.MustCompile(`(?i)^\s*(verse|chorus)\s*:\s*[A-G]`),
		regexp.MustCompile(`(?i)^\s*tempo\s*:\s*\d+\s*bpm`),
		regexp.MustCompile(`(?i)^\s*key\s*:\s*[A-G][#b]?\s*(major|minor)?`),
		regexp.MustCompile(`(?i)^\s*time\s*signature\s*:\s*\d+\/\d+`),
		regexp.MustCompile(`(?i)^\s*(bpm|tempo)\s*[:=]\s*\d+`),
		// Drum patterns
		regexp.MustCompile(`(?i)^\s*[xXoO\-\|]+\s*$`), // Drum notation
		regexp.MustCompile(`(?i)^\s*(kick|snare|hihat|cymbal)\s*[:=]\s*[xXoO\-\|]+`),
		// Pattern structures
		regexp.MustCompile(`(?i)^\s*pattern\s*\d*\s*[:\-=]`),
		regexp.MustCompile(`(?i)^\s*loop\s*\d*\s*[:\-=]`),
		regexp.MustCompile(`(?i)^\s*beat\s*\d*\s*[:\-=]`),
		regexp.MustCompile(`(?i)^\s*rhythm\s*\d*\s*[:\-=]`),
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
	lyricScore := 0
	patternScore := 0
	totalLines := len(recentLines)

	for _, line := range recentLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for lyric patterns
		for _, pattern := range cd.lyricPatterns {
			if pattern.MatchString(line) {
				lyricScore++
				break // Only count once per line
			}
		}

		// Check for pattern patterns
		for _, pattern := range cd.patternPatterns {
			if pattern.MatchString(line) {
				patternScore++
				break // Only count once per line
			}
		}
	}

	// Normalize scores
	lyricRatio := float64(lyricScore) / float64(totalLines)
	patternRatio := float64(patternScore) / float64(totalLines)

	// Determine content type based on ratios
	if lyricRatio == 0 && patternRatio == 0 {
		return ContentTypeUnknown
	}

	if lyricRatio > 0.7 && patternRatio < 0.3 {
		return ContentTypeLyrics
	}

	if patternRatio > 0.7 && lyricRatio < 0.3 {
		return ContentTypePatterns
	}

	// If both have significant presence, it's mixed
	if lyricRatio >= cd.mixedThreshold && patternRatio >= cd.mixedThreshold {
		return ContentTypeMixed
	}

	// Default to the dominant type
	if lyricRatio > patternRatio {
		return ContentTypeLyrics
	}

	return ContentTypePatterns
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