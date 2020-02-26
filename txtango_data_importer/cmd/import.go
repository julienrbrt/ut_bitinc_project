package cmd

import (
	"sync"
	"tx2db/database"

	"github.com/spf13/cobra"
)

//WaitGroup is used to wait for all the goroutines launched here to finish
var wg sync.WaitGroup

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import TX-TANGO database into better-driving database",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		wg.Add(1)
		go func() {
			err = database.ImportDrivers(&wg)
		}()

		wg.Add(1)
		go func() {
			err = database.ImportTrucks(&wg)
		}()

		//handle only one error
		if err != nil {
			return err
		}
		wg.Wait()

		err = database.ImportToursData()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
