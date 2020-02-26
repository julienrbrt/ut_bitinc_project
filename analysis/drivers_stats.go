package analysis

import (
	"time"
)

//reportDuration is the duration of a report
const monthlyDuration = 30 * 24 * time.Hour
const weeklyDuration = 7 * 24 * time.Hour

//DriverKmDrove calculte how many km a driver drove in the report duration
func DriverKmDrove(duration time.Duration) (float32, error) {

	return 0, nil
}
