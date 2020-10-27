package fiware

import (
	"strings"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Device is a Fiware entity
type Device struct {
	ngsi.BaseEntity
	Value          *ngsi.TextProperty     `json:"value"`
	DateCreated    *ngsi.DateTimeProperty `json:"dateCreated,omitempty"`
	DateModified   *ngsi.DateTimeProperty `json:"dateModified,omitempty"`
	Location       *ngsi.GeoJSONProperty  `json:"location,omitempty"`
	RefDeviceModel *ngsi.Relationship     `json:"refDeviceModel,omitempty"`
}

//NewDevice creates a new Device from given ID and Value
func NewDevice(id string, value string) *Device {
	if strings.HasPrefix(id, "urn:ngsi-ld:Device:") == false {
		id = "urn:ngsi-ld:Device:" + id
	}

	return &Device{
		Value: ngsi.NewTextProperty(value),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "Device",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
