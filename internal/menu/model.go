package menu

import (
	"btui/internal/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the main menu state
type Model struct {
	List       list.Model
	Choice     *ActionType
	Quitting   bool
	Width      int
	Height     int
	SubProgram tea.Model
	InSubMenu  bool
}

// NewModel creates a new menu model
func NewModel() Model {
	actions := []list.Item{
		ui.GenericItem{
			Title:       "List Devices",
			Description: "Browse available Bluetooth devices",
			Value:       ListDevicesAction,
		},
		ui.GenericItem{
			Title:       "Scan & Connect",
			Description: "Scan for nearby devices and connect",
			Value:       ScanAction,
		},
		ui.GenericItem{
			Title:       "Connect",
			Description: "Connect to a Bluetooth device",
			Value:       ConnectAction,
		},
		ui.GenericItem{
			Title:       "Disconnect",
			Description: "Disconnect from a Bluetooth device",
			Value:       DisconnectAction,
		},
		ui.GenericItem{
			Title:       "Quit",
			Description: "Exit btui",
			Value:       QuitAction,
		},
	}

	return Model{
		List: ui.NewList(actions, "btui - Bluetooth Manager", 80, 14),
	}
}
