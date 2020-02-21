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

//DriverEcoMonitorReport represents the eco monitor report of a driver
//EcoMonitorReport trip is determined from contact ON to contact OFF
type DriverEcoMonitorReport struct {
	gorm.Model
	DriverID                                           uint
	TourID                                             uint
	TransicsID                                         int
	Distance                                           float32
	DurationDriving                                    float32
	FuelConsumption                                    float32
	FuelConsumptionAverage                             float32
	RpmAverage                                         float32
	EmissionAverage                                    float32
	SpeedAverage                                       float32
	FuelConsumptionIdling                              float32
	DurationIdling                                     float32
	NumberIdling                                       int
	DurationOverSpeeding                               float32
	NumberOverSpeeding                                 int
	DistanceCoasting                                   float32
	DurationCoasting                                   float32
	NumberOfStops                                      int
	NumberOfBrakes                                     int
	NumberOfPanicBrakes                                int
	DistanceByBrakes                                   float32
	DurationByBrakes                                   float32
	DurationByRetarder                                 float32
	DurationHighRPMnoFuel                              float32
	DurationHighRPM                                    float32
	NumberOfHarshAccelerations                         float32
	DurationHarshAcceleration                          float32
	DistanceGreenSpot                                  float32
	DurationGreenSpot                                  float32
	FuelConsumptionGreenSpot                           float32
	NumberOfGearChanges                                int
	NumberOfGearChangesUp                              int
	PositionOfThrottleAverage                          float32
	PositionOfThrottleMaximum                          float32
	NumberOfPto                                        int
	FuelConsumptionPtoDriving                          float32
	FuelConsumptionPtoStandStill                       float32
	DurationPtoDriving                                 float32
	DurationPtoStandStill                              float32
	DistanceOnCruiseControl                            float32
	DurationOnCruiseControl                            float32
	AvgFuelConsumptionCruiseControlInLiterPerHundredKm float32
	AvgFuelConsumptionCruiseControlInkmPerLiter        float32
	StartTime                                          time.Time
	EndTime                                            time.Time
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
