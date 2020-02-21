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

var (
	loadingDataFromTransics = "Loading data from Transics TX-TANGO... this could take a while..."
	errParsingTransicsID    = "Error when parsing TransicsID"
	errParsingDate          = "Error while parsing date from Transics TX-TANGO"
	errDatabaseConnection   = "Connection error to the database"
)

//Truck represents trucks
type Truck struct {
	gorm.Model
	TruckGroupID        uint
	TruckActivityReport []TruckActivityReport `gorm:"foreignkey:TruckID"`
	Tour                []Tour                `gorm:"foreignkey:TruckID"`
	TransicsID          uint
	LicensePlate        string
	Inactive            bool
	LastModified        time.Time
}

//TruckGroup represents group of truck
type TruckGroup struct {
	gorm.Model
	Name  string
	Truck []Truck `gorm:"foreignkey:TruckGroupID"`
}

//Trailer represents a trailer
type Trailer struct {
	gorm.Model
	Tour         []Tour `gorm:"foreignkey:TrailerID"`
	TransicsID   uint
	LicensePlate string
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

//ImportTrucks imports all the trucks from TX-Tango and fill the database
func ImportTrucks(wg *sync.WaitGroup) error {
	//notify WaitGroup that we're done
	defer wg.Done()

	//import data from transics
	fmt.Println(loadingDataFromTransics)
	txVehicle, err := txtango.GetVehicle()
	if err != nil {
		return err
	}

	//check and return error
	if txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Errors.Error.CodeExplenation != "" {
		return errors.New(txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Errors.Error.CodeExplenation)
	}

	//check and print warning
	if txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Warnings.Warning.CodeExplenation != "" {
		fmt.Printf("WARNING: %s\n", txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Warnings.Warning.CodeExplenation)
	}

	for i, data := range txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Vehicles.InterfaceVehicleResultV13 {
		//import trailer of a vehicle asynchronously
		go addTrailer(data.Trailer)

		transicsID, err := strconv.ParseUint(data.VehicleTransicsID, 10, 64)
		if err != nil {
			return errors.Wrap(err, errParsingTransicsID)
		}

		//parse modified date into time.Time if existing
		modifiedDate, err := time.Parse("2006-01-02T15:04:05", data.Modified)
		if err != nil {
			fmt.Println(errParsingDate)
			modifiedDate = time.Time{}
		}

		newTruck := Truck{
			TransicsID:   uint(transicsID),
			LicensePlate: data.LicensePlate,
			Inactive:     data.Inactive,
			LastModified: modifiedDate,
		}

		//add or update truck
		var truck Truck
		status := "Skipped"

		if err = db.Where(Truck{TransicsID: newTruck.TransicsID}).First(&truck).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, errDatabaseConnection)
			}
			// add truck
			status = "Importing"
			db.Create(&newTruck)
		} else if truck.LastModified.Before(newTruck.LastModified) {
			// update truck
			status = "Updated"
			db.Model(&truck).Where(Truck{TransicsID: newTruck.TransicsID}).Update(newTruck)
		}

		//add truck
		fmt.Printf("(%d / %d) %s truck %d\n", i+1, len(txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Vehicles.InterfaceVehicleResultV13), status, newTruck.TransicsID)

		//add truck to group
		addGroup(newTruck, data.Groups.TxConnectGroups.ConnectGroups.ConnectGroup[0].SubGroup)
	}

	return nil
}

//add the trailer of a truck
//as error are not extremely important for this sub-category, there no error handling
func addTrailer(txTrailer txtango.TXTrailer) {
	transicsID, err := strconv.ParseUint(txTrailer.TransicsID, 10, 64)
	if err != nil || transicsID == 0 {
		return
	}

	trailer := Trailer{
		TransicsID:   uint(transicsID),
		LicensePlate: txTrailer.LicensePlate,
	}

	db.FirstOrCreate(&trailer, trailer)
}

//assign a group to a truck
//as error are not extremely important for this sub-category, there no error handling
func addGroup(truck Truck, groupName string) {
	truckGroup := TruckGroup{Name: groupName}
	db.FirstOrCreate(&truckGroup, truckGroup)

	if err := db.Model(&truckGroup).Association("Truck").Find(&truck).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}

		fmt.Printf("Truck %d added to TruckGroup %s (now containing %d trucks)\n", truck.TransicsID, groupName, db.Model(&truckGroup).Association("Truck").Count())
		db.Model(&truckGroup).Association("Truck").Append(truck)
	}
}
