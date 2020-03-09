package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	reasonQueueNoData = "No data found during import"
	tar               = "tar" // Truck Activity Report
	emr               = "emr" // Eco Monitor Report
)

//TourQueue represents the database tour queue
type TourQueue struct {
	gorm.Model
	TourID     uint
	ReportType string // should only be tar or emr
	ImportFrom time.Time
	Reason     string
}
