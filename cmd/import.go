package cmd

import (
	"log"
	"sync"
	"tx2db/database"

	"github.com/spf13/cobra"
)

//WaitGroup is used to wait for all the goroutines launched here to finish
var wg sync.WaitGroup

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import TX-TANGO database",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		log.Print("Connecting to database...")
		//connect to database
		err = database.InitDB()
		if err != nil {
			panic(err)
		}
		defer database.DB().Close()

		//connect to redis
		log.Print("Connecting to redis...")
		err = database.InitRedis()
		if err != nil {
			panic(err)
		}
		defer database.RDB().Close()

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
