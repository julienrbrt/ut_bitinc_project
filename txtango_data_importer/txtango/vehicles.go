package txtango

var getVehiculeTemplate = `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        <Get_Vehicles_V13 xmlns="http://transics.org">
			{{ template "login" }}
            <VehicleSelection>
                <Identifiers />
				<IncludePosition>{{.IncludePosition}}<IncludePosition>
				<IncludeActivity>{{.IncludeActivity}}<IncludeActivity>
				<IncludeDrivers>{{.IncludeDrivers}}<IncludeDrivers>
				<IncludeObcInfo>{{.IncludeObcInfo}}<IncludeObcInfo>
				<IncludeETAInfo>{{.IncludeETAInfo}}<IncludeETAInfo>
				<IncludeTemperatureInfo>{{.IncludeTemperatureInfo}}<IncludeTemperatureInfo>
				<IncludeInfoFields>{{.IncludeInfoFields}}<IncludeInfoFields>
				<IncludeUpdateDates>{{.IncludeUpdateDates}}<IncludeUpdateDates>
				<IncludeInactive>{{.IncludeInactive}}<IncludeInactive>
				<IncludeLastTextMessageInboxOutbox>{{.IncludeLastTextMessageInboxOutbox}}<IncludeLastTextMessageInboxOutbox>
				<IncludeLastAlarmMessage>{{.IncludeLastAlarmMessage}}<IncludeLastAlarmMessage>
				<IncludeBlockedVehicleInfo>{{.IncludeBlockedVehicleInfo}}<IncludeBlockedVehicleInfo>
				<IncludeCompanyCardInfo>{{.IncludeCompanyCardInfo}}<IncludeCompanyCardInfo>
				<IncludeVehicleProfile>{{.IncludeVehicleProfile}}<IncludeVehicleProfile>
				<IncludeNextStopInfo>{{.IncludeNextStopInfo}}<IncludeNextStopInfo>
				<DiagnosticFilter />
				<IncludeExtraTruckInfo>{{.IncludeExtraTruckInfo}}<IncludeExtraTruckInfo>
				<IncludeGroups>{{.IncludeGroups}}</IncludeGroups>
				<IncludeRented>{{.IncludeRented}}</IncludeRented>
                <IncludePlanningInfo xsi:nil="true" />
                <GenericFilter>
                    <GenericFilterItem />
                </GenericFilter>
            </VehicleSelection>
        </Get_Vehicles_V13>
    </soap:Body>
</soap:Envelope>
`

//GetVehicle implements GetVehicle_13 from TX-TANGO
type GetVehicle struct {
}
