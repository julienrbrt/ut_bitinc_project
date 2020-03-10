package database

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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
	ImportOn   time.Time
	Reason     string
	Trial      int //number of time the element of the queue has been tried to be imported
}

//ImportQueuedToursData imports the data from the queue
func ImportQueuedToursData(handleError bool) error {
	var queue []TourQueue

	//get only element from queue where the import_on date is older than 3 days and older date first
	db.Where("? > import_on", time.Now().AddDate(0, 0, -3)).Order("import_on asc").Find(&queue)

	log.Println("Checking & importing tour from queue")
	for i, data := range queue {
		var tour Tour

		log.Printf("(%d / %d) Checking & importing tour from queue\n", i+1, len(queue))
		err := db.Model(&tour).Where("id = ?", data.TourID).First(&tour).Error
		if err != nil {
			//if a tour of the queue cannot be gotten, skip it
			continue
		}

		//caculate elapsed time between last import and day to import
		diff := int(tour.LastImport.Sub(data.ImportOn).Hours() / 24)

		switch data.ReportType {
		case emr:
			err = importEcoMoniorReport(&tour, diff)
		case tar:
			err = importActivityReport(&tour, diff)
		}
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			if handleError {
				return err
			}
		}

		// if data.NbTrial == new.NbTrial
		var newQueue TourQueue
		err = db.Model(&TourQueue{}).Where(TourQueue{TourID: data.TourID}).First(&newQueue).Error
		if err != nil {
			return err
		}

		//element of the queue has been fetched, remove it permanently
		if data.Trial == newQueue.Trial {
			db.Unscoped().Where(data).Delete(&TourQueue{})
		}

	}

	return nil
}

//addTourToQueue add tours into the tour queue
func addTourToQueue(tour *Tour, importOn time.Time, reportType, reason string) error {
	var tourQueue TourQueue
	data := &TourQueue{
		TourID:     tour.ID,
		ReportType: reportType,
		ImportOn:   importOn,
		Reason:     reason,
	}

	if err := db.Model(&tourQueue).Where(TourQueue{TourID: tour.ID}).First(&tourQueue).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, errDatabaseConnection)
		}

		//add tour to queue
		db.Create(&data)
	} else {
		//update the number of trial for that queue
		tourQueue.Trial = tourQueue.Trial + 1
		db.Save(&tourQueue)
	}

	return nil
}
