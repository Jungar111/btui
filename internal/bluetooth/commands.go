package bluetooth

import (
	"os/exec"
	"strings"

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
		cmd := exec.Command("bluetoothctl", "connect", device.MacAddress)
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
		cmd := exec.Command("bluetoothctl", "disconnect", device.MacAddress)
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