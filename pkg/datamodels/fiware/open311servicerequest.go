package fiware

import (
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Open311ServiceRequest is a ...
type Open311ServiceRequest struct {
	ngsi.BaseEntity
	RequestedDateTime *ngsi.DateTimeProperty  `json:"requested_datetime,omitempty"`
	Location          geojson.GeoJSONProperty `json:"location"`
	ServiceCode       ngsi.NumberProperty     `json:"service_code"`
}

//NewOpen311ServiceRequest creates a new service request
func NewOpen311ServiceRequest(latitude float64, longitude float64, reportedType int, reportedTimestamp string) *Open311ServiceRequest {
	dateTimeValue := ngsi.CreateDateTimeProperty(reportedTimestamp)

	id := Open311ServiceRequestIDPrefix + reportedTimestamp

	return &Open311ServiceRequest{
		RequestedDateTime: dateTimeValue,
		Location:          *geojson.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		ServiceCode:       *ngsi.NewNumberPropertyFromInt(reportedType),
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
