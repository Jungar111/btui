package scan

import (
	"btui/internal/bluetooth"
	"btui/internal/ui"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}

func TestDeviceToListItem(t *testing.T) {
	device := bluetooth.BluetoothDevice{
		MacAddress: "AA:BB:CC:DD:EE:FF",
		Name:       "Test Device",
		Connected:  false,
		Paired:     true,
	}

	item := deviceToListItem(device, nil, nil)
	deviceItem, ok := item.(ui.DeviceItem)
	if !ok {
		t.Fatal("Expected item to be ui.DeviceItem")
	}

	if deviceItem.Title() != "Test Device" {
		t.Errorf("Expected title 'Test Device', got %q", deviceItem.Title())
	}

	if deviceItem.Description() != "Paired • AA:BB:CC:DD:EE:FF" {
		t.Errorf("Expected description 'Paired • AA:BB:CC:DD:EE:FF', got %q", deviceItem.Description())
	}
}

func TestDiscoveredDeviceToListItem(t *testing.T) {
	discovered := bluetooth.DiscoveredDevice{
		BluetoothDevice: bluetooth.BluetoothDevice{
			MacAddress: "BB:CC:DD:EE:FF:AA",
			Name:       "Discovered Device",
			Connected:  false,
			Paired:     false,
		},
		RSSI: -72,
	}

	item := discoveredDeviceToListItem(discovered, nil, nil)
	deviceItem, ok := item.(ui.DeviceItem)
	if !ok {
		t.Fatal("Expected item to be ui.DeviceItem")
	}

	if deviceItem.Title() != "Discovered Device" {
		t.Errorf("Expected title 'Discovered Device', got %q", deviceItem.Title())
	}

	expectedDesc := "Discovered • RSSI: -72 • BB:CC:DD:EE:FF:AA"
	if deviceItem.Description() != expectedDesc {
		t.Errorf("Expected description %q, got %q", expectedDesc, deviceItem.Description())
	}
}

func TestCombineDevicesToListItems(t *testing.T) {
	pairedDevices := []bluetooth.BluetoothDevice{
		{
			MacAddress: "AA:BB:CC:DD:EE:FF",
			Name:       "Paired Device",
			Connected:  true,
			Paired:     true,
		},
	}

	discoveredDevices := []bluetooth.DiscoveredDevice{
		{
			BluetoothDevice: bluetooth.BluetoothDevice{
				MacAddress: "BB:CC:DD:EE:FF:AA",
				Name:       "Discovered Device",
				Connected:  false,
				Paired:     false,
			},
			RSSI: -65,
		},
		{
			// Duplicate MAC - should be ignored in favor of paired device
			BluetoothDevice: bluetooth.BluetoothDevice{
				MacAddress: "AA:BB:CC:DD:EE:FF",
				Name:       "Should be ignored",
				Connected:  false,
				Paired:     false,
			},
			RSSI: -80,
		},
	}

	items := combineDevicesToListItems(pairedDevices, discoveredDevices, nil, nil)

	// Should have 2 items: 1 paired, 1 discovered (duplicate ignored)
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestConnectingDisconnectingStatus(t *testing.T) {
	device := bluetooth.BluetoothDevice{
		MacAddress: "AA:BB:CC:DD:EE:FF",
		Name:       "Test Device",
		Connected:  false,
		Paired:     true,
	}

	connectingDevice := &bluetooth.BluetoothDevice{
		MacAddress: "AA:BB:CC:DD:EE:FF",
		Name:       "Test Device",
	}

	disconnectingDevice := &bluetooth.BluetoothDevice{
		MacAddress: "AA:BB:CC:DD:EE:FF", 
		Name:       "Test Device",
	}

	// Test connecting status
	item := deviceToListItem(device, connectingDevice, nil)
	deviceItem := item.(ui.DeviceItem)
	if !strings.Contains(deviceItem.Description(), "Connecting...") {
		t.Errorf("Expected connecting status, got: %s", deviceItem.Description())
	}

	// Test disconnecting status  
	item = deviceToListItem(device, nil, disconnectingDevice)
	deviceItem = item.(ui.DeviceItem)
	if !strings.Contains(deviceItem.Description(), "Disconnecting...") {
		t.Errorf("Expected disconnecting status, got: %s", deviceItem.Description())
	}

	// Test normal status when no operations in progress
	item = deviceToListItem(device, nil, nil)
	deviceItem = item.(ui.DeviceItem)
	if !strings.Contains(deviceItem.Description(), "Paired") {
		t.Errorf("Expected paired status, got: %s", deviceItem.Description())
	}
}

func TestUIUpdateHandling(t *testing.T) {
	model := NewModel()
	
	// Set up connecting state
	device := bluetooth.BluetoothDevice{
		MacAddress: "AA:BB:CC:DD:EE:FF",
		Name:       "Test Device",
	}
	model.ConnectingTo = &device
	
	// Simulate UI update message
	msg := bluetooth.UIUpdateMsg{}
	newModelInterface, cmd := model.Update(msg)
	
	// Should continue periodic updates when operation is in progress
	if cmd == nil {
		t.Error("Expected UIUpdateCmd to be returned when operation is in progress")
	}
	
	// Cast back to Model type and clear connecting state
	newModel := newModelInterface.(Model)
	newModel.ConnectingTo = nil
	newModel.DisconnectingFrom = nil
	
	// Simulate another UI update message
	finalModel, finalCmd := newModel.Update(msg)
	
	// Should stop periodic updates when no operations are in progress
	if finalCmd != nil {
		t.Error("Expected no command when no operations are in progress")
	}
	
	_ = finalModel // Use the variable to avoid compiler warning
}

func TestUpdateWindowSize(t *testing.T) {
	model := NewModel()

	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(windowMsg)
	m := updatedModel.(Model)

	if m.Width != 100 || m.Height != 50 {
		t.Errorf("Expected dimensions 100x50, got %dx%d", m.Width, m.Height)
	}
}

func TestUpdateDevicesMsg(t *testing.T) {
	model := NewModel()

	devicesMsg := bluetooth.DevicesMsg{
		Devices: []string{
			"Device AA:BB:CC:DD:EE:FF Test Device",
		},
		ConnectedDevices: []string{
			"Device AA:BB:CC:DD:EE:FF Test Device",
		},
		Err: nil,
	}

	updatedModel, _ := model.Update(devicesMsg)
	m := updatedModel.(Model)

	if m.Loading {
		t.Error("Expected loading to be false after DevicesMsg")
	}

	if len(m.PairedDevices) != 1 {
		t.Errorf("Expected 1 paired device, got %d", len(m.PairedDevices))
	}

	if m.PairedDevices[0].MacAddress != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("Expected MAC address AA:BB:CC:DD:EE:FF, got %s", m.PairedDevices[0].MacAddress)
	}
}

func TestUpdateDiscoveryUpdateMsg(t *testing.T) {
	model := NewModel()
	// Initialize the list first
	model.List = ui.NewList([]list.Item{}, "Test", 80, 10)

	discoveryMsg := bluetooth.DiscoveryUpdateMsg{
		Devices: []bluetooth.DiscoveredDevice{
			{
				BluetoothDevice: bluetooth.BluetoothDevice{
					MacAddress: "CC:DD:EE:FF:AA:BB",
					Name:       "New Device",
				},
				RSSI: -70,
			},
		},
		Err: nil,
	}

	updatedModel, cmd := model.Update(discoveryMsg)
	m := updatedModel.(Model)

	if len(m.DiscoveredDevices) != 1 {
		t.Errorf("Expected 1 discovered device, got %d", len(m.DiscoveredDevices))
	}

	// Should not return a command if not actively scanning
	if cmd != nil && m.ScanState != ScanActive {
		t.Error("Should not return command when not actively scanning")
	}
}
