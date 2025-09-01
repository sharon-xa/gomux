/*
Copyright © 2025 Ali Jabar <ali93456@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a tmux script",
	Run: func(cmd *cobra.Command, args []string) {
		// a := actions.NewActions()
		// err := a.CreateScript()
		// if err != nil {
		// fmt.Println(err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
