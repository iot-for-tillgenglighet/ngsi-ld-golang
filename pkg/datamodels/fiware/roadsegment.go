package fiware

import (
	"strings"
	"time"

	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//RoadSurfaceType contains a surface type and a probability
type RoadSurfaceType struct {
	ngsi.TextProperty
	Probability float64 `json:"probability"`
}

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
	SurfaceType     *RoadSurfaceType         `json:"surfaceType,omitempty"`
}

//NewRoadSegment creates a new instance of RoadSegment
func NewRoadSegment(id, roadSegmentName, roadID string, coords [][2]float64, modified *time.Time) *RoadSegment {
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

	rs := &RoadSegment{
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

	if modified != nil {
		rs.DateModified = ngsi.CreateDateTimeProperty(modified.Format(time.RFC3339))
	}

	return rs
}

//WithSurfaceType takes a string surfaceType and a probability and returns the road segment instance
func (rs *RoadSegment) WithSurfaceType(surfaceType string, probability float64) *RoadSegment {

	rs.SurfaceType = &RoadSurfaceType{
		TextProperty: ngsi.TextProperty{
			Property: ngsi.Property{Type: "Property"},
			Value:    surfaceType,
		},
		Probability: probability,
	}

	return rs
}
