//Package bolk handles the better-driving database migration
package bolk

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql" // driver mssql
)

var db *gorm.DB

//InitDB initialize the database
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
	defer conn.Close()
	if err != nil {
		return err
	}

	db = conn
	//Database migration
	db.Debug().AutoMigrate(&Driver{}, &Truck{}, &Tour{})

	return nil
}

//DB returns a handle to the DB object
func DB() *gorm.DB {
	return db
}
