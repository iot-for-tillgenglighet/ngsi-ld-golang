package fiware

import (
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//WeatherObserved is an observation of weather conditions at a certain place and time.
type Open311ServiceType struct {
	ngsi.BaseEntity
	description  ngsi.TextProperty   `json:"description"`
	service_code ngsi.NumberProperty `json:"service_code"`
}

//NewWeatherObserved creates a new instance of WeatherObserved
func NewOpen311ServiceType(label string, reportType string) *Open311ServiceType {

	id := "urn:ngsi-ld:Open311ServiceType:" + label + ":" + reportType

	return &Open311ServiceType{
		description:  *ngsi.NewTextProperty(label),
		service_code: *ngsi.NewNumberPropertyFromString(reportType),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "Open311ServiceType",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
