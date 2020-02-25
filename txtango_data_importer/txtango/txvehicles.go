package txtango

import (
	"encoding/xml"
)

// http://integratorsprod.transics.com/Administration/Get_Vehicles.html

//getVehiculeTemplate implements Get_Vehicles_V13
var getVehiculeTemplate = `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
		<Get_Vehicles_V13 xmlns="http://transics.org">
			{{ template "login" .Login}}
            <VehicleSelection>
				<IncludePosition>false</IncludePosition>
				<IncludeActivity>false</IncludeActivity>
				<IncludeDrivers>true</IncludeDrivers>
				<IncludeObcInfo>false</IncludeObcInfo>
				<IncludeETAInfo>true</IncludeETAInfo>
				<IncludeTemperatureInfo>true</IncludeTemperatureInfo>
				<IncludeInfoFields>false</IncludeInfoFields>
				<IncludeUpdateDates>true</IncludeUpdateDates>
				<IncludeInactive>true</IncludeInactive>
				<IncludeCompanyCardInfo>true</IncludeCompanyCardInfo>
				<IncludeVehicleProfile>false</IncludeVehicleProfile>
				<IncludeNextStopInfo>true</IncludeNextStopInfo>
				<IncludeExtraTruckInfo>true</IncludeExtraTruckInfo>
				<IncludeGroups>true</IncludeGroups>
            </VehicleSelection>
        </Get_Vehicles_V13>
    </soap:Body>
</soap:Envelope>
`

//GetVehicleRequest implements Get_Vehicles_V13
type GetVehicleRequest struct {
	// every request must implement the login
	Login Login
}

//GetVehicleResponse parses the response from TX-TANGO -- written using https://www.onlinetool.io/xmltogo/
type GetVehicleResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                   string `xml:",chardata"`
		GetVehiclesV13Response struct {
			Text                 string `xml:",chardata"`
			Xmlns                string `xml:"xmlns,attr"`
			GetVehiclesV13Result struct {
				Text          string    `xml:",chardata"`
				Executiontime string    `xml:"Executiontime,attr"`
				Errors        TXError   `xml:"Errors"`
				Warnings      TXWarning `xml:"Warnings"`
				Vehicles      struct {
					Text                      string `xml:",chardata"`
					InterfaceVehicleResultV13 []struct {
						Text               string `xml:",chardata"`
						VehicleFleetNumber string `xml:"VehicleFleetNumber"`
						Groups             struct {
							Text            string `xml:",chardata"`
							TxConnectGroups struct {
								Text          string `xml:",chardata"`
								ConnectGroups struct {
									Text         string `xml:",chardata"`
									ConnectGroup []struct {
										Text     string `xml:",chardata"`
										Group    string `xml:"Group"`
										SubGroup string `xml:"SubGroup"`
									} `xml:"ConnectGroup"`
								} `xml:"ConnectGroups"`
							} `xml:"TxConnectGroups"`
						} `xml:"Groups"`
						ExtraTruckInfo struct {
							Text    string `xml:",chardata"`
							Country string `xml:"Country"`
							InDuty  struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"InDuty"`
							OutDuty struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"OutDuty"`
							VinNumber string `xml:"VinNumber"`
							Category  string `xml:"Category"`
						} `xml:"ExtraTruckInfo"`
						VehicleID           string `xml:"VehicleID"`
						VehicleExternalCode string `xml:"VehicleExternalCode"`
						LicensePlate        string `xml:"LicensePlate"`
						Inactive            bool   `xml:"Inactive"`
						CanBusConnection    struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"CanBusConnection"`
						Trailer           TXTrailer `xml:"Trailer"`
						VehicleTransicsID string    `xml:"VehicleTransicsID"`
						Modified          string    `xml:"Modified"`
						CurrentKms        struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"CurrentKms"`
						FuelLevel struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"FuelLevel"`
						FuelLevelIndex struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"FuelLevelIndex"`
						RefrigeratorIndex struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"RefrigeratorIndex"`
						Speed             string `xml:"Speed"`
						ActivityCompleted struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"ActivityCompleted"`
						Driver struct {
							Text          string `xml:",chardata"`
							ID            string `xml:"ID"`
							TransicsID    string `xml:"TransicsID"`
							Code          string `xml:"Code"`
							Filter        string `xml:"Filter"`
							LastName      string `xml:"LastName"`
							FirstName     string `xml:"FirstName"`
							FormattedName string `xml:"FormattedName"`
						} `xml:"Driver"`
						ETAInfo struct {
							Text                string `xml:",chardata"`
							PositionDestination struct {
								Text      string  `xml:",chardata"`
								Longitude float32 `xml:"Longitude"`
								Latitude  float32 `xml:"Latitude"`
							} `xml:"PositionDestination"`
							PrevETA struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"PrevETA"`
							ETAStatus struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"ETAStatus"`
							DistanceETA struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"DistanceETA"`
							ETA struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"ETA"`
							PositionInfoDestination string `xml:"PositionInfoDestination"`
							EtaRestIncluded         struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"EtaRestIncluded"`
						} `xml:"ETAInfo"`
						UpdateDatesList struct {
							Text            string `xml:",chardata"`
							UpdateDatesItem []struct {
								Text           string `xml:",chardata"`
								Name           string `xml:"Name"`
								DateLastUpdate struct {
									Text string `xml:",chardata"`
									Nil  string `xml:"nil,attr"`
								} `xml:"DateLastUpdate"`
							} `xml:"UpdateDatesItem"`
						} `xml:"UpdateDatesList"`
						Maintenance     string `xml:"Maintenance"`
						Remaining       string `xml:"Remaining"`
						AutoFilter      string `xml:"AutoFilter"`
						CompanyCardInfo struct {
							Text     string `xml:",chardata"`
							CardID   string `xml:"CardId"`
							CardName string `xml:"CardName"`
						} `xml:"CompanyCardInfo"`
						FormattedName   string `xml:"FormattedName"`
						LastTrailerCode string `xml:"LastTrailerCode"`
					} `xml:"InterfaceVehicleResult_V13"`
				} `xml:"Vehicles"`
			} `xml:"Get_Vehicles_V13Result"`
		} `xml:"Get_Vehicles_V13Response"`
	} `xml:"Body"`
}

//TXTrailer represent a trailer from GetVehicleResponse
type TXTrailer struct {
	Text          string `xml:",chardata"`
	ID            string `xml:"ID"`
	TransicsID    string `xml:"TransicsID"`
	Code          string `xml:"Code"`
	Filter        string `xml:"Filter"`
	LicensePlate  string `xml:"LicensePlate"`
	FormattedName string `xml:"FormattedName"`
}

//GetVehicle wraps SAOPCall to make a Get_Vehicles_V13 request
func GetVehicle() (*GetVehicleResponse, error) {
	//make an authenticated request
	params := &GetVehicleRequest{
		Login: *authenticate(),
	}
	resp, err := soapCall(params, "GetVehicle", getVehiculeTemplate)
	if err != nil {
		return nil, err
	}

	//unmarshal json
	data := &GetVehicleResponse{}
	err = xml.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
