package diwise

import (
	"fmt"
	"strings"

	fiware "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/datamodels/fiware"
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//RoadSurfaceObserved is a Diwise entity
type RoadSurfaceObserved struct {
	ngsi.BaseEntity
	Location       ngsi.GeoJSONProperty          `json:"location"`
	SurfaceType    fiware.RoadSurfaceType        `json:"surfaceType"`
	RefRoadSegment *ngsi.MultiObjectRelationship `json:"refRoadSegment,omitempty"`
	DateObserved   *ngsi.DateTimeProperty        `json:"dateObserved,omitempty"`
}

//NewRoadSurfaceObserved creates a new instance of RoadSurfaceObserved
func NewRoadSurfaceObserved(id string, surfaceType string, probability float64, latitude float64, longitude float64) *RoadSurfaceObserved {
	if !strings.HasPrefix(id, RoadSurfaceObservedIDPrefix) {
		id = RoadSurfaceObservedIDPrefix + id
	}

	return &RoadSurfaceObserved{
		Location: ngsi.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		SurfaceType: fiware.RoadSurfaceType{
			TextProperty: ngsi.TextProperty{
				Property: ngsi.Property{Type: "Property"},
				Value:    surfaceType,
			},
			Probability: probability,
		},
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "RoadSurfaceObserved",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}

//WithRoadSegment creates a reference to a road segment.
func (rso *RoadSurfaceObserved) WithRoadSegment(segmentID string) (*RoadSurfaceObserved, error) {

	if strings.HasPrefix(segmentID, fiware.RoadSegmentIDPrefix) {
		relationship := ngsi.NewMultiObjectRelationship([]string{segmentID})
		rso.RefRoadSegment = &relationship
	} else {
		return nil, fmt.Errorf("unable to create a RoadSegmentRelationship from invalid segment ID: %s", segmentID)
	}

	return rso, nil
}
