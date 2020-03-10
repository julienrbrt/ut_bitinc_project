package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var (
	//WaitGroup used to wait for all the goroutines launched here to finish
	wg sync.WaitGroup
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tx2db",
	Short: "Import Transics data in a database",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
