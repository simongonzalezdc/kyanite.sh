package editor

import "github.com/Kyanite/noise/internal/ui/dimension"

// PreviewScroll manages scroll state for preview pane
type PreviewScroll struct {
	offset       int
	totalLines   int
	visibleLines int
	smoothScroll bool
}

// NewPreviewScroll creates a new scroll manager
func NewPreviewScroll() *PreviewScroll {
	return &PreviewScroll{
		offset:       0,
		smoothScroll: true,
	}
}

// ScrollUp scrolls up by the given number of lines
func (s *PreviewScroll) ScrollUp(lines int) {
	if lines <= 0 {
		lines = dimension.ScrollStepSize
	}

	s.offset -= lines
	if s.offset < 0 {
		s.offset = 0
	}
}

// ScrollDown scrolls down by the given number of lines
func (s *PreviewScroll) ScrollDown(lines int) {
	if lines <= 0 {
		lines = dimension.ScrollStepSize
	}

	maxOffset := s.totalLines - s.visibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}

	s.offset += lines
	if s.offset > maxOffset {
		s.offset = maxOffset
	}
}

// PageUp scrolls up by one page
func (s *PreviewScroll) PageUp() {
	s.ScrollUp(s.visibleLines)
}

// PageDown scrolls down by one page
func (s *PreviewScroll) PageDown() {
	s.ScrollDown(s.visibleLines)
}

// ScrollToTop scrolls to the top
func (s *PreviewScroll) ScrollToTop() {
	s.offset = 0
}

// ScrollToBottom scrolls to the bottom
func (s *PreviewScroll) ScrollToBottom() {
	maxOffset := s.totalLines - s.visibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}
	s.offset = maxOffset
}

// ScrollToLine scrolls to make the given line visible
func (s *PreviewScroll) ScrollToLine(lineNum int) {
	if lineNum < s.offset {
		// Line is above viewport
		s.offset = lineNum
	} else if lineNum >= s.offset+s.visibleLines {
		// Line is below viewport
		s.offset = lineNum - s.visibleLines + 1
	}

	// Clamp to valid range
	s.clampOffset()
}

// UpdateDimensions updates the scroll dimensions
func (s *PreviewScroll) UpdateDimensions(totalLines, visibleLines int) {
	s.totalLines = totalLines
	s.visibleLines = visibleLines
	s.clampOffset()
}

// clampOffset ensures offset is within valid range
func (s *PreviewScroll) clampOffset() {
	if s.offset < 0 {
		s.offset = 0
	}

	maxOffset := s.totalLines - s.visibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}

	if s.offset > maxOffset {
		s.offset = maxOffset
	}
}

// Offset returns the current scroll offset
func (s *PreviewScroll) Offset() int {
	return s.offset
}

// CanScrollUp returns whether scrolling up is possible
func (s *PreviewScroll) CanScrollUp() bool {
	return s.offset > 0
}

// CanScrollDown returns whether scrolling down is possible
func (s *PreviewScroll) CanScrollDown() bool {
	maxOffset := s.totalLines - s.visibleLines
	return s.offset < maxOffset && maxOffset >= 0
}

// ScrollPercentage returns the current scroll position as a percentage
func (s *PreviewScroll) ScrollPercentage() float64 {
	if s.totalLines <= s.visibleLines {
		return 0
	}

	maxOffset := s.totalLines - s.visibleLines
	if maxOffset == 0 {
		return 0
	}

	return float64(s.offset) / float64(maxOffset) * 100
}
