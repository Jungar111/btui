package bluetooth

// BluetoothDevice represents a Bluetooth device
type BluetoothDevice struct {
	MacAddress string
	Name       string
	RawLine    string
	Connected  bool
}

// DevicesMsg represents the result of scanning for devices
type DevicesMsg struct {
	Devices          []string
	ConnectedDevices []string
	Err              error
}