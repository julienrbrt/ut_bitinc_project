package analysis

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"text/template"
	"time"
	"tx2db/util"

	"github.com/kardianos/osext"
	"github.com/pkg/errors"
)

//DriverReportData contains the data of a report
type DriverReportData struct {
	FullName         string
	TransicsID       string
	TruckDriven      []string
	DrivenKm         string
	CruiseControl    string
	PanicBrakes      string
	FuelConsumption  string
	VisitedCountries []string
	PersonalJoke     string
}

//getReportTemplate gets the template
func getReportTemplate() (*template.Template, error) {
	//get program path
	wd, err := osext.ExecutableFolder()
	if err != nil {
		return &template.Template{}, err
	}

	//path of templates used to build the driver report
	var reportPath = path.Join("analysis", "driver_report.html")

	tmpl, err := template.ParseFiles(path.Join(wd, reportPath))
	if err != nil {
		return &template.Template{}, err
	}

	return tmpl, nil
}

//runR launch the R analysis
func runR(startTime, endTime string) error {
	//path of the analysis
	analysis := path.Join("analysis", "analysis.R")

	//Run the analysis
	r := exec.Command("Rscript", analysis, startTime, endTime)
	//display error and output
	r.Stdout = os.Stdout
	r.Stderr = os.Stderr

	if err := r.Run(); err != nil {
		return errors.Wrap(err, "r.Run() failed")
	}

	return nil
}

//printReport will save a template from disk
func printReport(template, transicsID, startTime string) {
	_ = fmt.Sprintf("driver_%s_report_%s.png", transicsID, startTime)

}

//BuildDriverReport builds a report aimed at drivers
func BuildDriverReport() error {
	log.Print("Building drivers reports...")

	//reportRange defines the number of days a report contains -  according to the .env file
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
	driverName, err := getDriverName(driverList)
	if err != nil {
		return err
	}
	truckDriven, err := getTruckDriven(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	panicBrakes, err := getTotalPanicBrakes(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	vistedCountries, err := getVisitedCountries(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	cruiseControl, err := getCruiseControlRatio(driverList, startTime, endTime)
	if err != nil {
		return err
	}
	fuelConsumption, err := getFuelConsumption(driverList, startTime, endTime)
	if err != nil {
		return err
	}

	//runR analysis
	// if err := runR(startTime.Format("2006-01-02"), endTime.Format("2006-01-02")); err != nil {
	// 	return err
	// }

	//store report template
	tmpl, err := getReportTemplate()
	if err != nil {
		return err
	}

	//Create report based on someone who drove
	//Assumes that all request are sorted
	for i, driverDrivenKm := range drivenKm {
		var data DriverReportData

		//assign metric to report data
		data.DrivenKm = driverDrivenKm.Metric
		data.TransicsID = driverDrivenKm.TransicsID

		if driverName[i].TransicsID == data.TransicsID {
			data.FullName = driverName[i].Metric
		}

		for _, truck := range truckDriven {
			if truck.TransicsID != data.TransicsID {
				break
			}
			data.TruckDriven = append(data.TruckDriven, truck.Metric)
		}

		if panicBrakes[i].TransicsID == data.TransicsID {
			data.PanicBrakes = panicBrakes[i].Metric
		}

		for _, country := range vistedCountries {
			if country.TransicsID != data.TransicsID {
				break
			}
			data.VisitedCountries = append(data.VisitedCountries, country.Metric)
		}

		if cruiseControl[i].TransicsID == data.TransicsID {
			value, _ := strconv.ParseFloat(cruiseControl[i].Metric, 32)
			data.CruiseControl = fmt.Sprintf("%.2f", math.Round(value*100))
		}

		if fuelConsumption[i].TransicsID == data.TransicsID {
			data.FuelConsumption = fuelConsumption[i].Metric
		}

		//get personal joke
		data.PersonalJoke = util.GetJoke()

		//fill in template
		report := tmpl
		err = report.Execute(os.Stdout, data)

		//TODO print template to png
		break
	}

	return nil
}
