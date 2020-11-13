package fiware

import (
	"strings"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//RoadSegment is a Fiware entity
type RoadSegment struct {
	ngsi.BaseEntity
	Name            *ngsi.TextProperty       `json:"name"`
	DateCreated     *ngsi.DateTimeProperty   `json:"dateCreated,omitempty"`
	DateModified    *ngsi.DateTimeProperty   `json:"dateModified,omitempty"`
	Location        ngsi.RoadSegmentLocation `json:"location,omitempty"`
	EndPoint        ngsi.GeoJSONProperty     `json:"endPoint"`
	StartPoint      ngsi.GeoJSONProperty     `json:"startPoint"`
	RefRoad         *ngsi.RoadRelationship   `json:"refRoad,omitempty"`
	TotalLaneNumber *ngsi.NumberProperty     `json:"totalLaneNumber"`
}

//NewRoadSegment creates a new instance of RoadSegment
func NewRoadSegment(id, roadSegmentName, roadID string, coords [][2]float64) *RoadSegment {
	if strings.HasPrefix(id, "urn:ngsi-ld:RoadSegment:") == false {
		id = "urn:ngsi-ld:RoadSegment:" + id
	}

	if strings.HasPrefix(roadID, "urn:ngsi-ld:Road:") == false {
		roadID = "urn:ngsi-ld:Road:" + roadID
	}

	name := ngsi.NewTextProperty(roadSegmentName)
	refRoad := ngsi.NewRoadRelationship(roadID)
	startPoint := coords[0]
	endPoint := coords[len(coords)-1]

	return &RoadSegment{
		Name:            name,
		EndPoint:        ngsi.CreateGeoJSONPropertyFromWGS84(endPoint[0], endPoint[1]),
		StartPoint:      ngsi.CreateGeoJSONPropertyFromWGS84(startPoint[0], startPoint[1]),
		RefRoad:         refRoad,
		Location:        ngsi.NewRoadSegmentLocation(coords),
		TotalLaneNumber: ngsi.NewNumberPropertyFromInt(1),
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
