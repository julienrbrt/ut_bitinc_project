package database

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"tx2db/txtango"

	"github.com/pkg/errors"
)

var (
	loadingDataFromTransics = "Loading data from Transics TX-TANGO... this could take a while..."
	errParsingTransicsID    = "Error when parsing TransicsID"
	errParsingDate          = "Error while parsing date from Transics TX-TANGO"
)

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

	for i, data := range txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Persons.InterfacePersonResultV9 {
		transicsID, err := strconv.Atoi(data.PersonTransicsID)
		if err != nil {
			return errors.Wrap(err, errParsingTransicsID)
		}

		driver := Driver{
			TransicsID: transicsID,
			Name:       data.FormattedName,
			Language:   data.Languages.WorkingLanguage,
		}

		//add driver
		fmt.Printf("(%d / %d) Importing driver %d\n", i, len(txDrivers.Body.GetDriversV9Response.GetDriversV9Result.Persons.InterfacePersonResultV9), driver.TransicsID)
		db.FirstOrCreate(&driver)
	}

	return nil
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

	for i, data := range txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Vehicles.InterfaceVehicleResultV13 {
		transicsID, err := strconv.Atoi(data.VehicleTransicsID)
		if err != nil {
			return errors.Wrap(err, errParsingTransicsID)
		}

		modifiedDate, err := time.Parse(time.RFC3339, data.Modified)
		if err != nil {
			return errors.Wrap(err, errParsingDate)
		}

		truck := Truck{
			TransicsID:   transicsID,
			LicensePlate: data.LicensePlate,
			Inactive:     data.Inactive,
			LastModified: modifiedDate,
		}

		//add driver
		fmt.Printf("(%d / %d) Importing truck %d\n", i, len(txVehicle.Body.GetVehiclesV13Response.GetVehiclesV13Result.Vehicles.InterfaceVehicleResultV13), truck.TransicsID)
		db.FirstOrCreate(&truck)
	}

	return nil
}
