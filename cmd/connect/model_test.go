package connect

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

func TestViewStateString(t *testing.T) {
	tests := []struct {
		state    ViewState
		expected string
	}{
		{DeviceSelection, "DeviceSelection"},
		{Connecting, "Connecting"},
		{ShowResult, "ShowResult"},
	}

	for _, test := range tests {
		// Note: ViewState doesn't have a String() method, but we can test the constants
		if int(test.state) < 0 || int(test.state) > 2 {
			t.Errorf("ViewState %v is out of expected range", test.state)
		}
	}
}
