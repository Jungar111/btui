package scan

import (
	"btui/internal/bluetooth"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

// scanKeyMap defines the key bindings for the scan interface
type scanKeyMap struct {
	Up          key.Binding
	Down        key.Binding
	ViUp        key.Binding
	ViDown      key.Binding
	Enter       key.Binding
	Scan        key.Binding
	Connect     key.Binding
	Disconnect  key.Binding
	Refresh     key.Binding
	Quit        key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k scanKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Scan, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k scanKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.ViUp, k.Down, k.ViDown}, // navigation
		{k.Enter, k.Connect, k.Disconnect}, // actions
		{k.Scan, k.Refresh, k.Quit}, // controls
	}
}

// scanKeys defines the key map for scan interface
var scanKeys = scanKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "↑"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "↓"),
		key.WithHelp("↓", "move down"),
	),
	ViUp: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "move up"),
	),
	ViDown: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "smart connect/disconnect"),
	),
	Scan: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "toggle scan"),
	),
	Connect: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "connect"),
	),
	Disconnect: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "disconnect"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model represents the state of the scan view
type Model struct {
	List              list.Model
	ScanState         ScanState
	Loading           bool
	Quitting          bool
	Err               error
	Width             int
	Height            int
	ConnectingTo      *bluetooth.BluetoothDevice
	DisconnectingFrom *bluetooth.BluetoothDevice
	StatusMessage     string
	DiscoveryScanner  *bluetooth.DiscoveryScanner
	PairedDevices     []bluetooth.BluetoothDevice
	DiscoveredDevices []bluetooth.DiscoveredDevice
}

// NewModel creates a new model for the scan command
func NewModel() Model {
	return Model{
		ScanState:        ScanStopped,
		Loading:          true,
		DiscoveryScanner: bluetooth.NewDiscoveryScanner(),
	}
}
