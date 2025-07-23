package bluetooth

import (
	"context"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ConnectResult represents the result of a connect operation
type ConnectResult struct {
	Device  BluetoothDevice
	Success bool
	Output  string
	Err     error
}

// ConnectMsg is sent when a connect operation completes
type ConnectMsg ConnectResult

// ConnectCmd returns a command that connects to a Bluetooth device
func ConnectCmd(device BluetoothDevice) tea.Cmd {
	return func() tea.Msg {
		// Create a context with timeout to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "bluetoothctl", "connect", device.MacAddress)
		output, err := cmd.CombinedOutput()

		result := ConnectResult{
			Device: device,
			Output: strings.TrimSpace(string(output)),
			Err:    err,
		}

		// Check if connection was successful
		if err == nil && strings.Contains(result.Output, "Connection successful") {
			result.Success = true
		}

		return ConnectMsg(result)
	}
}

// DisconnectResult represents the result of a disconnect operation
type DisconnectResult struct {
	Device  BluetoothDevice
	Success bool
	Output  string
	Err     error
}

// DisconnectMsg is sent when a disconnect operation completes
type DisconnectMsg DisconnectResult

// DisconnectCmd returns a command that disconnects from a Bluetooth device
func DisconnectCmd(device BluetoothDevice) tea.Cmd {
	return func() tea.Msg {
		// Create a context with timeout to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "bluetoothctl", "disconnect", device.MacAddress)
		output, err := cmd.CombinedOutput()

		result := DisconnectResult{
			Device: device,
			Output: strings.TrimSpace(string(output)),
			Err:    err,
		}

		// Check if disconnection was successful
		if err == nil && strings.Contains(result.Output, "Successful disconnected") {
			result.Success = true
		}

		return DisconnectMsg(result)
	}
}

// ScanResult represents the result of a scan operation
type ScanResult struct {
	Success bool
	Output  string
	Err     error
}

// ScanStartMsg is sent when a scan operation starts
type ScanStartMsg ScanResult

// ScanStopMsg is sent when a scan operation stops
type ScanStopMsg ScanResult

// StartScanCmd returns a command that starts Bluetooth scanning
func StartScanCmd() tea.Cmd {
	return func() tea.Msg {
		// Create a context with timeout to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "bluetoothctl", "scan", "on")
		output, err := cmd.CombinedOutput()

		result := ScanResult{
			Output: strings.TrimSpace(string(output)),
			Err:    err,
		}

		// Check if scan start was successful
		if err == nil && (strings.Contains(result.Output, "Discovery started") || strings.Contains(result.Output, "SetDiscoveryFilter success")) {
			result.Success = true
		}

		return ScanStartMsg(result)
	}
}

// StopScanCmd returns a command that stops Bluetooth scanning
func StopScanCmd() tea.Cmd {
	return func() tea.Msg {
		// Create a context with timeout to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "bluetoothctl", "scan", "off")
		output, err := cmd.CombinedOutput()

		result := ScanResult{
			Output: strings.TrimSpace(string(output)),
			Err:    err,
		}

		// Check if scan stop was successful
		if err == nil && strings.Contains(result.Output, "Discovery stopped") {
			result.Success = true
		}

		return ScanStopMsg(result)
	}
}

// TickMsg is sent periodically during scanning to update device list
type TickMsg time.Time

// TickCmd returns a command that sends periodic ticks for updating device list
func TickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
