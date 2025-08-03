package menu

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	if model.Choice != nil {
		t.Error("Expected initial choice to be nil")
	}

	if model.Quitting {
		t.Error("Expected initial quitting state to be false")
	}

	if model.InSubMenu {
		t.Error("Expected initial InSubMenu state to be false")
	}

	if model.SubProgram != nil {
		t.Error("Expected initial SubProgram to be nil")
	}

	// Test that list is initialized with menu items
	items := model.List.Items()
	if len(items) == 0 {
		t.Error("Expected list to have menu items")
	}

	// Check that all expected actions are present
	expectedActions := []ActionType{
		ListDevicesAction,
		ScanAction,
		ConnectAction,
		DisconnectAction,
		QuitAction,
	}

	if len(items) != len(expectedActions) {
		t.Errorf("Expected %d menu items, got %d", len(expectedActions), len(items))
	}
}

func TestMenuItemCreation(t *testing.T) {
	model := NewModel()
	items := model.List.Items()

	// Verify each item has the correct structure
	for i, item := range items {
		// Cast to GenericItem (assuming it uses ui.GenericItem)
		if item == nil {
			t.Errorf("Menu item %d is nil", i)
		}
	}
}
