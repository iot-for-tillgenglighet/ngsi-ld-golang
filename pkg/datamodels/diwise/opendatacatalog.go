package diwise

import (
	"strings"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//OpenDataCatalog is a Diwise entity
type OpenDataCatalog struct {
	ngsi.BaseEntity
	CatalogId   *ngsi.TextProperty             `json:"id"`
	Title       *ngsi.TextProperty             `json:"title"`
	Description *ngsi.TextProperty             `json:"description"`
	Publisher   *ngsi.TextProperty             `json:"publisher"`
	License     *ngsi.TextProperty             `json:"license"`
	Dataset     *ngsi.SingleObjectRelationship `json:"dataset"`
}

//NewOpenDataCatalog creates a new instance of OpenDataCatalog
func NewOpenDataCatalog(id, title, description, publisher, license, dataset string) *OpenDataCatalog {
	if !strings.HasPrefix(id, OpenDataCatalogIDPrefix) {
		id = OpenDataCatalogIDPrefix + id
	}

	return &OpenDataCatalog{
		Title:       ngsi.NewTextProperty(title),
		Description: ngsi.NewTextProperty(description),
		Publisher:   ngsi.NewTextProperty(publisher),
		License:     ngsi.NewTextProperty(license),
		Dataset:     ngsi.NewSingleObjectRelationship(dataset),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "OpenDataCatalog",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}

}
