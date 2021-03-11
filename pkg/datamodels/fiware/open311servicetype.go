package fiware

import (
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Open311ServiceType is a ...
type Open311ServiceType struct {
	ngsi.BaseEntity
	Description ngsi.TextProperty   `json:"description"`
	ServiceCode ngsi.NumberProperty `json:"service_code"`
}

//NewOpen311ServiceType creates a new Open311ServiceType
func NewOpen311ServiceType(label string, reportType string) *Open311ServiceType {

	id := Open311ServiceTypeIDPrefix + label + ":" + reportType

	return &Open311ServiceType{
		Description: *ngsi.NewTextProperty(label),
		ServiceCode: *ngsi.NewNumberPropertyFromString(reportType),
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
