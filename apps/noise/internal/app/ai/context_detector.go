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
	// 1. Define all patterns as local variables for clarity and safety
	pLyricsHeaders := regexp.MustCompile(`(?i)^\s*(verse|chorus|bridge|intro|outro|pre-chorus)\s*[:\-\d]*`)
	pLyricsBrackets := regexp.MustCompile(`(?i)^\s*\[.*\]\s*$`)
	pLyricsSentences := regexp.MustCompile(`(?m)^\s*[A-Z][^.!?]*[.!?]`)
	pLyricsContractions := regexp.MustCompile(`(?i)\b(I'm|you're|we're|they're|can't|won't|don't)\b`)
	pLyricsEmotional := regexp.MustCompile(`(?i)\b(love|heart|night|day|sky|rain|sun|moon|stars|tonight|care|gone|back|time|life|feeling|lights)\b`)
	pLyricsRhyming := regexp.MustCompile(`(?m)^\s*[A-Z][a-z]+\s+[a-z]+\s+[a-z]+\s+`)
	pLyricsEndearment := regexp.MustCompile(`(?i)\b(baby|honey|darling|sweetheart)\b`)
	pLyricsShortLines := regexp.MustCompile(`(?m)^\s*\w+(?:\s+\w+){1,5}\s*$`)

	pPatternChordSeq := regexp.MustCompile(`(?i)^\s*([A-G][#b]?(m|maj|min|dim|aug)?\s*[-|/]\s*)+[A-G][#b]?(m|maj|min|dim|aug)?\s*$`)
	pPatternChordSpace := regexp.MustCompile(`(?i)^\s*(?:[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?)(?:\s+(?:[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?))*\s*$`)
	pPatternChordLabel := regexp.MustCompile(`(?i)^\s*(?:verse|chorus|bridge|intro|outro|pre-chorus|prechorus)\s*:\s*(?:[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?)(?:\s*[-|/ ]\s*[A-G][#b]?(?:maj|min|m|dim|aug|sus\d*|add\d*|m7|7|maj7|dim7|aug7)?)*\s*$`)
	pPatternRomanSeq := regexp.MustCompile(`(?i)^\s*([IVX]+[\/]?\s*)+[IVX]+\s*$`)
	pPatternRomanComplex := regexp.MustCompile(`(?i)^\s*(?:[IVX]+[ivx]?)(?:\s*[-|/ ]\s*(?:[IVX]+[ivx]?))*\s*$`)
	pPatternRomanLabel := regexp.MustCompile(`(?i)^\s*(?:verse|chorus|bridge|intro|outro|pre-chorus|prechorus)\s*:\s*(?:[IVX]+[ivx]?)(?:\s*[-|/ ]\s*(?:[IVX]+[ivx]?))*\s*$`)
	pPatternRomanGrid := regexp.MustCompile(`(?i)^\s*\|(?:\s*[IVX]+[ivx]?\s*\|)+\s*$`)
	pPatternTempo := regexp.MustCompile(`(?i)^\s*tempo\s*:\s*\d+\s*bpm`)
	pPatternKey := regexp.MustCompile(`(?i)^\s*key\s*:\s*[A-G][#b]?\s*(major|minor)?`)
	pPatternTimeSig := regexp.MustCompile(`(?i)^\s*time\s*signature\s*:\s*\d+\/\d+`)
	pPatternBPM := regexp.MustCompile(`(?i)^\s*(bpm|tempo)\s*[:=]\s*\d+`)
	pPatternDrumFlat := regexp.MustCompile(`(?i)^\s*[xXoO\-\|]+\s*$`)
	pPatternDrumMapped := regexp.MustCompile(`(?i)^\s*(kick|snare|hihat|hi-hat|ride|tom|floor tom|cymbal)\s*[:=]\s*[xXoO\-\|]+`)
	pPatternStructure := regexp.MustCompile(`(?i)^\s*(pattern|loop|beat|rhythm)\s*\d*\s*[:\-=]`)

	// 2. Initialize detector
	detector := &ContextDetector{
		mixedThreshold: 0.45,
		analysisWindow: 20,
	}

	detector.lyricPatterns = []*regexp.Regexp{
		pLyricsHeaders, pLyricsBrackets, pLyricsSentences, pLyricsContractions,
		pLyricsEmotional, pLyricsRhyming, pLyricsEndearment, pLyricsShortLines,
	}

	detector.patternPatterns = []*regexp.Regexp{
		pPatternChordSeq, pPatternChordSpace, pPatternChordLabel,
		pPatternRomanSeq, pPatternRomanComplex, pPatternRomanLabel, pPatternRomanGrid,
		pPatternTempo, pPatternKey, pPatternTimeSig, pPatternBPM,
		pPatternDrumFlat, pPatternDrumMapped, pPatternStructure,
	}

	// 3. Structural patterns
	detector.lyricStructuralPatterns = map[*regexp.Regexp]bool{
		pLyricsHeaders:  true,
		pLyricsBrackets: true,
	}

	// 4. Strong patterns
	detector.lyricStrongPatterns = map[*regexp.Regexp]bool{
		pLyricsSentences:    true,
		pLyricsContractions: true,
		pLyricsEmotional:    true,
		pLyricsRhyming:      true,
		pLyricsEndearment:   true,
	}

	// 5. Lexical patterns
	detector.lyricLexicalPatterns = map[*regexp.Regexp]bool{
		pLyricsEmotional:    true,
		pLyricsEndearment:   true,
		pLyricsContractions: true,
	}

	// 6. Strong pattern patterns
	detector.patternStrongPatterns = make(map[*regexp.Regexp]bool)
	for _, p := range detector.patternPatterns {
		detector.patternStrongPatterns[p] = true
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
			// Skip lyric scoring if it's a structural match AND also a pattern match
			// AND NOT a strong lyric match (which would indicate mixed intent)
			if linePatternMatched && lineLyricStructural && !lineLyricStrong {
				continue
			}

			// CRITICAL: If it's a strong pattern match, treat as pure pattern unless it's a strong lyric
			if linePatternStrong && !lineLyricStrong {
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
	} else {
		return ContentTypeUnknown
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
		// For very short content, be more lenient with pattern detection
		if totalLines == 1 && rawPatternRatio > 0.5 {
			return ContentTypePatterns
		}
	}

	if lyricScore > 0 && lyricLexicalHits == 0 && strongLyricLines == 0 && patternScore == 0 {
		return ContentTypeUnknown
	}

	if strongLyricLines == 0 && patternScore == 0 && rawLyricRatio < 0.5 {
		return ContentTypeUnknown
	}

	if strongPatternLines == 0 && lyricScore == 0 && rawPatternRatio < 0.5 {
		return ContentTypeUnknown
	}

	// 1. High-confidence single type detection (pure or largely pure content)
	if patternRatio == 0 && lyricRatio >= 0.25 {
		return ContentTypeLyrics
	}
	if lyricRatio == 0 && patternRatio >= 0.25 {
		return ContentTypePatterns
	}

	// 2. Strong dominance check (significant lead even if both types present)
	// If one type is more than 3x the other, it dominates unless it's very sparse
	if lyricRatio >= 3.0*patternRatio && lyricRatio >= 0.4 {
		return ContentTypeLyrics
	}
	if patternRatio >= 3.0*lyricRatio && patternRatio >= 0.4 {
		return ContentTypePatterns
	}

	// 2b. Strong pattern lines indicate pure pattern content
	// When most lines are strong pattern matches and there are few strong lyric indicators
	if strongPatternLines >= 3 && strongLyricLines == 0 && patternRatio > lyricRatio {
		return ContentTypePatterns
	}
	if strongLyricLines >= 3 && strongPatternLines == 0 && lyricRatio > patternRatio {
		return ContentTypeLyrics
	}

	// 3. Balanced mixed detection (both have meaningful presence with strong indicators)
	// Only classify as mixed if both types have substantial strong matches
	if lyricRatio >= 0.25 && patternRatio >= 0.25 && strongLyricLines > 0 && strongPatternLines > 0 {
		return ContentTypeMixed
	}

	// 4. Weaker mixed detection for cases with clear presence of both
	if lyricRatio >= 0.3 && patternRatio >= 0.3 {
		return ContentTypeMixed
	}

	// 5. Use absolute scores to resolve close calls for low-presence content
	if patternScore > lyricScore {
		return ContentTypePatterns
	}
	if lyricScore > patternScore {
		return ContentTypeLyrics
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
	var lyricRatio, patternRatio float64
	if totalLines > 0 {
		lyricRatio = float64(lyricMatches) / float64(totalLines)
		patternRatio = float64(patternMatches) / float64(totalLines)
	} else {
		return &ContextAnalysis{
			ContentType: ContentTypeUnknown,
			Confidence:  0.0,
			Details:     "No lines to analyze",
		}
	}

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
		ContentType:            contentType,
		Confidence:             confidence,
		TotalLines:             totalLines,
		LyricMatches:           lyricMatches,
		PatternMatches:         patternMatches,
		EmptyLines:             emptyLines,
		MatchedLyricPatterns:   matchedLyricPatterns,
		MatchedPatternPatterns: matchedPatternPatterns,
		Details:                fmt.Sprintf("Analyzed %d lines, found %d lyric matches and %d pattern matches", totalLines, lyricMatches, patternMatches),
	}
}

// ContextAnalysis provides detailed information about content analysis
type ContextAnalysis struct {
	ContentType            ContentType
	Confidence             float64
	TotalLines             int
	LyricMatches           int
	PatternMatches         int
	EmptyLines             int
	MatchedLyricPatterns   []string
	MatchedPatternPatterns []string
	Details                string
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
