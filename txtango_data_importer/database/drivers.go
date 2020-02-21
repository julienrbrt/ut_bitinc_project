package database

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"tx2db/txtango"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

//Driver represents driver of a truck
type Driver struct {
	gorm.Model
	DriverEcoMonitorReportID []DriverEcoMonitorReport `gorm:"foreignkey:DriverID"`
	Tour                     []Tour                   `gorm:"foreignkey:DriverID"`
	TransicsID               uint
	Name                     string
	Language                 string
	Inactive                 bool
	LastModified             time.Time
}

//ImportDrivers imports all the driver from TX-Tango and fill the database
func ImportDrivers(wg *sync.WaitGroup) error {
	//notify WaitGroup that we're done
	defer wg.Done()

	//import data from transics
	fmt.Println(loadingDataFromTransics)
	txDrivers, err := txtango.GetDrivers()
	if err != nil {
		return err
	}

	//check and return error
	if txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Errors.Error.CodeExplenation != "" {
		return errors.New(txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Errors.Error.CodeExplenation)
	}

	//check and print warning
	if txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Warnings.Warning.CodeExplenation != "" {
		fmt.Printf("WARNING: %s\n", txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Warnings.Warning.CodeExplenation)
	}

	for i, data := range txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Persons.InterfacePersonResultV9 {
		transicsID, err := strconv.ParseUint(data.PersonTransicsID, 10, 64)
		if err != nil {
			return errors.Wrap(err, errParsingTransicsID)
		}

		//parse modified date into time.Time if existing
		modifiedDate, err := time.Parse("2006-01-02T15:04:05", data.UpdateDatesList.UpdateDatesItem.DateLastUpdate)
		if err != nil {
			fmt.Println(errParsingDate)
			modifiedDate = time.Time{}
		}

		newDriver := Driver{
			TransicsID:   uint(transicsID),
			Name:         data.FormattedName,
			Language:     data.Languages.WorkingLanguage,
			Inactive:     data.Inactive,
			LastModified: modifiedDate,
		}

		//add or update driver
		var driver Driver
		status := "Skipped"

		if err = db.Where(Driver{TransicsID: newDriver.TransicsID}).First(&driver).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, errDatabaseConnection)
			}
			// add driver
			status = "Importing"
			db.Create(&newDriver)
		} else if driver.LastModified.Before(newDriver.LastModified) {
			// update driver
			status = "Updated"
			db.Model(&driver).Where(Driver{TransicsID: newDriver.TransicsID}).Update(newDriver)
		}

		fmt.Printf("(%d / %d) %s driver %d\n", i+1, len(txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Persons.InterfacePersonResultV9), status, newDriver.TransicsID)
	}

	return nil
}
