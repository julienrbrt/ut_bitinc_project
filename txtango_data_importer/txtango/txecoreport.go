package txtango

import (
	"encoding/xml"
	"time"
)

// http://integratorsprod.transics.com/Eco/Get_EcoPerformance.html

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
	DriverTransicsID int
	StartDate        string
	EndDate          string
}

//GetEcoReportResponse parses the response from TX-TANGO
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
						TransicsID      string `xml:"TransicsID"`
						Scope           string `xml:"Scope"`
						IsConfidentData string `xml:"IsConfidentData"`
						TripReference   string `xml:"TripReference"`
						Vehicle         struct {
							Text         string `xml:",chardata"`
							ID           string `xml:"ID"`
							TransicsID   string `xml:"TransicsID"`
							Code         string `xml:"Code"`
							Filter       string `xml:"Filter"`
							LicensePlate string `xml:"LicensePlate"`
						} `xml:"Vehicle"`
						Trainer string `xml:"Trainer"`
						Driver  struct {
							Text       string `xml:",chardata"`
							ID         string `xml:"ID"`
							TransicsID string `xml:"TransicsID"`
							Code       string `xml:"Code"`
							Filter     string `xml:"Filter"`
							LastName   string `xml:"LastName"`
							FirstName  string `xml:"FirstName"`
						} `xml:"Driver"`
						BeginDate  string `xml:"BeginDate"`
						EndDate    string `xml:"EndDate"`
						DataResult struct {
							Text                   string `xml:",chardata"`
							Distance               string `xml:"Distance"`
							Duration               string `xml:"Duration"`
							DurationDriving        string `xml:"DurationDriving"`
							FuelConsumption        string `xml:"FuelConsumption"`
							FuelConsumptionAverage struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"FuelConsumptionAverage"`
							RpmAverage         string `xml:"RpmAverage"`
							Co2EmissionAverage struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"Co2EmissionAverage"`
							SpeedAverage string `xml:"SpeedAverage"`
						} `xml:"DataResult"`
						IdlingResult struct {
							Text                     string `xml:",chardata"`
							NumberOfLongIdling       string `xml:"NumberOfLongIdling"`
							FuelConsumptionIdling    string `xml:"FuelConsumptionIdling"`
							DurationIdling           string `xml:"DurationIdling"`
							DurationIdlingPercentage struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"DurationIdlingPercentage"`
						} `xml:"IdlingResult"`
						OverSpeedingResult struct {
							Text                 string `xml:",chardata"`
							DurationOverSpeeding string `xml:"DurationOverSpeeding"`
							NumberOfOverSpeeding string `xml:"NumberOfOverSpeeding"`
						} `xml:"OverSpeedingResult"`
						CoastingResult struct {
							Text                string `xml:",chardata"`
							DistanceCoasting    string `xml:"DistanceCoasting"`
							DurationCoasting    string `xml:"DurationCoasting"`
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
							Text                       string `xml:",chardata"`
							NumberOfStops              string `xml:"NumberOfStops"`
							NumberOfBrakes             string `xml:"NumberOfBrakes"`
							NumberOfPanicBrakes        string `xml:"NumberOfPanicBrakes"`
							DistanceByBrakes           string `xml:"DistanceByBrakes"`
							DurationByBrakes           string `xml:"DurationByBrakes"`
							DurationByRetarder         string `xml:"DurationByRetarder"`
							DurationHighRPMnoFuel      string `xml:"DurationHighRPMnoFuel"`
							DurationHighRPM            string `xml:"DurationHighRPM"`
							DistanceByRetarder         string `xml:"DistanceByRetarder"`
							DistanceHighRPMnoFuel      string `xml:"DistanceHighRPMnoFuel"`
							NumberOfHarshAccelerations string `xml:"NumberOfHarshAccelerations"`
							DurationHarshAcceleration  string `xml:"DurationHarshAcceleration"`
						} `xml:"AnticipationResult"`
						GreenSpotResult struct {
							Text                     string `xml:",chardata"`
							DistanceGreenSpot        string `xml:"DistanceGreenSpot"`
							DurationGreenSpot        string `xml:"DurationGreenSpot"`
							FuelConsumptionGreenSpot string `xml:"FuelConsumptionGreenSpot"`
						} `xml:"GreenSpotResult"`
						GearingResult struct {
							Text                      string `xml:",chardata"`
							NumberOfGearChanges       string `xml:"NumberOfGearChanges"`
							NumberOfGearChangesUp     string `xml:"NumberOfGearChangesUp"`
							PositionOfThrottleAverage string `xml:"PositionOfThrottleAverage"`
							PositionOfThrottleMaximum string `xml:"PositionOfThrottleMaximum"`
						} `xml:"GearingResult"`
						PtoResult struct {
							Text                         string `xml:",chardata"`
							NumberOfPto                  string `xml:"NumberOfPto"`
							FuelConsumptionPtoDriving    string `xml:"FuelConsumptionPtoDriving"`
							FuelConsumptionPtoStandStill string `xml:"FuelConsumptionPtoStandStill"`
							DurationPtoDriving           string `xml:"DurationPtoDriving"`
							DurationPtoStandStill        string `xml:"DurationPtoStandStill"`
						} `xml:"PtoResult"`
						CruisingResult struct {
							Text                                           string `xml:",chardata"`
							DistanceOnCruiseControl                        string `xml:"DistanceOnCruiseControl"`
							DurationOnCruiseControl                        string `xml:"DurationOnCruiseControl"`
							DistanceOnCruiseControlPercentage              string `xml:"DistanceOnCruiseControlPercentage"`
							AvgFuelConsumptionCruiseControlInLiterPer100km string `xml:"AvgFuelConsumptionCruiseControlInLiterPer100km"`
							AvgFuelConsumptionCruiseControlInkmPerLiter    string `xml:"AvgFuelConsumptionCruiseControlInkmPerLiter"`
						} `xml:"CruisingResult"`
					} `xml:"EcoMonitorReportItem_V3"`
				} `xml:"EcoMonitorReportItems"`
			} `xml:"Get_EcoMonitor_Report_V4Result"`
		} `xml:"Get_EcoMonitor_Report_V4Response"`
	} `xml:"Body"`
}

//GetEcoReport wraps SAOPCall to make a Get_EcoMonitor_Report_V4 request
//the date argument is used to get the report of a specific date
func GetEcoReport(driverTransicsID int, date time.Time) (*GetEcoReportResponse, error) {
	startDate := date.Format("2006-01-02")
	// add a day to find the enddate
	endDate := date.Add(time.Hour * 24).Format("2006-01-02")

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
