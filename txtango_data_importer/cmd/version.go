package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print tx2db version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tx2db v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
