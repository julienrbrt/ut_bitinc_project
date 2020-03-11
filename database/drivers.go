package database

import (
	"log"
	"sync"
	"time"
	"tx2db/txtango"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	loadingDataFromTransics = "Loading data from Transics TX-TANGO... this could take a while..."
	errParsingTransicsID    = "Error when parsing TransicsID"
	errParsingDate          = "Error while parsing date from Transics TX-TANGO"
	errParsingCoordinates   = "Error parsing destination coordinates"
)

//Driver represents driver of a truck
type Driver struct {
	gorm.Model
	DriverEcoMonitorReportID []DriverEcoMonitorReport `gorm:"foreignkey:DriverTransicsID"`
	Tour                     []Tour                   `gorm:"foreignkey:DriverTransicsID"`
	TransicsID               uint
	PersonID                 string // identifier used within bolk and not by transics
	Name                     string
	Language                 string
	Inactive                 bool
	LastModified             time.Time
}

//ImportDrivers imports all the driver from Transics and fill the database
func ImportDrivers(wg *sync.WaitGroup) error {
	//notify WaitGroup that we're done
	defer wg.Done()

	//import data from transics
	log.Println(loadingDataFromTransics)
	txDrivers, err := txtango.GetDrivers()
	if err != nil {
		return err
	}

	//check and return error
	if txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Errors.Error != (txtango.TXError{}).Error {
		log.Printf("ERROR: %s - %s\n", txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Errors.Error.Code, txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Errors.Error.Value)
	}

	//check and print warning
	if txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Warnings.Warning != (txtango.TXWarning{}).Warning {
		log.Printf("WARNING: %s - %s\n", txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Errors.Error.Code, txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Warnings.Warning.Value)
	}

	for i, data := range txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Persons.InterfacePersonResultV9 {
		//parse modified date into time.Time if existing
		modifiedDate, err := time.Parse("2006-01-02T15:04:05", data.UpdateDatesList.UpdateDatesItem.DateLastUpdate)
		if err != nil {
			log.Println(errParsingDate)
			modifiedDate = time.Time{}
		}

		newDriver := Driver{
			TransicsID:   data.PersonTransicsID,
			PersonID:     data.PersonExternalCode,
			Name:         data.FormattedName,
			Language:     data.Languages.WorkingLanguage,
			Inactive:     data.Inactive,
			LastModified: modifiedDate,
		}

		//add or update driver
		var driver Driver
		status := "Skipped"

		if err = DB.Where(Driver{TransicsID: newDriver.TransicsID}).First(&driver).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, ErrorDB)
			}
			// add driver
			status = "Importing"
			DB.Create(&newDriver)
		} else if driver.LastModified.Before(newDriver.LastModified) {
			// update driver
			status = "Updated"
			DB.Model(&driver).Where(Driver{TransicsID: newDriver.TransicsID}).Update(newDriver)
		}

		log.Printf("(%d / %d) %s driver %d\n", i+1, len(txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Persons.InterfacePersonResultV9), status, newDriver.TransicsID)
	}

	return nil
}
