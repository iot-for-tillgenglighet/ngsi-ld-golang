package fiware

import (
	"strings"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Beach is a Fiware entity
type Beach struct {
	ngsi.BaseEntity
	Name         *ngsi.TextProperty     `json:"name,omitempty"`
	Description  *ngsi.TextProperty     `json:"description"`
	Location     ngsi.GeoJSONProperty   `json:"location,omitempty"`
	DateCreated  *ngsi.DateTimeProperty `json:"dateCreated,omitempty"`
	DateModified *ngsi.DateTimeProperty `json:"dateModified,omitempty"`
}

//NewBeach creates a new Beach from given ID and name
func NewBeach(id, name string, lat, lon float64) *Beach {
	if strings.HasPrefix(id, BeachIDPrefix) == false {
		id = BeachIDPrefix + id
	}

	return &Beach{
		Name:     ngsi.NewTextProperty(name),
		Location: ngsi.CreateGeoJSONPropertyFromWGS84(lon, lat),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "Beach",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

//WithDescription adds a text property named Deescription to this Beach instance
func (b *Beach) WithDescription(description string) *Beach {
	b.Description = ngsi.NewTextProperty(description)
	return b
}
