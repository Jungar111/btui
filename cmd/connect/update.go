package connect

import (
	"btui/internal/bluetooth"

	tea "github.com/charmbracelet/bubbletea"
)

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return m.DevicePicker.Init()
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case DeviceSelection:
		return m.updateDeviceSelection(msg)
	case Connecting:
		return m.updateConnecting(msg)
	case ShowResult:
		return m.updateShowResult(msg)
	default:
		return m, nil
	}
}

// updateDeviceSelection handles updates during device selection
func (m Model) updateDeviceSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	// Update the device picker
	var cmd tea.Cmd
	updatedModel, updateCmd := m.DevicePicker.Update(msg)
	m.DevicePicker = updatedModel.(bluetooth.PickerModel)
	cmd = updateCmd

	// Check if a device was selected
	if m.DevicePicker.Choice != nil {
		m.State = Connecting
		return m, bluetooth.ConnectCmd(*m.DevicePicker.Choice)
	}

	// Check if user quit device selection
	if m.DevicePicker.Quitting {
		return m, tea.Quit
	}

	return m, cmd
}

// updateConnecting handles updates during connection attempt
func (m Model) updateConnecting(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case bluetooth.ConnectMsg:
		result := bluetooth.ConnectResult(msg)
		m.Result = &result
		m.State = ShowResult
		return m, nil
	}

	return m, nil
}

// updateShowResult handles updates when showing connection result
func (m Model) updateShowResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}