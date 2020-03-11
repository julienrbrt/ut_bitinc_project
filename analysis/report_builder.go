package analysis

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/kardianos/osext"
	"github.com/pkg/errors"
)

// interesting read to build report
// https://github.com/chromedp/chromedp
// https://github.com/chromedp/examples/blob/master/screenshot/main.go

var (
	//name of the R setup analysis file
	analysis = "./analysis/driver_graphs.R"
	//path of templates used to build the reports
	reportPath = path.Join("analysis", "driver_report.html")
)

//getReportTemplate read template file to a string object
func getReportTemplate(path string) (string, error) {
	//get working directory
	//get program path
	wd, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	html, err := ioutil.ReadFile(wd + path)
	if err != nil {
		return "", err
	}

	return string(html), nil
}

//getReportEndTime builds the report range according to the .env file
func getReportEndTime(start time.Time) time.Time {
	//reportRange defines the number of days a report contains
	reportRange, err := strconv.Atoi(os.Getenv("REPORT_RANGE"))
	if err != nil {
		// a report range is necessary, no recovery possible
		panic(err)
	}

	return start.AddDate(0, 0, reportRange)
}

//initR launches the R analysis
func initR() error {
	log.Println("Starting R analysis")

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
}

//BuildDriverReport builds a report aimed at drivers
func BuildDriverReport() error {
	log.Print("Building drivers reports for drivers...")

	//get metrics
	_, err := getTruckDriven(time.Now().AddDate(0, 0, -14))
	if err != nil {
		return err
	}

	//initialize R
	// if err := initR(); err != nil {
	// 	return err
	// }

	return nil
}
