package listdevices

import (
	"btui/internal/ui"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}

func TestParseDeviceLine(t *testing.T) {
	connectedMacs := map[string]bool{
		"AA:BB:CC:DD:EE:FF": true,
	}

	tests := []struct {
		name              string
		line              string
		expectedMAC       string
		expectedName      string
		expectedConnected bool
	}{
		{
			name:              "Basic device with name",
			line:              "Device AA:BB:CC:DD:EE:FF Test Device",
			expectedMAC:       "AA:BB:CC:DD:EE:FF",
			expectedName:      "Test Device",
			expectedConnected: true,
		},
		{
			name:              "Device without name",
			line:              "Device BB:CC:DD:EE:FF:AA",
			expectedMAC:       "BB:CC:DD:EE:FF:AA",
			expectedName:      "",
			expectedConnected: false,
		},
		{
			name:              "Invalid line",
			line:              "Invalid",
			expectedMAC:       "",
			expectedName:      "",
			expectedConnected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			device := parseDeviceLine(test.line, connectedMacs)

			if device.MacAddress != test.expectedMAC {
				t.Errorf("Expected MAC %s, got %s", test.expectedMAC, device.MacAddress)
			}

			if device.Name != test.expectedName {
				t.Errorf("Expected name %s, got %s", test.expectedName, device.Name)
			}

			if device.Connected != test.expectedConnected {
				t.Errorf("Expected connected %v, got %v", test.expectedConnected, device.Connected)
			}

			if device.RawLine != test.line {
				t.Errorf("Expected raw line %s, got %s", test.line, device.RawLine)
			}
		})
	}
}

func TestDeviceToListItem(t *testing.T) {
	device := BluetoothDevice{
		MacAddress: "AA:BB:CC:DD:EE:FF",
		Name:       "Test Device",
		Connected:  true,
		RawLine:    "Device AA:BB:CC:DD:EE:FF Test Device",
	}

	item := deviceToListItem(device)
	genericItem, ok := item.(ui.GenericItem)
	if !ok {
		t.Fatal("Expected item to be ui.GenericItem")
	}

	if genericItem.Title != "🔗 Test Device" {
		t.Errorf("Expected title '🔗 Test Device', got %q", genericItem.Title)
	}

	if genericItem.Description != "AA:BB:CC:DD:EE:FF (Connected)" {
		t.Errorf("Expected description 'AA:BB:CC:DD:EE:FF (Connected)', got %q", genericItem.Description)
	}

	// Test device value
	deviceValue, ok := genericItem.Value.(BluetoothDevice)
	if !ok {
		t.Fatal("Expected value to be BluetoothDevice")
	}

	if deviceValue.MacAddress != device.MacAddress {
		t.Errorf("Expected MAC %s, got %s", device.MacAddress, deviceValue.MacAddress)
	}
}

func TestDevicesToListItems(t *testing.T) {
	deviceLines := []string{
		"Device AA:BB:CC:DD:EE:FF Connected Device",
		"Device BB:CC:DD:EE:FF:AA Disconnected Device",
	}

	connectedDeviceLines := []string{
		"Device AA:BB:CC:DD:EE:FF Connected Device",
	}

	items := devicesToListItems(deviceLines, connectedDeviceLines)

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
		return
	}

	// Check first item (connected)
	item1, ok := items[0].(ui.GenericItem)
	if !ok {
		t.Fatal("Expected first item to be ui.GenericItem")
	}

	if !strings.Contains(item1.Title, "🔗") {
		t.Error("Expected first item to have connected indicator")
	}

	// Check second item (not connected)
	item2, ok := items[1].(ui.GenericItem)
	if !ok {
		t.Fatal("Expected second item to be ui.GenericItem")
	}

	if strings.Contains(item2.Title, "🔗") {
		t.Error("Expected second item to not have connected indicator")
	}
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

	devicesMsg := DevicesMsg{
		Devices: []string{
			"Device AA:BB:CC:DD:EE:FF Test Device",
		},
		ConnectedDevices: []string{},
		Err:              nil,
	}

	updatedModel, _ := model.Update(devicesMsg)
	m := updatedModel.(Model)

	if m.Loading {
		t.Error("Expected loading to be false after DevicesMsg")
	}

	if m.List.Items() == nil {
		t.Error("Expected list items to be initialized")
	}
}

func TestUpdateDevicesMsgWithError(t *testing.T) {
	model := NewModel()

	devicesMsg := DevicesMsg{
		Devices:          nil,
		ConnectedDevices: nil,
		Err:              fmt.Errorf("test error"),
	}

	updatedModel, _ := model.Update(devicesMsg)
	m := updatedModel.(Model)

	if m.Loading {
		t.Error("Expected loading to be false after DevicesMsg with error")
	}

	if m.Err == nil {
		t.Error("Expected error to be set")
	}
}

func TestUpdateQuitKey(t *testing.T) {
	model := NewModel()

	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := model.Update(quitMsg)
	m := updatedModel.(Model)

	if !m.Quitting {
		t.Error("Expected quitting to be true after 'q' key")
	}

	if cmd == nil {
		t.Error("Expected quit command")
	}
}
