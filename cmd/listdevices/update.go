// Package listdevices contains the update logic for the list devices command
package listdevices

import (
	"btui/internal/ui"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func FetchDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		// Fetch all devices
		allDevicesCmd := exec.Command("bluetoothctl", "devices")
		allDevicesOutput, err := allDevicesCmd.Output()
		if err != nil {
			return DevicesMsg{Err: fmt.Errorf("issue with bluetoothctl devices: %w", err)}
		}

		// Split output into lines and filter empty ones
		allDevicesLines := strings.Split(strings.TrimSpace(string(allDevicesOutput)), "\n")
		var devices []string
		for _, line := range allDevicesLines {
			if strings.TrimSpace(line) != "" {
				devices = append(devices, line)
			}
		}

		// Fetch connected devices
		connectedDevicesCmd := exec.Command("bluetoothctl", "devices", "Connected")
		connectedDevicesOutput, err := connectedDevicesCmd.Output()
		if err != nil {
			return DevicesMsg{Err: fmt.Errorf("issue with bluetoothctl devices Connected: %w", err)}
		}

		// Split output into lines and filter empty ones
		connectedDevicesLines := strings.Split(strings.TrimSpace(string(connectedDevicesOutput)), "\n")
		var connectedDevices []string
		for _, line := range connectedDevicesLines {
			if strings.TrimSpace(line) != "" {
				connectedDevices = append(connectedDevices, line)
			}
		}

		return DevicesMsg{Devices: devices, ConnectedDevices: connectedDevices}
	}
}

// parseDeviceLine parses a bluetoothctl device line into a BluetoothDevice
func parseDeviceLine(line string, connectedMacs map[string]bool) BluetoothDevice {
	// bluetoothctl output format: "Device MAC_ADDRESS NAME"
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return BluetoothDevice{RawLine: line}
	}

	macAddress := parts[1]
	name := ""
	if len(parts) > 2 {
		name = strings.Join(parts[2:], " ")
	}

	return BluetoothDevice{
		MacAddress: macAddress,
		Name:       name,
		RawLine:    line,
		Connected:  connectedMacs[macAddress],
	}
}

// deviceToListItem converts a BluetoothDevice to a generic list item
func deviceToListItem(d BluetoothDevice) list.Item {
	title := d.Name
	if title == "" {
		title = "Unknown Device"
	}

	// Add connection status indicator to title
	if d.Connected {
		title = "🔗 " + title
	}

	description := d.MacAddress
	if d.Connected {
		description += " (Connected)"
	}

	return ui.GenericItem{
		Title:       title,
		Description: description,
		Value:       d, // Store the full device info
	}
}

// devicesToListItems converts a slice of device lines to list items
func devicesToListItems(deviceLines []string, connectedDeviceLines []string) []list.Item {
	// Create a map of connected MAC addresses for quick lookup
	connectedMacs := make(map[string]bool)
	for _, line := range connectedDeviceLines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			macAddress := parts[1]
			connectedMacs[macAddress] = true
		}
	}

	items := make([]list.Item, len(deviceLines))
	for i, line := range deviceLines {
		device := parseDeviceLine(line, connectedMacs)
		items[i] = deviceToListItem(device)
	}
	return items
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return FetchDevicesCmd()
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if m.List.Items() != nil {
			m.List.SetWidth(msg.Width)
		}
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			if !m.Loading && len(m.List.Items()) > 0 {
				selectedItem := m.List.SelectedItem()
				if genericItem, ok := selectedItem.(ui.GenericItem); ok {
					if device, ok := genericItem.Value.(BluetoothDevice); ok {
						m.Choice = &device
					}
				}
			}
			return m, tea.Quit
		}

	case DevicesMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}

		// Convert device strings to list items
		items := devicesToListItems(msg.Devices, msg.ConnectedDevices)

		// Create the list with stored dimensions
		width := m.Width
		height := m.Height
		if width == 0 {
			width = 80
		}
		if height == 0 {
			height = 14
		}

		m.List = ui.NewList(items, "Bluetooth Devices", width, height)
		return m, nil
	}

	// Only update the list if it's been initialized
	if m.List.Items() != nil {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}

	return m, nil
}
