package cmd

import (
	"fmt"
	"os"

	"gomux/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gomux",
	Short: "Gomux - tmux script manager and runner",
	Run: func(cmd *cobra.Command, args []string) {
		// If no args or flags, launch TUI
		if len(os.Args) == 1 {
			p := tea.NewProgram(tui.InitialModel())
			if _, err := p.Run(); err != nil {
				fmt.Println("Error running TUI:", err)
				os.Exit(1)
			}
			return
		}
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
