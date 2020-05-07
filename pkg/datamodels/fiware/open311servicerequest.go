package fiware

import (
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//WeatherObserved is an observation of weather conditions at a certain place and time.
type Open311ServiceRequest struct {
	ngsi.BaseEntity
	requested_datetime *ngsi.DateTimeProperty `json:"requested_datetime,omitempty"`
	Location           ngsi.GeoJSONProperty   `json:"location"`
	service_code       ngsi.NumberProperty    `json:"service_code"`
}

//NewWeatherObserved creates a new instance of WeatherObserved
func NewOpen311ServiceRequest(latitude float64, longitude float64, reportedType int, reportedTimestamp string) *Open311ServiceRequest {
	dateTimeValue := ngsi.CreateDateTimeProperty(reportedTimestamp)

	id := "urn:ngsi-ld:Open311ServiceRequest:" + reportedTimestamp 

	return &Open311ServiceRequest{
		requested_datetime : dateTimeValue,
		Location:     ngsi.CreateGeoJSONPropertyFromWGS84(latitude, longitude),
		service_code: *ngsi.NewNumberPropertyFromInt(reportedType)
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "Open311ServiceRequest",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
