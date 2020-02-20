package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Driver represents driver of a truck
type Driver struct {
	gorm.Model
	TransicsID int
	Name       string
	Language   string
}

//DriverEcoMonitorReport represents the eco monitor report of a driver
//EcoMonitorReport trip is determined from contact ON to contact OFF
type DriverEcoMonitorReport struct {
	gorm.Model
	Driver                                             Driver `gorm:"foreignkey:TransicsID"`
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

//Truck represents trucks
type Truck struct {
	gorm.Model
	TransicsID   int
	LicensePlate string
	Inactive     bool
	LastModified time.Time
}

//TruckGroup represents group of truck
type TruckGroup struct {
	gorm.Model
	Name  string
	Truck []Truck `gorm:"foreignkey:TransicsID"`
}

//Trailer represents a trailer
type Trailer struct {
	gorm.Model
	TransicsID   int
	LicensePlate string
}

//TruckActivityReport represents the activity report of a specific truck
type TruckActivityReport struct {
	gorm.Model
	TransicsID   int
	Truck        Truck `gorm:"foreignkey:TransicsID"`
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

//Tour represents information data about truck tours
//A tour is a period of driving connected to one driver
//Example Driver A and Driver B in the same trip will result in 2 Tours
type Tour struct {
	gorm.Model
	Truck                  Truck                    `gorm:"foreignkey:TransicsID"`
	Driver                 Driver                   `gorm:"foreignkey:TransicsID"`
	Trailer                Trailer                  `gorm:"foreignkey:TransicsID"`
	TruckActivityReport    []TruckActivityReport    `gorm:"foreignkey:TransicsID"`
	DriverEcoMonitorReport []DriverEcoMonitorReport `gorm:"foreignkey:TransicsID"`
	StartTime              time.Time
	EndTime                time.Time
}
