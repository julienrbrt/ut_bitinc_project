package cmd

import (
	"tx2db/bolk"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports TX-TANGO database into better-driving database",
	RunE: func(cmd *cobra.Command, args []string) error {
		// connect to better-driver db
		err := bolk.InitDB()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
