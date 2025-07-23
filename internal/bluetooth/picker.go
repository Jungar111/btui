package bluetooth

import (
	"btui/internal/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// PickerModel represents a device picker interface
type PickerModel struct {
	List     list.Model
	Choice   *BluetoothDevice
	Quitting bool
	Loading  bool
	Err      error
	Width    int
	Height   int
}

// NewPickerModel creates a new device picker model
func NewPickerModel() PickerModel {
	return PickerModel{
		Loading: true,
	}
}

// DeviceToListItem converts a BluetoothDevice to a generic list item
func DeviceToListItem(d BluetoothDevice) list.Item {
	title := d.Name
	if title == "" {
		title = "Unknown Device"
	}

	// Add connection status indicator to title
	if d.Connected {
		title = "🔗 " + title
	}

	description := d.MacAddress
	if d.Connected {
		description += " (Connected)"
	}

	return ui.GenericItem{
		Title:       title,
		Description: description,
		Value:       d, // Store the full device info
	}
}

// DevicesToListItems converts a slice of BluetoothDevice to list items
func DevicesToListItems(devices []BluetoothDevice) []list.Item {
	items := make([]list.Item, len(devices))
	for i, device := range devices {
		items[i] = DeviceToListItem(device)
	}
	return items
}

// Init implements tea.Model
func (m PickerModel) Init() tea.Cmd {
	return FetchDevicesCmd()
}

// Update implements tea.Model
func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if m.List.Items() != nil {
			m.List.SetWidth(msg.Width)
		}
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			if !m.Loading && len(m.List.Items()) > 0 {
				selectedItem := m.List.SelectedItem()
				if genericItem, ok := selectedItem.(ui.GenericItem); ok {
					if device, ok := genericItem.Value.(BluetoothDevice); ok {
						m.Choice = &device
					}
				}
			}
			return m, tea.Quit
		}

	case DevicesMsg:
		m.Loading = false
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}

		// Parse devices and convert to list items
		devices := ParseDevices(msg.Devices, msg.ConnectedDevices)
		items := DevicesToListItems(devices)

		// Create the list with stored dimensions
		width := m.Width
		height := m.Height
		if width == 0 {
			width = 80
		}
		if height == 0 {
			height = 14
		}

		m.List = ui.NewList(items, "Select Bluetooth Device", width, height)
		return m, nil
	}

	// Only update the list if it's been initialized
	if m.List.Items() != nil {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model
func (m PickerModel) View() string {
	if m.Quitting {
		return ""
	}

	if m.Err != nil {
		return "Error: " + m.Err.Error() + "\n\nPress q to quit."
	}

	if m.Loading {
		return "Scanning for Bluetooth devices...\n\nPress q to quit."
	}

	if m.List.Items() == nil || len(m.List.Items()) == 0 {
		return "No Bluetooth devices found.\n\nPress q to quit."
	}

	return m.List.View()
}