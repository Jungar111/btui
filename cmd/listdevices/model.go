// Package listdevices contains the model for the list devices command
package listdevices

import (
	"github.com/charmbracelet/bubbles/list"
)

// Model represents the state of the list devices view
type Model struct {
	List     list.Model
	Choice   *BluetoothDevice
	Quitting bool
	Loading  bool
	Err      error
	Width    int
	Height   int
}

// NewModel creates a new model for the list devices command
func NewModel() Model {
	return Model{
		Loading: true,
	}
}
