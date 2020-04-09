package txtango

import (
	"encoding/xml"
	"time"
)

// http://integratorsprod.transics.com/Eco/Get_EcoMonitor_Report.html

//getEcoReport implements Get_EcoMonitor_Report_V4
//the requests filters by driver using their transics_id
var getEcoReport = `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        <Get_EcoMonitor_Report_V4 xmlns="http://transics.org">
			{{ template "login" .Login}}
            <EcoMonitorReportSelection>
                <Drivers>
                    <Identifier>
                        <IdentifierType>TRANSICS_ID</IdentifierType>
                        <Id>{{.DriverTransicsID}}</Id>
                    </Identifier>
                </Drivers>
                <DateTimeRangeSelection>
                    <StartDate>{{.StartDate}}</StartDate>
                    <EndDate>{{.EndDate}}</EndDate>
                </DateTimeRangeSelection>
                <IncludeRecordsWithoutDriver>false</IncludeRecordsWithoutDriver>
            </EcoMonitorReportSelection>
        </Get_EcoMonitor_Report_V4>
    </soap:Body>
</soap:Envelope>`

//GetEcoReportRequest implements Get_EcoMonitor_Report_V4
type GetEcoReportRequest struct {
	// every request must implement the login
	Login            Login
	DriverTransicsID uint
	StartDate        string
	EndDate          string
}

//GetEcoReportResponse parses the response from Transics
type GetEcoReportResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                          string `xml:",chardata"`
		GetEcoMonitorReportV4Response struct {
			Text                        string `xml:",chardata"`
			Xmlns                       string `xml:"xmlns,attr"`
			GetEcoMonitorReportV4Result struct {
				Text                  string    `xml:",chardata"`
				Executiontime         string    `xml:"Executiontime,attr"`
				Errors                TXError   `xml:"Errors"`
				Warnings              TXWarning `xml:"Warnings"`
				EcoMonitorReportItems struct {
					Text                   string `xml:",chardata"`
					EcoMonitorReportItemV3 []struct {
						Text            string `xml:",chardata"`
						TransicsID      uint   `xml:"TransicsID"`
						Scope           string `xml:"Scope"`
						IsConfidentData string `xml:"IsConfidentData"`
						TripReference   string `xml:"TripReference"`
						Vehicle         struct {
							Text         string `xml:",chardata"`
							ID           string `xml:"ID"`
							TransicsID   uint   `xml:"TransicsID"`
							Code         string `xml:"Code"`
							Filter       string `xml:"Filter"`
							LicensePlate string `xml:"LicensePlate"`
						} `xml:"Vehicle"`
						Trainer string `xml:"Trainer"`
						Driver  struct {
							Text       string `xml:",chardata"`
							ID         string `xml:"ID"`
							TransicsID uint   `xml:"TransicsID"`
							Code       string `xml:"Code"`
							Filter     string `xml:"Filter"`
							LastName   string `xml:"LastName"`
							FirstName  string `xml:"FirstName"`
						} `xml:"Driver"`
						BeginDate  string `xml:"BeginDate"`
						EndDate    string `xml:"EndDate"`
						DataResult struct {
							Text                   string  `xml:",chardata"`
							Distance               float32 `xml:"Distance"`
							Duration               float32 `xml:"Duration"`
							DurationDriving        float32 `xml:"DurationDriving"`
							FuelConsumption        float32 `xml:"FuelConsumption"`
							FuelConsumptionAverage struct {
								Text float32 `xml:",chardata"`
								Nil  string  `xml:"nil,attr"`
							} `xml:"FuelConsumptionAverage"`
							RpmAverage         float32 `xml:"RpmAverage"`
							Co2EmissionAverage struct {
								Text float32 `xml:",chardata"`
								Nil  string  `xml:"nil,attr"`
							} `xml:"Co2EmissionAverage"`
							SpeedAverage float32 `xml:"SpeedAverage"`
						} `xml:"DataResult"`
						IdlingResult struct {
							Text                     string  `xml:",chardata"`
							NumberOfLongIdling       int     `xml:"NumberOfLongIdling"`
							FuelConsumptionIdling    float32 `xml:"FuelConsumptionIdling"`
							DurationIdling           float32 `xml:"DurationIdling"`
							DurationIdlingPercentage struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"DurationIdlingPercentage"`
						} `xml:"IdlingResult"`
						OverSpeedingResult struct {
							Text                 string  `xml:",chardata"`
							DurationOverSpeeding float32 `xml:"DurationOverSpeeding"`
							NumberOfOverSpeeding int     `xml:"NumberOfOverSpeeding"`
						} `xml:"OverSpeedingResult"`
						CoastingResult struct {
							Text                string  `xml:",chardata"`
							DistanceCoasting    float32 `xml:"DistanceCoasting"`
							DurationCoasting    float32 `xml:"DurationCoasting"`
							DistanceEcoRollInKm struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"DistanceEcoRollInKm"`
							DurationEcoRollInSec struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"DurationEcoRollInSec"`
							DistanceEcoRollInPercentage struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"DistanceEcoRollInPercentage"`
						} `xml:"CoastingResult"`
						AnticipationResult struct {
							Text                       string  `xml:",chardata"`
							NumberOfStops              int     `xml:"NumberOfStops"`
							NumberOfBrakes             int     `xml:"NumberOfBrakes"`
							NumberOfPanicBrakes        int     `xml:"NumberOfPanicBrakes"`
							DistanceByBrakes           float32 `xml:"DistanceByBrakes"`
							DurationByBrakes           float32 `xml:"DurationByBrakes"`
							DurationByRetarder         float32 `xml:"DurationByRetarder"`
							DurationHighRPMnoFuel      float32 `xml:"DurationHighRPMnoFuel"`
							DurationHighRPM            float32 `xml:"DurationHighRPM"`
							DistanceByRetarder         float32 `xml:"DistanceByRetarder"`
							DistanceHighRPMnoFuel      float32 `xml:"DistanceHighRPMnoFuel"`
							NumberOfHarshAccelerations int     `xml:"NumberOfHarshAccelerations"`
							DurationHarshAcceleration  float32 `xml:"DurationHarshAcceleration"`
						} `xml:"AnticipationResult"`
						GreenSpotResult struct {
							Text                     string  `xml:",chardata"`
							DistanceGreenSpot        float32 `xml:"DistanceGreenSpot"`
							DurationGreenSpot        float32 `xml:"DurationGreenSpot"`
							FuelConsumptionGreenSpot float32 `xml:"FuelConsumptionGreenSpot"`
						} `xml:"GreenSpotResult"`
						GearingResult struct {
							Text                      string  `xml:",chardata"`
							NumberOfGearChanges       int     `xml:"NumberOfGearChanges"`
							NumberOfGearChangesUp     int     `xml:"NumberOfGearChangesUp"`
							PositionOfThrottleAverage float32 `xml:"PositionOfThrottleAverage"`
							PositionOfThrottleMaximum float32 `xml:"PositionOfThrottleMaximum"`
						} `xml:"GearingResult"`
						PtoResult struct {
							Text                         string  `xml:",chardata"`
							NumberOfPto                  int     `xml:"NumberOfPto"`
							FuelConsumptionPtoDriving    float32 `xml:"FuelConsumptionPtoDriving"`
							FuelConsumptionPtoStandStill float32 `xml:"FuelConsumptionPtoStandStill"`
							DurationPtoDriving           float32 `xml:"DurationPtoDriving"`
							DurationPtoStandStill        float32 `xml:"DurationPtoStandStill"`
						} `xml:"PtoResult"`
						CruisingResult struct {
							Text                                           string  `xml:",chardata"`
							DistanceOnCruiseControl                        float32 `xml:"DistanceOnCruiseControl"`
							DurationOnCruiseControl                        float32 `xml:"DurationOnCruiseControl"`
							DistanceOnCruiseControlPercentage              float32 `xml:"DistanceOnCruiseControlPercentage"`
							AvgFuelConsumptionCruiseControlInLiterPer100km float32 `xml:"AvgFuelConsumptionCruiseControlInLiterPer100km"`
							AvgFuelConsumptionCruiseControlInkmPerLiter    float32 `xml:"AvgFuelConsumptionCruiseControlInkmPerLiter"`
						} `xml:"CruisingResult"`
					} `xml:"EcoMonitorReportItem_V3"`
				} `xml:"EcoMonitorReportItems"`
			} `xml:"Get_EcoMonitor_Report_V4Result"`
		} `xml:"Get_EcoMonitor_Report_V4Response"`
	} `xml:"Body"`
}

//GetEcoReport wraps SAOPCall to make a Get_EcoMonitor_Report_V4 request
//the date argument is used to get the report of a specific date
func GetEcoReport(driverTransicsID uint, start, end time.Time) (*GetEcoReportResponse, error) {
	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")

	//make an authenticated request
	params := &GetEcoReportRequest{
		Login:            *authenticate(),
		DriverTransicsID: driverTransicsID,
		// parse the date to transics format
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := soapCall(params, "GetEcoReport", getEcoReport)
	if err != nil {
		return nil, err
	}

	//unmarshal json
	data := &GetEcoReportResponse{}
	err = xml.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
