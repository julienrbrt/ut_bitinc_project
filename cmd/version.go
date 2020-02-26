package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//version number, should always be a float
var versionNb = 1.0

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tx2db v%f\n", versionNb)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
