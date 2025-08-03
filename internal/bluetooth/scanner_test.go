package bluetooth

import (
	"testing"
	"time"
)

func TestFetchDevicesCmd(t *testing.T) {
	// Test the FetchDevicesCmd function
	cmd := FetchDevicesCmd()
	msg := cmd()

	switch m := msg.(type) {
	case DevicesMsg:
		if m.Err != nil {
			t.Logf("Expected error on systems without Bluetooth: %v", m.Err)
			return // This is expected on systems without Bluetooth
		}

		t.Logf("Found %d devices", len(m.Devices))
		t.Logf("Found %d connected devices", len(m.ConnectedDevices))

		// Basic validation
		for i, device := range m.Devices {
			if device == "" {
				t.Errorf("Device %d is empty", i)
			}
			t.Logf("Device %d: %s", i, device)
		}

		for i, device := range m.ConnectedDevices {
			if device == "" {
				t.Errorf("Connected device %d is empty", i)
			}
			t.Logf("Connected device %d: %s", i, device)
		}

	default:
		t.Errorf("Expected DevicesMsg, got %T", msg)
	}
}

func TestParseDeviceLine(t *testing.T) {
	tests := []struct {
		name              string
		line              string
		connectedMacs     map[string]bool
		expectedMAC       string
		expectedName      string
		expectedConnected bool
	}{
		{
			name:              "Basic device line",
			line:              "Device 4C:87:5D:28:86:DD Bose NC 700 Headphones",
			connectedMacs:     map[string]bool{},
			expectedMAC:       "4C:87:5D:28:86:DD",
			expectedName:      "Bose NC 700 Headphones",
			expectedConnected: false,
		},
		{
			name:              "Connected device",
			line:              "Device DC:2C:26:09:D0:0C Keychron K4",
			connectedMacs:     map[string]bool{"DC:2C:26:09:D0:0C": true},
			expectedMAC:       "DC:2C:26:09:D0:0C",
			expectedName:      "Keychron K4",
			expectedConnected: true,
		},
		{
			name:              "Device without name",
			line:              "Device AA:BB:CC:DD:EE:FF",
			connectedMacs:     map[string]bool{},
			expectedMAC:       "AA:BB:CC:DD:EE:FF",
			expectedName:      "",
			expectedConnected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := ParseDeviceLine(tt.line, tt.connectedMacs)

			if device.MacAddress != tt.expectedMAC {
				t.Errorf("Expected MAC %s, got %s", tt.expectedMAC, device.MacAddress)
			}

			if device.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, device.Name)
			}

			if device.Connected != tt.expectedConnected {
				t.Errorf("Expected connected %v, got %v", tt.expectedConnected, device.Connected)
			}

			if device.RawLine != tt.line {
				t.Errorf("Expected raw line %s, got %s", tt.line, device.RawLine)
			}
		})
	}
}

func TestDiscoveryScanner(t *testing.T) {
	scanner := NewDiscoveryScanner()

	if scanner.IsScanning() {
		t.Error("Scanner should not be scanning initially")
	}

	// Test starting discovery
	err := scanner.StartDiscovery()
	if err != nil {
		t.Logf("Discovery start failed (expected on systems without Bluetooth): %v", err)
		return
	}

	if !scanner.IsScanning() {
		t.Error("Scanner should be scanning after start")
	}

	// Let it run briefly
	time.Sleep(2 * time.Second)

	devices := scanner.GetDiscoveredDevices()
	t.Logf("Discovered %d devices after 2 seconds", len(devices))

	// Test stopping discovery
	err = scanner.StopDiscovery()
	if err != nil {
		t.Errorf("Failed to stop discovery: %v", err)
	}

	if scanner.IsScanning() {
		t.Error("Scanner should not be scanning after stop")
	}
}
