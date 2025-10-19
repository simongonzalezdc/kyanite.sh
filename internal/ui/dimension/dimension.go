// Package dimension provides shared interfaces and helpers for UI components
// that respond to layout dimension changes.
package dimension

// Component defines the common interface that all dimension-aware UI components
// should implement.
type Component interface {
	// SetDimensions updates the component's width and height.
	SetDimensions(width, height int)
	// GetDimensions returns the current width and height.
	GetDimensions() (int, int)
}

// Bounds describes optional minimum and maximum constraints that can be applied
// to width and height values. A zero value for any bound disables that constraint.
type Bounds struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

// Set assigns width and height to the provided destination pointers. A nil
// pointer indicates that the corresponding dimension should be ignored.
func Set(widthPtr, heightPtr *int, width, height int) {
	if widthPtr != nil {
		*widthPtr = width
	}
	if heightPtr != nil {
		*heightPtr = height
	}
}

// ClampWidth ensures that width respects the provided minimum and maximum bounds.
// Use zero to disable either bound.
func ClampWidth(width, min, max int) int {
	if min > 0 && width < min {
		width = min
	}
	if max > 0 && width > max {
		width = max
	}
	return width
}

// ClampHeight ensures that height respects the provided minimum and maximum bounds.
// Use zero to disable either bound.
func ClampHeight(height, min, max int) int {
	if min > 0 && height < min {
		height = min
	}
	if max > 0 && height > max {
		height = max
	}
	return height
}

// ApplyBounds returns width and height values adjusted to respect the supplied bounds.
// Bounds that are zero are ignored.
func ApplyBounds(width, height int, bounds Bounds) (int, int) {
	width = ClampWidth(width, bounds.MinWidth, bounds.MaxWidth)
	height = ClampHeight(height, bounds.MinHeight, bounds.MaxHeight)
	return width, height
}

// MatchDimensions copies the dimensions from the source component to the target.
func MatchDimensions(target Component, source Component) {
	if target == nil || source == nil {
		return
	}
	width, height := source.GetDimensions()
	target.SetDimensions(width, height)
}