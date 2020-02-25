package database

import (
	"log"
	"strconv"
	"time"

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
	TransicsID   uint
	KmBegin      int
	KmEnd        int
	Consumption  float32
	LoadedStatus string
	Activity     string
	SpeedAvg     int
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

//checkTour handles tour import and creation flow
func checkTour(truck *Truck, driverTransicsID, trailerTransicsID, long, lat, tourStatus string) (Tour, error) {
	// if transics id not set, then do not create tour
	if driverTransicsID == "" || trailerTransicsID == "" {
		return Tour{}, nil
	}

	//get driver
	var driver Driver
	transicsID, err := strconv.ParseUint(driverTransicsID, 10, 64)
	if err != nil {
		return Tour{}, errors.Wrap(err, errParsingTransicsID)
	}

	driver.TransicsID = uint(transicsID)
	if err := db.Where(&driver).Find(&driver).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return Tour{}, err
		}
	}

	//get trailer
	var trailer Trailer
	transicsID, err = strconv.ParseUint(trailerTransicsID, 10, 64)
	if err != nil {
		return Tour{}, errors.Wrap(err, errParsingTransicsID)
	}

	trailer.TransicsID = uint(transicsID)
	if err := db.Where(&trailer).Find(&trailer).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return Tour{}, err
		}
	}

	//get destination coordinates
	longitude, err := strconv.ParseFloat(long, 32)
	if err != nil {
		return Tour{}, errors.Wrap(err, errParsingCoordinates)
	}

	latitude, err := strconv.ParseFloat(lat, 32)
	if err != nil {
		return Tour{}, errors.Wrap(err, errParsingCoordinates)
	}

	//initialize tour
	//keep in mind that if a driver keep doing the same tour with the same destination, his tour will never be finished.
	var tour Tour
	newTour := Tour{
		TruckID:              truck.ID,
		DriverID:             driver.ID,
		TrailerID:            trailer.ID,
		DestinationLongitude: float32(longitude),
		DestinationLatitude:  float32(latitude),
	}

	status := "Skipped"
	if err := db.Where(newTour).First(&tour).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return Tour{}, errors.Wrap(err, errDatabaseConnection)
		}

		//check how many tour has a truck
		var count int
		if err = db.Model(&tour).Where(Tour{TruckID: truck.ID}).Count(&count).Error; err != nil {
			return Tour{}, errors.Wrap(err, errDatabaseConnection)
		}

		if count > 0 {
			now := time.Now()
			// Get old tour
			var oldTour Tour
			if err := db.Model(&tour).Where(Tour{TruckID: truck.ID}).Last(&oldTour).Error; err != nil {
				return Tour{}, errors.Wrap(err, errDatabaseConnection)
			}

			// Start new tour and old tour using now date
			db.Model(&tour).Where(oldTour).Update(Tour{EndTime: now})
			newTour.StartTime = now
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

	return newTour, nil
}

//importData import DriverEcoMonitorReport and TruckActivityReport for a certain tour
func importData(tour Tour) error {
	return nil
}
