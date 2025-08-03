package scan

import (
	"btui/internal/ui"
	"fmt"
)

// View implements tea.Model
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	if m.Err != nil {
		return ui.AppStyle.Render(fmt.Sprintf("Error: %v\n\nPress q to quit", m.Err))
	}

	// Show loading state
	if m.Loading {
		return ui.AppStyle.Render("Loading devices...")
	}

	// Render the list and add status message area below
	if m.List.Items() != nil {
		listView := m.List.View()
		
		// Add status message area below the list (always present to prevent jumping)
		statusLine := ""
		if m.StatusMessage != "" {
			statusLine = m.StatusMessage
		} else {
			statusLine = " " // blank line to maintain consistent spacing
		}
		
		// Combine list view with status message area
		fullView := listView + "\n" + statusLine
		return ui.AppStyle.Render(fullView)
	}

	return ui.AppStyle.Render("No devices found")
}
