package txtango

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
                        <Id>{.VehicleTransicsID}}</Id>
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
