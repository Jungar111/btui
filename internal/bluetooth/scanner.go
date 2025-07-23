package bluetooth

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// FetchDevicesCmd returns a command that scans for Bluetooth devices
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

// ParseDeviceLine parses a bluetoothctl device line into a BluetoothDevice
func ParseDeviceLine(line string, connectedMacs map[string]bool) BluetoothDevice {
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

// ParseDevices converts device lines into BluetoothDevice structs
func ParseDevices(deviceLines []string, connectedDeviceLines []string) []BluetoothDevice {
	// Create a map of connected MAC addresses for quick lookup
	connectedMacs := make(map[string]bool)
	for _, line := range connectedDeviceLines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			macAddress := parts[1]
			connectedMacs[macAddress] = true
		}
	}

	devices := make([]BluetoothDevice, len(deviceLines))
	for i, line := range deviceLines {
		devices[i] = ParseDeviceLine(line, connectedMacs)
	}
	return devices
}