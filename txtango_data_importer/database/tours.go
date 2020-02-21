package database

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
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
	StartTime              time.Time
	EndTime                time.Time
}

//CreateTour creates a tour
func CreateTour(wg *sync.WaitGroup) error {

	return nil
}
