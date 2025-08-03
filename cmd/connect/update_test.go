package connect

import (
	"btui/internal/bluetooth"
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

func TestUpdateDeviceSelection(t *testing.T) {
	model := NewModel()

	// Test window resize
	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(windowMsg)
	m := updatedModel.(Model)

	if m.Width != 100 || m.Height != 50 {
		t.Errorf("Expected dimensions 100x50, got %dx%d", m.Width, m.Height)
	}

	// Test quit key
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(quitMsg)

	// Should return tea.Quit command
	if cmd == nil {
		t.Error("Expected quit command when 'q' is pressed")
	}
}

func TestUpdateConnecting(t *testing.T) {
	model := NewModel()
	model.State = Connecting

	// Test ConnectMsg
	testDevice := bluetooth.BluetoothDevice{
		MacAddress: "AA:BB:CC:DD:EE:FF",
		Name:       "Test Device",
	}

	connectResult := bluetooth.ConnectResult{
		Device:  testDevice,
		Success: true,
		Output:  "Connection successful",
	}

	connectMsg := bluetooth.ConnectMsg(connectResult)
	updatedModel, cmd := model.Update(connectMsg)
	m := updatedModel.(Model)

	if m.State != ShowResult {
		t.Errorf("Expected state to be ShowResult, got %v", m.State)
	}

	if m.Result == nil {
		t.Error("Expected result to be set")
	} else if !m.Result.Success {
		t.Error("Expected result to be successful")
	}

	// Should return a tick command for auto-quit
	if cmd == nil {
		t.Error("Expected a command for auto-quit timer")
	}
}

func TestUpdateShowResult(t *testing.T) {
	model := NewModel()
	model.State = ShowResult

	// Test quit key
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(quitMsg)

	if cmd == nil {
		t.Error("Expected quit command when 'q' is pressed in ShowResult state")
	}

	// Test QuitMsg
	quitMessage := tea.QuitMsg{}
	_, cmd = model.Update(quitMessage)

	if cmd == nil {
		t.Error("Expected quit command when QuitMsg is received")
	}
}
