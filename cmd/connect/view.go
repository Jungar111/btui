package connect

import (
	"fmt"

	"btui/internal/ui"
)

// View implements tea.Model
func (m Model) View() string {
	switch m.State {
	case DeviceSelection:
		return m.DevicePicker.View()
	case Connecting:
		return m.viewConnecting()
	case ShowResult:
		return m.viewResult()
	default:
		return "Unknown state"
	}
}

// viewConnecting renders the connecting state
func (m Model) viewConnecting() string {
	if m.DevicePicker.Choice == nil {
		return "Connecting...\n\nPress Ctrl+C to cancel."
	}

	deviceName := m.DevicePicker.Choice.Name
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	return fmt.Sprintf("Connecting to %s (%s)...\n\nPress Ctrl+C to cancel.",
		deviceName,
		m.DevicePicker.Choice.MacAddress)
}

// viewResult renders the connection result
func (m Model) viewResult() string {
	if m.Result == nil {
		return "No result to display.\n\nPress any key to exit."
	}

	deviceName := m.Result.Device.Name
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	style := ui.SuccessStyle()
	message := "✓ Successfully connected"

	if !m.Result.Success {
		style = ui.ErrorStyle()
		message = "✗ Failed to connect"
	}

	header := style.Render(fmt.Sprintf("%s to %s", message, deviceName))

	output := ""
	if m.Result.Output != "" {
		output = fmt.Sprintf("\nOutput: %s", m.Result.Output)
	}

	if m.Result.Err != nil {
		output += fmt.Sprintf("\nError: %s", m.Result.Err.Error())
	}

	return fmt.Sprintf("%s\nMAC: %s%s\n\nPress any key to exit.",
		header,
		m.Result.Device.MacAddress,
		output)
}
