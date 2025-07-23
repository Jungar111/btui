// Package ui contains user interface components and styling
package ui

import "github.com/charmbracelet/lipgloss"

// Styles using terminal default colors
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Faint(true).
			MarginBottom(2)

	HelpStyle = lipgloss.NewStyle().
			Faint(true).
			MarginTop(1)

	// List styles
	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Bold(true)

	ConnectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				Foreground(lipgloss.Color("10")) // Green color for connected devices

	ConnectedSelectedItemStyle = lipgloss.NewStyle().
					PaddingLeft(2).
					Bold(true).
					Foreground(lipgloss.Color("10")) // Green color for connected devices

	PaginationStyle = lipgloss.NewStyle().
			Faint(true)
)

// SuccessStyle returns the success style
func SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)
}

// ErrorStyle returns the error style
func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)
}
