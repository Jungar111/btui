// Package listdevices contains types for the list devices command
package listdevices

// BluetoothDevice represents a Bluetooth device
type BluetoothDevice struct {
	MacAddress string
	Name       string
	RawLine    string
	Connected  bool
}

type DevicesMsg struct {
	Devices          []string
	ConnectedDevices []string
	Err              error
}
