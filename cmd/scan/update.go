package scan

import (
	"btui/internal/bluetooth"
	"btui/internal/ui"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// deviceToListItem converts a BluetoothDevice to a device list item
func deviceToListItem(d bluetooth.BluetoothDevice, connectingTo *bluetooth.BluetoothDevice, disconnectingFrom *bluetooth.BluetoothDevice) list.Item {
	title := d.Name
	if title == "" {
		title = "Unknown Device"
	}

	// Create status description with device status and MAC address, applying colors
	// Priority: Connecting/Disconnecting > Connected > Paired > Discovered
	var status string
	if connectingTo != nil && connectingTo.MacAddress == d.MacAddress {
		status = ui.ConnectingStatusStyle.Render("Connecting...")
	} else if disconnectingFrom != nil && disconnectingFrom.MacAddress == d.MacAddress {
		status = ui.DisconnectingStatusStyle.Render("Disconnecting...")
	} else if d.Connected {
		status = ui.ConnectedStatusStyle.Render("Connected")
	} else if d.Paired {
		status = ui.PairedStatusStyle.Render("Paired")
	} else {
		status = ui.DiscoveredStatusStyle.Render("Discovered")
	}

	// Description includes colored status and muted MAC address
	description := status + " • " + ui.MacAddressStyle.Render(d.MacAddress)

	return ui.NewDeviceItem(title, description, d)
}

// discoveredDeviceToListItem converts a DiscoveredDevice to a device list item
func discoveredDeviceToListItem(d bluetooth.DiscoveredDevice, connectingTo *bluetooth.BluetoothDevice, disconnectingFrom *bluetooth.BluetoothDevice) list.Item {
	title := d.Name
	if title == "" {
		title = "Unknown Device"
	}

	// Create description with colored status, RSSI, and MAC address
	// Priority: Connecting/Disconnecting > Discovered
	var status string
	if connectingTo != nil && connectingTo.MacAddress == d.MacAddress {
		status = ui.ConnectingStatusStyle.Render("Connecting...")
	} else if disconnectingFrom != nil && disconnectingFrom.MacAddress == d.MacAddress {
		status = ui.DisconnectingStatusStyle.Render("Disconnecting...")
	} else {
		status = ui.DiscoveredStatusStyle.Render("Discovered")
	}

	description := status
	if d.RSSI != 0 {
		description += " • " + ui.RSSIStyle.Render(fmt.Sprintf("RSSI: %d", d.RSSI))
	}
	description += " • " + ui.MacAddressStyle.Render(d.MacAddress)

	return ui.NewDeviceItem(title, description, d.BluetoothDevice)
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
		items[i] = deviceToListItem(device, nil, nil)
	}
	return items
}

// combineDevicesToListItems combines paired and discovered devices into list items
func combineDevicesToListItems(pairedDevices []bluetooth.BluetoothDevice, discoveredDevices []bluetooth.DiscoveredDevice, connectingTo *bluetooth.BluetoothDevice, disconnectingFrom *bluetooth.BluetoothDevice) []list.Item {
	// Create a map to avoid duplicates (prioritize paired devices)
	deviceMap := make(map[string]list.Item)

	// Add paired devices first (they take priority)
	for _, device := range pairedDevices {
		deviceMap[device.MacAddress] = deviceToListItem(device, connectingTo, disconnectingFrom)
	}

	// Add discovered devices (only if not already paired)
	for _, device := range discoveredDevices {
		if _, exists := deviceMap[device.MacAddress]; !exists {
			deviceMap[device.MacAddress] = discoveredDeviceToListItem(device, connectingTo, disconnectingFrom)
		}
	}

	// Convert map to slice and separate by connection status: connected -> paired -> discovered
	var connectedItems []list.Item
	var pairedItems []list.Item
	var discoveredItems []list.Item

	for _, item := range deviceMap {
		if deviceItem, ok := item.(ui.DeviceItem); ok {
			if device, ok := deviceItem.Device().(bluetooth.BluetoothDevice); ok {
				if device.Connected {
					connectedItems = append(connectedItems, item)
				} else if device.Paired {
					pairedItems = append(pairedItems, item)
				} else {
					discoveredItems = append(discoveredItems, item)
				}
			}
		}
	}

	// Sort connected devices alphabetically by name
	sort.Slice(connectedItems, func(i, j int) bool {
		itemI := connectedItems[i].(ui.DeviceItem)
		itemJ := connectedItems[j].(ui.DeviceItem)
		deviceI := itemI.Device().(bluetooth.BluetoothDevice)
		deviceJ := itemJ.Device().(bluetooth.BluetoothDevice)
		return strings.ToLower(deviceI.Name) < strings.ToLower(deviceJ.Name)
	})

	// Sort paired devices alphabetically by name
	sort.Slice(pairedItems, func(i, j int) bool {
		itemI := pairedItems[i].(ui.DeviceItem)
		itemJ := pairedItems[j].(ui.DeviceItem)
		deviceI := itemI.Device().(bluetooth.BluetoothDevice)
		deviceJ := itemJ.Device().(bluetooth.BluetoothDevice)
		return strings.ToLower(deviceI.Name) < strings.ToLower(deviceJ.Name)
	})

	// Sort discovered devices alphabetically by name
	sort.Slice(discoveredItems, func(i, j int) bool {
		itemI := discoveredItems[i].(ui.DeviceItem)
		itemJ := discoveredItems[j].(ui.DeviceItem)
		deviceI := itemI.Device().(bluetooth.BluetoothDevice)
		deviceJ := itemJ.Device().(bluetooth.BluetoothDevice)
		return strings.ToLower(deviceI.Name) < strings.ToLower(deviceJ.Name)
	})

	// Combine: connected first, then paired, then discovered
	items := make([]list.Item, 0, len(connectedItems)+len(pairedItems)+len(discoveredItems))
	items = append(items, connectedItems...)
	items = append(items, pairedItems...)
	items = append(items, discoveredItems...)

	return items
}

// updateDeviceList updates or creates the device list with proper dimensions and preserves position
func (m *Model) updateDeviceList(items []list.Item) {
	width := m.Width
	height := m.Height
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 14
	}

	// Create clean title without status messages (status messages now shown below list)
	title := "Bluetooth Devices"
	if m.ScanState == ScanActive {
		title += " - Scanning..."
	} else if m.ScanState == ScanStopped {
		title += " - Ready"
	}

	if m.List.Items() == nil {
		// Create new list if it doesn't exist yet
		// Account for padding (2 rows) + status line (1 row) = 3 rows total
		m.List = ui.NewListWithKeys(items, title, width-4, height-3, scanKeys)
	} else {
		// Preserve cursor position during updates
		currentIndex := m.List.Index()
		
		// Update existing list
		m.List.SetItems(items)
		m.List.SetWidth(width - 4)   // Account for padding (4 columns left/right)
		m.List.SetHeight(height - 3) // Account for padding (2 rows) + status line (1 row)
		m.List.Title = title         // Update title with current status
		
		// Restore cursor position if it's still valid
		if currentIndex < len(items) && currentIndex >= 0 {
			m.List.Select(currentIndex)
		}
	}
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
			m.List.SetWidth(msg.Width - 4) // Account for padding
			m.List.SetHeight(msg.Height - 3) // Account for padding (2 rows) + status line (1 row)
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

		case "enter":
			// Smart connect/disconnect: if connected, disconnect; otherwise, connect
			if !m.Loading && len(m.List.Items()) > 0 && m.ConnectingTo == nil && m.DisconnectingFrom == nil {
				selectedItem := m.List.SelectedItem()
				if deviceItem, ok := selectedItem.(ui.DeviceItem); ok {
					if device, ok := deviceItem.Device().(bluetooth.BluetoothDevice); ok {
						if device.Connected {
							// Device is connected, so disconnect it
							m.DisconnectingFrom = &device
							m.StatusMessage = "Disconnecting from " + device.Name + "..."
							// Immediately refresh UI to show disconnecting status
							items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
							m.updateDeviceList(items)
							return m, tea.Batch(
								func() tea.Msg { return DisconnectingMsg{Device: device} },
								bluetooth.DisconnectCmd(device),
								bluetooth.UIUpdateCmd(),
							)
						} else {
							// Device is not connected, so connect to it
							m.ConnectingTo = &device
							m.StatusMessage = "Connecting to " + device.Name + "..."
							// Immediately refresh UI to show connecting status
							items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
							m.updateDeviceList(items)
							return m, tea.Batch(
								func() tea.Msg { return ConnectingMsg{Device: device} },
								bluetooth.ConnectCmd(device),
								bluetooth.UIUpdateCmd(),
							)
						}
					}
				}
			}

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
					// Update title to reflect new state
					if m.List.Items() != nil {
						items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
						m.updateDeviceList(items)
					}
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
				// Update title to reflect new state
				if m.List.Items() != nil {
					items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
					m.updateDeviceList(items)
				}
			}

		case "c":
			// Connect to selected device
			if !m.Loading && len(m.List.Items()) > 0 && m.ConnectingTo == nil {
				selectedItem := m.List.SelectedItem()
				if deviceItem, ok := selectedItem.(ui.DeviceItem); ok {
					if device, ok := deviceItem.Device().(bluetooth.BluetoothDevice); ok {
						if !device.Connected {
							m.ConnectingTo = &device
							m.StatusMessage = "Connecting to " + device.Name + "..."
							// Immediately refresh UI to show connecting status
							items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
							m.updateDeviceList(items)
							return m, tea.Batch(
								func() tea.Msg { return ConnectingMsg{Device: device} },
								bluetooth.ConnectCmd(device),
								bluetooth.UIUpdateCmd(),
							)
						}
					}
				}
			}

		case "d":
			// Disconnect from selected device
			if !m.Loading && len(m.List.Items()) > 0 && m.DisconnectingFrom == nil {
				selectedItem := m.List.SelectedItem()
				if deviceItem, ok := selectedItem.(ui.DeviceItem); ok {
					if device, ok := deviceItem.Device().(bluetooth.BluetoothDevice); ok {
						if device.Connected {
							m.DisconnectingFrom = &device
							m.StatusMessage = "Disconnecting from " + device.Name + "..."
							// Immediately refresh UI to show disconnecting status
							items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
							m.updateDeviceList(items)
							return m, tea.Batch(
								func() tea.Msg { return DisconnectingMsg{Device: device} },
								bluetooth.DisconnectCmd(device),
								bluetooth.UIUpdateCmd(),
							)
						}
					}
				}
			}

		case "r":
			// Refresh device list
			m.Loading = true
			return m, bluetooth.FetchDevicesCmd()

		case "j":
			// Move down in list (vi-style navigation)
			if m.List.Items() != nil {
				var cmd tea.Cmd
				m.List, cmd = m.List.Update(tea.KeyMsg{Type: tea.KeyDown})
				return m, cmd
			}

		case "k":
			// Move up in list (vi-style navigation)
			if m.List.Items() != nil {
				var cmd tea.Cmd
				m.List, cmd = m.List.Update(tea.KeyMsg{Type: tea.KeyUp})
				return m, cmd
			}
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
		items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
		m.updateDeviceList(items)
		return m, nil

	case bluetooth.DiscoveryUpdateMsg:
		if msg.Err != nil {
			m.StatusMessage = "Discovery error: " + msg.Err.Error()
			return m, nil
		}

		// Update discovered devices
		m.DiscoveredDevices = msg.Devices

		// Update the list with combined devices
		items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
		m.updateDeviceList(items)

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

	case bluetooth.UIUpdateMsg:
		// Refresh UI during operations to show connecting/disconnecting status
		if m.ConnectingTo != nil || m.DisconnectingFrom != nil {
			items := combineDevicesToListItems(m.PairedDevices, m.DiscoveredDevices, m.ConnectingTo, m.DisconnectingFrom)
			m.updateDeviceList(items)
			// Continue periodic updates while operations are in progress
			return m, bluetooth.UIUpdateCmd()
		}
		// No operations in progress, stop periodic updates
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
