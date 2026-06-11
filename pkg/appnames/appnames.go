// Package appnames provides canonical app name constants for the kyanite.sh suite.
// Use these instead of hardcoded strings to avoid drift when apps are added or renamed.
package appnames

const (
	Focus  = "focus"
	Noise  = "noise"
	Syntax = "syntax"
	Prism  = "prism"
)

// All returns every app name in stable order.
var All = []string{Focus, Noise, Syntax, Prism}

// IsValid reports whether the given name is a recognized app.
func IsValid(app string) bool {
	for _, a := range All {
		if a == app {
			return true
		}
	}
	return false
}

// Others returns all app names except the given one.
// If app is not a valid app name, it returns all names.
func Others(app string) []string {
	others := make([]string, 0, len(All)-1)
	for _, a := range All {
		if a != app {
			others = append(others, a)
		}
	}
	return others
}
