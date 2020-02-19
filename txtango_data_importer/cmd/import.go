package cmd

import (
	"tx2db/database"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports TX-TANGO database into better-driving database",
	RunE: func(cmd *cobra.Command, args []string) error {
		//TODO use go routine and wait groups
		err := database.ImportTrucks()
		if err != nil {
			return err
		}
		err = database.ImportDrivers()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
