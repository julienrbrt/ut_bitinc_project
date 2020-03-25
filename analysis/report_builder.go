package analysis

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
	"tx2db/util"

	"github.com/kardianos/osext"
	"github.com/pkg/errors"
)

//DriverReportData contains the data of a report
type DriverReportData struct {
	FullName         string
	PersonID         string
	TransicsID       string
	TruckDriven      []string
	DrivenKm         string
	CruiseControl    string
	PanicBrakes      string
	FuelConsumption  string
	VisitedCountries []string
	PersonalJoke     string
	StartTime        string
	EndTime          string
}

var (
	//GraphPath is where are stored the R graph
	GraphPath = path.Join("analysis", "assets", "graph")
	//ReportPath is where are stored the reports
	ReportPath = path.Join("analysis", "assets", "report")
	//path of the analysis
	reportTemplatePath = path.Join("analysis", "driver_report.html")
	analysisPath       = path.Join("analysis", "analysis.R")
	phantomPath        = path.Join("analysis", "html2png.js") //path of the html2png.js
	phantomGenPath     = path.Join("analysis", "html2png_gen.js")
)

//getReportTemplate gets the template
func getReportTemplate(wd string) (*template.Template, error) {
	//path of template used to build the driver report
	tmpl, err := template.ParseFiles(path.Join(wd, reportTemplatePath))
	if err != nil {
		return &template.Template{}, err
	}

	return tmpl, nil
}

//runR launch the R analysis
func runR(ignoreCache bool, wd, startTime, endTime string) error {
	//only verify if cache used
	if !ignoreCache {
		//verify that analysis has not already been performed (i.e. no files saved with an endtime)
		graph, err := os.Open(GraphPath)
		if err != nil {
			log.Fatal(err)
		}

		//get the names in only one slice
		names, err := graph.Readdirnames(0)
		for _, f := range names {
			if strings.Contains(f, endTime) {
				log.Println("Skipping analysis as graphs are already generated.")
				return nil
			}
		}
	}
	//Run the analysis
	r := exec.Command("Rscript", analysisPath, startTime, endTime)
	//display error and output
	r.Stdout = os.Stdout
	r.Stderr = os.Stderr

	if err := r.Run(); err != nil {
		return errors.Wrap(err, "r.Run() failed")
	}

	return nil
}

//runPhamtom runs phantomjs to take a convert a html template to png
func runPhantom(wd, reportPath string) error {
	//fill in template
	tmpl, err := template.ParseFiles(path.Join(wd, phantomPath))
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, reportPath)

	if err := ioutil.WriteFile(phantomGenPath, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	//run phantomjs
	phantom := exec.Command("phantomjs", path.Join(wd, phantomGenPath))
	//display error and output
	phantom.Stdout = os.Stdout
	phantom.Stderr = os.Stderr

	if err := phantom.Run(); err != nil {
		return errors.Wrap(err, "phantom.Run() failed")
	}

	log.Printf("Report generated in %s.png\n", reportPath)

	return nil
}

//BuildDriverReport builds a report aimed at drivers
func BuildDriverReport(ignoreCache bool) error {
	log.Print("Building drivers reports...")

	//reportRange defines the number of days a report contains - according to the .env file
	reportRange, err := strconv.Atoi(os.Getenv("REPORT_RANGE"))
	if err != nil {
		// a report range is necessary, no recovery possible
		panic(err)
	}

	//get date
	startTime := time.Now().AddDate(0, 0, -7)
	endTime := startTime.AddDate(0, 0, reportRange)

	//get metrics
	drivenKm, err := getDrivenKm(startTime, endTime)
	if err != nil {
		return err
	}
	//get list of which driver report to build
	var driverList []string
	for _, driver := range drivenKm {
		driverList = append(driverList, driver.TransicsID)
	}
	//get metrics
	driverName, err := getDriverName(driverList)
	if err != nil {
		return err
	}
	//get metrics
	personID, err := getPersonID(driverList)
	if err != nil {
		return err
	}
	//get metrics
	truckDriven, err := getTruckDriven(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	//get metrics
	panicBrakes, err := getTotalPanicBrakes(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	//get metrics
	vistedCountries, err := getVisitedCountries(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	//get metrics
	cruiseControl, err := getCruiseControlRatio(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	//get metrics
	fuelConsumption, err := getFuelConsumption(driverList, startTime, endTime)
	if err != nil {
		return err
	}

	//get program path
	wd, err := osext.ExecutableFolder()
	if err != nil {
		return err
	}

	//runR analysis
	if err := runR(ignoreCache, wd, startTime.Format("2006-01-02"), endTime.Format("2006-01-02")); err != nil {
		return err
	}

	//store report template
	tmpl, err := getReportTemplate(wd)
	if err != nil {
		return err
	}

	//Fill in template -- Assumes that all results are sorted
	for i, driverDrivenKm := range drivenKm {
		var data DriverReportData

		//assign metric to report data
		kms, _ := strconv.ParseFloat(driverDrivenKm.Metric, 32)
		data.DrivenKm = fmt.Sprintf("%.1f", kms)
		data.TransicsID = driverDrivenKm.TransicsID
		//used to get the right image
		data.StartTime = startTime.Format("2006-01-02")
		data.EndTime = endTime.Format("2006-01-02")

		if driverName[i].TransicsID == data.TransicsID {
			data.FullName = strings.ToUpper(driverName[i].Metric)
		}

		if personID[i].TransicsID == data.TransicsID {
			data.PersonID = personID[i].Metric
		}

		for _, truck := range truckDriven {
			if truck.TransicsID == data.TransicsID {
				data.TruckDriven = append(data.TruckDriven, truck.Metric)
			}
		}

		if panicBrakes[i].TransicsID == data.TransicsID {
			data.PanicBrakes = fmt.Sprintf("%sx", panicBrakes[i].Metric)
		}

		for _, country := range vistedCountries {
			if country.TransicsID == data.TransicsID {
				if country.Metric != "" {
					data.VisitedCountries = append(data.VisitedCountries, country.Metric)
				}
			}
		}

		if cruiseControl[i].TransicsID == data.TransicsID {
			value, _ := strconv.ParseFloat(cruiseControl[i].Metric, 32)
			data.CruiseControl = fmt.Sprintf("%.1f%%", value*100)
		}

		if fuelConsumption[i].TransicsID == data.TransicsID {
			data.FuelConsumption = fmt.Sprintf("%sL", fuelConsumption[i].Metric)
		}

		//get personal joke (short only)
		for len(data.PersonalJoke) == 0 || len(data.PersonalJoke) > 300 {
			data.PersonalJoke = util.GetJoke()
		}

		//fill in template
		report := tmpl
		buf := &bytes.Buffer{}
		err = report.Execute(buf, data)

		//save template to disk
		genReportPath := path.Join(wd, ReportPath, fmt.Sprintf("driver_%s_report_%s", data.TransicsID, endTime.Format("2006-01-02")))
		if err := ioutil.WriteFile(genReportPath+".html", buf.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}

		//save template to png
		if err := runPhantom(wd, genReportPath); err != nil {
			log.Fatal(err)
		}

	}

	return nil
}
