package txtango

import (
	"encoding/xml"
)

//http://integratorsprod.transics.com/Administration/Get_Drivers.html

//getDriversTemplate implements Get_Drivers_V9
var getDriversTemplate = `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Get_Drivers_V9 xmlns="http://transics.org">
	  {{ template "login" .Login}}
      <DriverSelection>
		<IncludeContactInfo>false</IncludeContactInfo>
		<IncludeGroups>false</IncludeGroups>
		<IncludeInactiveDrivers>true</IncludeInactiveDrivers>
		<IncludeLastVehicleInfo>false</IncludeLastVehicleInfo>
		<IncludeLicenseInfo>false</IncludeLicenseInfo>
		<IncludeTachoCardInfo>true</IncludeTachoCardInfo>
		<IncludeUpdateDates>true</IncludeUpdateDates>
		<IncludeHRInfo>false</IncludeHRInfo>
      </DriverSelection>
    </Get_Drivers_V9>
  </soap:Body>
</soap:Envelope>
`

//GetDriversRequest implements Get_Drivers_V9
type GetDriversRequest struct {
	// every request must implement the login
	Login Login
}

//GetDriversResponse parses the response from TX-TANGO
type GetDriversResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                 string `xml:",chardata"`
		GetDriversV9Response struct {
			Text               string `xml:",chardata"`
			Xmlns              string `xml:"xmlns,attr"`
			GetDriversV9Result struct {
				Text          string    `xml:",chardata"`
				Executiontime string    `xml:"Executiontime,attr"`
				Errors        TXError   `xml:"Errors"`
				Warnings      TXWarning `xml:"Warnings"`
				Persons       struct {
					Text                    string `xml:",chardata"`
					InterfacePersonResultV9 []struct {
						Text            string `xml:",chardata"`
						SmartCardNumber struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"SmartCardNumber"`
						PersonID           string `xml:"PersonId"`
						PersonExternalCode string `xml:"PersonExternalCode"`
						Lastname           string `xml:"Lastname"`
						Firstname          string `xml:"Firstname"`
						Filter             string `xml:"Filter"`
						Inactive           bool   `xml:"Inactive"`
						Description        string `xml:"Description"`
						Languages          struct {
							Text                      string `xml:",chardata"`
							WorkingLanguage           string `xml:"WorkingLanguage"`
							ObcUILanguage             string `xml:"ObcUILanguage"`
							ObcInstructionSetLanguage string `xml:"ObcInstructionSetLanguage"`
						} `xml:"Languages"`
						PersonTransicsID string `xml:"PersonTransicsId"` // important as identifier used thorough the code
						UpdateDatesList  struct {
							Text            string `xml:",chardata"`
							UpdateDatesItem struct {
								Text           string `xml:",chardata"`
								Name           string `xml:"Name"`
								DateLastUpdate string `xml:"DateLastUpdate"`
							} `xml:"UpdateDatesItem"`
						} `xml:"UpdateDatesList"`
						Modification struct {
							Text string `xml:",chardata"`
							Nil  string `xml:"nil,attr"`
						} `xml:"Modification"`
						TachoCardInfo struct {
							Text           string `xml:",chardata"`
							CardID         string `xml:"CardId"`
							CountryOfIssue struct {
								Text        string `xml:",chardata"`
								CountryCode string `xml:"CountryCode"`
								CountryName string `xml:"CountryName"`
							} `xml:"CountryOfIssue"`
							TachoLoginEnabled bool   `xml:"TachoLoginEnabled"`
							RenewalIndex      string `xml:"RenewalIndex"`
							ReplacementIndex  string `xml:"ReplacementIndex"`
							StartOfValidity   struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"StartOfValidity"`
							EndOfValidity struct {
								Text string `xml:",chardata"`
								Nil  string `xml:"nil,attr"`
							} `xml:"EndOfValidity"`
							RespectedDeadLine string `xml:"RespectedDeadLine"`
						} `xml:"TachoCardInfo"`
						FormattedName string `xml:"FormattedName"`
					} `xml:"InterfacePersonResult_V9"`
				} `xml:"Persons"`
			} `xml:"Get_Drivers_V9Result"`
		} `xml:"Get_Drivers_V9Response"`
	} `xml:"Body"`
}

//GetDrivers wraps SAOPCall to make a Get_Drivers_V9 request
func GetDrivers() (*GetDriversResponse, error) {
	//make an authenticated request
	params := &GetDriversRequest{
		Login: *authenticate(),
	}
	resp, err := soapCall(params, "GetDrivers", getDriversTemplate)
	if err != nil {
		return nil, err
	}

	//unmarshal json
	data := &GetDriversResponse{}
	err = xml.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
