package cmd

import (
	"log"
	"sync"
	"tx2db/database"

	"github.com/spf13/cobra"
)

var (
	//WaitGroup used to wait for all the goroutines launched here to finish
	wg sync.WaitGroup
	//ignoreLastImport will ignore the last import date and reimport all the data
	ignoreLastImport bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "fetch data from Transics and import it into a database",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		log.Print("Connecting to database...")
		//connect to database
		err = database.InitDB()
		if err != nil {
			panic(err)
		}
		defer database.DB().Close()

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

		err = database.ImportToursData(ignoreLastImport)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	//--ignoreLastImport flag
	importCmd.PersistentFlags().BoolVar(&ignoreLastImport, "ignoreLastImport", false, "Ignore the last import date and refetch everything")
	rootCmd.AddCommand(importCmd)
}
