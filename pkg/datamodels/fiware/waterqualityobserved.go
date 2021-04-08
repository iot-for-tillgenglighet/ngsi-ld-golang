package fiware

import (
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//WaterQualityObserved is intended to represent water quality parameters at a certain water mass (river, lake, sea, etc.) section.
type WaterQualityObserved struct {
	ngsi.BaseEntity
	DateCreated        *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified       *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObserved       ngsi.DateTimeProperty          `json:"dateObserved"`
	Location           geojson.GeoJSONProperty        `json:"location"`
	RefDevice          *ngsi.SingleObjectRelationship `json:"refDevice,omitempty"`
	RefPointOfInterest *ngsi.SingleObjectRelationship `json:"refPointOfInterest,omitempty"`
	Temperature        *ngsi.NumberProperty           `json:"temperature,omitempty"`
}

//NewWaterQualityObserved creates a new instance of WaterQualityObserved
func NewWaterQualityObserved(device string, latitude float64, longitude float64, observedAt string) *WaterQualityObserved {
	dateTimeValue := ngsi.CreateDateTimeProperty(observedAt)
	refDevice := CreateDeviceRelationshipFromDevice(device)

	if refDevice == nil {
		device = "manual"
	}

	id := WaterQualityObservedIDPrefix + device + ":" + observedAt

	return &WaterQualityObserved{
		DateObserved: *dateTimeValue,
		Location:     *geojson.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		RefDevice:    refDevice,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "WaterQualityObserved",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
