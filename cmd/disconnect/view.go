package disconnect

import (
	"fmt"

	"btui/internal/ui"
)

// View implements tea.Model
func (m Model) View() string {
	switch m.State {
	case DeviceSelection:
		return m.DevicePicker.View()
	case Disconnecting:
		return m.viewDisconnecting()
	case ShowResult:
		return m.viewResult()
	default:
		return "Unknown state"
	}
}

// viewDisconnecting renders the disconnecting state
func (m Model) viewDisconnecting() string {
	if m.DevicePicker.Choice == nil {
		return "Disconnecting...\n\nPress Ctrl+C to cancel."
	}

	deviceName := m.DevicePicker.Choice.Name
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	return fmt.Sprintf("Disconnecting from %s (%s)...\n\nPress Ctrl+C to cancel.",
		deviceName,
		m.DevicePicker.Choice.MacAddress)
}

// viewResult renders the disconnection result
func (m Model) viewResult() string {
	if m.Result == nil {
		return "No result to display.\n\nPress any key to exit."
	}

	deviceName := m.Result.Device.Name
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	style := ui.SuccessStyle()
	message := "✓ Successfully disconnected"
	
	if !m.Result.Success {
		style = ui.ErrorStyle()
		message = "✗ Failed to disconnect"
	}

	header := style.Render(fmt.Sprintf("%s from %s", message, deviceName))
	
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