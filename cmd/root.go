package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sharon-xa/gomux/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gomux",
	Short: "Gomux - tmux script manager and runner",
	Run: func(cmd *cobra.Command, args []string) {
		// If no args or flags, launch TUI
		if len(os.Args) == 1 {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				log.Fatalf("error: %w", err)
			}
			defer f.Close()

			p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				log.Fatalln("Error running TUI:", err)
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
