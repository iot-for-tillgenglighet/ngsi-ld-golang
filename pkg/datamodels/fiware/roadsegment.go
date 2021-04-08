package fiware

import (
	"strings"
	"time"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
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
	Name            *ngsi.TextProperty             `json:"name"`
	DateCreated     *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified    *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	Location        ngsi.RoadSegmentLocation       `json:"location,omitempty"`
	EndPoint        geojson.GeoJSONProperty        `json:"endPoint"`
	StartPoint      geojson.GeoJSONProperty        `json:"startPoint"`
	RefRoad         *ngsi.SingleObjectRelationship `json:"refRoad,omitempty"`
	TotalLaneNumber *ngsi.NumberProperty           `json:"totalLaneNumber"`
	SurfaceType     *RoadSurfaceType               `json:"surfaceType,omitempty"`
}

//NewRoadSegment creates a new instance of RoadSegment
func NewRoadSegment(id, roadSegmentName, roadID string, coords [][2]float64, modified *time.Time) *RoadSegment {
	if !strings.HasPrefix(id, RoadSegmentIDPrefix) {
		id = RoadSegmentIDPrefix + id
	}

	if !strings.HasPrefix(roadID, RoadIDPrefix) {
		roadID = RoadIDPrefix + roadID
	}

	name := ngsi.NewTextProperty(roadSegmentName)
	refRoad := ngsi.NewSingleObjectRelationship(roadID)
	startPoint := coords[0]
	endPoint := coords[len(coords)-1]

	rs := &RoadSegment{
		Name:            name,
		EndPoint:        *geojson.CreateGeoJSONPropertyFromWGS84(endPoint[0], endPoint[1]),
		StartPoint:      *geojson.CreateGeoJSONPropertyFromWGS84(startPoint[0], startPoint[1]),
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

	if len(surfaceType) > 0 && probability >= 0.0 && probability <= 1.0 {
		rs.SurfaceType = &RoadSurfaceType{
			TextProperty: ngsi.TextProperty{
				Property: ngsi.Property{Type: "Property"},
				Value:    surfaceType,
			},
			Probability: probability,
		}
	}

	return rs
}
