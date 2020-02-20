//Package database handles the better-driving database migration
package database

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql" // driver mssql
)

var db *gorm.DB

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

	db = conn
	//Database migration
	db.Debug().AutoMigrate(&Driver{}, &DriverEcoMonitorReport{}, &Truck{}, &TruckGroup{}, &TruckActivityReport{}, &Tour{})

	return nil
}

//DB returns a handle to the DB object
func DB() *gorm.DB {
	return db
}
