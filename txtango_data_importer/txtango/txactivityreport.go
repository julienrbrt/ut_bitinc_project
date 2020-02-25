package txtango

import (
	"encoding/xml"
	"time"
)

// http://integratorsprod.transics.com/Reporting/Get_ActivityReport.html

//getActivityReportTemplate implements Get_ActivityReport_V11
//the requests filters by vehicle using their transics_id
var getActivityReportTemplate = `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        <Get_ActivityReport_V11 xmlns="http://transics.org">
			{{ template "login" .Login}}
            <ActivityReportSelection>
                <Vehicles>
                    <IdentifierVehicle>
                        <IdentifierVehicleType>TRANSICS_ID</IdentifierVehicleType>
                        <Id>{{.VehicleTransicsID}}</Id>
                    </IdentifierVehicle>
				</Vehicles>
			<DateTimeRangeSelection>
				<DateTypeSelection>STARTED</DateTypeSelection>
				<StartDate>{{.StartDate}}</StartDate>
				<EndDate>{{.EndDate}}</EndDate>
			</DateTimeRangeSelection>
            </ActivityReportSelection>
        </Get_ActivityReport_V11>
    </soap:Body>
</soap:Envelope>
`

//GetActivityReportRequest implements Get_ActivityReport_V11
type GetActivityReportRequest struct {
	// every request must implement the login
	Login             Login
	VehicleTransicsID int
	StartDate         string
	EndDate           string
}

//GetActivityReportResponse parses the response from TX-TANGO
type GetActivityReportResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                         string `xml:",chardata"`
		GetActivityReportV11Response struct {
			Text                       string `xml:",chardata"`
			Xmlns                      string `xml:"xmlns,attr"`
			GetActivityReportV11Result struct {
				Text                string    `xml:",chardata"`
				Executiontime       string    `xml:"Executiontime,attr"`
				Errors              TXError   `xml:"Errors"`
				Warnings            TXWarning `xml:"Warnings"`
				ActivityReportItems struct {
					Text                  string `xml:",chardata"`
					ActivityReportItemV11 []struct {
						Text        string `xml:",chardata"`
						ID          string `xml:"ID"`
						WorkingCode struct {
							Text        string `xml:",chardata"`
							Code        string `xml:"Code"`
							Description string `xml:"Description"`
						} `xml:"WorkingCode"`
						Vehicle struct {
							Text          string `xml:",chardata"`
							ID            string `xml:"ID"`
							TransicsID    string `xml:"TransicsID"`
							Code          string `xml:"Code"`
							Filter        string `xml:"Filter"`
							LicensePlate  string `xml:"LicensePlate"`
							FormattedName string `xml:"FormattedName"`
						} `xml:"Vehicle"`
						Trailer struct {
							Text          string `xml:",chardata"`
							ID            string `xml:"ID"`
							TransicsID    string `xml:"TransicsID"`
							Code          string `xml:"Code"`
							Filter        string `xml:"Filter"`
							LicensePlate  string `xml:"LicensePlate"`
							FormattedName string `xml:"FormattedName"`
						} `xml:"Trailer"`
						BeginDate    string  `xml:"BeginDate"`
						EndDate      string  `xml:"EndDate"`
						KmBegin      int     `xml:"KmBegin"`
						KmEnd        int     `xml:"KmEnd"`
						Consumption  float32 `xml:"Consumption"`
						LoadedStatus string  `xml:"LoadedStatus"`
						Activity     struct {
							Text                  string `xml:",chardata"`
							ID                    string `xml:"ID"`
							Name                  string `xml:"Name"`
							IsPlanning            string `xml:"IsPlanning"`
							ActivityType          string `xml:"ActivityType"`
							InstructionSetVersion string `xml:"InstructionSetVersion"`
						} `xml:"Activity"`
						SpeedAvg  float32 `xml:"SpeedAvg"`
						Reference string  `xml:"Reference"`
						Position  struct {
							Text                        string  `xml:",chardata"`
							Longitude                   float32 `xml:"Longitude"`
							Latitude                    float32 `xml:"Latitude"`
							AddressInfo                 string  `xml:"AddressInfo"`
							DistanceFromCapitol         string  `xml:"DistanceFromCapitol"`
							DistanceFromLargeCity       string  `xml:"DistanceFromLargeCity"`
							DistanceFromSmallCity       string  `xml:"DistanceFromSmallCity"`
							DistanceFromPointOfInterest string  `xml:"DistanceFromPointOfInterest"`
							CountryCode                 string  `xml:"CountryCode"`
						} `xml:"Position"`
						ModificationDate string `xml:"ModificationDate"`
						RegistrationID   string `xml:"RegistrationID"`
						POI              struct {
							Text        string `xml:",chardata"`
							PoiID       string `xml:"PoiID"`
							Active      bool   `xml:"Active"`
							Name        string `xml:"Name"`
							StreetLine1 string `xml:"StreetLine1"`
							StreetLine2 string `xml:"StreetLine2"`
							StreetLine3 string `xml:"StreetLine3"`
							Number      string `xml:"Number"`
							POBox       string `xml:"POBox"`
							ZipCode     string `xml:"ZipCode"`
							City        string `xml:"City"`
							Country     string `xml:"Country"`
							Position    struct {
								Text      string `xml:",chardata"`
								Longitude string `xml:"Longitude"`
								Latitude  string `xml:"Latitude"`
							} `xml:"Position"`
						} `xml:"POI"`
						ModificationID string `xml:"ModificationID"`
						Soucre         string `xml:"Soucre"`
						Active         string `xml:"Active"`
						IsValidated    string `xml:"IsValidated"`
					} `xml:"ActivityReportItem_V11"`
				} `xml:"ActivityReportItems"`
				MaximumModificationID   string `xml:"MaximumModificationID"`
				MaximumModificationDate string `xml:"MaximumModificationDate"`
				IsMoreDataPresent       struct {
					Text string `xml:",chardata"`
					Nil  string `xml:"nil,attr"`
				} `xml:"IsMoreDataPresent"`
			} `xml:"Get_ActivityReport_V11Result"`
		} `xml:"Get_ActivityReport_V11Response"`
	} `xml:"Body"`
}

//GetActivityReport wraps SAOPCall to make a Get_ActivityReport_V11 request
//the date argument is used to get the report of a specific date
func GetActivityReport(vehicleTransicsID int, date time.Time) (*GetActivityReportResponse, error) {
	startDate := date.Format("2006-01-02")
	// add a day to find the end date
	endDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	//make an authenticated request
	params := &GetActivityReportRequest{
		Login:             *authenticate(),
		VehicleTransicsID: vehicleTransicsID,
		// parse the date to transics format
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := soapCall(params, "GetActivityReport", getActivityReportTemplate)
	if err != nil {
		return nil, err
	}

	//unmarshal json
	data := &GetActivityReportResponse{}
	err = xml.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
