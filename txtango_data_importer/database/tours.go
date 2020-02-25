package database

import (
	"log"
	"strconv"
	"time"
	"tx2db/txtango"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

//Tour represents information data about truck tours
//A tour is a period of driving connected to one driver
//Example Driver A and Driver B in the same trip will result in 2 Tours
type Tour struct {
	gorm.Model
	TruckID                uint
	DriverID               uint
	TrailerID              uint
	TruckActivityReport    []TruckActivityReport    `gorm:"foreignkey:TourID"`
	DriverEcoMonitorReport []DriverEcoMonitorReport `gorm:"foreignkey:TourID"`
	DestinationLongitude   float32
	DestinationLatitude    float32
	Status                 string
	StartTime              time.Time
	EndTime                time.Time
}

//TruckActivityReport represents the activity report of a specific truck
type TruckActivityReport struct {
	gorm.Model
	TruckID      uint
	TourID       uint
	KmBegin      int
	KmEnd        int
	Consumption  float32
	LoadedStatus string
	Activity     string
	SpeedAvg     float32
	Longitude    float32
	Latitude     float32
	AddressInfo  string
	CountryCode  string
	Reference    string
	StartTime    time.Time
	EndTime      time.Time
}

//DriverEcoMonitorReport represents the eco monitor report of a driver
//EcoMonitorReport trip is determined from contact ON to contact OFF
type DriverEcoMonitorReport struct {
	gorm.Model
	DriverID                                           uint
	TourID                                             uint
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

//checkTour handles tour import and creation flow
func checkTour(truck *Truck, driverTransicsID, trailerTransicsID, tourStatus string, long, lat float32) error {
	// if transics id not set, then do not create tour
	if driverTransicsID == "" {
		return nil
	}

	//get driver
	var driver Driver
	transicsID, err := strconv.ParseUint(driverTransicsID, 10, 64)
	if err != nil {
		return errors.Wrap(err, errParsingTransicsID)
	}

	driver.TransicsID = uint(transicsID)
	if err := db.Where(&driver).Find(&driver).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, errDatabaseConnection)
		}
	}

	//get trailer
	var trailer Trailer
	if trailerTransicsID != "" {
		transicsID, err = strconv.ParseUint(trailerTransicsID, 10, 64)
		if err != nil {
			return errors.Wrap(err, errParsingTransicsID)
		}

		trailer.TransicsID = uint(transicsID)
		if err := db.Where(&trailer).Find(&trailer).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, errDatabaseConnection)
			}
		}
	}

	//initialize tour
	//keep in mind that if a driver keep doing the same tour with the same destination, his tour will never be finished.
	var tour Tour
	newTour := Tour{
		TruckID:              truck.ID,
		DriverID:             driver.ID,
		TrailerID:            trailer.ID,
		DestinationLongitude: long,
		DestinationLatitude:  lat,
	}

	status := "Skipped"
	if err := db.Where(newTour).First(&tour).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, errDatabaseConnection)
		}

		//check how many tour has a truck
		var count int
		if err = db.Model(&tour).Where(Tour{TruckID: truck.ID}).Count(&count).Error; err != nil {
			return errors.Wrap(err, errDatabaseConnection)
		}

		if count > 0 {
			now := time.Now()
			// Get old tour
			var oldTour Tour
			if err := db.Model(&tour).Where(Tour{TruckID: truck.ID}).Last(&oldTour).Error; err != nil {
				return errors.Wrap(err, errDatabaseConnection)
			}

			// Start new tour and old tour using now date
			db.Model(&tour).Where(oldTour).Update(Tour{EndTime: now})
			newTour.StartTime = now
		} else {
			//set startTime first ever tour imported to yesterday date
			newTour.StartTime = time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		}

		// create tour
		status = "Creating a new"
		newTour.Status = tourStatus
		db.Create(&newTour)
	} else if truck.LastModified.After(tour.UpdatedAt) { // update tour
		status = "Updating"
		db.Model(&tour).Where(newTour).Update(Tour{Status: tourStatus})
	}

	//add tour
	log.Printf("%s tour of driver %s in truck %d\n", status, driverTransicsID, truck.TransicsID)

	return nil
}

//ImportToursData import DriverEcoMonitorReport and TruckActivityReport
//Yet we import only yesterday data, so if the program quit for x days, x days of data will not be imported
//It would nice to add a check for that in a future version
func ImportToursData() error {
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)

	var trucks []Truck
	err := db.Where("last_modified > ?", yesterday).Find(&trucks).Error
	if err != nil {
		return errors.Wrap(err, errDatabaseConnection)
	}

	for i, truck := range trucks {
		//import data from transics
		log.Printf("(%d / %d) %s\n", i+1, len(trucks), loadingDataFromTransics)
		txTruckActivity, err := txtango.GetActivityReport(int(truck.TransicsID), yesterday)
		if err != nil {
			return err
		}

		//check and return error
		if txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Errors.Error.CodeExplenation != "" {
			log.Printf("ERROR: %s\n", txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Errors.Error.CodeExplenation)
		}

		//check and print warning
		if txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Warnings.Warning.CodeExplenation != "" {
			log.Printf("WARNING: %s\n", txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Warnings.Warning.CodeExplenation)
		}

		var tour Tour
		//get matching tour
		//TODO check end date too
		err = db.Where("start_time >= ? AND truck_id = ?", yesterday, truck.ID).First(&tour).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, errDatabaseConnection)
			}
			log.Printf("No tour has been found for truck %d. Skip.", truck.TransicsID)
			continue
		}

		for _, data := range txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.ActivityReportItems.ActivityReportItemV11 {
			//parse begin and end date into time.Time
			startTime, err := time.Parse("2006-01-02T15:04:05", data.BeginDate)
			if err != nil {
				log.Println(errParsingDate)
				startTime = time.Time{}
			}

			endTime, err := time.Parse("2006-01-02T15:04:05", data.EndDate)
			if err != nil {
				log.Println(errParsingDate)
				endTime = time.Time{}
			}

			newTruckActivity := TruckActivityReport{
				TourID:       tour.ID,
				TruckID:      tour.TruckID,
				KmBegin:      data.KmBegin,
				KmEnd:        data.KmEnd,
				Consumption:  data.Consumption,
				LoadedStatus: data.LoadedStatus,
				Activity:     data.Activity.Name,
				SpeedAvg:     data.SpeedAvg,
				Longitude:    data.Position.Longitude,
				Latitude:     data.Position.Latitude,
				AddressInfo:  data.Position.AddressInfo,
				CountryCode:  data.Position.CountryCode,
				Reference:    data.Reference,
				StartTime:    startTime,
				EndTime:      endTime,
			}

			//add truck activty
			log.Printf("Activity added for Truck %d\n", truck.TransicsID)
			db.Create(&newTruckActivity)
		}
	}

	return nil
}
