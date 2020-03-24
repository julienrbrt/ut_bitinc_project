package analysis

import (
	"time"
	"tx2db/database"

	"github.com/pkg/errors"
)

//DriverMetric defines a driver metric
type DriverMetric struct {
	TransicsID string
	Metric     string
}

//getDrivenKm gets the number of kilometers driven
func getDrivenKm(start, end time.Time) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, FLOOR(SUM(distance)) as metric 
	FROM driver_eco_monitor_reports demr
	INNER JOIN tours t
	ON demr.tour_id = t.id
	WHERE distance > 0
	AND t.start_time >= ?
	AND t.end_time <= ?
	GROUP BY demr.driver_transics_id
	ORDER BY demr.driver_transics_id asc`,
		start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getDriverName gets a driver name
func getDriverName(driversList []string) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT transics_id, name as metric
	FROM drivers
	WHERE transics_id IN (?)
	ORDER BY transics_id asc`,
		driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getPersonID gets a person-id
func getPersonID(driversList []string) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT transics_id, person_id as metric
	FROM drivers
	WHERE transics_id IN (?)
	ORDER BY transics_id asc`,
		driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getTruckDriven gets the trucks that a driver has been driving
func getTruckDriven(driversList []string, start, end time.Time) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT t.driver_transics_id as transics_id, trucks.license_plate as metric
	FROM tours t
	INNER JOIN trucks
	ON t.truck_transics_id = trucks.transics_id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND t.driver_transics_id IN (?)
	GROUP BY t.driver_transics_id, trucks.license_plate
	ORDER BY t.driver_transics_id asc`,
		start.Format("2006-01-02"), end.Format("2006-01-02"), driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getTotalPanicBrakes gets the number of panic brakes performed
func getTotalPanicBrakes(driversList []string, start, end time.Time) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(number_of_panic_brakes) as metric
	FROM driver_eco_monitor_reports demr
	INNER JOIN tours t
	ON demr.tour_id = t.id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND t.driver_transics_id IN (?)
	GROUP BY demr.driver_transics_id
	ORDER BY demr.driver_transics_id asc`,
		start.Format("2006-01-02"), end.Format("2006-01-02"), driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getVisitedCountries gets the country list where drivers have been
func getVisitedCountries(driversList []string, start, end time.Time) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT d.transics_id as transics_id, tar.country_code as metric
	FROM tours t
	INNER JOIN truck_activity_reports tar
	ON t.id = tar.tour_id
	INNER JOIN drivers d
	ON d.transics_id = t.driver_transics_id 
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND t.driver_transics_id IN (?)
	GROUP BY transics_id, tar.country_code
	ORDER BY transics_id asc`,
		start.Format("2006-01-02"), end.Format("2006-01-02"), driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getRollOutRatio gets ratio of rolling out
func getRollOutRatio(driversList []string, start, end time.Time) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(demr.distance_coasting) / SUM(demr.distance) as metric 
	FROM driver_eco_monitor_reports demr
	LEFT JOIN tours t
	ON demr.tour_id = t.id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND t.driver_transics_id IN (?)
	AND distance > 0
	GROUP BY demr.driver_transics_id
	ORDER BY demr.driver_transics_id asc`,
		start.Format("2006-01-02"), end.Format("2006-01-02"), driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getCruiseControlRatio gets ratio of cruise control usage
func getCruiseControlRatio(driversList []string, start, end time.Time) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(demr.distance_on_cruise_control) / SUM(demr.distance) as metric 
	FROM driver_eco_monitor_reports demr
	LEFT JOIN tours t
	ON demr.tour_id = t.id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND distance > 0
	AND t.driver_transics_id IN (?)
	GROUP BY demr.driver_transics_id
	ORDER BY demr.driver_transics_id asc`,
		start.Format("2006-01-02"), end.Format("2006-01-02"), driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}

//getFuelConsumption gets the fuel consumption of a driver
func getFuelConsumption(driversList []string, start, end time.Time) ([]DriverMetric, error) {
	var result []DriverMetric
	if err := database.DB.Raw(`
	SELECT demr.driver_transics_id as transics_id, SUM(fuel_consumption) as metric 
	FROM driver_eco_monitor_reports demr
	LEFT JOIN tours t
	ON demr.tour_id = t.id
	WHERE t.start_time >= ?
	AND t.end_time <= ?
	AND t.driver_transics_id IN (?)
	GROUP BY demr.driver_transics_id
	ORDER BY demr.driver_transics_id asc`,
		start.Format("2006-01-02"), end.Format("2006-01-02"), driversList).Scan(&result).Error; err != nil {
		return result, errors.Wrap(err, database.ErrorDB)
	}

	return result, nil
}
