// Package examples provides example plugins demonstrating the plugin API.
package examples

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/plugins"
)

// ChordAnalyzerPlugin is an example plugin that adds chord analysis functionality
type ChordAnalyzerPlugin struct {
	metadata *plugins.PluginMetadata
	enabled  bool
}

// NewChordAnalyzerPlugin creates a new chord analyzer plugin
func NewChordAnalyzerPlugin() *ChordAnalyzerPlugin {
	return &ChordAnalyzerPlugin{
		metadata: &plugins.PluginMetadata{
			ID:          "chord_analyzer",
			Name:        "Chord Analyzer",
			Version:     "1.0.0",
			Description: "Analyzes chords in lyrics and provides theory insights",
			Author:      "noise.sh Team",
			License:     "MIT",
			Capabilities: []plugins.Capability{
				plugins.CapabilityTheoryTool,
				plugins.CapabilityEditorTool,
			},
			Enabled: true,
		},
		enabled: false,
	}
}

// Metadata returns the plugin metadata
func (p *ChordAnalyzerPlugin) Metadata() *plugins.PluginMetadata {
	return p.metadata
}

// Initialize sets up the plugin
func (p *ChordAnalyzerPlugin) Initialize(ctx *plugins.PluginContext) error {
	p.metadata.LoadTime = time.Now()
	p.enabled = p.metadata.Enabled
	return nil
}

// Cleanup performs cleanup when the plugin is unloaded
func (p *ChordAnalyzerPlugin) Cleanup() error {
	return nil
}

// Enable activates the plugin
func (p *ChordAnalyzerPlugin) Enable() error {
	p.enabled = true
	return nil
}

// Disable deactivates the plugin
func (p *ChordAnalyzerPlugin) Disable() error {
	p.enabled = false
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (p *ChordAnalyzerPlugin) IsEnabled() bool {
	return p.enabled
}

// AnalyzeChords analyzes chords in the given content
func (p *ChordAnalyzerPlugin) AnalyzeChords(content string) (*ChordAnalysis, error) {
	chords := p.extractChords(content)

	analysis := &ChordAnalysis{
		TotalChords:       len(chords),
		UniqueChords:      p.getUniqueChords(chords),
		ChordProgressions: p.identifyProgressions(chords),
		KeySignature:      p.detectKeySignature(chords),
		Complexity:        p.calculateComplexity(chords),
	}

	return analysis, nil
}

// extractChords extracts chord names from content using regex
func (p *ChordAnalyzerPlugin) extractChords(content string) []string {
	// Common chord patterns (major, minor, 7th, etc.)
	chordRegex := regexp.MustCompile(`\b([A-G][#b]?)(maj|min|m|aug|dim|sus|add)?(\d+)?\b`)
	matches := chordRegex.FindAllStringSubmatch(content, -1)

	var chords []string
	for _, match := range matches {
		if len(match) > 1 && match[1] != "" {
			chord := match[1]
			if len(match) > 2 && match[2] != "" {
				chord += match[2]
			}
			if len(match) > 3 && match[3] != "" {
				chord += match[3]
			}
			chords = append(chords, chord)
		}
	}

	return chords
}

// getUniqueChords returns unique chord names
func (p *ChordAnalyzerPlugin) getUniqueChords(chords []string) []string {
	seen := make(map[string]bool)
	var unique []string

	for _, chord := range chords {
		if !seen[chord] {
			seen[chord] = true
			unique = append(unique, chord)
		}
	}

	return unique
}

// identifyProgressions identifies common chord progressions
func (p *ChordAnalyzerPlugin) identifyProgressions(chords []string) []string {
	var progressions []string

	// Simple progression detection (I-V-vi-IV is very common)
	if len(chords) >= 4 {
		// This is a simplified example - real progression detection would be more sophisticated
		progressions = append(progressions, "I-V-vi-IV pattern detected")
	}

	return progressions
}

// detectKeySignature attempts to detect the key signature
func (p *ChordAnalyzerPlugin) detectKeySignature(chords []string) string {
	// Simple key detection based on most common chords
	// In reality, this would use music theory algorithms
	if len(chords) == 0 {
		return "Unknown"
	}

	// Count chord frequencies
	chordCount := make(map[string]int)
	for _, chord := range chords {
		chordCount[chord]++
	}

	// Find most frequent chord (likely the tonic)
	var mostFrequent string
	var maxCount int
	for chord, count := range chordCount {
		if count > maxCount {
			maxCount = count
			mostFrequent = chord
		}
	}

	return fmt.Sprintf("Key of %s (estimated)", mostFrequent)
}

// calculateComplexity calculates the complexity score of the chord progression
func (p *ChordAnalyzerPlugin) calculateComplexity(chords []string) int {
	if len(chords) == 0 {
		return 0
	}

	uniqueChords := p.getUniqueChords(chords)
	// Simple complexity score based on number of unique chords
	return len(uniqueChords) * 10
}

// ChordAnalysis represents the result of chord analysis
type ChordAnalysis struct {
	TotalChords       int      `json:"total_chords"`
	UniqueChords      []string `json:"unique_chords"`
	ChordProgressions []string `json:"chord_progressions"`
	KeySignature      string   `json:"key_signature"`
	Complexity        int      `json:"complexity"`
}

// GetEditorTool returns the editor tool definition for this plugin
func (p *ChordAnalyzerPlugin) GetEditorTool() *plugins.EditorTool {
	return &plugins.EditorTool{
		ID:          "chord_analyzer",
		Name:        "Analyze Chords",
		Description: "Analyze chords and provide theory insights",
		Icon:        "â™ª",
		Handler: func(content string) (string, error) {
			analysis, err := p.AnalyzeChords(content)
			if err != nil {
				return "", err
			}

			result := "Chord Analysis Results:\n\n"
			result += "Total Chords: " + fmt.Sprintf("%d", analysis.TotalChords) + "\n"
			result += "Unique Chords: " + strings.Join(analysis.UniqueChords, ", ") + "\n"
			result += "Key Signature: " + analysis.KeySignature + "\n"
			result += "Complexity Score: " + fmt.Sprintf("%d", analysis.Complexity) + "\n"

			if len(analysis.ChordProgressions) > 0 {
				result += fmt.Sprintf("Progressions: %s\n", strings.Join(analysis.ChordProgressions, ", "))
			}

			return result, nil
		},
		Shortcut: "ctrl+a",
		Enabled:  p.enabled,
	}
}

// GetManifest returns the plugin manifest for file-based loading
func (p *ChordAnalyzerPlugin) GetManifest() ([]byte, error) {
	return json.MarshalIndent(p.metadata, "", "  ")
}
