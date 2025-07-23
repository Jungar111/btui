// Package ui contains user interface components
package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// GenericItem represents any item that can be displayed in a list
type GenericItem struct {
	Title       string
	Description string
	Value       interface{} // Store any data you need
}

// FilterValue implements list.Item
func (i GenericItem) FilterValue() string {
	return i.Title + " " + i.Description
}

// Custom item delegate for styling
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(GenericItem)
	if !ok {
		return
	}

	var str string
	if i.Description != "" {
		str = fmt.Sprintf("%s (%s)", i.Title, i.Description)
	} else {
		str = i.Title
	}

	// Check if this is a connected device by looking for the connection indicator
	isConnected := strings.Contains(i.Title, "🔗") || strings.Contains(i.Description, "(Connected)")

	fn := ItemStyle.Render
	if index == m.Index() {
		if isConnected {
			fn = func(s ...string) string {
				return ConnectedSelectedItemStyle.Render("> " + strings.Join(s, " "))
			}
		} else {
			fn = func(s ...string) string {
				return SelectedItemStyle.Render("> " + strings.Join(s, " "))
			}
		}
	} else if isConnected {
		fn = ConnectedItemStyle.Render
	}

	fmt.Fprint(w, fn(str))
}

// NewList creates a new generic list
func NewList(items []list.Item, title string, width, height int) list.Model {
	l := list.New(items, itemDelegate{}, width, height)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle
	return l
}
