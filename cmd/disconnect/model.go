package disconnect

import (
	"btui/internal/bluetooth"
)

// ViewState represents the current view state
type ViewState int

const (
	DeviceSelection ViewState = iota
	Disconnecting
	ShowResult
)

// Model represents the state of the disconnect command
type Model struct {
	State        ViewState
	DevicePicker bluetooth.PickerModel
	Result       *bluetooth.DisconnectResult
	Width        int
	Height       int
}

// NewModel creates a new model for the disconnect command
func NewModel() Model {
	return Model{
		State:        DeviceSelection,
		DevicePicker: bluetooth.NewPickerModel(),
	}
}