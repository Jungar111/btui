package connect

import (
	"btui/internal/bluetooth"
)

// ViewState represents the current view state
type ViewState int

const (
	DeviceSelection ViewState = iota
	Connecting
	ShowResult
)

// Model represents the state of the connect command
type Model struct {
	State        ViewState
	DevicePicker bluetooth.PickerModel
	Result       *bluetooth.ConnectResult
	Width        int
	Height       int
}

// NewModel creates a new model for the connect command
func NewModel() Model {
	return Model{
		State:        DeviceSelection,
		DevicePicker: bluetooth.NewPickerModel(),
	}
}
