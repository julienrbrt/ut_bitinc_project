package txtango

//TXError parses errors message in response
type TXError struct {
	Text  string `xml:",chardata"`
	Error struct {
		Code            string `xml:"ErrorCode"`
		CodeExplenation string `xml:"ErrorCodeExplenation"`
		Field           string `xml:"Field"`
		Value           string `xml:"Value"`
	} `xml:"Error"`
}

//TXWarning parses warning message in response
type TXWarning struct {
	Text    string `xml:",chardata"`
	Warning struct {
		Code            string `xml:"WarningCode"`
		CodeExplenation string `xml:"WarningCodeExplenation"`
		Field           string `xml:"Field"`
		Value           string `xml:"Value"`
	} `xml:"Warning"`
}
