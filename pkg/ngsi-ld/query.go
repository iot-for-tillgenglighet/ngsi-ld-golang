package ngsi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//Query is an interface to be used when passing queries to context registries and sources
type Query interface {
	HasDeviceReference() bool
	Device() string

	IsGeoQuery() bool
	Geo() *GeoQuery

	EntityAttributes() []string
	EntityTypes() []string

	Request() *http.Request
}

const (
	//GeoSpatialRelationNearPoint describes a relation as a max or min distance from a Point
	GeoSpatialRelationNearPoint = "near"
)

//GeoQuery contains information about a geo-query that may be used for subscriptions
//or when querying entitites
type GeoQuery struct {
	Geometry    string    `json:"geometry"`
	Coordinates []float64 `json:"coordinates"`
	GeoRel      string    `json:"georel"`
	GeoProperty *string   `json:"geoproperty,omitempty"`

	distance uint32
}

//Distance returns the required distance in meters from a near Point and a boolean
//flag indicating if it is inclusive or exclusive
func (gq *GeoQuery) Distance() (uint32, bool) {
	return gq.distance, true
}

//Point extracts the first two coordinates of the enclosed geometry
func (gq *GeoQuery) Point() (float64, float64, error) {
	if len(gq.Coordinates) == 2 {
		return gq.Coordinates[0], gq.Coordinates[1], nil
	}

	return 0, 0, errors.New("invalid number of coordinates in GeoQuery for a Point geometry")
}

func newQueryFromParameters(req *http.Request, types []string, attributes []string, q string) (Query, error) {

	var err error

	const refDevicePrefix string = "refDevice==\""

	qw := &queryWrapper{request: req, types: types, attributes: attributes}

	if strings.HasPrefix(q, refDevicePrefix) {
		splitElems := strings.Split(q, "\"")
		qw.device = &splitElems[1]
	}

	georel := req.URL.Query().Get("georel")
	if len(georel) > 0 {
		qw.geoQuery, err = newGeoQueryFromHTTPRequest(georel, req)
	}

	return qw, err
}

func newGeoQueryFromHTTPRequest(georel string, req *http.Request) (*GeoQuery, error) {

	if georel == GeoSpatialRelationNearPoint {
		geoQuery := &GeoQuery{Geometry: "Point"}

		if req.URL.Query().Get("geometry") != "Point" {
			return nil, errors.New("The geospatial relationship near is only defined for the geometry type Point")
		}

		distanceString := req.URL.Query().Get("maxDistance")
		if len(distanceString) < 2 || strings.HasPrefix(distanceString, "=") == false {
			return nil, errors.New("Required parameter maxDistance missing or invalid")
		}

		distanceString = distanceString[1:]
		distance, err := strconv.Atoi(distanceString)

		if err != nil {
			return nil, errors.New("Failed to parse distance: " + err.Error())
		}

		if distance < 0 {
			return nil, errors.New("Distance value must be non negative")
		}

		geoQuery.distance = uint32(distance)

		geoQuery.Coordinates, err = parseGeometryCoordinates(req.URL.Query().Get("coordinates"))
		if err != nil {
			return nil, err
		}

		return geoQuery, nil
	}

	return nil, errors.New("Only the geo-spatial relationship \"near\" is supported at this time")
}

func parseGeometryCoordinates(coordparameter string) ([]float64, error) {
	if strings.HasPrefix(coordparameter, "[") == false || strings.HasSuffix(coordparameter, "]") == false {
		return nil, errors.New("Geometry coordinates must be enclosed in []")
	}

	coordparameter = coordparameter[1:]
	coordparameter = coordparameter[0 : len(coordparameter)-1]

	tokens := strings.Split(coordparameter, ",")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("expected two coordinates for a Point geometry, but got %d", len(tokens))
	}

	coordinates := []float64{}

	for idx, c := range tokens {
		f64, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coordinate %d as a decimal value: %s", idx, err.Error())
		}

		coordinates = append(coordinates, f64)
	}

	return coordinates, nil
}

type queryWrapper struct {
	request    *http.Request
	types      []string
	attributes []string
	device     *string

	geoQuery *GeoQuery
}

func (q *queryWrapper) HasDeviceReference() bool {
	return q.device != nil
}

func (q *queryWrapper) IsGeoQuery() bool {
	return q.Geo() != nil
}

func (q *queryWrapper) Geo() *GeoQuery {
	return q.geoQuery
}

func (q *queryWrapper) Device() string {
	return *q.device
}

func (q *queryWrapper) EntityAttributes() []string {
	return q.attributes
}

func (q *queryWrapper) EntityTypes() []string {
	return q.types
}

func (q *queryWrapper) Request() *http.Request {
	return q.request
}
