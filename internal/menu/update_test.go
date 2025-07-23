package menu

import (
	"btui/cmd/connect"
	"btui/cmd/disconnect"
	"btui/cmd/listdevices"
	"btui/cmd/scan"
	"btui/internal/bluetooth"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMenuInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()
	
	// Menu init should return nil command
	if cmd != nil {
		t.Error("Expected menu Init to return nil command")
	}
}

func TestMenuWindowResize(t *testing.T) {
	model := NewModel()
	
	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(windowMsg)
	m := updatedModel.(Model)
	
	if m.Width != 100 || m.Height != 50 {
		t.Errorf("Expected dimensions 100x50, got %dx%d", m.Width, m.Height)
	}
}

func TestMenuQuitKey(t *testing.T) {
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

func TestHandleActions(t *testing.T) {
	model := NewModel()
	
	tests := []struct {
		action       ActionType
		expectedType interface{}
	}{
		{ListDevicesAction, listdevices.Model{}},
		{ScanAction, scan.Model{}},
		{ConnectAction, connect.Model{}},
		{DisconnectAction, disconnect.Model{}},
	}
	
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			updatedModel, cmd := model.handleAction(test.action)
			m := updatedModel.(Model)
			
			if !m.InSubMenu {
				t.Error("Expected to be in sub-menu after action")
			}
			
			if m.SubProgram == nil {
				t.Error("Expected SubProgram to be set")
			}
			
			if cmd == nil {
				t.Error("Expected command to be returned")
			}
		})
	}
}

func TestHandleQuitAction(t *testing.T) {
	model := NewModel()
	
	updatedModel, cmd := model.handleAction(QuitAction)
	m := updatedModel.(Model)
	
	if !m.Quitting {
		t.Error("Expected quitting to be true after QuitAction")
	}
	
	if cmd == nil {
		t.Error("Expected quit command")
	}
}

func TestShouldReturnToMenu(t *testing.T) {
	tests := []struct {
		name     string
		model    tea.Model
		expected bool
	}{
		{
			name:     "listdevices quitting",
			model:    listdevices.Model{Quitting: true},
			expected: true,
		},
		{
			name:     "listdevices with choice",
			model:    listdevices.Model{Choice: &listdevices.BluetoothDevice{}},
			expected: true,
		},
		{
			name:     "scan quitting",
			model:    scan.Model{Quitting: true},
			expected: true,
		},
		{
			name:     "connect quitting",
			model:    connect.Model{DevicePicker: bluetooth.PickerModel{Quitting: true}},
			expected: true,
		},
		{
			name:     "disconnect quitting",
			model:    disconnect.Model{DevicePicker: bluetooth.PickerModel{Quitting: true}},
			expected: true,
		},
		{
			name:     "listdevices normal state",
			model:    listdevices.Model{Quitting: false, Choice: nil},
			expected: false,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := shouldReturnToMenu(test.model)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMenuEnterKey(t *testing.T) {
	model := NewModel()
	
	// Simulate selecting the first menu item (List Devices)
	// First we need to set up a proper item selection
	items := model.List.Items()
	if len(items) > 0 {
		// Create a mock enter key press
		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
		
		// The exact behavior depends on the list implementation
		// This is a basic test to ensure enter key is handled
		_, cmd := model.Update(enterMsg)
		
		// We expect some response, though the exact command depends on list state
		// This test mainly ensures the key is processed without panic
		_ = cmd // Command might be nil if no item is selected
	}
}

func TestUpdateSubMenu(t *testing.T) {
	model := NewModel()
	model.InSubMenu = true
	model.SubProgram = listdevices.NewModel()
	
	// Test escape key to return to main menu
	escMsg := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel, _ := model.Update(escMsg)
	m := updatedModel.(Model)
	
	if m.InSubMenu {
		t.Error("Expected to exit sub-menu after escape key")
	}
	
	if m.SubProgram != nil {
		t.Error("Expected SubProgram to be nil after escape")
	}
}