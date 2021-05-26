package fiware

import (
	"strings"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	ngsi "github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//TrafficFlowObserved is a Fiware entity
type TrafficFlowObserved struct {
	ngsi.BaseEntity
	DateObserved   *ngsi.DateTimeProperty         `json:"dateObserved,omitempty"`
	DateCreated    *ngsi.DateTimeProperty         `json:"dateCreated,omitempty"`
	DateModified   *ngsi.DateTimeProperty         `json:"dateModified,omitempty"`
	Location       geojson.GeoJSONProperty        `json:"location,omitempty"`
	LaneID         *ngsi.NumberProperty           `json:"laneID,omitempty"`
	RefRoadSegment *ngsi.SingleObjectRelationship `json:"refRoadSegment,omitempty"`
}

//NewTrafficFlowObserved creates a new TrafficFlowObserved from given ID
func NewTrafficFlowObserved(id string, latitude float64, longitude float64, observedAt string, laneID int) *TrafficFlowObserved {
	if !strings.HasPrefix(id, TrafficFlowObservedIDPrefix) {
		id = TrafficFlowObservedIDPrefix + id
	}

	dateTimeValue := ngsi.CreateDateTimeProperty(observedAt)
	lane := ngsi.NewNumberPropertyFromInt(laneID)

	return &TrafficFlowObserved{
		DateObserved: dateTimeValue,
		Location:     *geojson.CreateGeoJSONPropertyFromWGS84(longitude, latitude),
		LaneID:       lane,
		BaseEntity: ngsi.BaseEntity{
			ID:   id,
			Type: "TrafficFlowObserved",
			Context: []string{
				"https://schema.lab.fiware.org/ld/context",
				"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
			},
		},
	}
}
