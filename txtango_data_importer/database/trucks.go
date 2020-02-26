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
	errDatabaseConnection   = "Connection error to the database"
)

//Truck represents trucks
type Truck struct {
	gorm.Model
	TruckGroupID        uint
	TruckActivityReport []TruckActivityReport `gorm:"foreignkey:TruckTransicsID"`
	Tour                []Tour                `gorm:"foreignkey:TruckTransicsID"`
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
	Tour         []Tour `gorm:"foreignkey:TrailerTransicsID"`
	TransicsID   uint
	LicensePlate string
}

//ImportTrucks imports all the trucks from TX-Tango and fill the database
func ImportTrucks(wg *sync.WaitGroup) error {
	//notify WaitGroup that we're done
	defer wg.Done()

	//import data from transics
	log.Println(loadingDataFromTransics)
	txVehicle, err := txtango.GetVehicle()
	if err != nil {
		return err
	}

	//check and return error
	if txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Errors.Error != (txtango.TXError{}).Error {
		log.Printf("ERROR: %s - %s\n", txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Errors.Error.Code, txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Errors.Error.Value)
	}

	//check and print warning
	if txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Warnings.Warning != (txtango.TXWarning{}).Warning {
		log.Printf("WARNING: %s - %s\n", txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Warnings.Warning.Code, txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Warnings.Warning.Value)
	}

	for i, data := range txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Vehicles.InterfaceVehicleResultV13 {
		//import trailer of a vehicle asynchronously
		go addTrailer(&data.Trailer)

		//parse modified date into time.Time if existing
		modifiedDate, err := time.Parse("2006-01-02T15:04:05", data.Modified)
		if err != nil {
			log.Println(errParsingDate)
			modifiedDate = time.Time{}
		}

		newTruck := Truck{
			TransicsID:   data.VehicleTransicsID,
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
			truck = newTruck
		} else if truck.LastModified.Before(newTruck.LastModified) {
			// update truck
			status = "Updated"
			db.Model(&truck).Where(Truck{TransicsID: newTruck.TransicsID}).Update(newTruck)
		}

		//add truck
		log.Printf("(%d / %d) %s truck %d\n", i+1, len(txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Vehicles.InterfaceVehicleResultV13), status, newTruck.TransicsID)

		//start tour flow
		err = buildTour(&truck, data.Driver.TransicsID, data.Trailer.TransicsID, data.ETAInfo.ETAStatus.Text, data.ETAInfo.PositionDestination.Longitude, data.ETAInfo.PositionDestination.Latitude)
		if err != nil {
			// TODO add proper error handling
			log.Print(err)
		}

		//add truck to group
		addGroup(&newTruck, data.Groups.TxConnectGroups.ConnectGroups.ConnectGroup[0].SubGroup)
	}

	return nil
}

//add the trailer of a truck
//as error are not important for this sub-category, there no error handling
func addTrailer(txTrailer *txtango.TXTrailer) {
	// do not create unexisting trailer
	if txTrailer.TransicsID == 0 {
		return
	}

	trailer := Trailer{
		TransicsID:   txTrailer.TransicsID,
		LicensePlate: txTrailer.LicensePlate,
	}

	db.FirstOrCreate(&trailer, trailer)
}

//assign a group to a truck
//as error are not important for this sub-category, there no error handling
func addGroup(truck *Truck, groupName string) {
	truckGroup := TruckGroup{Name: groupName}
	db.FirstOrCreate(&truckGroup, truckGroup)

	if err := db.Model(&truckGroup).Association("Truck").Find(&truck).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}

		log.Printf("Truck %d added to TruckGroup %s (now containing %d trucks)\n", truck.TransicsID, groupName, db.Model(&truckGroup).Association("Truck").Count())
		db.Model(&truckGroup).Association("Truck").Append(truck)
	}
}
