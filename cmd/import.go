package cmd

import (
	"log"
	"tx2db/database"

	"github.com/spf13/cobra"
)

var (
	//ignoreLastImport will ignore the last import date and reimport all the data
	ignoreLastImport bool
	//importFromOnlyQueue will skip the import process and import it only from the queue
	importFromQueueOnly bool
	//cleanTourQueue will delete the entiere queue
	cleanTourQueue bool
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
			return err
		}
		defer database.DB.Close()

		//cleanTourQueue if requested
		if cleanTourQueue {
			database.DB.Unscoped().Delete(&[]database.TourQueue{})
			log.Print("Sucessfully cleaned tour queue")
		}

		wg.Add(1)
		go func() {
			//import drivers concurrently
			err = database.ImportDrivers(&wg)
		}()

		wg.Add(1)
		go func() {
			//import trucks concurrently and create tours
			err = database.ImportTrucks(&wg)
		}()

		//handle only one error
		if err != nil {
			return err
		}
		wg.Wait()

		if importFromQueueOnly {
			//import tours data from queue
			err = database.ImportQueuedToursData(true)
		} else {
			//import tours data
			err = database.ImportToursData(ignoreLastImport)
		}
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	//--ignoreLastImport flag
	importCmd.PersistentFlags().BoolVar(&ignoreLastImport, "ignoreLastImport", false, "Ignore the last import date and refetch everything")
	//--importFromOnlyQueue
	importCmd.PersistentFlags().BoolVar(&importFromQueueOnly, "importFromQueueOnly", false, "Import only missing data from the queue")
	//--cleanTourQueue
	importCmd.PersistentFlags().BoolVar(&cleanTourQueue, "cleanTourQueue", false, "Empty the tour queue")
	rootCmd.AddCommand(importCmd)
}
