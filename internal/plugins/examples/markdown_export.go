package examples

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/plugins"
)

// MarkdownExportPlugin is an example plugin that adds Markdown export functionality
type MarkdownExportPlugin struct {
	metadata *plugins.PluginMetadata
	enabled  bool
}

// NewMarkdownExportPlugin creates a new Markdown export plugin
func NewMarkdownExportPlugin() *MarkdownExportPlugin {
	return &MarkdownExportPlugin{
		metadata: &plugins.PluginMetadata{
			ID:          "markdown_export",
			Name:        "Markdown Export",
			Version:     "1.0.0",
			Description: "Export lyrics as formatted Markdown",
			Author:      "noise.sh Team",
			License:     "MIT",
			Capabilities: []plugins.Capability{
				plugins.CapabilityExportFormat,
			},
			Enabled: true,
		},
		enabled: false,
	}
}

// Metadata returns the plugin metadata
func (p *MarkdownExportPlugin) Metadata() *plugins.PluginMetadata {
	return p.metadata
}

// Initialize sets up the plugin
func (p *MarkdownExportPlugin) Initialize(ctx *plugins.PluginContext) error {
	p.metadata.LoadTime = time.Now()
	p.enabled = p.metadata.Enabled

	// Register the export format with the plugin manager
	// In a real implementation, this would be done through the plugin context
	if ctx != nil {
		// This is a simplified approach - in reality, the plugin manager
		// would provide a registration mechanism
	}

	return nil
}

// Cleanup performs cleanup when the plugin is unloaded
func (p *MarkdownExportPlugin) Cleanup() error {
	return nil
}

// Enable activates the plugin
func (p *MarkdownExportPlugin) Enable() error {
	p.enabled = true
	return nil
}

// Disable deactivates the plugin
func (p *MarkdownExportPlugin) Disable() error {
	p.enabled = false
	return nil
}

// IsEnabled returns whether the plugin is currently enabled
func (p *MarkdownExportPlugin) IsEnabled() bool {
	return p.enabled
}

// exportToMarkdown converts content to Markdown format
func (p *MarkdownExportPlugin) exportToMarkdown(content string) ([]byte, error) {
	lines := strings.Split(content, "\n")
	var markdown strings.Builder

	// Add header
	markdown.WriteString("# Exported Lyrics\n\n")
	markdown.WriteString(fmt.Sprintf("Exported on: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	markdown.WriteString("---\n\n")

	// Process content
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines at the beginning of sections
		if trimmed == "" {
			if i > 0 && strings.TrimSpace(lines[i-1]) != "" {
				markdown.WriteString("\n")
			}
			continue
		}

		// Check if this looks like a chorus (common patterns)
		if p.isChorusLine(trimmed) {
			markdown.WriteString(fmt.Sprintf("**%s**\n\n", trimmed))
		} else {
			markdown.WriteString(fmt.Sprintf("%s\n", trimmed))
		}
	}

	return []byte(markdown.String()), nil
}

// isChorusLine checks if a line appears to be a chorus
func (p *MarkdownExportPlugin) isChorusLine(line string) bool {
	chorusIndicators := []string{
		"chorus", "refrain", "(chorus)", "[chorus]",
		"hook", "(hook)", "[hook]",
	}

	lineLower := strings.ToLower(line)
	for _, indicator := range chorusIndicators {
		if strings.Contains(lineLower, indicator) {
			return true
		}
	}

	return false
}

// GetExportFormat returns the export format definition for this plugin
func (p *MarkdownExportPlugin) GetExportFormat() *plugins.ExportFormat {
	return &plugins.ExportFormat{
		ID:          "markdown",
		Name:        "Markdown",
		Description: "Export as Markdown with formatting",
		Extension:   ".md",
		MimeType:    "text/markdown",
		Handler: func(content string) ([]byte, error) {
			return p.exportToMarkdown(content)
		},
		Enabled: p.enabled,
	}
}

// GetManifest returns the plugin manifest for file-based loading
func (p *MarkdownExportPlugin) GetManifest() ([]byte, error) {
	return json.MarshalIndent(p.metadata, "", "  ")
}
