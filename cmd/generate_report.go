package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//name of the analysis file
const analysis = "./analysis/setup_analysis.R"

var reportCmd = &cobra.Command{
	Use:   "gen-report",
	Short: "Build analysis report",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Starting analysis...")

		//run R analysis
		r := exec.Command("Rscript", analysis)
		//display error and output
		r.Stdout = os.Stdout
		r.Stderr = os.Stderr

		err := r.Run()
		if err != nil {
			return errors.Wrap(err, "r.Run() failed")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
