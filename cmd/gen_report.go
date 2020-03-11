package cmd

import (
	"log"
	"tx2db/analysis"
	"tx2db/database"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "gen-report",
	Short: "Generate driver reports aimed at drivers only",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		log.Print("Connecting to database...")
		//connect to database
		err = database.InitDB()
		if err != nil {
			panic(err)
		}
		defer database.DB.Close()

		err = analysis.BuildDriverReport()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
