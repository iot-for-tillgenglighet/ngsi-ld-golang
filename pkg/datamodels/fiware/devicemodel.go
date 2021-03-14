package fiware

import (
	"strings"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//DeviceModel is a Fiware entity
type DeviceModel struct {
	ngsi.BaseEntity
	Category           *ngsi.TextListProperty `json:"category"`
	ModelName          *ngsi.TextProperty     `json:"modelName,omitempty"`
	Name               *ngsi.TextProperty     `json:"name,omitempty"`
	BrandName          *ngsi.TextProperty     `json:"brandName,omitempty"`
	ManufacturerName   *ngsi.TextProperty     `json:"manufacturerName,omitempty"`
	ControlledProperty *ngsi.TextListProperty `json:"controlledProperty,omitempty"`
	DateCreated        *ngsi.DateTimeProperty `json:"dateCreated,omitempty"`
	DateModified       *ngsi.DateTimeProperty `json:"dateModified,omitempty"`
}

//NewDeviceModel creates a new DeviceModel from given ID and Value
func NewDeviceModel(id string, categories []string) *DeviceModel {
	if !strings.HasPrefix(id, DeviceModelIDPrefix) {
		id = DeviceModelIDPrefix + id
	}

	return &DeviceModel{
		Category: ngsi.NewTextListProperty(categories),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "DeviceModel",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

//NewDeviceModelRelationship creates a single object relationship to a DeviceModelID
func NewDeviceModelRelationship(deviceModelID string) (*ngsi.SingleObjectRelationship, error) {
	if !strings.HasPrefix(deviceModelID, DeviceModelIDPrefix) {
		deviceModelID = DeviceModelIDPrefix + deviceModelID
	}

	rel := ngsi.NewSingleObjectRelationship(deviceModelID)
	return rel, nil
}
