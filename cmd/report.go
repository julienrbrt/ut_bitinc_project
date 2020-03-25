package cmd

import (
	"log"
	"tx2db/analysis"
	"tx2db/database"

	"github.com/spf13/cobra"
)

var (
	//ignoreCache will ignore the already generated graph and generate them agaub
	ignoreCache bool
)

var genReportCmd = &cobra.Command{
	Use:   "gen-report",
	Short: "Generate driver reports aimed at drivers only",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Print("Connecting to database...")
		//connect to database
		err := database.InitDB()
		if err != nil {
			panic(err)
		}
		defer database.DB.Close()

		err = analysis.BuildDriverReport(ignoreCache)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	//--ignoreCache flag
	genReportCmd.PersistentFlags().BoolVar(&ignoreCache, "ignoreCache", false, "Do not use already generated report graphs")
	rootCmd.AddCommand(genReportCmd)
}
