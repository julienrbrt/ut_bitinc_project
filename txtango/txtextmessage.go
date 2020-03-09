package txtango

import (
	"encoding/xml"
)

//http://integratorsprod.transics.com/Messaging/Send_TextMessage.html

//sendTextMessageTemplate implements Send_TextMessage
var sendTextMessageTemplate = `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
  <Send_TextMessage xmlns="http://transics.org">
  	{{ template "login" .Login}}
	<TextMessageSend>
	  <Vehicles>
		<IdentifierVehicle>
		  <IdentifierVehicleType>TRANSICS_ID</IdentifierVehicleType>
		  <Id>{{.VehicleTransicsID}}</Id>
		</IdentifierVehicle>
	  </Vehicles>
	  <VehicleType>NONE</VehicleType>
	  <ForceOBCWakeUp>true</ForceOBCWakeUp>
	  <Message>{{.Text}}</Message>
	</TextMessageSend>
  </Send_TextMessage>
</soap:Body>
</soap:Envelope>`

//SentTextMessageRequest implements Send_TextMessage
type SentTextMessageRequest struct {
	// every request must implement the login
	Login             Login
	VehicleTransicsID uint
	Message           string
}

//SentTextMessageResponse parses the response from Transics
type SentTextMessageResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                    string `xml:",chardata"`
		SendTextMessageResponse struct {
			Text                  string `xml:",chardata"`
			Xmlns                 string `xml:"xmlns,attr"`
			SendTextMessageResult struct {
				Text                       string    `xml:",chardata"`
				Executiontime              string    `xml:"Executiontime,attr"`
				Errors                     TXError   `xml:"Errors"`
				Warnings                   TXWarning `xml:"Warnings"`
				SendTextMessageResultInfos struct {
					Text                      string `xml:",chardata"`
					SendTextMessageResultInfo struct {
						Text    string `xml:",chardata"`
						Vehicle struct {
							Text         string `xml:",chardata"`
							ID           string `xml:"ID"`
							TransicsID   string `xml:"TransicsID"`
							Code         string `xml:"Code"`
							Filter       string `xml:"Filter"`
							LicensePlate string `xml:"LicensePlate"`
						} `xml:"Vehicle"`
						MessageID string `xml:"MessageId"`
					} `xml:"SendTextMessageResultInfo"`
				} `xml:"SendTextMessageResultInfos"`
			} `xml:"Send_TextMessageResult"`
		} `xml:"Send_TextMessageResponse"`
	} `xml:"Body"`
}

//SendMessage wraps SAOPCall to make a Send_TextMessage request
func SendMessage(vehicleTransicsID uint, text string) (*SentTextMessageResponse, error) {
	//make an authenticated request
	params := &SentTextMessageRequest{
		Login:             *authenticate(),
		VehicleTransicsID: vehicleTransicsID,
		Message:           text,
	}
	resp, err := soapCall(params, "SendMessage", sendTextMessageTemplate)
	if err != nil {
		return nil, err
	}

	//unmarshal json
	data := &SentTextMessageResponse{}
	err = xml.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
