package editor

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
)

// PreviewRenderer handles markdown rendering with caching
type PreviewRenderer struct {
	renderer *glamour.TermRenderer
	cache    *PreviewCache
	width    int
}

// NewPreviewRenderer creates a new preview renderer
func NewPreviewRenderer(width int, themeName string) (*PreviewRenderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	return &PreviewRenderer{
		renderer: r,
		cache:    NewPreviewCache(100),
		width:    width,
	}, nil
}

// Render renders markdown content with caching
func (r *PreviewRenderer) Render(content string) (string, error) {
	// Generate cache key
	key := r.cacheKey(content)

	// Check cache
	if cached, ok := r.cache.Get(key); ok {
		return cached, nil
	}

	// Render
	rendered, err := r.renderer.Render(content)
	if err != nil {
		return "", err
	}

	// Cache result
	r.cache.Set(key, content, rendered)

	return rendered, nil
}

// cacheKey generates a cache key for content
func (r *PreviewRenderer) cacheKey(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x-%d", hash, r.width)
}

// SetWidth updates the renderer width
func (r *PreviewRenderer) SetWidth(width int) error {
	if width == r.width {
		return nil
	}

	r.width = width
	r.cache.Clear() // Clear cache as width affects rendering

	// Recreate renderer with new width
	newRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return err
	}

	r.renderer = newRenderer
	return nil
}

// CacheStats returns cache statistics
func (r *PreviewRenderer) CacheStats() (hits, misses int64, size int) {
	return r.cache.Stats()
}

// SyntaxHighlight applies syntax highlighting to code blocks
func (r *PreviewRenderer) SyntaxHighlight(content, language string) string {
	// For now, just wrap in code fence
	// In production, would use actual syntax highlighter
	return fmt.Sprintf("```%s\n%s\n```", language, content)
}

// FormatChords formats chord notation
func (r *PreviewRenderer) FormatChords(line string) string {
	// Basic chord detection and formatting
	// Chords typically appear in brackets [C] or uppercase
	formatted := strings.ReplaceAll(line, "[", "**[")
	formatted = strings.ReplaceAll(formatted, "]", "]**")
	return formatted
}
