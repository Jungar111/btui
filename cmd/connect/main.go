package connect

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// New creates a new cobra command for connecting to devices  
func New() *cobra.Command {
	c := &cobra.Command{}
	c.Use = "connect"
	c.Short = "Connect to a Bluetooth device"
	c.Long = "Select and connect to a Bluetooth device"
	c.Run = run
	return c
}

// run executes the connect command
func run(cmd *cobra.Command, args []string) {
	m := NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}