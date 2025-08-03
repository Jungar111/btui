// Package cmd contains all commands for btui
package cmd

import (
	"btui/cmd/connect"
	"btui/cmd/disconnect"
	"btui/cmd/listdevices"
	"btui/cmd/scan"
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "btui",
		Short: "A TUI for interacting with bluetoothctl",
		Long:  "btui provides a terminal user interface for managing Bluetooth devices using bluetoothctl",
		Run: func(cmd *cobra.Command, args []string) {
			// Launch the scan interface directly
			m := scan.NewModel()
			p := tea.NewProgram(m, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error running program: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// Keep individual commands for direct CLI access if needed
	rootCmd.AddCommand(listdevices.New())
	rootCmd.AddCommand(connect.New())
	rootCmd.AddCommand(disconnect.New())
	rootCmd.AddCommand(scan.New())

	return rootCmd
}

func Execute(ctx context.Context, cmd *cobra.Command) error {
	_, err := cmd.ExecuteContextC(ctx)
	if err != nil {
		return fmt.Errorf("command failed: %v", err)
	}
	return nil
}
