package fiware

import (
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//WeatherObserved is an observation of weather conditions at a certain place and time.
type WeatherObserved struct {
	ngsi.BaseEntity
	DateCreated  *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObserved ngsi.DateTimeProperty          `json:"dateObserved"`
	Location     ngsi.GeoJSONProperty           `json:"location"`
	RefDevice    *ngsi.SingleObjectRelationship `json:"refDevice,omitempty"`
	SnowHeight   *ngsi.NumberProperty           `json:"snowHeight,omitempty"`
	Temperature  *ngsi.NumberProperty           `json:"temperature,omitempty"`
}

//NewWeatherObserved creates a new instance of WeatherObserved
func NewWeatherObserved(device string, latitude float64, longitude float64, observedAt string) *WeatherObserved {
	dateTimeValue := ngsi.CreateDateTimeProperty(observedAt)
	refDevice := CreateDeviceRelationshipFromDevice(device)

	if refDevice == nil {
		device = "manual"
	}

	id := WeatherObservedIDPrefix + device + ":" + observedAt

	return &WeatherObserved{
		DateObserved: *dateTimeValue,
		Location:     *ngsi.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		RefDevice:    refDevice,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "WeatherObserved",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
