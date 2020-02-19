package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Driver represents driver of a truck
type Driver struct {
	gorm.Model
	Name  string
	Email string
	Phone string
}

//Truck represents trucks
type Truck struct {
	gorm.Model
	Brand     string
	Type      string
	BuildYear int
}

//Tour represents information data about truck tours
type Tour struct {
	gorm.Model
	Truck       Truck  `gorm:"foreignkey:TruckRefer"`
	Driver      Driver `gorm:"foreignkey:DriverRefer"`
	TruckRefer  uint
	DriverRefer uint
	StartTime   time.Time
	EndTime     time.Time
}
