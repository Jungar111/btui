package scan

import (
	"btui/internal/ui"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	if m.Err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit", m.Err)
	}

	var content strings.Builder

	// Header with scan status
	header := fmt.Sprintf("Bluetooth Scanner - Status: %s", m.ScanState.String())
	content.WriteString(ui.TitleStyle.Render(header))
	content.WriteString("\n\n")

	// Show loading state
	if m.Loading {
		content.WriteString("Loading devices...\n")
		return content.String()
	}

	// Show device list if available
	if m.List.Items() != nil {
		content.WriteString(m.List.View())
		content.WriteString("\n")
	}

	// Status message
	if m.StatusMessage != "" {
		statusStyle := lipgloss.NewStyle().
			Italic(true)
		content.WriteString(statusStyle.Render(m.StatusMessage))
		content.WriteString("\n")
	}

	// Connection status indicators
	if m.ConnectingTo != nil {
		content.WriteString(fmt.Sprintf("🔄 Connecting to %s...\n", m.ConnectingTo.Name))
	}
	if m.DisconnectingFrom != nil {
		content.WriteString(fmt.Sprintf("🔄 Disconnecting from %s...\n", m.DisconnectingFrom.Name))
	}

	// Help text
	help := "\nControls:\n"
	help += "  s - Start/Stop scanning\n"
	help += "  c - Connect to selected device\n"
	help += "  d - Disconnect from selected device\n"
	help += "  r - Refresh device list\n"
	help += "  ↑/↓ - Navigate list\n"
	help += "  q - Quit\n"

	content.WriteString(ui.HelpStyle.Render(help))

	return content.String()
}
