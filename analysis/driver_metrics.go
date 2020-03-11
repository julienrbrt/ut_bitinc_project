package analysis

import (
	"time"
	"tx2db/database"

	"github.com/pkg/errors"
)

//DriverMetric defines a driver metric
type DriverMetric struct {
	TransicsID uint
	Metric     int
}

//getPanicBrakesNb gets the number of panic brakes performed in a specific time range
func getPanicBrakesNb(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT d.transics_id, SUM(number_of_panic_brakes) as metric 
	FROM driver_eco_monitor_reports demr
	INNER JOIN tours t
	ON demr.tour_id = t.id
	INNER JOIN drivers d
	ON t.driver_transics_id = d.transics_id 
	WHERE number_of_panic_brakes > 0
	AND t.start_time >= ?
	AND t.end_time <= ?
	GROUP BY transics_id`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}

//getDrivenKm gets the number of kilometers driven in a specific time range
func getDrivenKm(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT d.transics_id, FLOOR(SUM(distance)) as metric 
	FROM driver_eco_monitor_reports demr
	INNER JOIN tours t
	ON demr.tour_id = t.id
	INNER JOIN drivers d
	on t.driver_transics_id = d.transics_id 
	WHERE distance > 0
	AND t.start_time >= ?
	AND t.end_time <= ?
	GROUP BY transics_id`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}
