package template

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/kardianos/osext"
	"github.com/pkg/errors"
)

// interesting read to build report
// https://github.com/chromedp/chromedp
// https://github.com/chromedp/examples/blob/master/screenshot/main.go

//name of the R analysis file
const analysis = "./analysis/setup_analysis.R"

var (
	//path of templates used to build the reports
	truckReportPath      = path.Join("template", "html", "truck_report_self.html")
	driverReportPath     = path.Join("template", "html", "tmpl_driver_report.html")
	driverReportSelfPath = path.Join("template", "html", "tmpl_driver_report_self.html")
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

//InitR launches the R analysis
func InitR() error {
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

//BuildTruckReport builds a report aimed at the operation manager
func BuildTruckReport(wg *sync.WaitGroup) error {
	//notify WaitGroup that we're done
	defer wg.Done()

	return nil
}

//BuildDriverReport builds a report aimed at the driver instructor,
func BuildDriverReport(wg *sync.WaitGroup) error {
	//notify WaitGroup that we're done
	defer wg.Done()

	return nil
}

//BuildDriverSelfReport builds a report aimed at the operation manager
func BuildDriverSelfReport(wg *sync.WaitGroup) error {
	//notify WaitGroup that we're done
	defer wg.Done()

	return nil
}
