package fiware

import (
	"strings"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Road is a Fiware entity
type Road struct {
	ngsi.BaseEntity
	Name           *ngsi.TextProperty           `json:"name"`
	RoadClass      *ngsi.TextProperty           `json:"roadClass"`
	RefRoadSegment ngsi.RoadSegmentRelationship `json:"refRoadSegment"`
}

//NewRoad creates a new instance of Road
func NewRoad(id string, roadName string, roadClass string, roadSegmentIdentities []string) *Road {
	if strings.HasPrefix(id, "urn:ngsi-ld:Road:") == false {
		id = "urn:ngsi-ld:Road:" + id
	}

	return &Road{
		Name:           ngsi.NewTextProperty(roadName),
		RoadClass:      ngsi.NewTextProperty(roadClass),
		RefRoadSegment: ngsi.NewRoadSegmentRelationship(roadSegmentIdentities),
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "RoadSegment",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
