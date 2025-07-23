package scan

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// New creates a new cobra command for scanning devices
func New() *cobra.Command {
	c := &cobra.Command{}
	c.Use = "scan"
	c.Short = "Scan for and connect to Bluetooth devices"
	c.Long = "Scan for nearby Bluetooth devices and allow selection and connection"
	c.Run = run
	return c
}

// run executes the scan command
func run(cmd *cobra.Command, args []string) {
	m := NewModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
