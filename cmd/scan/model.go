package scan

import (
	"btui/internal/bluetooth"

	"github.com/charmbracelet/bubbles/list"
)

// Model represents the state of the scan view
type Model struct {
	List              list.Model
	ScanState         ScanState
	Loading           bool
	Quitting          bool
	Err               error
	Width             int
	Height            int
	ConnectingTo      *bluetooth.BluetoothDevice
	DisconnectingFrom *bluetooth.BluetoothDevice
	StatusMessage     string
	DiscoveryScanner  *bluetooth.DiscoveryScanner
	PairedDevices     []bluetooth.BluetoothDevice
	DiscoveredDevices []bluetooth.DiscoveredDevice
}

// NewModel creates a new model for the scan command
func NewModel() Model {
	return Model{
		ScanState:        ScanStopped,
		Loading:          true,
		DiscoveryScanner: bluetooth.NewDiscoveryScanner(),
	}
}
