package fiware

import (
	"strings"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Device is a Fiware entity
type Device struct {
	ngsi.BaseEntity
	Value                 *ngsi.TextProperty             `json:"value"`
	DateLastValueReported *ngsi.DateTimeProperty         `json:"dateLastValueReported,omitempty"`
	DateCreated           *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified          *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	Location              *geojson.GeoJSONProperty       `json:"location,omitempty"`
	RefDeviceModel        *ngsi.SingleObjectRelationship `json:"refDeviceModel,omitempty"`
}

//NewDevice creates a new Device from given ID and Value
func NewDevice(id string, value string) *Device {
	if !strings.HasPrefix(id, DeviceIDPrefix) {
		id = DeviceIDPrefix + id
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

//CreateDeviceRelationshipFromDevice create a DeviceRelationship from a Device
func CreateDeviceRelationshipFromDevice(device string) *ngsi.SingleObjectRelationship {
	if len(device) == 0 {
		return nil
	}

	if !strings.HasPrefix(device, DeviceIDPrefix) {
		device = DeviceIDPrefix + device
	}

	return ngsi.NewSingleObjectRelationship(device)
}
