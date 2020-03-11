package analysis

import (
	"time"
	"tx2db/database"

	"github.com/pkg/errors"
)

//DriverMetric defines a driver metric
type DriverMetric struct {
	TransicsID uint
	Metric     string
}

//getTotalPanicBrakes gets the number of panic brakes performed
func getTotalPanicBrakes(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(number_of_panic_brakes) as metric
	FROM driver_eco_monitor_reports demr
	INNER JOIN tours t
	ON demr.tour_id = t.id
	WHERE number_of_panic_brakes > 0
	AND t.start_time >= ?
	AND t.end_time <= ?
	GROUP BY demr.driver_transics_id`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}

//getDrivenKm gets the number of kilometers driven
func getDrivenKm(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, FLOOR(SUM(distance)) as metric 
	FROM driver_eco_monitor_reports demr
	INNER JOIN tours t
	ON demr.tour_id = t.id
	WHERE distance > 0
	AND t.start_time >= ?
	AND t.end_time <= ?
	GROUP BY demr.driver_transics_id`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}

//getVisitedCountries gets the country list where drivers have been
func getVisitedCountries(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT d.transics_id as transics_id, tar.country_code as metric
	FROM tours t
	INNER JOIN truck_activity_reports tar
	ON t.id = tar.tour_id
	INNER JOIN drivers d
	ON d.transics_id = t.driver_transics_id 
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	GROUP BY transics_id, tar.country_code`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}

//getRollOutRatio gets ratio of rolling out
func getRollOutRatio(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(demr.distance_coasting) / SUM(demr.distance) as metric 
	FROM driver_eco_monitor_reports demr
	LEFT JOIN tours t
	ON demr.tour_id = t.id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND distance > 0
	GROUP BY demr.driver_transics_id`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}

//getCruiseControlRatio gets ratio of cruise control usage
func getCruiseControlRatio(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(demr.distance_on_cruise_control) / SUM(demr.distance) as metric 
	FROM driver_eco_monitor_reports demr
	LEFT JOIN tours t
	ON demr.tour_id = t.id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND distance > 0
	GROUP BY demr.driver_transics_id`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}

//getConsumption gets the consumption of a driver
func getConsumption(start time.Time) ([]DriverMetric, error) {
	//build the endtime
	end := getReportEndTime(start)

	var driversMetrics []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(fuel_consumption) as metric 
	FROM driver_eco_monitor_reports demr
	LEFT JOIN tours t
	ON demr.tour_id = t.id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	GROUP BY demr.driver_transics_id`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&driversMetrics).Error; err != nil {
		return driversMetrics, errors.Wrap(err, database.ErrorDB)
	}

	return driversMetrics, nil
}
