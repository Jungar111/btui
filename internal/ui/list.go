// Package ui contains user interface components
package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// DeviceItem represents a Bluetooth device that can be displayed in a list
type DeviceItem struct {
	title       string
	description string
	device      any // Store the actual device data
}

// Title implements list.Item
func (i DeviceItem) Title() string { return i.title }

// Description implements list.Item  
func (i DeviceItem) Description() string { return i.description }

// FilterValue implements list.Item
func (i DeviceItem) FilterValue() string { return i.title }

// Device returns the stored device data
func (i DeviceItem) Device() any { return i.device }

// NewDeviceItem creates a new device item
func NewDeviceItem(title, description string, device any) DeviceItem {
	return DeviceItem{
		title:       title,
		description: description,
		device:      device,
	}
}

// GenericItem represents any item that can be displayed in a list (kept for backward compatibility)
type GenericItem struct {
	Title       string
	Description string
	Value       any
}

// FilterValue implements list.Item
func (i GenericItem) FilterValue() string {
	return i.Title + " " + i.Description
}

// newDeviceDelegate creates an extended default delegate with status-specific colors
func newDeviceDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	
	// Customize selection colors to use terminal colors that fit our scheme
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("15")).
    BorderLeftForeground(lipgloss.Color("2")).
		Bold(true).
    Italic(true)
	
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
    BorderLeftForeground(lipgloss.Color("2")).
		Foreground(lipgloss.Color("7")) // Terminal white

	
	return d
}

// NewList creates a new list with colored status indicators
func NewList(items []list.Item, title string, width, height int) list.Model {
	delegate := newDeviceDelegate()
	l := list.New(items, delegate, width, height)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle
	return l
}

// NewListWithKeys creates a new list with colored status indicators and custom key map
func NewListWithKeys(items []list.Item, title string, width, height int, keyMap interface{}) list.Model {
	delegate := newDeviceDelegate()
	l := list.New(items, delegate, width, height)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		if km, ok := keyMap.(interface{ ShortHelp() []key.Binding }); ok {
			return km.ShortHelp()
		}
		return []key.Binding{}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		if km, ok := keyMap.(interface{ FullHelp() [][]key.Binding }); ok {
			fullHelp := km.FullHelp()
			var allKeys []key.Binding
			for _, row := range fullHelp {
				allKeys = append(allKeys, row...)
			}
			return allKeys
		}
		return []key.Binding{}
	}
	return l
}
