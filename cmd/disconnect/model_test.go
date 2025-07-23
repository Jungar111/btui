package disconnect

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	model := NewModel()
	
	if model.State != DeviceSelection {
		t.Errorf("Expected initial state to be DeviceSelection, got %v", model.State)
	}
	
	if model.Result != nil {
		t.Error("Expected initial result to be nil")
	}
	
	if model.Width != 0 || model.Height != 0 {
		t.Error("Expected initial dimensions to be 0")
	}
}

func TestViewStateConstants(t *testing.T) {
	// Test that ViewState constants are properly defined
	if DeviceSelection != 0 {
		t.Errorf("Expected DeviceSelection to be 0, got %d", DeviceSelection)
	}
	
	if Disconnecting != 1 {
		t.Errorf("Expected Disconnecting to be 1, got %d", Disconnecting)
	}
	
	if ShowResult != 2 {
		t.Errorf("Expected ShowResult to be 2, got %d", ShowResult)
	}
}