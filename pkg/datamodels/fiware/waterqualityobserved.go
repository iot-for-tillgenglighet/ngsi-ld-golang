package fiware

import (
	"bytes"
	"encoding/json"

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

type waterQualityDTO struct {
	ngsi.BaseEntity
	DateCreated        *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified       *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	DateObserved       ngsi.DateTimeProperty          `json:"dateObserved"`
	Location           json.RawMessage                `json:"location"`
	RefDevice          *ngsi.SingleObjectRelationship `json:"refDevice,omitempty"`
	RefPointOfInterest *ngsi.SingleObjectRelationship `json:"refPointOfInterested,omitempty"`
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

func (wqo WaterQualityObserved) ToGeoJSONFeature(propertyName string, simplified bool) (geojson.GeoJSONFeature, error) {
	g := geojson.NewGeoJSONFeature(wqo.ID, wqo.Type, wqo.Location.GeoPropertyValue())

	if simplified {
		g.SetProperty(propertyName, wqo.Location.GeoPropertyValue())
		g.SetProperty("dateObserved", wqo.DateObserved.Value.Value)

		if wqo.Temperature != nil {
			g.SetProperty("temperature", wqo.Temperature.Value)
		}

		if wqo.RefDevice != nil {
			g.SetProperty("refDevice", wqo.RefDevice.Object)
		}

		if wqo.RefPointOfInterest != nil {
			g.SetProperty("refPointOfInterest", wqo.RefPointOfInterest.Object)
		}
	} else {
		g.SetProperty(propertyName, wqo.Location)
		g.SetProperty("dateObserved", wqo.DateObserved)

		g.SetProperty("temperature", wqo.Temperature)
		g.SetProperty("refDevice", wqo.RefDevice)
		g.SetProperty("refPointOfInterest", wqo.RefPointOfInterest)
	}

	return g, nil
}

func (wqo *WaterQualityObserved) UnmarshalJSON(data []byte) error {
	dto := &waterQualityDTO{}
	err := json.NewDecoder(bytes.NewReader(data)).Decode(dto)

	if err == nil {
		wqo.ID = dto.ID
		wqo.Type = dto.Type

		wqo.DateCreated = dto.DateCreated
		wqo.DateModified = dto.DateModified
		wqo.DateObserved = dto.DateObserved
		wqo.RefDevice = dto.RefDevice
		wqo.RefPointOfInterest = dto.RefPointOfInterest
		wqo.Temperature = dto.Temperature

		wqo.Context = dto.Context

		wqo.Location = *geojson.CreateGeoJSONPropertyFromJSON(dto.Location)

	}

	return err
}
