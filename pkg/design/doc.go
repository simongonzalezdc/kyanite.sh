// Package design provides a unified design system for the kyanite.sh TUI suite.
//
// It defines themes, style tokens, typography conventions, spacing constants,
// icon systems, and WCAG accessibility validation for terminal user interfaces
// built with the Charm stack (Bubble Tea, Lipgloss, Bubbles).
//
// # Theme Persistence
//
// Theme persistence (saving/loading the user's selected theme preference) is
// managed per-app, not by this shared module. Each application is responsible
// for persisting the user's theme choice to its own configuration file
// (e.g., ~/.config/<app>/theme.json or theme.toml).
//
// # Registration
//
// Built-in themes are registered via RegisterBuiltIn, which panics if the
// theme fails WCAG AA contrast validation. User-defined custom themes should
// use RegisterCustom, which returns an error instead of panicking.
package design
