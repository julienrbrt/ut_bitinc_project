//Package database handles the database migration
package database

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql" // driver mssql
)

//DB is the database object
var DB *gorm.DB

//ErrorDB specified an connection error to the database
var ErrorDB = "Connection error to the database"

//InitDB initialize the sql database
//We are using an GO ORM named GORM
func InitDB() error {
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	//encode the password so no failure in database connection
	dbPassword = url.QueryEscape(dbPassword)

	//build connection string
	dbURI := fmt.Sprintf("sqlserver://%s:%s@%s:1433?database=%s", dbUser, dbPassword, dbHost, dbName)

	//database connection
	conn, err := gorm.Open("mssql", dbURI)
	if err != nil {
		return err
	}

	DB = conn
	//Database migration
	DB.Debug().AutoMigrate(&Driver{}, &DriverEcoMonitorReport{}, &Truck{}, &TruckGroup{}, &TruckActivityReport{}, &Trailer{}, &Tour{}, &TourQueue{})

	return nil
}
