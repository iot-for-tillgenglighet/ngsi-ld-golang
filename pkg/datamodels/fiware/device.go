package fiware

import (
	"strings"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Device is a Fiware entity
type Device struct {
	ngsi.BaseEntity
	Value                 *ngsi.TextProperty       `json:"value"`
	DateLastValueReported *ngsi.DateTimeProperty   `json:"dateLastValueReported,omitempty"`
	DateCreated           *ngsi.DateTimeProperty   `json:"dateCreated,omitempty"`
	DateModified          *ngsi.DateTimeProperty   `json:"dateModified,omitempty"`
	Location              *ngsi.GeoJSONProperty    `json:"location,omitempty"`
	RefDeviceModel        *DeviceModelRelationship `json:"refDeviceModel,omitempty"`
}

//NewDevice creates a new Device from given ID and Value
func NewDevice(id string, value string) *Device {
	if strings.HasPrefix(id, DeviceIDPrefix) == false {
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

//DeviceRelationship stores information about an entity's relation to a certain Device
type DeviceRelationship struct {
	ngsi.Relationship
	Object string `json:"object"`
}

//CreateDeviceRelationshipFromDevice create a DeviceRelationship from a Device
func CreateDeviceRelationshipFromDevice(device string) *DeviceRelationship {
	if len(device) == 0 {
		return nil
	}

	if strings.HasPrefix(device, DeviceIDPrefix) == false {
		device = DeviceIDPrefix + device
	}

	return &DeviceRelationship{
		Relationship: ngsi.Relationship{Type: "Relationship"},
		Object:       device,
	}
}
