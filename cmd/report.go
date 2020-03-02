package cmd

import (
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "gen-report",
	Short: "Build analysis report",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
