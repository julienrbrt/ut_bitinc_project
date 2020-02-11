package txtango

import (
	"encoding/xml"
	"fmt"
)

var getVehiculeTemplate = `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
		<Get_Vehicles_V13 xmlns="http://transics.org">
			{{ template "login" .Login}}
            <VehicleSelection>
                <Identifiers />
				<IncludePosition>{{.IncludePosition}}</IncludePosition>
				<IncludeActivity>{{.IncludeActivity}}</IncludeActivity>
				<IncludeDrivers>{{.IncludeDrivers}}</IncludeDrivers>
				<IncludeObcInfo>{{.IncludeObcInfo}}</IncludeObcInfo>
				<IncludeETAInfo>{{.IncludeETAInfo}}</IncludeETAInfo>
				<IncludeTemperatureInfo>{{.IncludeTemperatureInfo}}</IncludeTemperatureInfo>
				<IncludeInfoFields>{{.IncludeInfoFields}}</IncludeInfoFields>
				<IncludeUpdateDates>{{.IncludeUpdateDates}}</IncludeUpdateDates>
				<IncludeInactive>{{.IncludeInactive}}</IncludeInactive>
				<IncludeCompanyCardInfo>{{.IncludeCompanyCardInfo}}</IncludeCompanyCardInfo>
				<IncludeVehicleProfile>{{.IncludeVehicleProfile}}</IncludeVehicleProfile>
				<IncludeNextStopInfo>{{.IncludeNextStopInfo}}</IncludeNextStopInfo>
				<DiagnosticFilter />
				<IncludeExtraTruckInfo>{{.IncludeExtraTruckInfo}}</IncludeExtraTruckInfo>
				<IncludeGroups>{{.IncludeGroups}}</IncludeGroups>
                <GenericFilter>
                    <GenericFilterItem />
                </GenericFilter>
            </VehicleSelection>
        </Get_Vehicles_V13>
    </soap:Body>
</soap:Envelope>
`

//GetVehicleRequest implements GetVehicle_13
type GetVehicleRequest struct {
	// every request must implement the login
	Login                  Login
	IncludePosition        bool
	IncludeActivity        bool
	IncludeDrivers         bool
	IncludeObcInfo         bool
	IncludeETAInfo         bool
	IncludeTemperatureInfo bool
	IncludeInfoFields      bool
	IncludeUpdateDates     bool
	IncludeInactive        bool
	IncludeCompanyCardInfo bool
	IncludeVehicleProfile  bool
	IncludeNextStopInfo    bool
	IncludeExtraTruckInfo  bool
	IncludeGroups          bool
}

//GetVehicleResponse parse the response from TX-TANGO -- written using https://www.onlinetool.io/xmltogo/
type GetVehicleResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                   string `xml:",chardata"`
		GetVehiclesV13Response struct {
			GetVehiclesV13Result struct {
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
							Statuses  struct {
								Text         string `xml:",chardata"`
								Company      string `xml:"Company"`
								Rental       string `xml:"Rental"`
								Loan         string `xml:"Loan"`
								Truck        string `xml:"Truck"`
								LightVehicle string `xml:"LightVehicle"`
								DigitalTacho string `xml:"DigitalTacho"`
								TractorUnit  string `xml:"TractorUnit"`
								Coach        string `xml:"Coach"`
							} `xml:"Statuses"`
						} `xml:"ExtraTruckInfo"`
						VehicleID           string `xml:"VehicleID"`
						VehicleExternalCode string `xml:"VehicleExternalCode"`
						LicensePlate        string `xml:"LicensePlate"`
						Inactive            string `xml:"Inactive"`
						CanBusConnection    struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"CanBusConnection"`
						Trailer struct {
							Text          string `xml:",chardata"`
							ID            string `xml:"ID"`
							TransicsID    string `xml:"TransicsID"`
							Code          string `xml:"Code"`
							Filter        string `xml:"Filter"`
							LicensePlate  string `xml:"LicensePlate"`
							FormattedName string `xml:"FormattedName"`
						} `xml:"Trailer"`
						VehicleProfile struct {
							Text               string `xml:",chardata"`
							ProfileName        string `xml:"ProfileName"`
							ProfileDescription string `xml:"ProfileDescription"`
							IsDefaultProfile   string `xml:"IsDefaultProfile"`
							ProfileType        string `xml:"ProfileType"`
							ProfileID          string `xml:"ProfileId"`
						} `xml:"VehicleProfile"`
						VehicleTransicsID string `xml:"VehicleTransicsID"`
						Modified          struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"Modified"`
						CurrentKms struct {
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
						Position struct {
							Text                  string `xml:",chardata"`
							Longitude             string `xml:"Longitude"`
							Latitude              string `xml:"Latitude"`
							AddressInfo           string `xml:"AddressInfo"`
							DistanceFromCapitol   string `xml:"DistanceFromCapitol"`
							DistanceFromLargeCity string `xml:"DistanceFromLargeCity"`
							DistanceFromSmallCity string `xml:"DistanceFromSmallCity"`
							CountryCode           string `xml:"CountryCode"`
							Heading               struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"Heading"`
							LocationSource              string `xml:"LocationSource"`
							DistanceFromPointOfInterest string `xml:"DistanceFromPointOfInterest"`
						} `xml:"Position"`
						Speed    string `xml:"Speed"`
						Activity struct {
							Text string `xml:",chardata"`
							ID   string `xml:"ID"`
							Name string `xml:"Name"`
						} `xml:"Activity"`
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
						ObcInfo struct {
							Text                  string `xml:",chardata"`
							SoftwareVersion       string `xml:"SoftwareVersion"`
							InstructionSetVersion string `xml:"InstructionSetVersion"`
							Planningfeatures      struct {
								Text           string `xml:",chardata"`
								EnableTrips    string `xml:"Enable_Trips"`
								EnableJobs     string `xml:"Enable_Jobs"`
								EnablePlaces   string `xml:"Enable_Places"`
								EnableProducts string `xml:"Enable_Products"`
							} `xml:"Planningfeatures"`
							UpdateProgress struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"UpdateProgress"`
							ResponseImage string `xml:"ResponseImage"`
							IBC           struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"IBC"`
							ModemChannelID string `xml:"ModemChannelID"`
							DeviceType     string `xml:"DeviceType"`
						} `xml:"ObcInfo"`
						ETAInfo struct {
							Text                string `xml:",chardata"`
							PositionDestination struct {
								Text      string `xml:",chardata"`
								Longitude string `xml:"Longitude"`
								Latitude  string `xml:"Latitude"`
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
						InfoFieldList struct {
							Text          string `xml:",chardata"`
							InfoFieldItem []struct {
								Text           string `xml:",chardata"`
								Name           string `xml:"Name"`
								Value          string `xml:"Value"`
								DateLastUpdate string `xml:"DateLastUpdate"`
							} `xml:"InfoFieldItem"`
						} `xml:"InfoFieldList"`
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

//GetVehicle wraps SAOPCall to make a GetVehicle_13 request
func GetVehicle(req *GetVehicleRequest) (*GetVehicleResponse, error) {
	// add authentication to request
	req.Login = *authenticate()
	resp, err := soapCall(&req, "GetVehicle", getVehiculeTemplate)
	if err != nil {
		return nil, err
	}

	//unmarshal json
	data := &GetVehicleResponse{}
	err = xml.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	fmt.Println(data.Body.GetVehiclesV13Response.GetVehiclesV13Result.Errors.Error.Code)

	return data, nil
}
