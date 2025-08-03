// Package ui contains user interface components and styling
package ui

import "github.com/charmbracelet/lipgloss"

// Styles using terminal default colors and fancy list styling
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")). // Terminal bright white
			Background(lipgloss.Color("2")).  // Terminal bright black (muted)
			Padding(0, 1).
			Bold(true)

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
				Foreground(lipgloss.Color("2")) // Terminal green for connected devices

	ConnectedSelectedItemStyle = lipgloss.NewStyle().
					PaddingLeft(2).
					Bold(true).
					Foreground(lipgloss.Color("2")) // Terminal green for connected devices

	PaginationStyle = lipgloss.NewStyle().
			Faint(true)

	// Status-specific styles using terminal colors that respect user's theme
	ConnectedStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("2")) // Terminal green

	PairedStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("3")) // Terminal yellow

	DiscoveredStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")) // Terminal cyan

	ConnectingStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("11")) // Terminal bright yellow/orange

	DisconnectingStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("13")) // Terminal bright magenta

	MacAddressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // Terminal bright black (muted)

	RSSIStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")) // Terminal white

	// Application-wide padding style for comfortable spacing
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2) // 1 row padding top/bottom, 2 column padding left/right
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
