// Package listdevices contains the view logic for the list devices command
package listdevices

import "fmt"

// View implements tea.Model
func (m Model) View() string {
	if m.Choice != nil {
		status := "Disconnected"
		if m.Choice.Connected {
			status = "Connected"
		}
		return fmt.Sprintf("\nSelected device: %s (%s) - %s\n", m.Choice.Name, m.Choice.MacAddress, status)
	}

	if m.Quitting {
		return "\nGoodbye!\n"
	}

	if m.Loading {
		return "\nLoading devices...\n"
	}

	if m.Err != nil {
		return fmt.Sprintf("\nError: %v\n", m.Err)
	}

	if len(m.List.Items()) == 0 {
		return "\nNo devices found.\n"
	}

	return "\n" + m.List.View()
}
