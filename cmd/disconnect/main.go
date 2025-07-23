package disconnect

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// New creates a new cobra command for disconnecting from devices  
func New() *cobra.Command {
	c := &cobra.Command{}
	c.Use = "disconnect"
	c.Short = "Disconnect from a Bluetooth device"
	c.Long = "Select and disconnect from a Bluetooth device"
	c.Run = run
	return c
}

// run executes the disconnect command
func run(cmd *cobra.Command, args []string) {
	m := NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}