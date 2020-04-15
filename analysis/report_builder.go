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
	Email            string
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
	r := exec.Command("Rscript", path.Join(wd, analysisPath), path.Join(wd, reportFolderPath), startTime, endTime)
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

	if err := ioutil.WriteFile(path.Join(wd, phantomGenPath), buf.Bytes(), 0644); err != nil {
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

//cleanAnalysis remove the uncessary analysis report files required only for its generation
func cleanAnalysis(wd string) error {
	//clean report files
	report, err := os.Open(path.Join(wd, reportFolderPath))
	if err != nil {
		return err
	}
	defer report.Close()

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
func BuildDriverReport(skipSendMail, skipSendDriverMail, skipUploadToFtp bool, startTime, endTime time.Time) error {
	//format start and end time
	formatedStartTime := startTime.Format("2006-01-02")
	formatedEndTime := endTime.Format("2006-01-02")

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

	log.Printf("Generating %d drivers reports for the period %s to %s\n", len(driverList), formatedStartTime, formatedEndTime)

	//get metrics
	driverData, err := getDriverData(driverList)
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

	//start (and clean) analysis
	if err := startAnalysis(wd, formatedStartTime, formatedEndTime); err != nil {
		return err
	}
	defer cleanAnalysis(wd)

	//genReportPathList contains the list of path of the generated reports
	var genReportPathList []string
	//fill in templates -- ASSUME THAT RESULTS ARE SORTED
	for i, driverDrivenKm := range drivenKm {
		var data DriverReportData

		//assign metrics to report data
		kms, _ := strconv.ParseFloat(driverDrivenKm.Metric, 32)
		data.DrivenKm = fmt.Sprintf("%.1f", kms)
		data.TransicsID = driverDrivenKm.TransicsID

		data.StartTime = formatedStartTime
		data.EndTime = formatedEndTime

		data.FullName = strings.ToUpper(driverData[i].Name)
		data.PersonID = driverData[i].PersonID
		data.Email = driverData[i].Email

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

		//add all report path a list
		genReportPathList = append(genReportPathList, genReportPath+".png")

		//send analysis mail to drivers
		if !skipSendMail && !skipSendDriverMail {
			//inform SYSTEM_ADMINISTATOR_EMAIL if no driver mail provided
			if data.Email == "" {
				if err := util.InformSystemAdministratorDriverEmailMissing(data.PersonID); err != nil {
					log.Printf("ERROR: System Administrator not informed of unexisting mail: %v\n", err)
				}
			} else {
				if err := util.InformDriver(data.Email, genReportPath+".png", formatedStartTime, formatedEndTime); err != nil {
					log.Printf("ERROR: Driver mail not informed of available report: %v\n", err)
				}
			}
		}
	}

	//create bulk reports
	if len(genReportPathList) > 0 {
		//build all reports to pdf
		pdfName := fmt.Sprintf("weekly_report_%s.pdf", formatedEndTime)
		pdfPath := path.Join(wd, reportFolderPath, pdfName)
		if err := util.BuildPDFFromImages(pdfPath, genReportPathList); err != nil {
			return err
		}

		//upload pdf to ftp
		if !skipUploadToFtp {
			if err := util.UploadToFTP(pdfName, pdfPath); err != nil {
				//inform system administator
				util.InformSystemAdministratorFTPError(pdfPath)
				return err
			}
			log.Println("Weekly report successfully uploaded to FTP")
		}

		//inform INSTRUCTOR_EMAIL that weekly analysis are available
		if !skipSendMail {
			if err := util.InformInstructor(formatedStartTime, formatedEndTime); err != nil {
				log.Fatalf("ERROR: Instructor not informed of new weekly driver analysis available: %v\n", err)
			}
		}
	}

	return nil
}
