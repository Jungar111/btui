package scan

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	if model.ScanState != ScanStopped {
		t.Errorf("Expected initial scan state to be ScanStopped, got %v", model.ScanState)
	}

	if !model.Loading {
		t.Error("Expected initial loading state to be true")
	}

	if model.DiscoveryScanner == nil {
		t.Error("Expected DiscoveryScanner to be initialized")
	}

	if model.PairedDevices != nil {
		t.Error("Expected initial paired devices to be nil")
	}

	if model.DiscoveredDevices != nil {
		t.Error("Expected initial discovered devices to be nil")
	}
}

func TestScanStateString(t *testing.T) {
	tests := []struct {
		state    ScanState
		expected string
	}{
		{ScanStopped, "Stopped"},
		{ScanStarting, "Starting..."},
		{ScanActive, "Scanning"},
		{ScanStopping, "Stopping..."},
	}

	for _, test := range tests {
		result := test.state.String()
		if result != test.expected {
			t.Errorf("Expected %s.String() to return %q, got %q",
				test.state, test.expected, result)
		}
	}
}

func TestScanStateUnknown(t *testing.T) {
	// Test unknown state
	unknownState := ScanState(999)
	result := unknownState.String()
	if result != "Unknown" {
		t.Errorf("Expected unknown state to return 'Unknown', got %q", result)
	}
}
