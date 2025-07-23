package scan

import (
	"btui/internal/bluetooth"
	"btui/internal/ui"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// deviceToListItem converts a BluetoothDevice to a generic list item
func deviceToListItem(d bluetooth.BluetoothDevice) list.Item {
	title := d.Name
	if title == "" {
		title = "Unknown Device"
	}

	// Add status indicators to title
	if d.Connected {
		title = "🔗 " + title
	} else if d.Paired {
		title = "📱 " + title
	} else {
		title = "🔍 " + title // New/discovered device
	}

	description := d.MacAddress
	if d.Connected {
		description += " (Connected)"
	} else if d.Paired {
		description += " (Paired)"
	} else {
		description += " (Discovered)"
	}

	return ui.GenericItem{
		Title:       title,
		Description: description,
		Value:       d,
	}
}

// discoveredDeviceToListItem converts a DiscoveredDevice to a generic list item
func discoveredDeviceToListItem(d bluetooth.DiscoveredDevice) list.Item {
	title := d.Name
	if title == "" {
		title = "Unknown Device"
	}
	title = "🔍 " + title // New/discovered device indicator

	description := d.MacAddress + " (Discovered)"
	if d.RSSI != 0 {
		description += fmt.Sprintf(" RSSI: %d", d.RSSI)
	}

	return ui.GenericItem{
		Title:       title,
		Description: description,
		Value:       d.BluetoothDevice, // Use the embedded BluetoothDevice
	}
}

// devicesToListItems converts device data to list items
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
		device := bluetooth.ParseDeviceLine(line, connectedMacs)
		// Check if device is paired by checking if it appears in regular device list
		device.Paired = true // All devices from bluetoothctl devices are paired
		items[i] = deviceToListItem(device)
	}
	return items
}

// combineDevicesToListItems combines paired and discovered devices into list items
func combineDevicesToListItems(pairedDevices []bluetooth.BluetoothDevice, discoveredDevices []bluetooth.DiscoveredDevice) []list.Item {
	// Create a map to avoid duplicates (prioritize paired devices)
	deviceMap := make(map[string]list.Item)
	
	// Add paired devices first (they take priority)
	for _, device := range pairedDevices {
		deviceMap[device.MacAddress] = deviceToListItem(device)
	}
	
	// Add discovered devices (only if not already paired)
	for _, device := range discoveredDevices {
		if _, exists := deviceMap[device.MacAddress]; !exists {
			deviceMap[device.MacAddress] = discoveredDeviceToListItem(device)
		}
	}
	
	// Convert map to slice
	items := make([]list.Item, 0, len(deviceMap))
	for _, item := range deviceMap {
		items = append(items, item)
	}
	
	return items
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	// Start by fetching existing devices
	return bluetooth.FetchDevicesCmd()
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
			// Stop discovery scanner if active
			if m.DiscoveryScanner != nil && m.DiscoveryScanner.IsScanning() {
				m.DiscoveryScanner.StopDiscovery()
			}
			return m, tea.Quit

		case "s":
			// Toggle scanning
			switch m.ScanState {
			case ScanStopped:
				m.ScanState = ScanStarting
				m.StatusMessage = "Starting real-time device discovery..."
				if err := m.DiscoveryScanner.StartDiscovery(); err != nil {
					m.StatusMessage = "Failed to start discovery: " + err.Error()
					m.ScanState = ScanStopped
				} else {
					m.ScanState = ScanActive
					m.StatusMessage = "Scanning for devices... Press 's' to stop"
					return m, bluetooth.DiscoveryTickCmd(m.DiscoveryScanner)
				}
			case ScanActive:
				m.ScanState = ScanStopping
				m.StatusMessage = "Stopping discovery..."
				if err := m.DiscoveryScanner.StopDiscovery(); err != nil {
					m.StatusMessage = "Failed to stop discovery: " + err.Error()
				} else {
					m.ScanState = ScanStopped
					m.StatusMessage = "Discovery stopped"
				}
			}

		case "c":
			// Connect to selected device
			if !m.Loading && len(m.List.Items()) > 0 && m.ConnectingTo == nil {
				selectedItem := m.List.SelectedItem()
				if genericItem, ok := selectedItem.(ui.GenericItem); ok {
					if device, ok := genericItem.Value.(bluetooth.BluetoothDevice); ok {
						if !device.Connected {
							m.ConnectingTo = &device
							m.StatusMessage = "Connecting to " + device.Name + "..."
							return m, tea.Batch(
								func() tea.Msg { return ConnectingMsg{Device: device} },
								bluetooth.ConnectCmd(device),
							)
						}
					}
				}
			}

		case "d":
			// Disconnect from selected device
			if !m.Loading && len(m.List.Items()) > 0 && m.DisconnectingFrom == nil {
				selectedItem := m.List.SelectedItem()
				if genericItem, ok := selectedItem.(ui.GenericItem); ok {
					if device, ok := genericItem.Value.(bluetooth.BluetoothDevice); ok {
						if device.Connected {
							m.DisconnectingFrom = &device
							m.StatusMessage = "Disconnecting from " + device.Name + "..."
							return m, tea.Batch(
								func() tea.Msg { return DisconnectingMsg{Device: device} },
								bluetooth.DisconnectCmd(device),
							)
						}
					}
				}
			}

		case "r":
			// Refresh device list
			m.Loading = true
			return m, bluetooth.FetchDevicesCmd()
		}

	case bluetooth.DevicesMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}

		// Parse paired devices
		connectedMacs := make(map[string]bool)
		for _, line := range msg.ConnectedDevices {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				macAddress := parts[1]
				connectedMacs[macAddress] = true
			}
		}

		m.PairedDevices = make([]bluetooth.BluetoothDevice, len(msg.Devices))
		for i, line := range msg.Devices {
			device := bluetooth.ParseDeviceLine(line, connectedMacs)
			device.Paired = true // All devices from bluetoothctl devices are paired
			m.PairedDevices[i] = device
		}

		// Update the list with combined devices
		items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices)

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

	case bluetooth.DiscoveryUpdateMsg:
		if msg.Err != nil {
			m.StatusMessage = "Discovery error: " + msg.Err.Error()
			return m, nil
		}

		// Update discovered devices
		m.DiscoveredDevices = msg.Devices

		// Update the list with combined devices
		items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices)

		// Update the list
		if m.List.Items() != nil {
			m.List.SetItems(items)
		}

		// Continue discovery updates if still scanning
		var cmd tea.Cmd
		if m.ScanState == ScanActive {
			cmd = bluetooth.DiscoveryTickCmd(m.DiscoveryScanner)
		}
		return m, cmd

	case ConnectingMsg:
		m.ConnectingTo = &msg.Device
		return m, nil

	case bluetooth.ConnectMsg:
		m.ConnectingTo = nil
		if msg.Success {
			m.StatusMessage = "Successfully connected to " + msg.Device.Name
		} else {
			m.StatusMessage = "Failed to connect to " + msg.Device.Name + ": " + msg.Output
		}
		// Refresh device list to show updated connection status
		return m, bluetooth.FetchDevicesCmd()

	case DisconnectingMsg:
		m.DisconnectingFrom = &msg.Device
		return m, nil

	case bluetooth.DisconnectMsg:
		m.DisconnectingFrom = nil
		if msg.Success {
			m.StatusMessage = "Successfully disconnected from " + msg.Device.Name
		} else {
			m.StatusMessage = "Failed to disconnect from " + msg.Device.Name + ": " + msg.Output
		}
		// Refresh device list to show updated connection status
		return m, bluetooth.FetchDevicesCmd()
	}

	// Only update the list if it's been initialized
	if m.List.Items() != nil {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}

	return m, nil
}
