package template

import (
	"io/ioutil"
	"path"

	"github.com/kardianos/osext"
)

var (
	//variables storing the path of each template used to build report
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

//BuildTruckReport builds a report aimed at the operation manager
func BuildTruckReport() {

}

//BuildDriverReport builds a report aimed at the driver instructor,
func BuildDriverReport() {

}

//BuildDriverSelfReport builds a report aimed at the operation manager
func BuildDriverSelfReport() {

}
