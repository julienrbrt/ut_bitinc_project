package database

import (
	"log"
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
	DriverTransicsID       uint
	TruckTransicsID        uint
	TrailerTransicsID      uint
	TruckActivityReport    []TruckActivityReport    `gorm:"foreignkey:TourID"`
	DriverEcoMonitorReport []DriverEcoMonitorReport `gorm:"foreignkey:TourID"`
	DestinationLongitude   float32
	DestinationLatitude    float32
	Status                 string
	StartTime              time.Time
	EndTime                time.Time `sql:"default: null"`
	LastImport             time.Time `sql:"default: null"`
}

//TruckActivityReport represents the activity report of a specific truck
type TruckActivityReport struct {
	gorm.Model
	TruckTransicsID uint
	TourID          uint
	KmBegin         int
	KmEnd           int
	Consumption     float32
	LoadedStatus    string
	Activity        string
	SpeedAvg        float32
	Longitude       float32
	Latitude        float32
	AddressInfo     string
	CountryCode     string
	Reference       string
	StartTime       time.Time
	EndTime         time.Time
}

//DriverEcoMonitorReport represents the eco monitor report of a driver
//EcoMonitorReport trip is determined from contact ON to contact OFF
type DriverEcoMonitorReport struct {
	gorm.Model
	TourID                                             uint
	DriverTransicsID                                   uint
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
	NumberOfHarshAccelerations                         int
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

//buildTour handles tour import and creation flow
func buildTour(truck *Truck, driverTransicsID, trailerTransicsID uint, tourStatus string, long, lat float32) error {
	// if transics id not set, then do not create tour
	if driverTransicsID == 0 || long == 0 || lat == 0 {
		return nil
	}

	//initialize tour
	//keep in mind that if a driver keep doing the same tour with the same destination, his tour will never be finished.
	var tour Tour
	newTour := Tour{
		TruckTransicsID:      truck.TransicsID,
		DriverTransicsID:     driverTransicsID,
		TrailerTransicsID:    trailerTransicsID,
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
		if err = db.Model(&tour).Where(Tour{TruckTransicsID: truck.TransicsID}).Count(&count).Error; err != nil {
			return errors.Wrap(err, errDatabaseConnection)
		}

		now := time.Now()
		if count > 0 {
			// Get old tour
			var oldTour Tour
			if err := db.Model(&tour).Where(Tour{TruckTransicsID: truck.TransicsID}).Last(&oldTour).Error; err != nil {
				return errors.Wrap(err, errDatabaseConnection)
			}

			// Start new tour and old tour using now date
			db.Model(&tour).Where(oldTour).Update(Tour{EndTime: now})
			newTour.StartTime = now
		} else {
			//set startTime first ever tour imported to yesterday date
			newTour.StartTime = now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
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
	log.Printf("%s tour of driver %d in truck %d\n", status, driverTransicsID, truck.TransicsID)

	return nil
}

//ImportToursData import TruckActivityReport and DriverEcoMonitorReport
func ImportToursData() error {
	log.Println("Importing tours data")

	now := time.Now()
	var tours []Tour
	//getting the lastest tours which has been end in the past 3 dats
	err := db.Where("last_import IS NULL OR end_time IS NULL OR end_time > last_import").Find(&tours).Error
	if err != nil {
		return errors.Wrap(err, errDatabaseConnection)
	}

	for i, tour := range tours {
		log.Printf("(%d / %d) %s\n", i+1, len(tours), loadingDataFromTransics)

		// caculate elapsed time betfore last import and queries missing days
		if (tour.LastImport == time.Time{}) {
			//if never has been imported, only import the previous day
			tour.LastImport = now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
		}

		//import activity report
		err = importActivityReport(&tour, now)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return err
		}

		//import eco monitor report
		err = importEcoMoniorReport(&tour)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return err

		}
	}
	return nil
}

//importActivityReport import the activity report of truck given a toru
func importActivityReport(tour *Tour, now time.Time) error {
	diff := int(now.Sub(tour.LastImport).Hours() / 24)

	for day := 0; day <= diff; day++ {
		//import data from transics
		txTruckActivity, err := txtango.GetActivityReport(tour.TruckTransicsID, tour.LastImport.AddDate(0, 0, day))
		if err != nil {
			return err
		}

		//check and return error
		if txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Errors.Error != (txtango.TXError{}).Error {
			log.Printf("ERROR: %s - %s\n", txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Errors.Error.Code, txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Errors.Error.Value)
		}

		//check and print warning
		if txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Warnings.Warning != (txtango.TXWarning{}).Warning {
			log.Printf("WARNING: %s - %s\n", txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Warnings.Warning.Code, txTruckActivity.Body.GetActivityReportV11Response.GetActivityReportV11Result.Warnings.Warning.Value)
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

			//check if date is contained in tour date boundaries
			if tour.StartTime.Before(startTime) && (tour.EndTime.After(endTime) || tour.EndTime == time.Time{}) {
				newTruckActivity := TruckActivityReport{
					TourID:          tour.ID,
					TruckTransicsID: tour.TruckTransicsID,
					KmBegin:         data.KmBegin,
					KmEnd:           data.KmEnd,
					Consumption:     data.Consumption,
					LoadedStatus:    data.LoadedStatus,
					Activity:        data.Activity.Name,
					SpeedAvg:        data.SpeedAvg,
					Longitude:       data.Position.Longitude,
					Latitude:        data.Position.Latitude,
					AddressInfo:     data.Position.AddressInfo,
					CountryCode:     data.Position.CountryCode,
					Reference:       data.Reference,
					StartTime:       startTime,
					EndTime:         endTime,
				}

				//add truck activty
				log.Printf("Activity added for Truck %d\n", tour.TruckTransicsID)
				db.Create(&newTruckActivity)

				//update last import of tour
				db.Model(&tour).Where("id = ?", tour.ID).Update(Tour{LastImport: now})
			}
		}
	}

	return nil
}

func importEcoMoniorReport(tour *Tour) error {
	//import data from transics
	txDriverEcoMonitor, err := txtango.GetEcoReport(tour.DriverTransicsID, tour.StartTime)
	if err != nil {
		return err
	}

	//check and return error
	if txDriverEcoMonitor.Body.GetEcoMonitorReportV4Response.GetEcoMonitorReportV4Result.Errors.Error != (txtango.TXError{}).Error {
		log.Printf("ERROR: %s - %s\n", txDriverEcoMonitor.Body.GetEcoMonitorReportV4Response.GetEcoMonitorReportV4Result.Errors.Error.Code, txDriverEcoMonitor.Body.GetEcoMonitorReportV4Response.GetEcoMonitorReportV4Result.Errors.Error.Value)

	}

	//check and print warning
	if txDriverEcoMonitor.Body.GetEcoMonitorReportV4Response.GetEcoMonitorReportV4Result.Warnings.Warning != (txtango.TXWarning{}).Warning {
		log.Printf("WARNING: %s - %s\n", txDriverEcoMonitor.Body.GetEcoMonitorReportV4Response.GetEcoMonitorReportV4Result.Warnings.Warning.Code, txDriverEcoMonitor.Body.GetEcoMonitorReportV4Response.GetEcoMonitorReportV4Result.Warnings.Warning.Value)
	}

	for _, data := range txDriverEcoMonitor.Body.GetEcoMonitorReportV4Response.GetEcoMonitorReportV4Result.EcoMonitorReportItems.EcoMonitorReportItemV3 {
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

		//check if date is contained in tour date boundaries
		if tour.StartTime.Before(startTime) && (tour.EndTime.After(endTime) || tour.EndTime == time.Time{}) {
			newEcoMonitor := DriverEcoMonitorReport{
				TourID:                       tour.ID,
				DriverTransicsID:             tour.DriverTransicsID,
				Distance:                     data.DataResult.Distance,
				DurationDriving:              data.DataResult.Duration,
				FuelConsumption:              data.DataResult.FuelConsumption,
				FuelConsumptionAverage:       data.DataResult.FuelConsumptionAverage.Text,
				RpmAverage:                   data.DataResult.RpmAverage,
				EmissionAverage:              data.DataResult.Co2EmissionAverage.Text,
				SpeedAverage:                 data.DataResult.SpeedAverage,
				FuelConsumptionIdling:        data.IdlingResult.FuelConsumptionIdling,
				DurationIdling:               data.IdlingResult.DurationIdling,
				NumberIdling:                 data.IdlingResult.NumberOfLongIdling,
				DurationOverSpeeding:         data.OverSpeedingResult.DurationOverSpeeding,
				NumberOverSpeeding:           data.OverSpeedingResult.NumberOfOverSpeeding,
				DistanceCoasting:             data.CoastingResult.DistanceCoasting,
				DurationCoasting:             data.CoastingResult.DurationCoasting,
				NumberOfStops:                data.AnticipationResult.NumberOfStops,
				NumberOfBrakes:               data.AnticipationResult.NumberOfBrakes,
				NumberOfPanicBrakes:          data.AnticipationResult.NumberOfPanicBrakes,
				DistanceByBrakes:             data.AnticipationResult.DistanceByBrakes,
				DurationByBrakes:             data.AnticipationResult.DurationByBrakes,
				DurationByRetarder:           data.AnticipationResult.DurationByRetarder,
				DurationHighRPMnoFuel:        data.AnticipationResult.DurationHighRPMnoFuel,
				DurationHighRPM:              data.AnticipationResult.DurationHighRPM,
				NumberOfHarshAccelerations:   data.AnticipationResult.NumberOfHarshAccelerations,
				DurationHarshAcceleration:    data.AnticipationResult.DurationHarshAcceleration,
				DistanceGreenSpot:            data.GreenSpotResult.DistanceGreenSpot,
				DurationGreenSpot:            data.GreenSpotResult.DurationGreenSpot,
				FuelConsumptionGreenSpot:     data.GreenSpotResult.FuelConsumptionGreenSpot,
				NumberOfGearChanges:          data.GearingResult.NumberOfGearChanges,
				NumberOfGearChangesUp:        data.GearingResult.NumberOfGearChangesUp,
				PositionOfThrottleAverage:    data.GearingResult.PositionOfThrottleAverage,
				PositionOfThrottleMaximum:    data.GearingResult.PositionOfThrottleMaximum,
				NumberOfPto:                  data.PtoResult.NumberOfPto,
				FuelConsumptionPtoDriving:    data.PtoResult.FuelConsumptionPtoDriving,
				FuelConsumptionPtoStandStill: data.PtoResult.FuelConsumptionPtoStandStill,
				DurationPtoDriving:           data.PtoResult.DurationPtoDriving,
				DurationPtoStandStill:        data.PtoResult.DurationPtoStandStill,
				DistanceOnCruiseControl:      data.CruisingResult.DistanceOnCruiseControl,
				DurationOnCruiseControl:      data.CruisingResult.DurationOnCruiseControl,
				AvgFuelConsumptionCruiseControlInLiterPerHundredKm: data.CruisingResult.AvgFuelConsumptionCruiseControlInLiterPer100km,
				AvgFuelConsumptionCruiseControlInkmPerLiter:        data.CruisingResult.AvgFuelConsumptionCruiseControlInkmPerLiter,
				StartTime: startTime,
				EndTime:   endTime,
			}

			//add eco monitor for driver
			log.Printf("EcoMonitor Report added for driver %d\n", tour.DriverTransicsID)
			db.Create(&newEcoMonitor)

			//wait 3 seconds to do not be blocked by TX-TANGO
			time.Sleep(3 * time.Second)
		}
	}

	return nil
}
