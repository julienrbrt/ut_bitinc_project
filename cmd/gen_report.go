package cmd

import (
	"tx2db/template"

	"github.com/spf13/cobra"
)

var (
	//flags specifying which report to generates
	truckReportOnly      bool
	driverReportOnly     bool
	driverReportSelfOnly bool
)

var reportCmd = &cobra.Command{
	Use:   "gen-report [OPTIONS]",
	Short: "Generate analysis report",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if err := template.InitR(); err != nil {
			return err
		}

		if truckReportOnly {
			wg.Add(1)
			go func() {
				err = template.BuildTruckReport(&wg)
			}()
		}

		if driverReportOnly {
			wg.Add(1)
			go func() {
				err = template.BuildDriverReport(&wg)
			}()
		}

		if driverReportSelfOnly {
			wg.Add(1)
			go func() {
				err = template.BuildDriverSelfReport(&wg)
			}()
		}

		//handle only one error
		if err != nil {
			return err
		}

		if !truckReportOnly && !driverReportOnly && !driverReportSelfOnly {

		}
		wg.Wait()

		return nil
	},
}

func init() {
	//--truckReportOnly flag
	reportCmd.PersistentFlags().BoolVar(&truckReportOnly, "truckReportOnly", false, "Generate truck reports only")
	//--driverReportOnly flag
	reportCmd.PersistentFlags().BoolVar(&driverReportOnly, "driverReportOnly", false, "Generate driver reports only")
	//--driverReportSelfOnly flag
	reportCmd.PersistentFlags().BoolVar(&driverReportSelfOnly, "driverReportSelfOnly", false, "Generate driver reports compared to themself only")
	rootCmd.AddCommand(reportCmd)
}
