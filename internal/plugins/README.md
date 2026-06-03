# noise.sh Plugin System

The noise.sh plugin system allows developers to extend the application's functionality without modifying the core codebase. This document explains how to create, install, and manage plugins.

## Overview

The plugin system provides a secure and stable foundation for extending noise.sh with:

- **UI Extensions**: Add new screens, menu items, or modify existing UI
- **Editor Tools**: Add new editing capabilities and tools
- **Export Formats**: Support for additional export formats
- **Theory Tools**: Music theory analysis and tools
- **Audio Features**: Audio processing and MIDI handling
- **Data Integration**: External data sources and exporters

## Architecture

### Core Components

1. **Plugin Interface** (`Plugin`): Main interface all plugins must implement
2. **Plugin Manager** (`PluginManager`): Handles plugin discovery, loading, and lifecycle
3. **Security Manager** (`SecurityManager`): Validates plugins and ensures safe execution
4. **Plugin Context** (`PluginContext`): Provides access to application services

### Plugin Lifecycle

1. **Discovery**: Plugin manager scans configured directories for plugins
2. **Validation**: Security manager validates plugin files and manifests
3. **Loading**: Valid plugins are loaded into memory
4. **Initialization**: Plugin's `Initialize` method is called
5. **Registration**: Plugin capabilities are registered with the system
6. **Execution**: Plugin responds to application events and user actions

## Creating a Plugin

### 1. Plugin Structure

Create a new Go file that implements the `Plugin` interface:

```go
package main

import (
    "github.com/Kyanite/noise/internal/plugins"
)

type MyPlugin struct {
    metadata *plugins.PluginMetadata
    enabled  bool
}

func (p *MyPlugin) Metadata() *plugins.PluginMetadata {
    return p.metadata
}

func (p *MyPlugin) Initialize(ctx *plugins.PluginContext) error {
    // Initialize plugin with context
    return nil
}

func (p *MyPlugin) Cleanup() error {
    // Cleanup resources
    return nil
}

func (p *MyPlugin) Enable() error {
    p.enabled = true
    return nil
}

func (p *MyPlugin) Disable() error {
    p.enabled = false
    return nil
}

func (p *MyPlugin) IsEnabled() bool {
    return p.enabled
}
```

### 2. Plugin Manifest

Create a `plugin.json` file with plugin metadata:

```json
{
  "id": "my_plugin",
  "name": "My Custom Plugin",
  "version": "1.0.0",
  "description": "A custom plugin for noise.sh",
  "author": "Your Name",
  "license": "MIT",
  "capabilities": [
    "export_format",
    "editor_tool"
  ],
  "config": {
    "custom_option": "value"
  },
  "enabled": true
}
```

### 3. Plugin Capabilities

#### Export Format Plugin

```go
func (p *MyPlugin) GetExportFormat() *plugins.ExportFormat {
    return &plugins.ExportFormat{
        ID:          "my_format",
        Name:        "My Format",
        Description: "Custom export format",
        Extension:   ".myf",
        MimeType:    "application/my-format",
        Handler: func(content string) ([]byte, error) {
            // Convert content to your format
            return []byte(processedContent), nil
        },
        Enabled:     true,
    }
}
```

#### Editor Tool Plugin

```go
func (p *MyPlugin) GetEditorTool() *plugins.EditorTool {
    return &plugins.EditorTool{
        ID:          "my_tool",
        Name:        "My Tool",
        Description: "Custom editor tool",
        Icon:        "âš¡",
        Handler: func(content string) (string, error) {
            // Process content with your tool
            return processedContent, nil
        },
        Shortcut:    "ctrl+m",
        Enabled:     true,
    }
}
```

#### Menu Item Plugin

```go
func (p *MyPlugin) GetMenuItems() []*plugins.MenuItem {
    return []*plugins.MenuItem{
        {
            ID:          "my_menu_item",
            Title:       "My Menu Item",
            Description: "Custom menu action",
            Handler: func() tea.Cmd {
                // Handle menu item selection
                return nil
            },
            Enabled:     true,
        },
    }
}
```

## Installation

### 1. Plugin Directories

Plugins are loaded from these directories (in order of precedence):

- `~/.noise.sh/plugins/` (User plugins)
- `./plugins/` (Project plugins)
- `./internal/plugins/examples/` (Built-in examples)

### 2. Manual Installation

1. Create your plugin files (`.go` and `.json`)
2. Copy them to one of the plugin directories
3. Restart noise.sh

### 3. Development Installation

For development, place your plugin in `./internal/plugins/examples/`:

```bash
cp my_plugin.go ./internal/plugins/examples/
cp my_plugin.json ./internal/plugins/examples/
```

## Security

The plugin system includes multiple security layers:

### Path Validation
- Plugins can only be loaded from approved directories
- System directories are blocked
- Path traversal attacks are prevented

### File Validation
- File size limits (10MB max)
- Extension whitelisting
- Permission checks (no world-writable files)

### Manifest Validation
- Required fields validation
- Capability validation
- Suspicious pattern detection

### Sandboxing
- Plugins run in restricted environment
- Resource usage monitoring
- Integrity checking

## Management

### Plugin Settings Interface

Access the plugin management interface through the application settings. The interface provides:

- **Plugin List**: View all loaded plugins
- **Enable/Disable**: Toggle plugin state
- **Details View**: Plugin information and capabilities
- **Security Report**: Plugin security status

### Command Line

```bash
# List plugins
noise.sh --list-plugins

# Enable plugin
noise.sh --enable-plugin plugin_id

# Disable plugin
noise.sh --disable-plugin plugin_id

# Plugin info
noise.sh --plugin-info plugin_id
```

## Best Practices

### 1. Plugin Development

- **Keep it simple**: Start with minimal functionality
- **Handle errors gracefully**: Always return meaningful errors
- **Document capabilities**: Clearly describe what your plugin does
- **Version carefully**: Use semantic versioning

### 2. Security

- **Validate inputs**: Never trust user data
- **Limit resource usage**: Be conscious of memory and CPU usage
- **Log appropriately**: Use the provided logger
- **Test thoroughly**: Test with various inputs and edge cases

### 3. Performance

- **Lazy initialization**: Initialize resources only when needed
- **Efficient algorithms**: Optimize for large content
- **Background processing**: Use goroutines for long operations
- **Memory management**: Clean up resources in `Cleanup()`

## Examples

### Example 1: Simple Export Plugin

See `internal/plugins/examples/markdown_export.go` for a complete example that adds Markdown export functionality.

### Example 2: Theory Tool Plugin

See `internal/plugins/examples/chord_analyzer.go` for an example that provides chord analysis capabilities.

## Troubleshooting

### Common Issues

1. **Plugin not loading**
   - Check file permissions
   - Verify manifest JSON syntax
   - Check noise.sh logs for errors

2. **Plugin not appearing in UI**
   - Ensure plugin is enabled
   - Check capability registration
   - Verify UI integration points

3. **Security validation errors**
   - Check plugin directory permissions
   - Verify manifest format
   - Review security logs

### Debug Mode

Enable debug logging to troubleshoot plugin issues:

```bash
noise.sh --debug
```

Check the log file at `~/.noise.sh/logs/noise.sh.log` for detailed plugin loading information.

## API Reference

### Plugin Interface

```go
type Plugin interface {
    Metadata() *PluginMetadata
    Initialize(ctx *PluginContext) error
    Cleanup() error
    Enable() error
    Disable() error
    IsEnabled() bool
}
```

### Plugin Context

```go
type PluginContext struct {
    Context  context.Context
    Config   *config.Config
    Metadata *PluginMetadata
    // Additional services will be added as needed
}
```

### Capabilities

- `CapabilityMenuItem`: Add menu items
- `CapabilityScreen`: Provide new screens
- `CapabilityUIExtension`: Extend existing UI
- `CapabilityEditorTool`: Add editor tools
- `CapabilityExportFormat`: Add export formats
- `CapabilitySyntax`: Provide syntax highlighting
- `CapabilityTheoryTool`: Add theory analysis tools
- `CapabilityChordLib`: Provide chord libraries
- `CapabilityAudioEffect`: Provide audio effects
- `CapabilityMIDIHandler`: Handle MIDI events
- `CapabilityDataProvider`: Provide external data
- `CapabilityExporter`: Export data
- `CapabilityHook`: Hook into system events
- `CapabilityConfig`: Add configuration options

## Contributing

When contributing plugins to the noise.sh project:

1. Follow the established patterns in existing plugins
2. Include comprehensive tests
3. Document all public APIs
4. Ensure security best practices
5. Update this documentation

## License

Plugins are licensed separately from the main noise.sh application. Ensure your plugin complies with the chosen license terms.
