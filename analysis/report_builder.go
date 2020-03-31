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
	//path of the analysis
	reportFolderPath     = path.Join("analysis", "assets", "report")
	reportTemplatePathDE = path.Join("analysis", "driver_report_de.html")
	reportTemplatePathEN = path.Join("analysis", "driver_report_en.html")
	reportTemplatePathFR = path.Join("analysis", "driver_report_fr.html")
	reportTemplatePathNL = path.Join("analysis", "driver_report_nl.html")
	analysisPath         = path.Join("analysis", "analysis.R")
	//path of the html2png.js
	phantomPath    = path.Join("analysis", "html2png.js")
	phantomGenPath = path.Join("analysis", "html2png_gen.js")
)

//startAnalysis launch the R analysis
func startAnalysis(wd, startTime, endTime string) error {
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

//saveReport runs phantomjs to take a convert a html template to png
func saveReport(wd, genReportPath string) error {
	//fill in template
	tmpl, err := template.ParseFiles(path.Join(wd, phantomPath))
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, genReportPath)

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

	log.Printf("Report successfully generated in %s.png\n", genReportPath)

	return nil
}

//cleanReportFiles remove the uncessary report files required only for its generation
func cleanReportFiles(wd string) error {
	//clean report files
	report, err := os.Open(reportFolderPath)
	if err != nil {
		return err
	}

	names, err := report.Readdirnames(0)
	for _, f := range names {
		if strings.Contains(f, ".html") || strings.Contains(f, "_graph_") {
			if err := os.Remove(path.Join(wd, reportFolderPath, f)); err != nil {
				return errors.Wrap(err, "Could not remove report files")
			}
		}
	}

	return nil
}

//BuildDriverReport builds a report aimed at drivers
func BuildDriverReport(skipSendMail bool, startTime, endTime time.Time) error {
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

	log.Printf("Generating %d drivers reports...", len(driverList))

	//get metrics
	driverName, err := getDriverName(driverList)
	if err != nil {
		return err
	}
	//get metrics
	personID, err := getDriverPersonID(driverList)
	if err != nil {
		return err
	}
	//get metrics
	driverLanguage, err := getDriverLanguage(driverList)
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

	//parse all templates
	tmplDE, err := template.ParseFiles(path.Join(wd, reportTemplatePathDE))
	if err != nil {
		return err
	}
	tmplEN, err := template.ParseFiles(path.Join(wd, reportTemplatePathEN))
	if err != nil {
		return err
	}
	tmplFR, err := template.ParseFiles(path.Join(wd, reportTemplatePathFR))
	if err != nil {
		return err
	}
	tmplNL, err := template.ParseFiles(path.Join(wd, reportTemplatePathNL))
	if err != nil {
		return err
	}

	//start analysis
	if err := startAnalysis(wd, startTime.Format("2006-01-02"), endTime.Format("2006-01-02")); err != nil {
		return err
	}

	//fill in templates -- ASSUME THAT RESULTS ARE SORTED
	for i, driverDrivenKm := range drivenKm {
		var data DriverReportData

		//assign metrics to report data
		kms, _ := strconv.ParseFloat(driverDrivenKm.Metric, 32)
		data.DrivenKm = fmt.Sprintf("%.1f", kms)
		data.TransicsID = driverDrivenKm.TransicsID
		data.StartTime = startTime.Format("2006-01-02")
		data.EndTime = endTime.Format("2006-01-02")
		data.FullName = strings.ToUpper(driverName[i].Metric)
		data.PersonID = personID[i].Metric
		data.PanicBrakes = fmt.Sprintf("%sx", panicBrakes[i].Metric)
		data.FuelConsumption = fmt.Sprintf("%sL", fuelConsumption[i].Metric)
		cc, _ := strconv.ParseFloat(cruiseControl[i].Metric, 32)
		data.CruiseControl = fmt.Sprintf("%.1f%%", cc*100)

		for _, truck := range truckDriven {
			if truck.TransicsID == data.TransicsID {
				data.TruckDriven = append(data.TruckDriven, truck.Metric)
			}
		}

		for _, country := range vistedCountries {
			if country.TransicsID == data.TransicsID {
				if country.Metric != "" {
					data.VisitedCountries = append(data.VisitedCountries, country.Metric)
				}
			}
		}

		//get personal joke (short only)
		for len(data.PersonalJoke) == 0 || len(data.PersonalJoke) > 500 {
			data.PersonalJoke = util.GetJoke(driverLanguage[i].Metric)
		}

		//fill in template (with right translation)
		var report *template.Template
		switch driverLanguage[i].Metric {
		case "DU":
			report = tmplDE
		case "FR":
			report = tmplFR
		case "NL":
			report = tmplNL
		default:
			report = tmplEN
		}
		buf := &bytes.Buffer{}
		err = report.Execute(buf, data)

		//save template to disk
		genReportPath := path.Join(wd, reportFolderPath, fmt.Sprintf("driver_%s_report_%s", data.PersonID, endTime.Format("2006-01-02")))
		if err := ioutil.WriteFile(genReportPath+".html", buf.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}

		//save template to png
		if err := saveReport(wd, genReportPath); err != nil {
			log.Fatal(err)
		}

		//send analysis mail
		if !skipSendMail {
			if err := util.SendReportMail("", genReportPath+".png", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), data.PersonID); err != nil {
				log.Fatalf("ERROR: Mail not sent: %v\n", err)
			}
		}

	}

	//clean report files
	if err := cleanReportFiles(wd); err != nil {
		return err
	}

	return nil
}
