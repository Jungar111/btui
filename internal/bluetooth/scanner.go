package bluetooth

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// FetchDevicesCmd returns a command that scans for Bluetooth devices
func FetchDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		// Create a context with timeout to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Fetch all devices
		allDevicesCmd := exec.CommandContext(ctx, "bluetoothctl", "devices")
		allDevicesOutput, err := allDevicesCmd.CombinedOutput()
		if err != nil {
			// Include the command output in the error for better debugging
			return DevicesMsg{Err: fmt.Errorf("issue with bluetoothctl devices: %w (output: %s)", err, string(allDevicesOutput))}
		}

		// Split output into lines and filter empty ones
		allDevicesLines := strings.Split(strings.TrimSpace(string(allDevicesOutput)), "\n")
		var devices []string
		for _, line := range allDevicesLines {
			if strings.TrimSpace(line) != "" {
				devices = append(devices, line)
			}
		}

		// Fetch connected devices - if this fails, we'll continue with empty connected list
		connectedDevicesCmd := exec.CommandContext(ctx, "bluetoothctl", "devices", "Connected")
		connectedDevicesOutput, err := connectedDevicesCmd.CombinedOutput()
		var connectedDevices []string
		if err != nil {
			// Log the error but don't fail the entire operation
			// Some systems might not support the "Connected" filter
			connectedDevices = []string{}
		} else {
			// Split output into lines and filter empty ones
			connectedDevicesLines := strings.Split(strings.TrimSpace(string(connectedDevicesOutput)), "\n")
			for _, line := range connectedDevicesLines {
				if strings.TrimSpace(line) != "" {
					connectedDevices = append(connectedDevices, line)
				}
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
