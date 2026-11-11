/*
Copyright © 2025 Ali Jabar <ali93456@gmail.com>
*/
package cmd

import (
	"github.com/sharon-xa/gomux/internal/workspace"
	"github.com/sharon-xa/gomux/pkg/actions"
	"github.com/spf13/cobra"
)

var (
	actionsInstance actions.IActions
	createCmd       = &cobra.Command{
		Use:   "create",
		Short: "Build a tmux script configuration",
		Long: `The create command lets you build a tmux script step by step 
or all at once. Use subcommands to add sessions, windows, and panes.

Examples:

  # Create a session only
  gomux create session --name dev --dir ~/projects --type development

  # Add a window with layout
  gomux create window --name editor --cmd "vim ." --layout main-horizontal

  # Add a pane to a window
  gomux create pane --window 1 --direction horizontal --cmd "npm run dev"

  # Create everything at once
  gomux create all --name dev --dir ~/projects`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			actionsInstance = actions.NewActions()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// return actionsInstance.NewScript()
			return nil
		},
	}
)

func newSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Create a new tmux session",
		RunE: func(cmd *cobra.Command, args []string) error {
			// name, _ := cmd.Flags().GetString("name")
			// dir, _ := cmd.Flags().GetString("dir")
			// typ, _ := cmd.Flags().GetString("type")

			// session := workspace.Session{
			// 	Name:      name,
			// 	Directory: dir,
			// 	Type:      typ,
			// }
			// return actionsInstance.AddSession(&session)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Session name")
	cmd.Flags().String("dir", "", "Session working directory")
	cmd.Flags().String("type", "", "Session type")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("dir")

	return cmd
}

func newWindowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "window",
		Short: "Add a window to the current session",
		RunE: func(cmd *cobra.Command, args []string) error {
			// name, _ := cmd.Flags().GetString("name")
			// number, _ := cmd.Flags().GetInt("number")
			// cmdStr, _ := cmd.Flags().GetString("cmd")
			// layout, _ := cmd.Flags().GetString("layout")
			// sync, _ := cmd.Flags().GetBool("sync")
			//
			// w := workspace.Window{
			// 	Name:              name,
			// 	Number:            number,
			// 	Command:           cmdStr,
			// 	Layout:            workspace.Layout(layout),
			// 	SynchronizedPanes: sync,
			// }
			//
			// return actionsInstance.AddWindow(&w)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Window name")
	cmd.Flags().Int("number", 1, "Window number")
	cmd.Flags().String("cmd", "", "Command to run in this window")
	cmd.Flags().String("layout", string(workspace.LayoutEvenVertical), "Layout type")
	cmd.Flags().Bool("sync", false, "Synchronize panes")

	return cmd
}

func newPaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pane",
		Short: "Add a pane to a window",
		RunE: func(cmd *cobra.Command, args []string) error {
			// id, _ := cmd.Flags().GetInt("id")
			// location, _ := cmd.Flags().GetInt("window")
			// dir, _ := cmd.Flags().GetString("direction")
			// cmdStr, _ := cmd.Flags().GetString("cmd")
			//
			// p := workspace.Pane{
			// 	ID:        id,
			// 	Location:  location,
			// 	Direction: workspace.Direction(dir),
			// 	Command:   cmdStr,
			// }
			//
			// return actionsInstance.AddPane(&p)
			return nil
		},
	}

	cmd.Flags().Int("id", 0, "Pane ID")
	cmd.Flags().Int("window", 0, "Window location (number)")
	cmd.Flags().String("direction", string(workspace.Vertical), "Pane split direction")
	cmd.Flags().String("cmd", "", "Command to run in pane")

	return cmd
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(newSessionCmd())
	createCmd.AddCommand(newWindowCmd())
	createCmd.AddCommand(newPaneCmd())
}
