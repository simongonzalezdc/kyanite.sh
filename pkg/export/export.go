// Package export provides a shared interface and helpers for file export
// across kyanite.sh apps. Each app (noise, prism, syntax) has its own
// domain-specific exporters that implement the Formatter interface.
//
// Usage:
//
//	exporter := export.New("palette", "css")
//	data, err := exporter.Export(paletteData)
//	err = exporter.WriteFile("palette.css", data)
package export

import (
	"fmt"
	"os"
	"path/filepath"
)

// Formatter converts domain data into a specific output format.
// Each app registers formatters for its export types.
type Formatter interface {
	// Format converts data to bytes in the target format.
	// data is the domain-specific struct (e.g., a palette, song, or document).
	Format(data any) ([]byte, error)

	// Extension returns the file extension without dot (e.g., "css", "json").
	Extension() string
}

// FormatterFunc is a convenience type for single-function formatters.
// Use NewFormatter to create a full Formatter with an extension.
type FormatterFunc func(data any) ([]byte, error)

func (f FormatterFunc) Format(data any) ([]byte, error) { return f(data) }

// formatterWrapper wraps a FormatterFunc with an extension.
type formatterWrapper struct {
	FormatterFunc
	ext string
}

func (w formatterWrapper) Extension() string { return w.ext }

// NewFormatter creates a Formatter from a function and file extension.
func NewFormatter(ext string, f func(data any) ([]byte, error)) Formatter {
	return formatterWrapper{FormatterFunc: FormatterFunc(f), ext: ext}
}

// Registry holds named formatters indexed by category and format.
type Registry struct {
	formatters map[string]map[string]Formatter // category → format → Formatter
}

// NewRegistry creates an empty formatter registry.
func NewRegistry() *Registry {
	return &Registry{formatters: make(map[string]map[string]Formatter)}
}

// Register adds a formatter under a category and format name.
func (r *Registry) Register(category, format string, f Formatter) {
	if r.formatters[category] == nil {
		r.formatters[category] = make(map[string]Formatter)
	}
	r.formatters[category][format] = f
}

// Lookup finds a formatter by category and format.
// Returns nil if not found.
func (r *Registry) Lookup(category, format string) Formatter {
	if cats, ok := r.formatters[category]; ok {
		return cats[format]
	}
	return nil
}

// Formats returns all registered format names for a category.
func (r *Registry) Formats(category string) []string {
	cats, ok := r.formatters[category]
	if !ok {
		return nil
	}
	formats := make([]string, 0, len(cats))
	for f := range cats {
		formats = append(formats, f)
	}
	return formats
}

// Export formats data using the named formatter and returns bytes.
func Export(r *Registry, category, format string, data any) ([]byte, error) {
	f := r.Lookup(category, format)
	if f == nil {
		return nil, fmt.Errorf("no formatter for %s/%s (available: %v)", category, format, r.Formats(category))
	}
	return f.Format(data)
}

// WriteFile writes formatted data to a file, creating parent directories.
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}
