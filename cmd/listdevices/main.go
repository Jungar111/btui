// Package listdevices contains the main entry point for the list devices command
package listdevices

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// New creates a new cobra command for listing devices
func New() *cobra.Command {
	c := &cobra.Command{}
	c.Use = "list-devices"
	c.Short = "List and select Bluetooth devices"
	c.Long = "Display a list of available Bluetooth devices and allow selection"
	c.Run = run
	return c
}

// run executes the list devices command
func run(cmd *cobra.Command, args []string) {
	m := NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
