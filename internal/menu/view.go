package menu

// View implements tea.Model
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	// If we're in a sub-menu, show that
	if m.InSubMenu && m.SubProgram != nil {
		return m.SubProgram.View()
	}

	// Show the main menu
	return m.List.View()
}