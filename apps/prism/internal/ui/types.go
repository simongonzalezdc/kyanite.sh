package ui

// NavigateMsg is sent to navigate to a different screen
type NavigateMsg struct {
	Screen int
}

// Navigate creates a navigation message
func Navigate(screen int) NavigateMsg {
	return NavigateMsg{Screen: screen}
}
