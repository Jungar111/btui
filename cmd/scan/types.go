package scan

import "btui/internal/bluetooth"

// ScanState represents the current scanning state
type ScanState int

const (
	ScanStopped ScanState = iota
	ScanStarting
	ScanActive
	ScanStopping
)

// String returns a string representation of the scan state
func (s ScanState) String() string {
	switch s {
	case ScanStopped:
		return "Stopped"
	case ScanStarting:
		return "Starting..."
	case ScanActive:
		return "Scanning"
	case ScanStopping:
		return "Stopping..."
	default:
		return "Unknown"
	}
}

// ConnectingMsg is sent when starting to connect to a device
type ConnectingMsg struct {
	Device bluetooth.BluetoothDevice
}

// DisconnectingMsg is sent when starting to disconnect from a device
type DisconnectingMsg struct {
	Device bluetooth.BluetoothDevice
}
