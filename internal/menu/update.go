package menu

import (
	"btui/cmd/connect"
	"btui/cmd/disconnect"
	"btui/cmd/listdevices"
	"btui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model  
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If we're in a sub-menu, handle it there
	if m.InSubMenu && m.SubProgram != nil {
		return m.updateSubMenu(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			selectedItem := m.List.SelectedItem()
			if genericItem, ok := selectedItem.(ui.GenericItem); ok {
				if actionType, ok := genericItem.Value.(ActionType); ok {
					return m.handleAction(actionType)
				}
			}
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

// handleAction processes the selected menu action
func (m Model) handleAction(action ActionType) (tea.Model, tea.Cmd) {
	switch action {
	case ListDevicesAction:
		subModel := listdevices.NewModel()
		m.SubProgram = subModel
		m.InSubMenu = true
		return m, subModel.Init()

	case ConnectAction:
		subModel := connect.NewModel()
		m.SubProgram = subModel
		m.InSubMenu = true
		return m, subModel.Init()

	case DisconnectAction:
		subModel := disconnect.NewModel()
		m.SubProgram = subModel
		m.InSubMenu = true
		return m, subModel.Init()

	case QuitAction:
		m.Quitting = true
		return m, tea.Quit

	default:
		return m, nil
	}
}

// updateSubMenu handles updates when in a sub-menu
func (m Model) updateSubMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle escape to go back to main menu
		if msg.String() == "esc" {
			m.InSubMenu = false
			m.SubProgram = nil
			return m, nil
		}
	}

	// Update the sub-program
	var cmd tea.Cmd
	m.SubProgram, cmd = m.SubProgram.Update(msg)

	// Check if sub-program wants to quit (return to main menu)
	if shouldReturnToMenu(m.SubProgram) {
		m.InSubMenu = false
		m.SubProgram = nil
		return m, nil
	}

	return m, cmd
}

// shouldReturnToMenu checks if we should return to the main menu
func shouldReturnToMenu(subProgram tea.Model) bool {
	// Check if listdevices model has quit
	if listModel, ok := subProgram.(listdevices.Model); ok {
		return listModel.Quitting || listModel.Choice != nil
	}
	
	// Check if connect model has quit or finished
	if connectModel, ok := subProgram.(connect.Model); ok {
		return connectModel.DevicePicker.Quitting || 
			   (connectModel.State == connect.ShowResult && connectModel.Result != nil)
	}
	
	// Check if disconnect model has quit or finished
	if disconnectModel, ok := subProgram.(disconnect.Model); ok {
		return disconnectModel.DevicePicker.Quitting || 
			   (disconnectModel.State == disconnect.ShowResult && disconnectModel.Result != nil)
	}

	return false
}