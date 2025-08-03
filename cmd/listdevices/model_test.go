package listdevices

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	if !model.Loading {
		t.Error("Expected initial loading state to be true")
	}

	if model.Choice != nil {
		t.Error("Expected initial choice to be nil")
	}

	if model.Quitting {
		t.Error("Expected initial quitting state to be false")
	}

	if model.Err != nil {
		t.Error("Expected initial error to be nil")
	}

	if model.Width != 0 || model.Height != 0 {
		t.Error("Expected initial dimensions to be 0")
	}
}
