package cmd

import (
	"log"
	"time"
	"tx2db/analysis"
	"tx2db/database"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	//ignoreCache will ignore the already generated graph and generate them agaub
	ignoreCache bool
	//startTime define the startTime of the report
	startTime string
	//reportRange defines the number of days a report contains
	reportRange int
)

var genReportCmd = &cobra.Command{
	Use: "gen-report",
	Example: `
	tx2db gen-report
	tx2db gen-report --startTime 2020-02-20`,
	Short: "Generate driver reports aimed at drivers only",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var reportTime time.Time

		//get report date
		if startTime == "" {
			reportTime = time.Now().AddDate(0, 0, -8)
		} else {
			//parse begin and end date into time.Time
			reportTime, err = time.Parse("2006-01-02", startTime)
			if err != nil {
				return errors.Wrap(err, "Wrong date format, should be in the format 2020-02-10")
			}
		}

		log.Print("Connecting to database...")
		//connect to database
		err = database.InitDB()
		if err != nil {
			return err
		}
		defer database.DB.Close()

		err = analysis.BuildDriverReport(ignoreCache, reportTime, reportTime.AddDate(0, 0, reportRange))
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	//--ignoreCache flag
	genReportCmd.PersistentFlags().BoolVar(&ignoreCache, "ignoreCache", false, "Ignore already generated graphs")
	//--startTime flags, define the startTime of the report
	genReportCmd.PersistentFlags().StringVar(&startTime, "startTime", "", "Define the start time of a report (default a week ago)")
	//--reportRange flag, default to 7 days
	genReportCmd.PersistentFlags().IntVar(&reportRange, "reportRange", 7, "Define a report range")
	rootCmd.AddCommand(genReportCmd)
}
